package types

import (
	"errors"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/st-chain/me-hub/types"
)

const (
	MaxOracleSize                  = 21
	DefaultOracleDelegateThreshold = 10
	OutgoingTxBatchSize            = 100
	MaxResults                     = 100
	MaxOracleSetRequestsResults    = 5
	MaxKeepEventSize               = 100
	DefBridgeCallTimeout           = 604_800_000 // 7 * 24 * 3600 * 1000
	MaxGasLimit                    = 30_000_000
	MaxLiquidationSize             = 100
)

var (
	// AttestationVotesPowerThreshold threshold of votes power to succeed
	AttestationVotesPowerThreshold = sdkmath.NewInt(66)

	AttestationRelayerChangePowerThreshold = sdkmath.NewInt(30)
)

func DefaultParams() Params {
	return Params{
		GravityId:                          "me-gravity",
		AverageBlockTime:                   7_000,
		AverageExternalBlockTime:           5_000,
		ExternalBatchTimeout:               12 * 3600 * 1000,
		SignedWindow:                       30_000,
		SlashFraction:                      sdk.NewDecWithPrec(8, 1), // 80%
		RelayerSetUpdatePowerChangePercent: sdk.NewDecWithPrec(1, 1), // 10%
		IbcTransferTimeoutHeight:           20_000,
		Relayers:                           nil,
	}
}

// ValidateBasic checks that the parameters have valid values.
// nolint:gocyclo
func (m *Params) ValidateBasic() error {
	if len(m.GravityId) == 0 {
		return fmt.Errorf("gravityId cannpt be empty")
	}
	if _, err := fxtypes.StrToByte32(m.GravityId); err != nil {
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
		return fmt.Errorf("invalid signed window too short")
	}
	if m.SlashFraction.IsNegative() {
		return fmt.Errorf("attempted to slash with a negative slash factor: %v", m.SlashFraction)
	}
	if m.SlashFraction.GT(sdk.OneDec()) {
		return fmt.Errorf("slash factor too large: %s", m.SlashFraction)
	}
	if m.IbcTransferTimeoutHeight <= 1 {
		return fmt.Errorf("invalid ibc transfer timeout too short")
	}
	if m.OracleSetUpdatePowerChangePercent.IsNegative() {
		return fmt.Errorf("attempted to powet change percent with a negative: %v", m.OracleSetUpdatePowerChangePercent)
	}
	if m.OracleSetUpdatePowerChangePercent.GT(sdk.OneDec()) {
		return fmt.Errorf("powet change percent too large: %s", m.OracleSetUpdatePowerChangePercent)
	}
	if !m.DelegateThreshold.IsValid() || !m.DelegateThreshold.IsPositive() {
		return fmt.Errorf("invalid delegate threshold")
	}
	if m.DelegateThreshold.Denom != fxtypes.DefaultDenom {
		return fmt.Errorf("oracle delegate denom must FX")
	}
	if m.DelegateMultiple <= 0 {
		return fmt.Errorf("invalid delegate multiple")
	}
	if len(m.Oracles) > 0 {
		return errors.New("deprecated oracles")
	}
	if m.BridgeCallTimeout <= 3_600_000 {
		return fmt.Errorf("invalid bridge call timeout")
	}
	return nil
}
