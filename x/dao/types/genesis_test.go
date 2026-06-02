package types

import (
	"bytes"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisStateValidateFreeGasAccounts(t *testing.T) {
	state := GenesisState{
		DaoAddresses:    testGenesisDaoAddresses(),
		FreeGasAccounts: []string{testGenesisAddress(5), testGenesisAddress(6)},
	}
	require.NoError(t, state.Validate())

	state.FreeGasAccounts = []string{testGenesisAddress(5), "not-a-bech32-address"}
	require.ErrorContains(t, state.Validate(), "invalid free gas account address")

	duplicate := testGenesisAddress(5)
	state.FreeGasAccounts = []string{duplicate, duplicate}
	require.ErrorContains(t, state.Validate(), "duplicate free gas account address")
}

func testGenesisDaoAddresses() DaoAddresses {
	return DaoAddresses{
		GlobalDao:      testGenesisAddress(1),
		MeidDao:        testGenesisAddress(2),
		DevOperator:    testGenesisAddress(3),
		AirdropAddress: testGenesisAddress(4),
	}
}

func testGenesisAddress(seed byte) string {
	return sdk.AccAddress(bytes.Repeat([]byte{seed}, 20)).String()
}
