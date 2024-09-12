package types

import (
	"fmt"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

const (
	KeyElectionPeriod               = "KeyElectionPeriod"
	KeyMinStakeAmount               = "KeyMinStakeAmount"
	KeySequencerNumber              = "KeySequencerNumber"
	KeyBackupNumber                 = "KeyBackupNumber"
	KeyFirstElectInterval           = "KeyFirstElectInterval"
	KeyApplyElectionTime            = "KeyApplyElectionTime"
	KeyElectionInterimTime          = "KeyElectionInterimTime"
	KeyDaFraudChallengeStake        = "KeyDaFraudChallengeStake"
	MecPrecision             uint64 = 100000000
)

var (
	//默认的选举间隔周期，单位为分钟
	defaultElectionPeriod  uint32 = 43200
	defaultMinStakeAmount  uint64 = 1000
	defaultSequencerNumber uint32 = 10
	defaultBackupNumber    uint32 = 3
	//默认首次选举的时间，单位为分钟
	defaultFirstElectionInterval uint32 = 120
	//默认允许申请参与质押的时间，默认2天，单位为分钟
	defaultAllowApplyElectionTime uint32 = 2880
	//默认选举后的过渡时间，单位为秒
	defaultElectionInterimTime uint32 = 300
	//默认的DA挑战者的质押金额
	defaultDaFraudChallengeStake uint32 = 100
)

func DefaultParams() Params {
	return Params{
		ElectionPeriod:         defaultElectionPeriod,
		MinStakeAmount:         defaultMinStakeAmount,
		SequencerNumber:        defaultSequencerNumber,
		BackupSequencerNumber:  defaultBackupNumber,
		FirstElectionInterval:  defaultFirstElectionInterval,
		AllowApplyElectionTime: defaultAllowApplyElectionTime,
		ElectionInterimTime:    defaultElectionInterimTime,
		DaFraudChallengeStake:  defaultDaFraudChallengeStake,
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
		paramtypes.NewParamSetPair([]byte(KeyElectionInterimTime), &p.ElectionInterimTime, validateElectionInterimTime),
		paramtypes.NewParamSetPair([]byte(KeyDaFraudChallengeStake), &p.DaFraudChallengeStake, validateDaFraudChallengeStake),
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
	if err := validateElectionInterimTime(p.ElectionInterimTime); err != nil {
		return err
	}
	if err := validateDaFraudChallengeStake(p.DaFraudChallengeStake); err != nil {
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
	if val < 1 {
		return fmt.Errorf("minStakeAmount <  1")
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

func validateElectionInterimTime(v interface{}) error {
	val, ok := v.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}
	if val < 1 {
		return fmt.Errorf("ElectionInterimTime error. val = %d", val)
	}
	return nil
}

func validateDaFraudChallengeStake(v interface{}) error {
	val, ok := v.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}
	if val < 1 {
		return fmt.Errorf("DaFraudChallengeStake error. val = %d", val)
	}
	return nil
}
