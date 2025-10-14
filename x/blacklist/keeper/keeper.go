package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/st-chain/me-hub/x/blacklist/types"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storetypes.StoreKey
		memKey   storetypes.StoreKey
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,
	}
}

// InitGenesis initializes the blacklist module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	// Initialize blacklist
	for _, address := range genState.Blacklist.Addresses {
		k.SetBlacklist(ctx, types.Blacklist{Addresses: []string{address}})
	}

}

// GetBlacklist returns the current blacklist
func (k Keeper) GetBlacklist(ctx sdk.Context) (blacklist types.Blacklist, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.BlacklistKeyPrefix
	bz := store.Get(key)
	if bz == nil {
		return blacklist, false
	}
	k.cdc.MustUnmarshal(bz, &blacklist)
	return blacklist, true
}

// SetBlacklist sets the blacklist
func (k Keeper) SetBlacklist(ctx sdk.Context, blacklist types.Blacklist) error {
	store := ctx.KVStore(k.storeKey)
	key := types.BlacklistKeyPrefix
	bz := k.cdc.MustMarshal(&blacklist)
	store.Set(key, bz)
	return nil
}

// ExportGenesis returns the blacklist module's exported genesis
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	blacklist, found := k.GetBlacklist(ctx)
	if !found {
		blacklist = types.Blacklist{Addresses: []string{}}
	}
	return &types.GenesisState{
		Blacklist: blacklist,
	}
}
