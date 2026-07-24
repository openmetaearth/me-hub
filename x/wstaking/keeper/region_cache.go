package keeper

import (
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (k *Keeper) InitRegionCache(ctx sdk.Context) {
	k.regionCacheOnce.Do(func() {
		regions := k.GetAllRegion(ctx)
		k.SetRegionsCache(ctx, regions)
		k.Logger(ctx).Info("region cache initialized", "count", len(regions))
	})
}

/*
SetRegionCache only use immutable fields changes
Usage Rules:
 1. the risk is come from simulate tx, if region cache changed by simulate, nodes will have different state and appHash
 2. do not use cache to calculate in tx, only use immutable and non-calculated fields.
 3. be ware of tx that may change region cache, such as SetRegion, RemoveRegion, GetAllRegion.
*/
func (k *Keeper) SetRegionsCache(ctx sdk.Context, regions []types.Region) {
	// Only update the cache during DeliverTx (not during simulation)
	if !ctx.IsCheckTx() {
		k.regions = new(sync.Map)
		for _, region := range regions {
			k.regions.Store(region.RegionId, region)
		}
		k.Logger(ctx).Debug("SetRegionsCache", "count", len(regions))
	}
}

// GetRegionsCache only use for immutable fields
func (k *Keeper) GetRegionsCache() map[string]types.Region {
	regionsMap := make(map[string]types.Region)
	k.regions.Range(func(key, value interface{}) bool {
		regionId, ok := key.(string)
		if !ok {
			return true
		}
		region, ok := value.(types.Region)
		if !ok {
			return true
		}
		regionsMap[regionId] = region
		return true
	})
	return regionsMap
}

// GetRegionCache only use for immutable fields in one block, if there are many txs in a block to get and set region,
// assure just use unchangeable region fields, you cannot use region cache to calculate then reset region in tx for any reason.
func (k *Keeper) GetRegionCache(regionId string) (types.Region, bool) {
	value, ok := k.regions.Load(regionId)
	if !ok {
		return types.Region{}, false
	}
	region, ok := value.(types.Region)
	return region, ok
}
