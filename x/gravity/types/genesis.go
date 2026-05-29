package types

import (
	errorsmod "cosmossdk.io/errors"
)

// ValidateBasic validates genesis state by looping through the params and
// calling their validation functions
func (m *GenesisState) ValidateBasic() error {
	if err := m.Params.ValidateBasic(); err != nil {
		return err
	}

	denoms := make(map[string]bool)
	contracts := make(map[string]bool)

	for _, token := range m.BridgeTokens {
		if err := token.ValidateBasic(); err != nil {
			return err
		}
		if denoms[token.Denom] {
			return errorsmod.Wrapf(ErrDuplicate, "duplicate bridge token denom: %s", token.Denom)
		}
		denoms[token.Denom] = true

		if contracts[token.ContractAddress] {
			return errorsmod.Wrapf(ErrDuplicate, "duplicate bridge token contract: %s", token.ContractAddress)
		}
		contracts[token.ContractAddress] = true
	}

	return nil
}
