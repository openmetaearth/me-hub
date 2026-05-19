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
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"time"
)

// Delegate defines a method for performing a delegation of coins from a delegator to a validator
func (k MsgServer) Delegate(goCtx context.Context, msg *stakingtypes.MsgDelegate) (*stakingtypes.MsgDelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	regionId := k.GetRegionIdByAccount(ctx, sdk.MustAccAddressFromBech32(msg.DelegatorAddress))
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

	delegation, isOK := k.GetDelegation(ctx, delegatorAddress, validator.GetOperator())
	rewards := sdk.ZeroDec()
	var regionTreasureAddr sdk.AccAddress
	if isOK {
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
		err = k.bankKeeper.Extend().SendCoinsWithTag(ctx, regionTreasureAddr, delegatorAddress, sdk.NewCoins(sdk.NewCoin(params.BaseDenom, rewards.TruncateInt())),
			fmt.Sprintf("Delegate_SendRewards_%s", region.RegionId),
		)
		if err != nil {
			return nil, err
		}
	}

	// NOTE: source funds are always UnBonded
	newShares, err := k.Keeper.Delegate(ctx, delegatorAddress, msg.Amount.Amount, stakingtypes.Unbonded, validator, delegation, valAddr)
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
