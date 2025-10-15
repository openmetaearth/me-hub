package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

// InitGenesis initializes the bank module's state from a given genesis state.
func (k BaseKeeperWrapper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) {
	k.SetParams(ctx, genState.Params)

	totalSupply := sdk.Coins{}
	genState.Balances = types.SanitizeGenesisBalances(genState.Balances)
	for _, balance := range genState.Balances {
		addr := balance.GetAddress()
		k.HasOrSetAccount(ctx, addr)

		if err := k.SetBalances(ctx, addr, balance.Coins); err != nil {
			panic(fmt.Errorf("error on setting balances %w", err))
		}

		totalSupply = totalSupply.Add(balance.Coins...)
	}

	if !genState.Supply.Empty() && !genState.Supply.IsEqual(totalSupply) {
		panic(fmt.Errorf("genesis supply is incorrect, expected %v, got %v", genState.Supply, totalSupply))
	}

	k.SetSupply(ctx, totalSupply)

	for _, meta := range genState.DenomMetadata {
		k.SetDenomMetaData(ctx, meta)
	}

	for _, se := range genState.GetAllSendEnabled() {
		k.SetSendEnabled(ctx, se.Denom, se.Enabled)
	}
}

func (k BaseKeeperWrapper) HasOrSetAccount(ctx sdk.Context, address sdk.AccAddress) {
	if !k.ak.HasAccount(ctx, address) {
		k.ak.SetAccount(ctx, k.ak.NewAccountWithAddress(ctx, address))
	}
}
