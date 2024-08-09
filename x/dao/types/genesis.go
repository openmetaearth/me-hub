package types

import (
	// this line is used by starport scaffolding # genesis/types/import
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		GlobalDao:   "",
		MeidDao:     "",
		DevOperator: "",
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	_, err := sdk.AccAddressFromBech32(gs.GlobalDao)
	if err != nil {
		return fmt.Errorf("invalid global dao address %s", gs.GlobalDao)
	}

	_, err = sdk.AccAddressFromBech32(gs.MeidDao)
	if err != nil {
		return fmt.Errorf("invalid dao address %s", gs.MeidDao)
	}

	_, err = sdk.AccAddressFromBech32(gs.DevOperator)
	if err != nil {
		return fmt.Errorf("invalid dev operator address %s", gs.DevOperator)
	}

	return nil
}
