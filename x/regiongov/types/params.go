package types

import (
	"fmt"
	"math"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	// DefaultMinVotingPeriod is one hour, expressed in seconds.
	DefaultMinVotingPeriod uint64 = 3600
	// MaxVotingPeriod is the largest duration time.Duration can represent, in seconds.
	MaxVotingPeriod              uint64 = uint64(math.MaxInt64 / int64(time.Second))
	DefaultPassThresholdPercent  uint32 = 51
	DefaultEquityWithdrawPercent uint32 = 1
)

// DefaultParams returns the default regiongov parameters.
func DefaultParams() Params {
	return Params{
		MinVotingPeriod:       DefaultMinVotingPeriod,
		PassThresholdPercent:  DefaultPassThresholdPercent,
		EquityWithdrawPercent: DefaultEquityWithdrawPercent,
	}
}

// Validate checks the parameter bounds agreed for region governance.
func (p Params) Validate() error {
	if p.MinVotingPeriod < 1 || p.MinVotingPeriod > MaxVotingPeriod {
		return fmt.Errorf("min_voting_period must be a number of seconds in [1,%d], got %d", MaxVotingPeriod, p.MinVotingPeriod)
	}
	if p.PassThresholdPercent < 1 || p.PassThresholdPercent > 100 {
		return fmt.Errorf("pass_threshold_percent must be in [1,100], got %d", p.PassThresholdPercent)
	}
	if p.EquityWithdrawPercent < 1 || p.EquityWithdrawPercent > 100 {
		return fmt.Errorf("equity_withdraw_percent must be in [1,100], got %d", p.EquityWithdrawPercent)
	}
	return nil
}

// String implements fmt.Stringer.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
