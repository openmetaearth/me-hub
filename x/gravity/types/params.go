package types

import (
	"fmt"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/utils"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	OutgoingTxBatchSize = 100
	MaxKeepEventSize    = 100
	MaxGasLimit         = 30_000_000
)

var (
	// AttestationVotesPowerThreshold threshold of votes power to succeed
	AttestationVotesPowerThreshold = sdkmath.NewInt(66)

	AttestationProposalRelayerChangePowerThreshold = sdkmath.NewInt(30)
)

func DefaultParams() Params {
	return Params{
		GravityId:                          "me-gravity",
		AverageBlockTime:                   7_000,
		ExternalBatchTimeout:               24 * 3600 * 1000, // 24 hours
		AverageExternalBlockTime:           5_000,            // 5 seconds
		SignedWindow:                       30_000,
		SlashFraction:                      sdk.NewDecWithPrec(8, 1), // 80%
		RelayerSetUpdatePowerChangePercent: sdk.NewDecWithPrec(2, 1), // 20%
		MaxRelayers:                        5,
		MinDelegate:                        sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100_000_000)),    // 1 MEC
		MaxDelegate:                        sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(10_000_000_000)), // 100 MEC
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
		return fmt.Errorf("powet change percent too large: %s", m.RelayerSetUpdatePowerChangePercent)
	}
	if !m.MinDelegate.IsValid() || !m.MinDelegate.IsPositive() {
		return fmt.Errorf("invalid delegate threshold")
	}
	if m.MinDelegate.Denom != params.BaseDenom {
		return fmt.Errorf("oracle delegate denom must FX")
	}
	if !m.MaxDelegate.IsValid() || !m.MaxDelegate.IsPositive() {
		return fmt.Errorf("invalid delegate threshold")
	}
	if m.MaxDelegate.Denom != params.BaseDenom {
		return fmt.Errorf("oracle delegate denom must FX")
	}
	return nil
}
