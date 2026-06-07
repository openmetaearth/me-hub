package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

// SetRegionWithdraw stores the authorized address for a region.
// Calling this again with a different address overwrites the previous grant.
func (k Keeper) SetRegionWithdraw(ctx sdk.Context, regionId, address string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RegionWithdrawKeyPrefix))
	store.Set([]byte(regionId), []byte(address))
}

// GetRegionWithdraw returns the authorized address for a region, if any.
func (k Keeper) GetRegionWithdraw(ctx sdk.Context, regionId string) (address string, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RegionWithdrawKeyPrefix))
	bz := store.Get([]byte(regionId))
	if bz == nil {
		return "", false
	}
	return string(bz), true
}

// DeleteRegionWithdraw removes the withdrawer for a region.
func (k Keeper) DeleteRegionWithdraw(ctx sdk.Context, regionId string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RegionWithdrawKeyPrefix))
	store.Delete([]byte(regionId))
}

// GetAllRegionWithdrawGrants returns every region withdraw grant in store.
func (k Keeper) GetAllRegionWithdrawGrants(ctx sdk.Context) (grants []types.RegionWithdrawGrant) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RegionWithdrawKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		grants = append(grants, types.RegionWithdrawGrant{
			RegionId: string(iterator.Key()),
			Address:  string(iterator.Value()),
		})
	}

	return grants
}

// CanRegionWithdraw returns true if address is the authorized
// withdrawer for regionId.
func (k Keeper) CanRegionWithdraw(ctx sdk.Context, address, regionId string) bool {
	granted, found := k.GetRegionWithdraw(ctx, regionId)
	if !found {
		return false
	}
	return granted == address
}
