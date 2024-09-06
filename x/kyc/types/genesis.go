package types

import (
	"fmt"
	didtypes "github.com/st-chain/me-hub/x/did/types"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Issuer: didtypes.DidInfo{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if len(gs.Issuer.Did) != 16 {
		return fmt.Errorf("DID length must be equal to 16")
	}

	if gs.Issuer.Pubkey == "" {
		return fmt.Errorf("the pubkey is empty")
	}

	if gs.Issuer.Status != didtypes.DID_STATUS_ACTIVE {
		return fmt.Errorf("DID status must be active")
	}

	return nil
}
