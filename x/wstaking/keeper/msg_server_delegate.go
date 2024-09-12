package keeper

import (
	"context"
	"errors"
	"fmt"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/x/wstaking/types"
	"strings"
	"time"
)

// Delegate defines a method for performing a delegation of coins from a delegator to a validator
func (k MsgServer) Delegate(goCtx context.Context, msg *stakingtypes.MsgDelegate) (*stakingtypes.MsgDelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	did, ok := k.KycKeeper.GetDID(ctx, sdk.MustAccAddressFromBech32(msg.DelegatorAddress))
	if !ok {
		return k.UnMeidDelegate(goCtx, msg)
	} else {
		return k.MeidDelegate(goCtx, msg, did)
	}
}

// MeidDelegate defines a method for performing a delegation of coins from a KYC to a validator
func (k MsgServer) MeidDelegate(goCtx context.Context, msg *stakingtypes.MsgDelegate, did string) (*stakingtypes.MsgDelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	kycData, _ := k.KycKeeper.GetKYC(ctx, did)
	regionId := string(kycData.Data)
	region, isFound := k.GetRegion(ctx, regionId)
	if !isFound {
		return nil, types.ErrRegionNotExist
	}
	msg.ValidatorAddress = region.OperatorAddress
	valAddr, valErr := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if valErr != nil {
		return nil, valErr
	}
	err := types.CheckMinDelegate(msg.Amount.Amount)
	if err != nil {
		return nil, err
	}
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return nil, stakingtypes.ErrNoValidatorFound
	}

	validator.DelegationAmount = validator.DelegationAmount.Add(msg.Amount.Amount)
	if validator.Tokens.LT(validator.DelegationAmount) {
		return nil, types.ErrNodeLimitExceeded
	}

	region.DelegateAmount = region.DelegateAmount.Add(msg.Amount.Amount)
	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	bondDenom := k.BondDenom(ctx)
	if msg.Amount.Denom != bondDenom {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", msg.Amount.Denom, bondDenom,
		)
	}

	del := k.Delegation(ctx, delegatorAddress, valAddr)
	rewards := sdk.ZeroDec()
	var regionTreasureAddr sdk.AccAddress
	if del != nil {
		delegation, isOK := del.(stakingtypes.Delegation)
		if !isOK {
			panic("withdrawDelegationRewards err:type Delegation assertion failed")
			return nil, types.ErrAssertionFailed
		}
		rewards, err = k.CalculateInterest(ctx, delegation.Amount.Add(delegation.UnMeidAmount).Add(delegation.Unmovable), delegation.StartHeight)
		if err != nil {
			return nil, types.ErrCalculateInterest.Wrap(err.Error())
		}
		regionTreasureAddr, err = sdk.AccAddressFromBech32(region.RegionTreasureAddr)
		if err != nil {
			return nil, err
		}
		if region.DelegateInterest.GTE(rewards) {
			region.DelegateInterest = region.DelegateInterest.Sub(rewards)
		} else {
			return nil, errors.New(fmt.Sprintf("region(%s) total interest not enough.need pay %s,only have %s",
				region.RegionId, rewards.String(), region.DelegateInterest.String()))
		}
		err = k.BankKeeper.SendCoins(ctx, regionTreasureAddr, delegatorAddress, sdk.NewCoins(sdk.NewCoin(params.BaseDenom, rewards.TruncateInt())))
		if err != nil {
			return nil, err
		}
	}

	// NOTE: source funds are always UnBonded
	newShares, err := k.Keeper.Delegate(ctx, delegatorAddress, msg.Amount.Amount, stakingtypes.Unbonded, validator, true)
	if err != nil {
		return nil, err
	}
	k.SetRegion(ctx, region)
	k.SetValidator(ctx, validator)

	if msg.Amount.Amount.IsInt64() {
		defer func() {
			telemetry.IncrCounter(1, types.ModuleName, "delegate")
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", msg.Type()},
				float32(msg.Amount.Amount.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", msg.Amount.Denom)},
			)
		}()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingtypes.EventTypeDelegate,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(types.AttributeKeyRegionId, region.RegionId),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(stakingtypes.AttributeKeyCompletionTime, ctx.BlockHeader().Time.Format(time.RFC3339)),
			sdk.NewAttribute(types.AttributeKeyRegionTreasure, regionTreasureAddr.String()),
			sdk.NewAttribute(types.AttributeKeyNewShares, newShares.String()),
			sdk.NewAttribute(types.AttributeKeyTotalAmountDelegate, validator.DelegationAmount.String()+params.BaseDenom),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress),
			sdk.NewAttribute(types.AttributeKeyRewards, rewards.TruncateInt().String()+params.BaseDenom),
		),
	})

	return &stakingtypes.MsgDelegateResponse{}, nil
}

