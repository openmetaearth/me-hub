package types

import (
	"fmt"
)

// GenesisState defines the gravity genesis state
type GenesisState struct {
	Params       Params        `json:"params" yaml:"params"`
	BridgeTokens []BridgeToken `json:"bridge_tokens" yaml:"bridge_tokens"`
}

// ValidateBasic performs basic validation of the genesis data returning an
// error for any failed validation criteria.
func (s GenesisState) ValidateBasic() error {
	if err := s.Params.ValidateBasic(); err != nil {
		return fmt.Errorf("params validation failed: %w", err)
	}

	denomSet := make(map[string]struct{})
	contractSet := make(map[string]struct{})
	for i, token := range s.BridgeTokens {
		if _, exists := denomSet[token.Denom]; exists {
			return fmt.Errorf("duplicate denom '%s' at index %d: bridge token denoms must be unique", token.Denom, i)
		}
		if _, exists := contractSet[token.ContractAddress]; exists {
			return fmt.Errorf("duplicate contract address '%s' at index %d: bridge token contracts must be unique", token.ContractAddress, i)
		}
		denomSet[token.Denom] = struct{}{}
		contractSet[token.ContractAddress] = struct{}{}
	}

	return nil
}
