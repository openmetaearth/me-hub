package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/did/types"
)

func (k Keeper) HasDidBySubAccount(ctx sdk.Context, subAccount string) bool {
	if _, found := k.GetSubAccountDidMap(ctx, subAccount); found {
		return true
	}
	return false
}

func (k Keeper) GetSubAccountDidMap(ctx sdk.Context, subAccount string) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetSubAccountKey(subAccount))
	if bz == nil {
		return "", false
	}

	return string(bz), true
}

func (k Keeper) SetSubAccountDidMap(ctx sdk.Context, subAccount, did string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetSubAccountKey(subAccount), []byte(did))
}

func (k Keeper) DeleteSubAccountDidMap(ctx sdk.Context, subAccount string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetSubAccountKey(subAccount))
}
