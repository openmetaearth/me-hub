package types

import (
	"fmt"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidateBasic validates genesis state by looping through the params and
// calling their validation functions
func (m *GenesisState) ValidateBasic() error {
	if err := m.Params.ValidateBasic(); err != nil {
		return err
	}
	var totalPower uint64
	for i := range m.Relayers {
		power, err := validateGenesisRelayer(m.Params, m.Relayers[i])
		if err != nil {
			return fmt.Errorf("invalid relayer %d: %w", i, err)
		}
		if !m.Relayers[i].Online || power == 0 {
			continue
		}
		if totalPower > math.MaxUint64-power {
			return fmt.Errorf("total online relayer power exceeds uint64")
		}
		totalPower += power
	}
	return nil
}

func validateGenesisRelayer(params Params, relayer Relayer) (uint64, error) {
	if _, err := sdk.AccAddressFromBech32(relayer.RelayerAddress); err != nil {
		return 0, fmt.Errorf("invalid relayer address: %w", err)
	}
	if relayer.DelegateAmount.IsNil() || !relayer.DelegateAmount.IsPositive() {
		return 0, fmt.Errorf("delegate amount must be positive")
	}
	if relayer.DelegateAmount.LT(params.MinDelegate) {
		return 0, ErrDelegateAmountBelowMinimum
	}
	if relayer.DelegateAmount.GT(params.MaxDelegate) {
		return 0, ErrDelegateAmountAboveMaximum
	}
	power := relayer.GetPower()
	if power.IsPositive() && !power.IsUint64() {
		return 0, fmt.Errorf("relayer power must fit uint64")
	}
	if !power.IsPositive() {
		return 0, nil
	}
	return power.Uint64(), nil
}
