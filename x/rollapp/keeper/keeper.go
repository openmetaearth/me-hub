package keeper

import (
	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/openmetaearth/me-hub/x/rollapp/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	hooks      types.MultiRollappHooks
	paramstore paramtypes.Subspace

	ibcClientKeeper types.IBCClientKeeper
	channelKeeper   types.ChannelKeeper

	registeredRollappDenoms collections.KeySet[collections.Pair[string, string]] // [ rollappID, denom ]

	finalizePending func(ctx sdk.Context, stateInfoIndex types.StateInfoIndex) error
	daoKeeper       types.DaoKeeper
	// sequencerKeeper types.SequencerKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	channelKeeper types.ChannelKeeper,
	ibcclientKeeper types.IBCClientKeeper,
	daoKeeper types.DaoKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	k := &Keeper{
		cdc:             cdc,
		storeKey:        storeKey,
		paramstore:      ps,
		hooks:           nil,
		channelKeeper:   channelKeeper,
		ibcClientKeeper: ibcclientKeeper,
		daoKeeper:       daoKeeper,
	}
	k.SetFinalizePendingFn(k.finalizePendingState)
	return k
}

func (k *Keeper) SetFinalizePendingFn(fn func(ctx sdk.Context, stateInfoIndex types.StateInfoIndex) error) {
	k.finalizePending = fn
}

/* -------------------------------------------------------------------------- */
/*                                    Hooks                                   */
/* -------------------------------------------------------------------------- */

func (k *Keeper) SetHooks(sh types.MultiRollappHooks) {
	if k.hooks != nil {
		panic("cannot set rollapp hooks twice")
	}
	k.hooks = sh
}

func (k *Keeper) GetHooks() types.MultiRollappHooks {
	return k.hooks
}
