package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	metypes "github.com/openmetaearth/me-hub/types"
)

// SetChainIDFromCosmos configures the EVM keeper from a Cosmos chain-id.
// Legacy ids such as "mechain" are mapped to the EIP-155 form (e.g. mechain_2404-1).
func (k *Keeper) SetChainIDFromCosmos(cosmosChainID string) {
	k.WithChainID(sdk.Context{}.WithChainID(metypes.ChainIdWithEIP155From(cosmosChainID)))
}
