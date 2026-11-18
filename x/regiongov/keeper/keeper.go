package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

// Keeper stores region governance state.
// At VotingEndTime, EndBlocker tallies the proposal and executes it when it
// passes (or marks Failed if execution errors); otherwise it is Rejected.
type Keeper struct {
	cdc           codec.BinaryCodec
	storeKey      storetypes.StoreKey
	stakingKeeper types.StakingKeeper
	daoKeeper     types.DaoKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	stakingKeeper types.StakingKeeper,
	daoKeeper types.DaoKeeper,
) Keeper {
	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		stakingKeeper: stakingKeeper,
		daoKeeper:     daoKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
