package keeper

import (
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
)

// Wrapper wraps the original mint keeper and intercepts its original methods if needed.
type Keeper struct {
	*evmkeeper.Keeper
}

func NewKeeper(ek *evmkeeper.Keeper) *Keeper {
	return &Keeper{
		Keeper: ek,
	}
}
