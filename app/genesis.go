package app

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// The genesis state of the blockchain is represented here as a map of raw json
// messages key'd by a identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState(cdc codec.JSONCodec) GenesisState {
	return NewDefaultGenesisStateWithBasicManager(cdc, nil)
}

// NewDefaultGenesisStateWithBasicManager generates the default state with a BasicManager.
func NewDefaultGenesisStateWithBasicManager(cdc codec.JSONCodec, basicManager module.BasicManager) GenesisState {
	if basicManager == nil {
		// Return empty genesis state if no basic manager provided
		return make(GenesisState)
	}
	return basicManager.DefaultGenesis(cdc)
}
