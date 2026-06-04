package rollapp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/rollapp/keeper"
	"github.com/openmetaearth/me-hub/x/rollapp/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the rollapp
	for _, elem := range genState.RollappList {
		k.SetRollapp(ctx, elem)
	}
	// Set all the stateInfo
	for _, elem := range genState.StateInfoList {
		k.SetStateInfo(ctx, elem)
	}
	// Set all the latestStateInfoIndex
	for _, elem := range genState.LatestStateInfoIndexList {
		k.SetLatestStateInfoIndex(ctx, elem)
	}
	// Set all the latestFinalizedStateIndex
	for _, elem := range genState.LatestFinalizedStateIndexList {
		k.SetLatestFinalizedStateIndex(ctx, elem)
	}
	// Set all the blockHeightToFinalizationQueue
	for _, elem := range genState.BlockHeightToFinalizationQueueList {
		k.SetBlockHeightToFinalizationQueue(ctx, elem)
	}
	for _, rollappID := range genState.SkipDelayRollappList {
		k.SetSkipDelayRollapp(ctx, rollappID, true)
	}
	k.SetParams(ctx, genState.Params)
	// this line is used by starport scaffolding # genesis/module/init
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.RollappList = k.GetAllRollapps(ctx)
	genesis.StateInfoList = k.GetAllStateInfo(ctx)
	genesis.LatestStateInfoIndexList = k.GetAllLatestStateInfoIndex(ctx)
	genesis.LatestFinalizedStateIndexList = k.GetAllLatestFinalizedStateIndex(ctx)
	genesis.BlockHeightToFinalizationQueueList = k.GetAllBlockHeightToFinalizationQueue(ctx)
	genesis.SkipDelayRollappList = k.GetSkipDelayRollapps(ctx)

	return genesis
}
