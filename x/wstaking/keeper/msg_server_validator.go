package keeper

import (
	"context"
	"strings"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/st-chain/me-hub/utils"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

// CreateValidator defines wrapped method for creating a new validator.
func (k MsgServer) CreateValidator(
	goCtx context.Context, msg *stakingtypes.MsgCreateValidator,
) (*stakingtypes.MsgCreateValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.daoKeeper.IsGlobalDao(ctx, msg.DelegatorAddress) {
		return nil, types.ErrCheckGlobalDao
	}

	_, err := utils.CheckRegionName(strings.ToUpper(msg.Description.RegionID))
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrRegionName, msg.Description.RegionID)
	}
	msg.Description.RegionID = strings.ToLower(msg.Description.RegionID)

	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	if msg.Commission.Rate.LT(k.MinCommissionRate(ctx)) {
		return nil, sdkerrors.Wrapf(stakingtypes.ErrCommissionLTMinRate, "cannot set validator commission to less than minimum rate of %s", k.MinCommissionRate(ctx))
	}

	// check to see if the pubkey or sender has been registered before
	if _, found := k.GetValidator(ctx, valAddr); found {
		return nil, stakingtypes.ErrValidatorOwnerExists
	}

	pk, ok := msg.Pubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Expecting cryptotypes.PubKey, got %T", pk)
	}

	if _, found := k.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk)); found {
		return nil, stakingtypes.ErrValidatorPubKeyExists
	}

	bondDenom := k.BondDenom(ctx)
	if msg.Value.Denom != bondDenom {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", msg.Value.Denom, bondDenom,
		)
	}

	if _, err := msg.Description.EnsureLength(); err != nil {
		return nil, err
	}

	cp := ctx.ConsensusParams()
	if cp != nil && cp.Validator != nil {
		pkType := pk.Type()
		hasKeyType := false
		for _, keyType := range cp.Validator.PubKeyTypes {
			if pkType == keyType {
				hasKeyType = true
				break
			}
		}
		if !hasKeyType {
			return nil, sdkerrors.Wrapf(
				stakingtypes.ErrValidatorPubKeyTypeNotSupported,
				"got: %s, expected: %s", pk.Type(), cp.Validator.PubKeyTypes,
			)
		}
	}

	validator, err := stakingtypes.NewValidator(valAddr, pk, msg.Description)
	if err != nil {
		return nil, err
	}

	commission := stakingtypes.NewCommissionWithTime(
		msg.Commission.Rate, msg.Commission.MaxRate,
		msg.Commission.MaxChangeRate, ctx.BlockHeader().Time,
	)

	validator, err = validator.SetInitialCommission(commission)
	if err != nil {
		return nil, err
	}

	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	validator.MinSelfDelegation = msg.MinSelfDelegation
	validator.OwnerAddress = sdk.AccAddress(valAddr).String()

	k.SetValidator(ctx, validator)
	k.SetValidatorByConsAddr(ctx, validator)
	k.SetNewValidatorByPowerIndex(ctx, validator)

	// call the after-creation hook
	if err := k.Hooks().AfterValidatorCreated(ctx, validator.GetOperator()); err != nil {
		return nil, err
	}

	// move coins from the msg.Address account to a (self-delegation) delegator account
	// the validator account and global shares are updated within here
	// NOTE source will always be from a wallet which are unbonded
	_, err = k.Keeper.Stake(ctx, delegatorAddress, msg.Value.Amount, stakingtypes.Unbonded, validator, true)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingtypes.EventTypeCreateValidator,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Value.String()),
		),
	})

	return &stakingtypes.MsgCreateValidatorResponse{}, nil
}

