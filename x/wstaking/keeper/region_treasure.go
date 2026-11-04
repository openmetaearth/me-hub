package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

// SpendableRegionTreasure returns the coins a region treasury can pay out.
func (k Keeper) SpendableRegionTreasure(ctx sdk.Context, regionID string) (sdk.Coins, error) {
	region, found := k.GetRegion(ctx, regionID)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrRegionNotExist, "region %s", regionID)
	}
	fromAddr, err := sdk.AccAddressFromBech32(region.RegionTreasureAddr)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrUnknownAccount, "region account %s: %s", region.RegionTreasureAddr, err)
	}
	return k.bankKeeper.SpendableCoins(ctx, fromAddr), nil
}

// WithdrawRegionTreasure sends coins from a region treasury to receiver.
func (k Keeper) WithdrawRegionTreasure(ctx sdk.Context, regionID, receiver string, amount sdk.Coins) error {
	region, found := k.GetRegion(ctx, regionID)
	if !found {
		return sdkerrors.Wrapf(types.ErrRegionNotExist, "region %s", regionID)
	}

	fromAddr, err := sdk.AccAddressFromBech32(region.RegionTreasureAddr)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrUnknownAccount, "region account %s: %s", region.RegionTreasureAddr, err)
	}
	toAddr, err := sdk.AccAddressFromBech32(receiver)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrUnknownAccount, "receiver account %s: %s", receiver, err)
	}
	if k.bankKeeper.Extend().BlockedAddr(toAddr) {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive region treasury withdrawals", toAddr)
	}

	err = k.bankKeeper.Extend().SendCoinsWithTag(
		ctx,
		fromAddr,
		toAddr,
		amount,
		fmt.Sprintf("RegionGovWithdraw_%s", region.RegionId),
	)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "region treasure %s does not have enough balance", region.RegionTreasureAddr)
	}
	return nil
}

// SetRegionOperator writes Region.OperatorAddress and updates the region group admin.
func (k Keeper) SetRegionOperator(ctx sdk.Context, regionID, newOperator string) error {
	region, found := k.GetRegion(ctx, regionID)
	if !found {
		return sdkerrors.Wrapf(types.ErrRegionNotExist, "region %s", regionID)
	}
	valAddr, err := sdk.ValAddressFromBech32(newOperator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid operator %s: %s", newOperator, err)
	}
	validator, ok := k.GetValidator(ctx, valAddr)
	if !ok {
		return sdkerrors.Wrapf(types.ErrNoValidatorExists, "%s", newOperator)
	}

	region.OperatorAddress = validator.OperatorAddress
	k.SetRegion(ctx, region)
	if k.groupKeeper != nil {
		k.groupKeeper.UpdateGroupAdmin(ctx, regionID, validator.OwnerAddress)
	}
	return nil
}
