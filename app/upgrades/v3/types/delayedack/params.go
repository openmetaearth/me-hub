package delayedack

import (
	"fmt"

	"cosmossdk.io/math"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const ModuleName = "delayedack"

// Params is the pre-v3 delayedack params stored in x/params.
type Params struct {
	EpochIdentifier         string         `json:"epoch_identifier" yaml:"epoch_identifier"`
	BridgingFee             math.LegacyDec `json:"bridging_fee" yaml:"bridging_fee"`
	DeletePacketsEpochLimit int32          `json:"delete_packets_epoch_limit" yaml:"delete_packets_epoch_limit"`
}

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyEpochIdentifier         = []byte("EpochIdentifier")
	KeyBridgeFee               = []byte("BridgeFee")
	KeyDeletePacketsEpochLimit = []byte("DeletePacketsEpochLimit")
)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyEpochIdentifier, &p.EpochIdentifier, validateEpochIdentifier),
		paramtypes.NewParamSetPair(KeyBridgeFee, &p.BridgingFee, validateBridgingFee),
		paramtypes.NewParamSetPair(KeyDeletePacketsEpochLimit, &p.DeletePacketsEpochLimit, validateDeletePacketsEpochLimit),
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

func validateBridgingFee(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNil() || v.IsNegative() || v.GTE(math.LegacyOneDec()) {
		return fmt.Errorf("invalid bridging fee: %s", v)
	}
	return nil
}

func validateDeletePacketsEpochLimit(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < 0 {
		return fmt.Errorf("delete packet epoch limit must not be negative: %d", v)
	}
	return nil
}
