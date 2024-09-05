package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/did/types"
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

func (k Keeper) SetService(ctx sdk.Context, sid string, svc types.Service) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetServiceKey(sid), k.cdc.MustMarshal(&svc))
}

func (k Keeper) DeleteService(ctx sdk.Context, sid string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetServiceKey(sid))
}
