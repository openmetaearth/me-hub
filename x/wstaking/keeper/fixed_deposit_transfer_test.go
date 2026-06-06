package keeper

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func TestValidateFixedDepositTransferConfigIgnoresInactiveUnrelatedTerms(t *testing.T) {
	rate := sdk.MustNewDecFromStr("0.1")
	configs := fixedDepositConfigByTerm([]types.FixedDepositCfg{
		{Term: 30, Rate: rate, Status: types.RegionFixedDepositCfgStatusActive},
		{Term: 180, Rate: sdk.MustNewDecFromStr("0.2"), Status: types.RegionFixedDepositCfgStatusInactive},
	})

	err := validateFixedDepositTransferConfig(configs, types.FixedDeposit{
		Term: 30,
		Rate: rate,
	})
	if err != nil {
		t.Fatalf("expected unrelated inactive term to be ignored, got %v", err)
	}
}

func TestValidateFixedDepositTransferConfigRejectsInactiveMatchingTerm(t *testing.T) {
	rate := sdk.MustNewDecFromStr("0.1")
	configs := fixedDepositConfigByTerm([]types.FixedDepositCfg{
		{Term: 30, Rate: rate, Status: types.RegionFixedDepositCfgStatusInactive},
	})

	err := validateFixedDepositTransferConfig(configs, types.FixedDeposit{
		Term: 30,
		Rate: rate,
	})
	if err == nil {
		t.Fatal("expected inactive matching term to be rejected")
	}
	if !strings.Contains(err.Error(), "fixed deposit cfg status is inactive") {
		t.Fatalf("unexpected error: %v", err)
	}
}
