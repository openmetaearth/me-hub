package dao

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/dao/keeper"
	"github.com/st-chain/me-hub/x/dao/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetGlobalDao(ctx, sdk.MustAccAddressFromBech32(genState.GlobalDao))
	k.SetMeidDao(ctx, sdk.MustAccAddressFromBech32(genState.MeidDao))
	k.SetDevOperator(ctx, sdk.MustAccAddressFromBech32(genState.DevOperator))
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.GlobalDao = k.GetGlobalDao(ctx).String()
	genesis.MeidDao = k.GetMeidDao(ctx).String()
	genesis.DevOperator = k.GetDevOperator(ctx).String()
	return genesis
}