func (k MsgServer) UnMeidDelegate(goCtx context.Context, msg *stakingtypes.MsgDelegate) (*stakingtypes.MsgDelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	region, isFound := k.GetRegion(ctx, strings.ToLower(types.ExperienceRegion))
	if !isFound {
		return nil, types.ErrRegionNotExist
	}
	msg.ValidatorAddress = region.OperatorAddress

	valAddr, valErr := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if valErr != nil {
		return nil, valErr
	}
	err := types.CheckMinDelegate(msg.Amount.Amount)
	if err != nil {
		return nil, err
	}

	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return nil, stakingtypes.ErrNoValidatorFound
	}

	if validator.Tokens.LT(validator.DelegationAmount) {
		return nil, types.ErrNodeLimitExceeded
	}
	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	validator.DelegationAmount = validator.DelegationAmount.Add(msg.Amount.Amount)
	region.DelegateAmount = region.DelegateAmount.Add(msg.Amount.Amount)

	bondDenom := k.BondDenom(ctx)
	if msg.Amount.Denom != bondDenom {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", msg.Amount.Denom, bondDenom,
		)
	}

	del := k.Delegation(ctx, delegatorAddress, sdk.ValAddress{})
	rewards := sdk.ZeroDec()
	var regionTreasureAddr sdk.AccAddress
	if del != nil {
		delegation, isOK := del.(stakingtypes.Delegation)
		if !isOK {
			return nil, types.ErrAssertionFailed
		}
		rewards, err = k.CalculateInterest(ctx, delegation.Amount.Add(delegation.UnMeidAmount).Add(delegation.Unmovable), delegation.StartHeight)
		if err != nil {
			return nil, types.ErrCalculateInterest.Wrap(err.Error())
		}
		regionTreasureAddr, err = sdk.AccAddressFromBech32(region.RegionTreasureAddr)
		if err != nil {
			return nil, err
		}
		err = k.BankKeeper.SendCoins(ctx, regionTreasureAddr, delegatorAddress, sdk.NewCoins(sdk.NewCoin(params.BaseDenom, rewards.TruncateInt())))
		if err != nil {
			return nil, err
		}
		if region.DelegateInterest.GTE(rewards) {
			region.DelegateInterest = region.DelegateInterest.Sub(rewards)
		}
	}

	// NOTE: source funds are always unbonded
	newShares, err := k.Keeper.UnMeidDelegate(ctx, delegatorAddress, msg.Amount.Amount, validator)
	if err != nil {
		return nil, err
	}
	k.SetRegion(ctx, region)
	k.SetValidator(ctx, validator)

	if msg.Amount.Amount.IsInt64() {
		defer func() {
			telemetry.IncrCounter(1, types.ModuleName, "delegate")
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", msg.Type()},
				float32(msg.Amount.Amount.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", msg.Amount.Denom)},
			)
		}()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnMeidDelegate,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyRewards, rewards.TruncateInt().String()+params.BaseDenom),
			sdk.NewAttribute(types.AttributeKeyRegionTreasure, regionTreasureAddr.String()),
			sdk.NewAttribute(types.AttributeKeyNewShares, newShares.String()),
			sdk.NewAttribute(types.AttributeKeyTotalAmountDelegate, validator.DelegationAmount.String()+params.BaseDenom),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress),
			sdk.NewAttribute(types.AttributeKeyRegionId, region.RegionId),
		),
	})
	return &stakingtypes.MsgDelegateResponse{}, nil
}
