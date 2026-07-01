package types_test

import (
	"strings"
	"testing"

	_ "github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/x/did/types"
)

const (
	testDID     = "1000000000001" // length == DidLength (13)
	testAddress = "me1kjnt3ypezt3yf58w8upujvejdtt5xsvkq5dpk4"
	testPubkey  = "{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"AjkBriaNQIyoihm/Op5a53ovjdThnbs8G3GhSdErW7Mt\"}"
)

func validGenesis() types.GenesisState {
	return types.GenesisState{
		Infos: []types.DidInfo{
			{
				Did:     testDID,
				Address: testAddress,
				Pubkey:  testPubkey,
				Status:  types.DID_STATUS_ACTIVE,
			},
		},
		Svcs: []types.Service{
			{
				Sid:         "kyc",
				Name:        "kyc",
				Description: "this is kyc test service.",
				Issuers:     []string{testDID},
				Status:      types.SERVICE_STATUS_ACTIVE,
			},
		},
		Vcs: []types.Credential{
			{
				Did:  testDID,
				Sid:  "kyc",
				Hash: "0000000000000000001",
				Uri:  "http://metaearth.com/files/0001.vc",
				Data: []byte("test"),
			},
		},
		Flogs: []types.FilterLogger{
			{
				Did: testDID,
				Sid: "kyc",
				Filters: [][]byte{
					[]byte("A0"),
				},
			},
		},
	}
}

func TestGenesisState_Validate_Valid(t *testing.T) {
	gs := validGenesis()
	require.NoError(t, gs.Validate())
}

func TestGenesisState_Validate_EmptyGenesis(t *testing.T) {
	gs := types.DefaultGenesis()
	require.NoError(t, gs.Validate())
}

func TestGenesisState_Validate_DuplicateDID(t *testing.T) {
	gs := validGenesis()
	gs.Infos = append(gs.Infos, gs.Infos[0]) // duplicate DID
	require.ErrorContains(t, gs.Validate(), "duplicate DID")
}

func TestGenesisState_Validate_DuplicateAddress(t *testing.T) {
	gs := validGenesis()
	// different DID but same address
	gs.Infos = append(gs.Infos, types.DidInfo{
		Did:     "9999999999999",
		Address: testAddress,
		Pubkey:  testPubkey,
		Status:  types.DID_STATUS_ACTIVE,
	})
	require.ErrorContains(t, gs.Validate(), "duplicate address")
}

func TestGenesisState_Validate_InvalidDIDLength(t *testing.T) {
	gs := validGenesis()
	gs.Infos[0].Did = "short"
	require.ErrorContains(t, gs.Validate(), "DID length")
}

func TestGenesisState_Validate_EmptyPubkey(t *testing.T) {
	gs := validGenesis()
	gs.Infos[0].Pubkey = ""
	require.ErrorContains(t, gs.Validate(), "pubkey must not be empty")
}

func TestGenesisState_Validate_InvalidAddress(t *testing.T) {
	gs := validGenesis()
	gs.Infos[0].Address = "not-a-valid-address"
	require.ErrorContains(t, gs.Validate(), "invalid address")
}

func TestGenesisState_Validate_DuplicateServiceSid(t *testing.T) {
	gs := validGenesis()
	gs.Svcs = append(gs.Svcs, gs.Svcs[0])
	require.ErrorContains(t, gs.Validate(), "duplicate service sid")
}

func TestGenesisState_Validate_InvalidServiceSidLength(t *testing.T) {
	gs := validGenesis()
	gs.Svcs[0].Sid = "a" // too short
	require.ErrorContains(t, gs.Validate(), "sid length")
}

func TestGenesisState_Validate_DuplicateCredential(t *testing.T) {
	gs := validGenesis()
	gs.Vcs = append(gs.Vcs, gs.Vcs[0])
	require.ErrorContains(t, gs.Validate(), "duplicate credential")
}

func TestGenesisState_Validate_CredentialDataTooLarge(t *testing.T) {
	gs := validGenesis()
	gs.Vcs[0].Data = []byte(strings.Repeat("x", 64*1024+1))
	require.ErrorContains(t, gs.Validate(), "data length exceeds")
}

func TestGenesisState_Validate_FlogOrphanCredential(t *testing.T) {
	gs := validGenesis()
	// flog references a credential that does not exist in Vcs
	gs.Flogs[0].Did = "9999999999999"
	require.ErrorContains(t, gs.Validate(), "not found in vcs")
}

func TestGenesisState_Validate_FlogInvalidDIDLength(t *testing.T) {
	gs := validGenesis()
	gs.Flogs[0].Did = "short"
	require.ErrorContains(t, gs.Validate(), "DID length")
}

func TestGenesisState_Validate_DuplicateFlog(t *testing.T) {
	gs := validGenesis()
	gs.Flogs = append(gs.Flogs, gs.Flogs[0])
	require.ErrorContains(t, gs.Validate(), "duplicate filter logger")
}
