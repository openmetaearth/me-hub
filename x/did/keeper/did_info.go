package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/did/types"
)

func (k Keeper) HasDidInfo(ctx sdk.Context, did string) bool {
	if _, found := k.GetDidInfo(ctx, did); found {
		return true
	}
	return false
}

func (k Keeper) GetDidInfo(ctx sdk.Context, did string) (info types.DidInfo, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetDidInfoKey(did))
	if bz == nil {
		return types.DidInfo{}, false
	}

	k.cdc.MustUnmarshal(bz, &info)
	return info, true
}

func (k Keeper) GetDidInfos(ctx sdk.Context) (infos []types.DidInfo) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DidInfoPrefix)
	defer iterator.Close() // nolint: errcheck

	for ; iterator.Valid(); iterator.Next() {
		var info types.DidInfo
		k.cdc.MustUnmarshal(iterator.Value(), &info)
		infos = append(infos, info)
	}

	return infos
}

func (k Keeper) SetDidInfo(ctx sdk.Context, did string, info types.DidInfo) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetDidInfoKey(did), k.cdc.MustMarshal(&info))
}

func (k Keeper) DeleteDidInfo(ctx sdk.Context, did string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDidInfoKey(did))
}
