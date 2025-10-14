package types

import (
	"fmt"
	"github.com/st-chain/me-hub/utils"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	OutgoingTxBatchSize = 100
	MaxKeepEventSize    = 20
	MaxGasLimit         = 30_000_000
)

var (
	// AttestationVotesPowerThreshold threshold of votes power to succeed
	AttestationVotesPowerThreshold                 = sdkmath.NewInt(66)
	AttestationProposalRelayerChangePowerThreshold = sdkmath.NewInt(34)
)

func DefaultParams() Params {
	return Params{
		GravityId:                          "me-gravity",
		AverageBlockTime:                   5_000,
		ExternalBatchTimeout:               24 * 3600 * 1000, // 24 hours
		AverageExternalBlockTime:           1_000,            // 5 seconds
		SignedWindow:                       30_000,
		SlashFraction:                      sdk.NewDecWithPrec(8, 1), // 80%
		RelayerSetUpdatePowerChangePercent: sdk.NewDecWithPrec(2, 1), // 20%
		MaxRelayers:                        10,
		MinDelegate:                        sdkmath.NewInt(100_000_000),    // 1 MEC
		MaxDelegate:                        sdkmath.NewInt(10_000_000_000), // 100 MEC
	}
}

// ValidateBasic checks that the parameters have valid values.
// nolint:gocyclo
func (m *Params) ValidateBasic() error {
	if len(m.GravityId) == 0 {
		return fmt.Errorf("gravityId cannpt be empty")
	}
	if _, err := utils.StrToByte32(m.GravityId); err != nil {
		return err
	}
	if m.AverageBlockTime < 100 {
		return fmt.Errorf("invalid average block time, too short for latency limitations")
	}
	if m.ExternalBatchTimeout < 60000 {
		return fmt.Errorf("invalid target batch timeout, less than 60 seconds is too short")
	}
	if m.AverageExternalBlockTime < 100 {
		return fmt.Errorf("invalid average external block time, too short for latency limitations")
	}
	if m.SignedWindow <= 1 {
		return fmt.Errorf("invalid signed window, too short")
	}
	if m.SlashFraction.IsNegative() {
		return fmt.Errorf("attempted to slash with a negative slash factor: %v", m.SlashFraction)
	}
	if m.SlashFraction.GT(sdk.OneDec()) {
		return fmt.Errorf("slash factor too large: %s", m.SlashFraction)
	}
	if m.MaxRelayers < 1 {
		return fmt.Errorf("invalid max relayers, too short")
	}
	if m.RelayerSetUpdatePowerChangePercent.IsNegative() {
		return fmt.Errorf("attempted to powet change percent with a negative: %v", m.RelayerSetUpdatePowerChangePercent)
	}
	if m.RelayerSetUpdatePowerChangePercent.GT(sdk.OneDec()) {
		return fmt.Errorf("power change percent too large: %s", m.RelayerSetUpdatePowerChangePercent)
	}
	if !m.MinDelegate.IsPositive() {
		return fmt.Errorf("invalid delegate threshold")
	}
	if !m.MaxDelegate.IsPositive() {
		return fmt.Errorf("invalid delegate threshold")
	}
	return nil
}
