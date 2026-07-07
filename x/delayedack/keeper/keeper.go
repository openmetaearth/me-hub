package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"

	"github.com/openmetaearth/me-hub/x/delayedack/types"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	hooks      types.MultiDelayedAckHooks
	paramstore paramtypes.Subspace

	rollappKeeper types.RollappKeeper
	porttypes.ICS4Wrapper
	channelKeeper types.ChannelKeeper
	types.EIBCKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	rollappKeeper types.RollappKeeper,
	ics4Wrapper porttypes.ICS4Wrapper,
	channelKeeper types.ChannelKeeper,
	eibcKeeper types.EIBCKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}
	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		paramstore:    ps,
		rollappKeeper: rollappKeeper,
		ICS4Wrapper:   ics4Wrapper,
		channelKeeper: channelKeeper,
		EIBCKeeper:    eibcKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) IsRollappsEnabled(ctx sdk.Context) bool {
	return k.rollappKeeper.GetParams(ctx).RollappsEnabled
}

func (k Keeper) IsSkipDelayRollapp(ctx sdk.Context, rollappID string) bool {
	return k.rollappKeeper.IsSkipDelayRollapp(ctx, rollappID)
}

func (k Keeper) getRollappFinalizedHeight(ctx sdk.Context, chainID string) (uint64, error) {
	// GetLatestFinalizedStateIndex
	latestFinalizedStateIndex, found := k.rollappKeeper.GetLatestFinalizedStateIndex(ctx, chainID)
	if !found {
		return 0, rollapptypes.ErrNoFinalizedStateYetForRollapp
	}

	stateInfo := k.rollappKeeper.MustGetStateInfo(ctx, chainID, latestFinalizedStateIndex.Index)
	return stateInfo.StartHeight + stateInfo.NumBlocks - 1, nil
}

/* -------------------------------------------------------------------------- */
/*                               Hooks handling                               */
/* -------------------------------------------------------------------------- */

func (k *Keeper) SetHooks(hooks types.MultiDelayedAckHooks) {
	if k.hooks != nil {
		panic("DelayedAckHooks already set")
	}
	k.hooks = hooks
}

func (k *Keeper) GetHooks() types.MultiDelayedAckHooks {
	return k.hooks
}

/* -------------------------------------------------------------------------- */
/*                                 ICS4Wrapper                                */
/* -------------------------------------------------------------------------- */

// LookupModuleByChannel wraps ChannelKeeper LookupModuleByChannel function.
func (k *Keeper) LookupModuleByChannel(ctx sdk.Context, portID, channelID string) (string, *capabilitytypes.Capability, error) {
	return k.channelKeeper.LookupModuleByChannel(ctx, portID, channelID)
}
