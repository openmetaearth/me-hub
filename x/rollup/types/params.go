package types

import (
	"fmt"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

const (
	KeyElectionPeriod            = "KeyElectionPeriod"
	KeyMinStakeAmount            = "KeyMinStakeAmount"
	KeySequencerNumber           = "KeySequencerNumber"
	KeyBackupNumber              = "KeyBackupNumber"
	KeyFirstElectInterval        = "KeyFirstElectInterval"
	KeyApplyElectionTime         = "KeyApplyElectionTime"
	MecPrecision          uint64 = 100000000
)

var (
	defaultElectionPeriod  uint32 = 30
	defaultMinStakeAmount  uint64 = 1000 * MecPrecision
	defaultSequencerNumber uint32 = 10
	defaultBackupNumber    uint32 = 3
)

func DefaultParams() Params {
	return Params{
		ElectionPeriod:        defaultElectionPeriod,
		MinStakeAmount:        defaultMinStakeAmount,
		SequencerNumber:       defaultSequencerNumber,
		BackupSequencerNumber: defaultBackupNumber,
	}
}

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair([]byte(KeyElectionPeriod), &p.ElectionPeriod, validateElectionPeriod),
		paramtypes.NewParamSetPair([]byte(KeyMinStakeAmount), &p.MinStakeAmount, validateMinStakeAmount),
		paramtypes.NewParamSetPair([]byte(KeySequencerNumber), &p.SequencerNumber, validateSequencerNumber),
		paramtypes.NewParamSetPair([]byte(KeyBackupNumber), &p.BackupSequencerNumber, validateBackupSequencerNumber),
		paramtypes.NewParamSetPair([]byte(KeyFirstElectInterval), &p.FirstElectionInterval, validateFirstElectInterval),
		paramtypes.NewParamSetPair([]byte(KeyApplyElectionTime), &p.AllowApplyElectionTime, validateAllowApplyElectionTime),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateElectionPeriod(p.ElectionPeriod); err != nil {
		return err
	}
	if err := validateMinStakeAmount(p.MinStakeAmount); err != nil {
		return err
	}
	if err := validateSequencerNumber(p.SequencerNumber); err != nil {
		return err
	}
	if err := validateFirstElectInterval(p.FirstElectionInterval); err != nil {
		return err
	}
	if err := validateAllowApplyElectionTime(p.AllowApplyElectionTime); err != nil {
		return err
	}
	if p.AllowApplyElectionTime >= p.ElectionPeriod {
		return fmt.Errorf("AllowApplyElectionTime(%d) must less than ElectionPeriod(%d)",
			p.AllowApplyElectionTime, p.ElectionPeriod)
	}
	return validateBackupSequencerNumber(p.BackupSequencerNumber)
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateElectionPeriod(v interface{}) error {
	val, ok := v.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}
	if val == 0 {
		return fmt.Errorf("ElectionPeriod be 0")
	}
	return nil
}

func validateMinStakeAmount(v interface{}) error {
	val, ok := v.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}
	if val < 100000000 {
		return fmt.Errorf("minStakeAmount <  100000000")
	}
	return nil
}

func validateSequencerNumber(v interface{}) error {
	val, ok := v.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}
	if val < 4 {
		return fmt.Errorf("Sequencer's number <  4")
	}
	return nil
}

func validateBackupSequencerNumber(v interface{}) error {
	_, ok := v.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	return nil
}

func validateFirstElectInterval(v interface{}) error {
	val, ok := v.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}
	if val < 1 {
		return fmt.Errorf("FirstElectInterval error. val = %d", val)
	}
	return nil
}

func validateAllowApplyElectionTime(v interface{}) error {
	val, ok := v.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}
	if val < 1 {
		return fmt.Errorf("AllowApplyElectionTime error. val = %d", val)
	}
	return nil
}
