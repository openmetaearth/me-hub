package types

import (
	"fmt"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Issuers: []didtypes.DidInfo{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	for _, issuer := range gs.Issuers {
		if len(issuer.Did) != didtypes.DidLength {
			return fmt.Errorf("DID length must be equal to %d", didtypes.DidLength)
		}

		if issuer.Pubkey == "" {
			return fmt.Errorf("the pubkey is empty")
		}

		if issuer.Status != didtypes.DID_STATUS_ACTIVE {
			return fmt.Errorf("DID status must be active")
		}
	}

	return nil
}
