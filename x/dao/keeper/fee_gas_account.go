package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
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

func (k Keeper) GetAllFreeGasAccounts(ctx sdk.Context) (list []string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FreeGasAddressePrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close() // nolint: errcheck

	for ; iterator.Valid(); iterator.Next() {
		list = append(list, string(iterator.Value()))
	}

	return list
}
