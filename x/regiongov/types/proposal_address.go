package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// RequiredOperatorAccount returns the snapshotted ballot signer, not the current
// staking operator. Legacy proposals may still encode this account as valoper.
func (p Proposal) RequiredOperatorAccount() (sdk.AccAddress, error) {
	return (Elector{Operator: p.RequiredOperator}).AccountAddress()
}

// NormalizeRequiredOperator changes only the address encoding, never the signer.
// Non-equity proposals have no required operator and retain the empty value.
func (p *Proposal) NormalizeRequiredOperator() error {
	if p.RequiredOperator == "" {
		return nil
	}
	account, err := p.RequiredOperatorAccount()
	if err != nil {
		return err
	}
	p.RequiredOperator = account.String()
	return nil
}
