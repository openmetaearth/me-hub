package types_test

import (
	"testing"

	_ "github.com/evmos/ethermint/crypto/ethsecp256k1"
	_ "github.com/openmetaearth/me-hub/app/params"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/x/did/types"
)

func TestGenesisState_Validate(t *testing.T) {
	test := types.GenesisState{
		Infos: []types.DidInfo{
			{
				Did:     "1000000000001",
				Address: "me1kjnt3ypezt3yf58w8upujvejdtt5xsvkq5dpk4",
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

	err := test.Validate()
	require.NoError(t, err)
}

func TestGenesisState_Validate_Errors(t *testing.T) {
	validDid := "1000000000001"
	validAddr := "me1kjnt3ypezt3yf58w8upujvejdtt5xsvkq5dpk4"
	validPubkey := "{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"AjkBriaNQIyoihm/Op5a53ovjdThnbs8G3GhSdErW7Mt\"}"

	tests := []struct {
		name string
		gen  types.GenesisState
	}{
		{
			name: "invalid did length",
			gen: types.GenesisState{
				Infos: []types.DidInfo{
					{Did: "short", Address: validAddr, Pubkey: validPubkey, Status: types.DID_STATUS_ACTIVE},
				},
			},
		},
		{
			name: "invalid status value",
			gen: types.GenesisState{
				Infos: []types.DidInfo{
					{Did: validDid, Address: validAddr, Pubkey: validPubkey, Status: 99},
				},
			},
		},
		{
			name: "invalid bech32 address",
			gen: types.GenesisState{
				Infos: []types.DidInfo{
					{Did: validDid, Address: "invalid_bech32", Pubkey: validPubkey, Status: types.DID_STATUS_ACTIVE},
				},
			},
		},
		{
			name: "duplicate did",
			gen: types.GenesisState{
				Infos: []types.DidInfo{
					{Did: validDid, Address: validAddr, Pubkey: validPubkey, Status: types.DID_STATUS_ACTIVE},
					{Did: validDid, Address: "me1kjnt3ypezt3yf58w8upujvejdtt5xsvkq5dpk5", Pubkey: validPubkey, Status: types.DID_STATUS_ACTIVE},
				},
			},
		},
		{
			name: "duplicate address",
			gen: types.GenesisState{
				Infos: []types.DidInfo{
					{Did: validDid, Address: validAddr, Pubkey: validPubkey, Status: types.DID_STATUS_ACTIVE},
					{Did: "1000000000002", Address: validAddr, Pubkey: validPubkey, Status: types.DID_STATUS_ACTIVE},
				},
			},
		},
		{
			name: "invalid service name length",
			gen: types.GenesisState{
				Infos: []types.DidInfo{
					{Did: validDid, Address: validAddr, Pubkey: validPubkey, Status: types.DID_STATUS_ACTIVE},
				},
				Svcs: []types.Service{
					{Sid: "kyc", Name: "", Description: "desc", Issuers: []string{validDid}, Status: types.SERVICE_STATUS_ACTIVE},
				},
			},
		},
		{
			name: "phantom issuer reference",
			gen: types.GenesisState{
				Infos: []types.DidInfo{
					{Did: validDid, Address: validAddr, Pubkey: validPubkey, Status: types.DID_STATUS_ACTIVE},
				},
				Svcs: []types.Service{
					{Sid: "kyc", Name: "kyc", Description: "desc", Issuers: []string{"1000000000009"}, Status: types.SERVICE_STATUS_ACTIVE},
				},
			},
		},
		{
			name: "phantom credential references from filter logger",
			gen: types.GenesisState{
				Infos: []types.DidInfo{
					{Did: validDid, Address: validAddr, Pubkey: validPubkey, Status: types.DID_STATUS_ACTIVE},
				},
				Svcs: []types.Service{
					{Sid: "kyc", Name: "kyc", Description: "desc", Issuers: []string{validDid}, Status: types.SERVICE_STATUS_ACTIVE},
				},
				Vcs: []types.Credential{
					{Did: validDid, Sid: "kyc", Hash: "some_hash", Uri: "some_uri"},
				},
				Flogs: []types.FilterLogger{
					{Did: validDid, Sid: "other_svc", Filters: [][]byte{[]byte("filter")}},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.gen.Validate()
			require.Error(t, err)
		})
	}
}
