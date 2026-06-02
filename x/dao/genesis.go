package dao

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/dao/keeper"
	"github.com/openmetaearth/me-hub/x/dao/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetDaoAddresses(ctx, genState.DaoAddresses)
	for _, address := range genState.FreeGasAccounts {
		k.SetFreeGasAccount(ctx, address)
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.DaoAddresses, _ = k.GetDaoAddresses(ctx)
	k.IterateFreeGasAccounts(ctx, func(address string) bool {
		genesis.FreeGasAccounts = append(genesis.FreeGasAccounts, address)
		return false
	})
	sort.Strings(genesis.FreeGasAccounts)
	return genesis
}
