package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/dao/types"
)

func (k Keeper) SetFreeGasAccount(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.FreeGasAddressePrefix, []byte(address)...)
	store.Set(key, []byte(address))

	acc := sdk.MustAccAddressFromBech32(address)
	if has := k.authKeeper.HasAccount(ctx, acc); !has {
		newAccount := k.authKeeper.NewAccountWithAddress(ctx, acc)
		k.authKeeper.SetAccount(ctx, newAccount)
	}
}

func (k Keeper) RemoveFreeGasAccount(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.FreeGasAddressePrefix, []byte(address)...)
	store.Delete(key)
}

func (k Keeper) CheckFreeGasAccount(ctx sdk.Context, address string) bool {
	store := ctx.KVStore(k.storeKey)
	key := append(types.FreeGasAddressePrefix, []byte(address)...)
	value := store.Get(key)
	if len(value) == 0 {
		return false
	}
	return true
}

func (k Keeper) IterateFreeGasAccounts(ctx sdk.Context, cb func(address string) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.FreeGasAddressePrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		if cb(string(iterator.Value())) {
			break
		}
	}
}
