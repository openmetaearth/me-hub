package keeper

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/dao/types"
)

func (k Keeper) SetFreeGasAccount(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(freeGasAccountKey(address), []byte(address))

	acc := sdk.MustAccAddressFromBech32(address)
	if has := k.authKeeper.HasAccount(ctx, acc); !has {
		newAccount := k.authKeeper.NewAccountWithAddress(ctx, acc)
		k.authKeeper.SetAccount(ctx, newAccount)
	}
}

func (k Keeper) RemoveFreeGasAccount(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(freeGasAccountKey(address))
}

func (k Keeper) CheckFreeGasAccount(ctx sdk.Context, address string) bool {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(freeGasAccountKey(address))
	if len(value) == 0 {
		return false
	}
	return true
}

func (k Keeper) IterateFreeGasAccounts(ctx sdk.Context, cb func(address string) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.FreeGasAddressPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		address := string(bytes.TrimPrefix(iterator.Key(), types.FreeGasAddressPrefix))
		if cb(address) {
			break
		}
	}
}

func freeGasAccountKey(address string) []byte {
	key := append([]byte(nil), types.FreeGasAddressPrefix...)
	return append(key, []byte(address)...)
}
