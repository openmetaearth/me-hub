package wdistri

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/wdistri/keeper"
	"github.com/st-chain/me-hub/x/wdistri/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
	k.Keeper.InitGenesis(ctx, genState.CosmosDistributionState)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	genesis.CosmosDistributionState = *k.Keeper.ExportGenesis(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
