package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/st-chain/me-hub/x/gravity/types"
)

func (k Keeper) GetBridgeTokenByContract(ctx sdk.Context, tokenContract string) (bridgeToken *types.BridgeToken, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBridgeTokenByContractKey(tokenContract))
	if len(bz) == 0 {
		return nil, types.ErrNotFound
	}
	k.cdc.MustUnmarshal(bz, bridgeToken)
	return
}

func (k Keeper) GetBridgeTokenByDenom(ctx sdk.Context, denom string) (bridgeToken *types.BridgeToken, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBridgeTokenByDenomKey(denom))
	if len(bz) == 0 {
		return nil, types.ErrNotFound
	}
	k.cdc.MustUnmarshal(bz, bridgeToken)
	return
}

func (k Keeper) HasBridgeToken(ctx sdk.Context, tokenContract string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetBridgeTokenByContractKey(tokenContract))
}

func (k Keeper) SetBridgeToken(ctx sdk.Context, token *types.BridgeToken) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetBridgeTokenByDenomKey(token.Denom), k.cdc.MustMarshal(token))
	store.Set(types.GetBridgeTokenByContractKey(token.Contract), k.cdc.MustMarshal(token))
}

func (k Keeper) IterateBridgeTokenToDenom(ctx sdk.Context, cb func(*types.BridgeToken) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.BridgeTokenByDenomKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		bridgeToken := &types.BridgeToken{}
		k.cdc.MustUnmarshal(iter.Value(), bridgeToken)
		// cb returns true to stop early
		if cb(bridgeToken) {
			break
		}
	}
}
