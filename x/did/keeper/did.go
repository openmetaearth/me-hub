package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/did/types"
)

func (k Keeper) HasDID(ctx sdk.Context, addr sdk.AccAddress) bool {
	if _, found := k.GetDID(ctx, addr); found {
		return true
	}
	return false
}

func (k Keeper) GetDID(ctx sdk.Context, addr sdk.AccAddress) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetDIDKey(addr))
	if bz == nil {
		return "", false
	}

	return string(bz), true
}

func (k Keeper) SetDID(ctx sdk.Context, addr sdk.AccAddress, did string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetDIDKey(addr), []byte(did))
}

func (k Keeper) DeleteDID(ctx sdk.Context, addr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDIDKey(addr))
}
