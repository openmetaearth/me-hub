package types

import (
	fmt "fmt"

	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		DistributionParams: types.Params{
			CommunityTax:        sdk.NewDecWithPrec(2, 2), // 2%
			BaseProposerReward:  sdk.NewDecWithPrec(1, 2), // 1%
			BonusProposerReward: sdk.NewDecWithPrec(4, 2), // 4%
			WithdrawAddrEnabled: true,
		},
	}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params
// FIXME: may need to refactor
func (p Params) Validate() error {
	if p.DistributionParams.BaseProposerReward.IsNegative() {
		return fmt.Errorf(
			"base proposer reward should be positive: %s", p.DistributionParams.BaseProposerReward,
		)
	}
	if p.DistributionParams.BonusProposerReward.IsNegative() {
		return fmt.Errorf(
			"bonus proposer reward should be positive: %s", p.DistributionParams.BonusProposerReward,
		)
	}
	if v := p.DistributionParams.BaseProposerReward.Add(p.DistributionParams.BonusProposerReward).Add(p.DistributionParams.CommunityTax); v.GT(sdk.OneDec()) {
		return fmt.Errorf(
			"sum of base, bonus proposer rewards, and community tax cannot be greater than one: %s", v,
		)
	}
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
