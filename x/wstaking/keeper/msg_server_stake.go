package keeper

import (
	"context"
	gomath "math"
	"time"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/hashicorp/go-metrics"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

// Stake defines a method for performing a stake of coins from stake_tokens_pool to a validator
func (k MsgServer) Stake(goCtx context.Context, msg *types.MsgStake) (*types.MsgStakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.daoKeeper.IsGlobalDao(ctx, msg.StakerAddress) {
		return nil, types.ErrCheckGlobalDao
	}

	valAddr, valErr := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if valErr != nil {
		return nil, valErr
	}

	validator, err := k.GetValidator(ctx, valAddr)
	if err != nil {
		return nil, stakingtypes.ErrNoValidatorFound
	}

	bondDenom, _ := k.BondDenom(ctx)
	if msg.Amount.Denom != bondDenom {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", msg.Amount.Denom, bondDenom,
		)
	}

	minSelfStake := math.NewInt(int64(gomath.Pow10(params.BaseDenomUnit)))
	if msg.Amount.Amount.Mod(minSelfStake).Int64() != 0 {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin amount: got %d, expected %d integer multiple", msg.Amount.Amount.Int64(), int64(gomath.Pow10(params.BaseDenomUnit)),
		)
	}
	// should before modified region shared
	err = k.WstakingHooks().BeforeValidatorStakingModified(ctx, valAddr)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrHooks, "before stake:%+v", err)
	}

	// NOTE: source funds are always unbonded
	newShares, err := k.Keeper.Stake(ctx, sdk.MustAccAddressFromBech32(msg.StakerAddress), msg.Amount.Amount, stakingtypes.Unbonded, validator, true, "stake_"+validator.Description.RegionID)
	if err != nil {
		return nil, err
	}

	// Update RegionShare after successful stake (must be after Stake to avoid BondRegion overwrite)
	region, found := k.Keeper.GetRegion(ctx, validator.Description.RegionID)
	if found {
		region.RegionShare = region.RegionShare.Add(msg.Amount.Amount)
		k.Keeper.SetRegion(ctx, region)
	}

	if msg.Amount.Amount.IsInt64() {
		defer func() {
			telemetry.IncrCounter(1, types.ModuleName, "stake")
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", msg.Type()},
				float32(msg.Amount.Amount.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", msg.Amount.Denom)},
			)
		}()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeStake,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyNewShares, newShares.String()),
			sdk.NewAttribute(types.AttributeKeyRegionId, validator.Description.RegionID),
		),
	})

	return &types.MsgStakeResponse{}, nil
}

// Unstake defines a method for performing an unstake from a stake and a validator
func (k MsgServer) Unstake(goCtx context.Context, msg *types.MsgUnstake) (*types.MsgUnstakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	stakerAddress, err := sdk.AccAddressFromBech32(msg.StakerAddress)
	if err != nil {
		return nil, err
	}

	if !k.daoKeeper.IsGlobalDao(ctx, msg.StakerAddress) {
		return nil, types.ErrCheckGlobalDao
	}

	addr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	shares, err := k.ValidateUnbondAmount(ctx, stakerAddress, addr, msg.Amount.Amount)
	if err != nil {
		return nil, err
	}

	bondDenom, _ := k.BondDenom(ctx)
	if msg.Amount.Denom != bondDenom {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", msg.Amount.Denom, bondDenom,
		)
	}
	err = k.WstakingHooks().BeforeValidatorStakingModified(ctx, addr)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrHooks, "before unStake :error :%+v", err)
	}

	// Update RegionShare before unstake completes
	unstakeValidator, valErr := k.GetValidator(ctx, addr)
	if valErr == nil {
		if unstakeRegion, foundRegion := k.Keeper.GetRegion(ctx, unstakeValidator.Description.RegionID); foundRegion {
			if unstakeRegion.RegionShare.GTE(msg.Amount.Amount) {
				unstakeRegion.RegionShare = unstakeRegion.RegionShare.Sub(msg.Amount.Amount)
			} else {
				unstakeRegion.RegionShare = math.ZeroInt()
			}
			k.Keeper.SetRegion(ctx, unstakeRegion)
		}
	}

	completionTime, err := k.Keeper.Unstake(ctx, stakerAddress, addr, shares)
	if err != nil {
		return nil, err
	}

	if msg.Amount.Amount.IsInt64() {
		defer func() {
			telemetry.IncrCounter(1, types.ModuleName, "unstake")
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", msg.Type()},
				float32(msg.Amount.Amount.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", msg.Amount.Denom)},
			)
		}()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnstake,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.UTC().Format(time.RFC3339)),
		),
	})

	return &types.MsgUnstakeResponse{
		CompletionTime: completionTime,
	}, nil
}
