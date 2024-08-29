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
	"strconv"
	"strings"
	"time"
)

// Undelegate defines a method for performing an undelegation from a delegate and a validator
func (k MsgServer) Undelegate(goCtx context.Context, msg *stakingtypes.MsgUndelegate) (*stakingtypes.MsgUndelegateResponse, error) {
	var region types.Region
	ctx := sdk.UnwrapSDKContext(goCtx)
	meid, isFound := k.GetMeid(ctx, msg.DelegatorAddress)
	if msg.IsMeid {
		if !isFound {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "is-meid is true, you're meid")
		}
		region, isFound = k.GetRegion(ctx, meid.RegionId)
		if !isFound {
			return nil, types.ErrRegionNotExist
		}
	} else {
		region, isFound = k.GetRegion(ctx, strings.ToLower(types.ExperienceRegion))
		if !isFound {
			return nil, types.ErrRegionNotExist
		}
	}
	msg.ValidatorAddress = region.OperatorAddress
	addr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

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

	val, isFound := k.GetValidator(ctx, addr)
	if isFound {
		if val.DelegationAmount.LT(sdk.ZeroInt()) {
			return nil, types.ErrValidatorDelegationAmount.Wrapf("validator amount: %s, requested value: %s",
				val.DelegationAmount.String(), msg.Amount.Amount.String())
		}

	}

	// current interest balance * personal withdrawal pledge limit / district total pledge limit
	//regionTreasureAddr, err := sdk.AccAddressFromBech32(region.RegionTreasureAddr)
	//if err != nil {
	//	return nil, err
	//}
	//person_dele_inte := region.DelegateInterest.Mul(sdk.NewDecFromInt(msg.Amount.Amount).Quo(sdk.NewDecFromInt(validator.DelegationAmount)))
	del := k.Delegation(ctx, delegatorAddress, addr)
	if del == nil {
		return nil, types.ErrEmptyDelegationDistInfo
	}
	delegation, isOK := del.(stakingtypes.Delegation)
	if !isOK {
		panic("Undelegate err:type Delegation assertion failed")
		return nil, types.ErrAssertionFailed
	}
	userTotalStaking := delegation.Amount.Add(delegation.UnMeidAmount).Add(delegation.Unmovable)
	rewards, err := k.CalculateInterest(ctx, userTotalStaking, delegation.StartHeight)
	if err != nil {
		return nil, types.ErrCalculateInterest.Wrap(err.Error())
	}
	if region.DelegateInterest.GTE(rewards) {
		region.DelegateInterest = region.DelegateInterest.Sub(rewards)
	} else {
		return nil, errors.New(fmt.Sprintf("undelegate err,region(%s) total interest not enough.need pay %s,only have %s",
			region.RegionId, rewards.String(), region.DelegateInterest.String()))
	}

	//TODO: send rewards in staking module
	//err = k.bankKeeper.SendCoins(ctx, regionTreasureAddr, delegatorAddress, sdk.NewCoins(sdk.NewCoin(sdk.BaseMEDenom, rewards.TruncateInt())))
	//if err != nil {
	//	return nil, err
	//}

	if msg.IsMeid {
		if delegation.Amount.LT(msg.Amount.Amount) {
			return nil, types.ErrNotEnoughDelegationAmount
		}
	} else {
		if delegation.UnMeidAmount.LT(msg.Amount.Amount) {
			return nil, types.ErrNotEnoughDelegationAmount
		}
	}

	completionTime, returnAmount, err := k.Keeper.Undelegate(ctx, delegatorAddress, addr, msg.IsMeid, msg.Amount.Amount)
	if err != nil {
		return nil, err
	}
	region.DelegateAmount = region.DelegateAmount.Sub(msg.Amount.Amount)
	k.SetRegion(ctx, region)
	val.DelegationAmount = val.DelegationAmount.Sub(msg.Amount.Amount)
	k.SetValidator(ctx, val)

	if msg.Amount.Amount.IsInt64() {
		defer func() {
			telemetry.IncrCounter(1, stakingtypes.ModuleName, "undelegate")
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", msg.Type()},
				float32(returnAmount.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", msg.Amount.Denom)},
			)
		}()
	}
	delegateTreasure := k.AuthKeeper.GetModuleAccount(ctx, stakingtypes.BondedPoolName)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnDelegate,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(types.AttributeKeyRegionId, region.RegionId),
			sdk.NewAttribute(sdk.AttributeKeyAmount, returnAmount.String()+params.BaseDenom),
			sdk.NewAttribute(stakingtypes.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
			sdk.NewAttribute(types.AttributeKeyAmountDelegateInterest, region.DelegateInterest.String()+params.BaseDenom),
			sdk.NewAttribute(stakingtypes.BondedPoolName, delegateTreasure.String()),
			sdk.NewAttribute(types.AttributeKeyRegionTreasure, region.RegionTreasureAddr),
			sdk.NewAttribute(types.AttributeKeyDelegatorAddress, delegatorAddress.String()),
			sdk.NewAttribute(types.AttributeKeyPersonalDelegateInterest, rewards.TruncateInt().String()+params.BaseDenom),
			sdk.NewAttribute(types.AttributeKeyIsMeid, strconv.FormatBool(msg.IsMeid)),
		),
	})

	return &stakingtypes.MsgUndelegateResponse{
		CompletionTime: completionTime,
	}, nil
}
