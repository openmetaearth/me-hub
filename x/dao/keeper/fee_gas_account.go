package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/dao/types"
)

func (k Keeper) SetFreeGasAccount(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.FreeGasAddressePrefix, []byte(address)...)
	store.Set(key, []byte(address))
}

func (k Keeper) RemoveFreeGasAccount(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.FreeGasAddressePrefix, []byte(address)...)
	store.Delete(key)
}

func (k Keeper) IsFreeGasAccount(ctx sdk.Context, address string) bool {
	store := ctx.KVStore(k.storeKey)
	key := append(types.FreeGasAddressePrefix, []byte(address)...)
	value := store.Get(key)
	if len(value) == 0 {
		return false
	}
	return true
}
