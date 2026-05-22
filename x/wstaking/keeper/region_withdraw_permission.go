package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

// SetRegionWithdrawPermission stores the authorized address for a region.
// Calling this again with a different address overwrites the previous grant.
func (k Keeper) SetRegionWithdrawPermission(ctx sdk.Context, regionId, address string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RegionWithdrawKeyPrefix))
	store.Set([]byte(regionId), []byte(address))
}

// GetRegionWithdrawPermission returns the authorized address for a region, if any.
func (k Keeper) GetRegionWithdrawPermission(ctx sdk.Context, regionId string) (address string, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RegionWithdrawKeyPrefix))
	bz := store.Get([]byte(regionId))
	if bz == nil {
		return "", false
	}
	return string(bz), true
}

// DeleteRegionWithdrawPermission removes the withdraw permission for a region.
func (k Keeper) DeleteRegionWithdrawPermission(ctx sdk.Context, regionId string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RegionWithdrawKeyPrefix))
	store.Delete([]byte(regionId))
}

// HasRegionWithdrawPermission returns true if address is the authorized
// withdrawer for regionId.
func (k Keeper) HasRegionWithdrawPermission(ctx sdk.Context, address, regionId string) bool {
	granted, found := k.GetRegionWithdrawPermission(ctx, regionId)
	if !found {
		return false
	}
	return granted == address
}
