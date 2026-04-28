package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/gravity/types"
)

func (k Keeper) GetBridgeTokenByContract(ctx sdk.Context, tokenContract string) (*types.BridgeToken, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBridgeTokenByContractKey(tokenContract))
	if len(bz) == 0 {
		return nil, types.ErrNotFound
	}
	bridgeToken := &types.BridgeToken{}
	k.cdc.MustUnmarshal(bz, bridgeToken)
	return bridgeToken, nil
}

// GetBridgeTokenByDenom retrieves a BridgeToken by its denom.
// Returns ErrNotFound if no token exists for the given denom.
func (k Keeper) GetBridgeTokenByDenom(ctx sdk.Context, denom string) (*types.BridgeToken, error) {
	if denom == "" {
		return nil, types.ErrInvalid
	}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBridgeTokenByDenomKey(denom))
	if len(bz) == 0 {
		return nil, types.ErrNotFound
	}
	bridgeToken := &types.BridgeToken{}
	k.cdc.MustUnmarshal(bz, bridgeToken)
	return bridgeToken, nil
}

// HasBridgeToken returns true if a BridgeToken exists for the given contract address.
func (k Keeper) HasBridgeToken(ctx sdk.Context, tokenContract string) bool {
	if tokenContract == "" {
		return false
	}
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetBridgeTokenByContractKey(tokenContract))
}

func (k Keeper) SetBridgeToken(ctx sdk.Context, token *types.BridgeToken) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetBridgeTokenByDenomKey(token.Denom), k.cdc.MustMarshal(token))
	store.Set(types.GetBridgeTokenByContractKey(token.ContractAddress), k.cdc.MustMarshal(token))
}

func (k Keeper) DelBridgeToken(ctx sdk.Context, token *types.BridgeToken) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetBridgeTokenByDenomKey(token.Denom))
	store.Delete(types.GetBridgeTokenByContractKey(token.ContractAddress))
}

func (k Keeper) IterateBridgeTokenByDenom(ctx sdk.Context, cb func(*types.BridgeToken) bool) {
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
