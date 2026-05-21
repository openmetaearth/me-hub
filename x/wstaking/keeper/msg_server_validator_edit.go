package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/utils"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"strings"
)

// UpdateValidator defines a method for editing an existing validator
func (k MsgServer) UpdateValidator(goCtx context.Context, msg *types.MsgUpdateValidator) (*types.MsgUpdateValidatorResponse, error) {
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
	oldRegionId := validator.Description.RegionID
	// replace all editable fields (clients should autofill existing values)
	description, err := validator.Description.UpdateDescription(msg.Description)
	if err != nil {
		return nil, err
	}
	if msg.Description.RegionID == stakingtypes.DoNotModifyDesc {
		description.RegionID = oldRegionId
		validator.Description = description
	} else {
		if _, err := utils.CheckRegionName(strings.ToUpper(msg.Description.RegionID)); err != nil {
			return nil, types.ErrRegionName
		}
		// remove duplication
		validators := k.GetAllValidators(ctx)
		for _, v := range validators {
			if v.Description.RegionID == msg.Description.RegionID {
				return nil, types.ErrValidatorRegionDuplication
			}
		}
		k.UnBondRegion(ctx, oldRegionId)
		description.RegionID = msg.Description.RegionID
		validator.Description = description
		region, f := k.GetRegion(ctx, msg.Description.RegionID)
		if f && region.OperatorAddress == "" {
			k.BondRegion(ctx, validator, validator.Tokens, true)
		}
	}

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
			types.EventTypeUpdateValidator,
			sdk.NewAttribute(stakingtypes.AttributeKeyCommissionRate, validator.Commission.String()),
			sdk.NewAttribute(stakingtypes.AttributeKeyMinSelfDelegation, validator.MinSelfDelegation.String()),
			sdk.NewAttribute(types.AttributeKeyRegionId, validator.Description.RegionID),
		),
	})
	return &types.MsgUpdateValidatorResponse{}, nil
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

	ctx.Logger().Info("==>old validator", "old validator", validator.String(), "old owner", validator.OwnerAddress)

	stake, found := k.GetStake(ctx, staker, validator.GetOperator())
	if !found {
		return sdkerrors.Wrapf(types.ErrNoStake, "stake(%s) for operator(%s) not found", staker, validator.GetOperator())
	}
	k.RemoveStake(ctx, stake)

	k.RemoveValidator(ctx, validator.GetOperator())
	k.DeleteLastValidatorPower(ctx, validator.GetOperator())
	if validator.Status == stakingtypes.Unbonding {
		k.DeleteValidatorQueue(ctx, validator)
	}

	stake.ValidatorAddress = newValOperAddr.String()
	k.SetStake(ctx, stake)

	validator.OperatorAddress = newValOperAddr.String()
	validator.OwnerAddress = newValAddr.String()

	err := k.SetValidatorByConsAddr(ctx, validator)
	if err != nil {
		return err
	}

	k.SetValidator(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)
	// bond region again
	k.BondRegion(ctx, validator, validator.Tokens, true)

	if validator.Status == stakingtypes.Unbonding {
		k.InsertUnbondingValidatorQueue(ctx, validator)
	}

	k.SetChangeDelegationValidator(ctx, validator.Description.RegionID)

	ctx.Logger().Info("==>new validator", "validator", validator.String(), "owner", validator.OwnerAddress)
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
