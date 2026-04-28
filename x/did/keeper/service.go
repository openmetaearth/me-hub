package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/did/types"
)

func (k Keeper) GetService(ctx sdk.Context, sid string) (svc types.Service, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetServiceKey(sid))
	if bz == nil {
		return types.Service{}, false
	}

	k.cdc.MustUnmarshal(bz, &svc)
	return svc, true
}

func (k Keeper) GetServices(ctx sdk.Context) (svcs []types.Service) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ServicePrefix)
	defer iterator.Close() // nolint: errcheck

	for ; iterator.Valid(); iterator.Next() {
		var svc types.Service
		k.cdc.MustUnmarshal(iterator.Value(), &svc)
		svcs = append(svcs, svc)
	}

	return svcs
}

func (k Keeper) SetService(ctx sdk.Context, sid string, svc types.Service) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetServiceKey(sid), k.cdc.MustMarshal(&svc))
}

func (k Keeper) DeleteService(ctx sdk.Context, sid string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetServiceKey(sid))
}
