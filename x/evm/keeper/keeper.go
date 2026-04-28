package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	metypes "github.com/openmetaearth/me-hub/types"
)

// Wrapper wraps the original mint keeper and intercepts its original methods if needed.
type Keeper struct {
	*evmkeeper.Keeper
}

func NewKeeper(ek *evmkeeper.Keeper) *Keeper {
	ek.WithChainID(sdk.Context{}.WithChainID(metypes.ChainIdWithEIP155()))
	return &Keeper{
		Keeper: ek,
	}
}
