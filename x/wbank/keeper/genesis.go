package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

// InitGenesis initializes the bank module's state from a given genesis state.
func (k BaseKeeperWrapper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) {
	k.BaseKeeper.InitGenesis(ctx, genState)
}

func (k BaseKeeperWrapper) HasOrSetAccount(ctx sdk.Context, address sdk.AccAddress) {
	if !k.ak.HasAccount(ctx, address) {
		k.ak.SetAccount(ctx, k.ak.NewAccountWithAddress(ctx, address))
	}
}
