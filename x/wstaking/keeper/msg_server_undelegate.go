package keeper

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

// Undelegate defines a method for performing an undelegation from a delegate and a validator
func (k MsgServer) Undelegate(goCtx context.Context, msg *stakingtypes.MsgUndelegate) (*stakingtypes.MsgUndelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	regionID := k.GetRegionIdByAccount(ctx, sdk.MustAccAddressFromBech32(msg.DelegatorAddress))
	region, isFound := k.GetRegion(ctx, regionID)
	if !isFound {
		return nil, types.ErrRegionNotExist
	}
	msg.ValidatorAddress = region.OperatorAddress
	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	regionTreasureAddr, err := sdk.AccAddressFromBech32(region.RegionTreasureAddr)
	if err != nil {
		return nil, err
	}
	bondDenom := k.BondDenom(ctx)
	if msg.Amount.Denom != bondDenom {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", msg.Amount.Denom, bondDenom,
		)
	}

	val, isFound := k.GetValidator(ctx, valAddr)
	if isFound {
		if val.DelegationAmount.LT(sdk.ZeroInt()) {
			return nil, types.ErrValidatorDelegationAmount.Wrapf("validator amount: %s, requested value: %s",
				val.DelegationAmount.String(), msg.Amount.Amount.String())
		}
	}

	// current interest balance * personal withdrawal pledge limit / district total pledge limit
	//person_dele_inte := region.DelegateInterest.Mul(sdk.NewDecFromInt(msg.Amount.Amount).Quo(sdk.NewDecFromInt(validator.DelegationAmount)))
	delegation, isOK := k.GetDelegation(ctx, delegatorAddress, val.GetOperator())
	if !isOK {
		return nil, types.ErrEmptyDelegationDistInfo
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

	isMeid := true
	if strings.ToLower(val.Description.RegionID) == strings.ToLower(types.ExperienceRegionName) {
		isMeid = false
	}

	completionTime, returnAmount, err := k.Keeper.Undelegate(ctx, delegatorAddress, valAddr, isMeid, msg.Amount.Amount, delegation)
	if err != nil {
		return nil, err
	}
	region.DelegateAmount = region.DelegateAmount.Sub(returnAmount)
	k.SetRegion(ctx, region)
	val.DelegationAmount = val.DelegationAmount.Sub(returnAmount)
	k.SetValidator(ctx, val)
	//send delegation rewards
	err = k.bankKeeper.Extend().SendCoinsWithTag(ctx, regionTreasureAddr, delegatorAddress, sdk.NewCoins(sdk.NewCoin(params.BaseDenom, rewards.TruncateInt())),
		fmt.Sprintf("Undelegate_SendRewards_%s", region.RegionId),
	)
	if err != nil {
		return nil, err
	}
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
	delegateTreasure := k.authKeeper.GetModuleAccount(ctx, stakingtypes.BondedPoolName)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnDelegate,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(types.AttributeKeyRegionId, region.RegionId),
			sdk.NewAttribute(sdk.AttributeKeyAmount, returnAmount.String()+params.BaseDenom),
			sdk.NewAttribute(stakingtypes.AttributeKeyCompletionTime, completionTime.UTC().Format(time.RFC3339)),
			sdk.NewAttribute(types.AttributeKeyAmountDelegateInterest, region.DelegateInterest.String()+params.BaseDenom),
			sdk.NewAttribute(stakingtypes.BondedPoolName, delegateTreasure.GetAddress().String()),
			sdk.NewAttribute(types.AttributeKeyRegionTreasure, region.RegionTreasureAddr),
			sdk.NewAttribute(types.AttributeKeyDelegatorAddress, delegatorAddress.String()),
			sdk.NewAttribute(types.AttributeKeyPersonalDelegateInterest, rewards.String()+params.BaseDenom),
			sdk.NewAttribute(types.AttributeKeyIsMeid, strconv.FormatBool(isMeid)),
		),
	})

	return &stakingtypes.MsgUndelegateResponse{
		CompletionTime: completionTime,
	}, nil
}
