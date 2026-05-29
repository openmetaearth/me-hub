package types

import (
	fmt "fmt"
)

// ValidateBasic validates genesis state by looping through the params and
// calling their validation functions
func (m *GenesisState) ValidateBasic() error {
	if err := m.Params.ValidateBasic(); err != nil {
		return err
	}
	if err := ValidateGenesisBridgeTokens(m.BridgeTokens); err != nil {
		return err
	}
	return nil
}

// ValidateGenesisBridgeTokens checks that no two BridgeTokens share the same
// denom or the same contract address.  Duplicate entries would silently
// overwrite one of the two KV indexes (denom→token or contract→token)
// during InitGenesis, splitting them and corrupting bridge asset routing.
func ValidateGenesisBridgeTokens(tokens []BridgeToken) error {
	seenDenom := make(map[string]string, len(tokens))    // denom -> contract
	seenContract := make(map[string]string, len(tokens)) // contract -> denom
	for _, bt := range tokens {
		if bt.Denom == "" {
			return fmt.Errorf("bridge token with empty denom (contract %s)", bt.ContractAddress)
		}
		if bt.ContractAddress == "" {
			return fmt.Errorf("bridge token with empty contract address (denom %s)", bt.Denom)
		}
		if existing, ok := seenDenom[bt.Denom]; ok {
			return ErrDuplicate.Wrapf(
				"bridge token denom %q already registered to contract %s, cannot also register contract %s",
				bt.Denom, existing, bt.ContractAddress,
			)
		}
		if existing, ok := seenContract[bt.ContractAddress]; ok {
			return ErrDuplicate.Wrapf(
				"bridge token contract %q already registered to denom %s, cannot also register denom %s",
				bt.ContractAddress, existing, bt.Denom,
			)
		}
		seenDenom[bt.Denom] = bt.ContractAddress
		seenContract[bt.ContractAddress] = bt.Denom
	}
	return nil
}
