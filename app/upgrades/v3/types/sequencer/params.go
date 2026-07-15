package sequencer

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const ModuleName = "sequencer"

// Params is the pre-v3 sequencer params stored in x/params.
type Params struct {
	MinBond       sdk.Coin     `json:"min_bond" yaml:"min_bond"`
	UnbondingTime time.Duration `json:"unbonding_time" yaml:"unbonding_time"`
}

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyMinBond       = []byte("MinBond")
	KeyUnbondingTime = []byte("UnbondingTime")
)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMinBond, &p.MinBond, validateMinBond),
		paramtypes.NewParamSetPair(KeyUnbondingTime, &p.UnbondingTime, validateUnbondingTime),
	}
}

func validateMinBond(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if !v.IsValid() {
		return fmt.Errorf("invalid coin: %s", v)
	}
	return nil
}

func validateUnbondingTime(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v <= 0 {
		return fmt.Errorf("unbonding time must be positive: %d", v)
	}
	return nil
}
