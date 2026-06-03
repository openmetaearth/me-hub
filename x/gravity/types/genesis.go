package types

import (
	"fmt"
	errorsmod "cosmossdk.io/errors"
	"github.com/openmetaearth/me-hub/utils"
)

// ValidateBasic validates genesis state by looping through the params and
// calling their validation functions
func (m *GenesisState) ValidateBasic() error {
	if err := m.Params.ValidateBasic(); err != nil {
		return err
	}

	denoms := make(map[string]bool)
	contracts := make(map[string]bool)
	for i, token := range m.BridgeTokens {
		if token.Denom == "" {
			return errorsmod.Wrapf(ErrEmpty, "bridge token %d denom is empty", i)
		}
		if token.ContractAddress == "" {
			return errorsmod.Wrapf(ErrEmpty, "bridge token %d contract address is empty", i)
		}
		if err := utils.ValidateEthereumAddress(token.ContractAddress); err != nil {
			return errorsmod.Wrapf(err, "bridge token %d contract address %s is invalid", i, token.ContractAddress)
		}
		if token.Supply.IsNegative() {
			return errorsmod.Wrapf(ErrInvalid, "bridge token %d supply %s cannot be negative", i, token.Supply.String())
		}
		if denoms[token.Denom] {
			return errorsmod.Wrapf(ErrDuplicate, "duplicate bridge token denom: %s", token.Denom)
		}
		if contracts[token.ContractAddress] {
			return errorsmod.Wrapf(ErrDuplicate, "duplicate bridge token contract address: %s", token.ContractAddress)
		}
		denoms[token.Denom] = true
		contracts[token.ContractAddress] = true
	}

	return nil
}
