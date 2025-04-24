package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

func (k *Keeper) DeleteApproveReward(ctx sdk.Context, address, regionId string) error {
	return k.stkKeeper.RemoveKycReward(ctx, sdk.MustAccAddressFromBech32(address), regionId)
}

func (k *Keeper) TransferKycRegion(ctx sdk.Context, address, issuer, fromRegionId, toRegionId string) error {
	if strings.EqualFold(fromRegionId, toRegionId) {
		return nil
	}
	return k.stkKeeper.TransferKycRegion(ctx, sdk.MustAccAddressFromBech32(address), issuer, fromRegionId, toRegionId)
}
