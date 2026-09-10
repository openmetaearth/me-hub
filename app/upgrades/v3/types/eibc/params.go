package eibc

import (
	"fmt"

	"cosmossdk.io/math"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const ModuleName = "eibc"

// Params is the pre-v3 eibc params stored in x/params.
type Params struct {
	EpochIdentifier string         `json:"epoch_identifier" yaml:"epoch_identifier"`
	TimeoutFee      math.LegacyDec `json:"timeout_fee" yaml:"timeout_fee"`
	ErrackFee       math.LegacyDec `json:"errack_fee" yaml:"errack_fee"`
}

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyEpochIdentifier = []byte("EpochIdentifier")
	KeyTimeoutFee      = []byte("TimeoutFee")
	KeyErrAckFee       = []byte("ErrAckFee")
)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyEpochIdentifier, &p.EpochIdentifier, validateEpochIdentifier),
		paramtypes.NewParamSetPair(KeyTimeoutFee, &p.TimeoutFee, validateFee),
		paramtypes.NewParamSetPair(KeyErrAckFee, &p.ErrackFee, validateFee),
	}
}

func validateEpochIdentifier(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if len(v) == 0 {
		return fmt.Errorf("epoch identifier cannot be empty")
	}
	return nil
}

func validateFee(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNil() || v.IsNegative() || v.GTE(math.LegacyOneDec()) {
		return fmt.Errorf("invalid fee: %s", v)
	}
	return nil
}
