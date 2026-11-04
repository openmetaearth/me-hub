package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

// WithdrawRegionTreasure sends coins from a region treasury to receiver.
// When maxPercent > 0, each denom must satisfy amount*100 <= spendable*maxPercent.
func (k Keeper) WithdrawRegionTreasure(ctx sdk.Context, regionID, receiver string, amount sdk.Coins, maxPercent uint32) error {
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

	if maxPercent > 0 {
		spendable := k.bankKeeper.SpendableCoins(ctx, fromAddr)
		for _, coin := range amount {
			balance := spendable.AmountOf(coin.Denom)
			if coin.Amount.MulRaw(100).GT(balance.MulRaw(int64(maxPercent))) {
				return sdkerrors.Wrapf(
					sdkerrors.ErrInsufficientFunds,
					"equity withdrawal %s exceeds %d%% of region %s treasury balance %s",
					coin, maxPercent, regionID, sdk.NewCoin(coin.Denom, balance),
				)
			}
		}
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
