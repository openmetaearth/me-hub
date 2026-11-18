package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// AccountAddress returns the ballot signer. New snapshots store account addresses;
// validator operator addresses remain supported for legacy country snapshots.
func (e Elector) AccountAddress() (sdk.AccAddress, error) {
	if account, err := sdk.AccAddressFromBech32(e.Operator); err == nil {
		return account, nil
	}
	operator, err := sdk.ValAddressFromBech32(e.Operator)
	if err != nil {
		return nil, err
	}
	return sdk.AccAddress(operator), nil
}
