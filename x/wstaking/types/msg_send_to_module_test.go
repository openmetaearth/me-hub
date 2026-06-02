package types

import (
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
)

func TestIsAllowedSendToModuleTarget(t *testing.T) {
	tests := []struct {
		name     string
		target   string
		expected bool
	}{
		{name: "stake pool", target: StakePoolName, expected: true},
		{name: "fixed deposit principal pool", target: FixedDepositPrincipalPool, expected: true},
		{name: "bridge fee pool", target: BridgeFeePool, expected: true},
		{name: "global dao fee pool", target: GlobalDaoFeePool, expected: false},
		{name: "fee collector", target: authtypes.FeeCollectorName, expected: false},
		{name: "distribution", target: distrtypes.ModuleName, expected: false},
		{name: "bonded pool", target: stakingtypes.BondedPoolName, expected: false},
		{name: "not bonded pool", target: stakingtypes.NotBondedPoolName, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, IsAllowedSendToModuleTarget(tt.target))
		})
	}
}
