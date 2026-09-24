package rollapp

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const ModuleName = "rollapp"

// Params is the pre-v3 rollapp params stored in x/params.
// Only fields needed for migration are retained.
type Params struct {
	DisputePeriodInBlocks uint64 `json:"dispute_period_in_blocks" yaml:"dispute_period_in_blocks"`
	RollappsEnabled       bool   `json:"rollapps_enabled" yaml:"rollapps_enabled"`
}

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyDisputePeriodInBlocks = []byte("DisputePeriodInBlocks")
	KeyRollappsEnabled       = []byte("RollappsEnabled")
)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDisputePeriodInBlocks, &p.DisputePeriodInBlocks, validateDisputePeriodInBlocks),
		paramtypes.NewParamSetPair(KeyRollappsEnabled, &p.RollappsEnabled, func(_ interface{}) error { return nil }),
	}
}

func validateDisputePeriodInBlocks(v interface{}) error {
	disputePeriodInBlocks, ok := v.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}
	if disputePeriodInBlocks < 1 {
		return fmt.Errorf("dispute period cannot be lower than 1 block")
	}
	return nil
}
