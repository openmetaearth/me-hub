package types_test

import (
	"bytes"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/evmos/ethermint/crypto/ethsecp256k1"
	_ "github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/did/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	test := validGenesis()

	err := test.Validate()
	require.NoError(t, err)
}

func TestGenesisState_ValidateRejectsCorruptState(t *testing.T) {
	tooLargeCredentialData := bytes.Repeat([]byte("x"), types.MaxCredentialDataSize+1)

	testCases := []struct {
		name   string
		mutate func(*types.GenesisState)
	}{
		{
			name: "duplicate DID",
			mutate: func(genesis *types.GenesisState) {
				duplicate := genesis.Infos[0]
				duplicate.Address = testAddress(2)
				genesis.Infos = append(genesis.Infos, duplicate)
			},
		},
		{
			name: "duplicate address",
			mutate: func(genesis *types.GenesisState) {
				duplicate := genesis.Infos[0]
				duplicate.Did = "1000000000002"
				genesis.Infos = append(genesis.Infos, duplicate)
			},
		},
		{
			name: "invalid DID status",
			mutate: func(genesis *types.GenesisState) {
				genesis.Infos[0].Status = types.DidStatus(99)
			},
		},
		{
			name: "service references unknown issuer DID",
			mutate: func(genesis *types.GenesisState) {
				genesis.Svcs[0].Issuers = []string{"9999999999999"}
			},
		},
		{
			name: "duplicate service",
			mutate: func(genesis *types.GenesisState) {
				genesis.Svcs = append(genesis.Svcs, genesis.Svcs[0])
			},
		},
		{
			name: "credential references unknown DID",
			mutate: func(genesis *types.GenesisState) {
				genesis.Vcs[0].Did = "9999999999999"
			},
		},
		{
			name: "credential references unknown service",
			mutate: func(genesis *types.GenesisState) {
				genesis.Vcs[0].Sid = "other"
			},
		},
		{
			name: "duplicate credential",
			mutate: func(genesis *types.GenesisState) {
				genesis.Vcs = append(genesis.Vcs, genesis.Vcs[0])
			},
		},
		{
			name: "oversized credential data",
			mutate: func(genesis *types.GenesisState) {
				genesis.Vcs[0].Data = tooLargeCredentialData
			},
		},
		{
			name: "filter logger references missing credential",
			mutate: func(genesis *types.GenesisState) {
				genesis.Flogs[0].Sid = "other"
			},
		},
		{
			name: "duplicate filter logger",
			mutate: func(genesis *types.GenesisState) {
				genesis.Flogs = append(genesis.Flogs, genesis.Flogs[0])
			},
		},
		{
			name: "oversized filter",
			mutate: func(genesis *types.GenesisState) {
				genesis.Flogs[0].Filters[0] = bytes.Repeat([]byte("x"), 1025)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			genesis := validGenesis()
			tc.mutate(&genesis)

			err := genesis.Validate()
			require.Error(t, err)
		})
	}
}

func validGenesis() types.GenesisState {
	return types.GenesisState{
		Infos: []types.DidInfo{
			{
				Did:     "1000000000001",
				Address: testAddress(1),
				Pubkey:  "{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"AjkBriaNQIyoihm/Op5a53ovjdThnbs8G3GhSdErW7Mt\"}",
				Status:  types.DID_STATUS_ACTIVE,
			},
		},
		Svcs: []types.Service{
			{
				Sid:         "kyc",
				Name:        "kyc",
				Description: "this is kyc test service.",
				Issuers:     []string{"1000000000001"},
				Status:      types.SERVICE_STATUS_ACTIVE,
			},
		},
		Vcs: []types.Credential{
			{
				Did:  "1000000000001",
				Sid:  "kyc",
				Hash: "0000000000000000001",
				Uri:  "http://metaearth.com/files/0001.vc",
				Data: []byte("test"),
			},
		},
		Flogs: []types.FilterLogger{
			{
				Did: "1000000000001",
				Sid: "kyc",
				Filters: [][]byte{
					[]byte("A0"),
				},
			},
		},
	}
}

func testAddress(seed byte) string {
	return sdk.AccAddress(bytes.Repeat([]byte{seed}, 20)).String()
}