// EditValidator defines a method for editing an existing validator
func (k MsgServer) EditValidator(goCtx context.Context, msg *stakingtypes.MsgEditValidator) (*stakingtypes.MsgEditValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.daoKeeper.IsGlobalDao(ctx, msg.StakerAddress) {
		return nil, types.ErrCheckGlobalDao
	}

	valAddr, err := sdk.ValAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return nil, stakingtypes.ErrNoValidatorFound
	}

	// replace all editable fields (clients should autofill existing values)
	description, err := validator.Description.UpdateDescription(msg.Description)
	if err != nil {
		return nil, err
	}
	if msg.Description.RegionID == stakingtypes.DoNotModifyDesc {
		description.RegionID = validator.Description.RegionID
	} else {
		description.RegionID = msg.Description.RegionID
	}
	validator.Description = description

	//region, f := k.GetRegion(ctx, validator.Description.RegionID)
	//if !f {
	//	return nil, sdkerrors.Wrapf(types.ErrRegionNotExist, "please set region first")
	//}
	//if region.OperatorAddress != validator.OperatorAddress {
	//	return nil, fmt.Errorf("region id already bound to another validator(%s), please set region first", region.OperatorAddress)
	//}

	if msg.CommissionRate != nil {
		commission, err := k.UpdateValidatorCommission(ctx, validator, *msg.CommissionRate)
		if err != nil {
			return nil, err
		}

		// call the before-modification hook since we're about to update the commission
		if err := k.Hooks().BeforeValidatorModified(ctx, valAddr); err != nil {
			return nil, err
		}

		validator.Commission = commission
	}

	if msg.OwnerAddress != "" {
		err = k.resetValidator(ctx,
			sdk.MustAccAddressFromBech32(msg.StakerAddress),
			sdk.MustAccAddressFromBech32(msg.OwnerAddress),
			validator)
		if err != nil {
			return nil, sdkerrors.Wrapf(types.ErrResetValidator, err.Error())
		}
	} else {
		k.SetValidator(ctx, validator)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingtypes.EventTypeEditValidator,
			sdk.NewAttribute(stakingtypes.AttributeKeyCommissionRate, validator.Commission.String()),
			sdk.NewAttribute(stakingtypes.AttributeKeyMinSelfDelegation, validator.MinSelfDelegation.String()),
		),
	})
	return &stakingtypes.MsgEditValidatorResponse{}, nil
}

func (k Keeper) resetValidator(goCtx context.Context, staker, newValAddr sdk.AccAddress, validator stakingtypes.Validator) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	oldValOperator := validator.GetOperator()

	acc := k.authKeeper.GetAccount(ctx, newValAddr)
	if acc != nil {
		_, ok := acc.(authtypes.ModuleAccountI)
		if ok {
			return types.ErrValidatorAddress
		}
	}

	newValOperAddr := sdk.ValAddress(newValAddr)
	_, exist := k.GetValidator(ctx, newValOperAddr)
	if exist {
		return types.ErrValidatorExist
	}

	ctx.Logger().Info("==>old validator", "old validator", oldValOperator, "old owner", validator.OwnerAddress)

	stake, found := k.GetStake(ctx, staker, validator.GetOperator())
	if !found {
		return sdkerrors.Wrapf(types.ErrNoStake, "stake(%s) for operator(%s) not found", staker, validator.GetOperator())
	}

	k.RemoveValidator(ctx, validator.GetOperator())
	k.DeleteLastValidatorPower(ctx, validator.GetOperator())
	if validator.Status == stakingtypes.Unbonding || validator.UnbondingHeight > 0 {
		k.DeleteValidatorQueue(ctx, validator)
	}

	stake.ValidatorAddress = newValOperAddr.String()
	k.SetStake(ctx, stake)

	k.IterateAllDelegations(ctx, func(delegation stakingtypes.Delegation) bool {
		if delegation.ValidatorAddress == validator.OperatorAddress {
			delegation.ValidatorAddress = newValOperAddr.String()
			k.SetDelegation(ctx, delegation)
		}
		return false
	})

	region, isFound := k.GetRegion(ctx, validator.Description.RegionID)
	if !isFound {
		return sdkerrors.Wrapf(types.ErrRegion, "region id(%s) not found", validator.Description.RegionID)
	}

	validator.OperatorAddress = newValOperAddr.String()
	validator.OwnerAddress = newValAddr.String()
	region.OperatorAddress = newValOperAddr.String()
	err := k.SetValidatorByConsAddr(ctx, validator)
	if err != nil {
		return err
	}

	k.SetValidator(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)
	k.SetRegion(ctx, region)

	if validator.Status == stakingtypes.Unbonding || validator.UnbondingHeight > 0 {
		k.InsertUnbondingValidatorQueue(ctx, validator)
	}

	ctx.Logger().Info("==>new validator", "validator", validator.OperatorAddress, "owner", validator.OwnerAddress)
	if err := k.Hooks().AfterValidatorCreated(ctx, validator.GetOperator()); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeResetValidator,
			sdk.NewAttribute(types.AttributeKeyValidator, oldValOperator.String()),
			sdk.NewAttribute(types.AttributeKeyNewValidator, newValOperAddr.String()),
			sdk.NewAttribute(types.AttributeKeyNewOwnerAddress, validator.OwnerAddress),
		),
	})
	return nil
}
