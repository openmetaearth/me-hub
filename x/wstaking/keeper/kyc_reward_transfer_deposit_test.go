package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"github.com/stretchr/testify/require"
)

func TestValidateTransferFixedDepositCfgAllowsHistoricalRateMismatch(t *testing.T) {
	fixed := types.FixedDeposit{
		Term: 30,
		Rate: sdk.NewDecWithPrec(1, 1),
	}

	cfgs := map[int64]types.FixedDepositCfg{
		30: {
			Term:   30,
			Rate:   sdk.NewDecWithPrec(2, 1),
			Status: types.RegionFixedDepositCfgStatusActive,
		},
	}

	require.NoError(t, validateTransferFixedDepositCfg(fixed, cfgs))
}

func TestValidateTransferFixedDepositCfgRejectsMissingOrInactiveTerm(t *testing.T) {
	fixed := types.FixedDeposit{Term: 30}

	require.Error(t, validateTransferFixedDepositCfg(fixed, map[int64]types.FixedDepositCfg{}))
	require.Error(t, validateTransferFixedDepositCfg(fixed, map[int64]types.FixedDepositCfg{
		30: {
			Term:   30,
			Status: types.RegionFixedDepositCfgStatusInactive,
		},
	}))
}
