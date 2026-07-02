package types

import (
	"strings"
	"testing"

	_ "github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/stretchr/testify/require"
)

// genesis_test.go is in package types (not types_test) so it shares the
// TestMain defined in message_vc_test.go, which configures the "me" bech32
// prefix before any test runs.

const (
	genesisTestDID        = "1000000000001" // length == DidLength (13)
	genesisTestAddress    = "me1kjnt3ypezt3yf58w8upujvejdtt5xsvkq5dpk4"
	genesisTestPubkey     = "{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"AjkBriaNQIyoihm/Op5a53ovjdThnbs8G3GhSdErW7Mt\"}"
	genesisTestShortValue = "short" // invalid DID/sid value used across multiple tests
)

func validGenesis() GenesisState {
	return GenesisState{
		Infos: []DidInfo{
			{
				Did:     genesisTestDID,
				Address: genesisTestAddress,
				Pubkey:  genesisTestPubkey,
				Status:  DID_STATUS_ACTIVE,
			},
		},
		Svcs: []Service{
			{
				Sid:         "kyc",
				Name:        "kyc",
				Description: "this is kyc test service.",
				Issuers:     []string{genesisTestDID},
				Status:      SERVICE_STATUS_ACTIVE,
			},
		},
		Vcs: []Credential{
			{
				Did:  genesisTestDID,
				Sid:  "kyc",
				Hash: "0000000000000000001",
				Uri:  "http://metaearth.com/files/0001.vc",
				Data: []byte("test"),
			},
		},
		Flogs: []FilterLogger{
			{
				Did: genesisTestDID,
				Sid: "kyc",
				Filters: [][]byte{
					[]byte("A0"),
				},
			},
		},
	}
}

// --- happy path ---

func TestGenesisState_Validate_Valid(t *testing.T) {
	gs := validGenesis()
	require.NoError(t, gs.Validate())
}

func TestGenesisState_Validate_EmptyGenesis(t *testing.T) {
	gs := DefaultGenesis()
	require.NoError(t, gs.Validate())
}

// --- Infos ---

func TestGenesisState_Validate_DuplicateDID(t *testing.T) {
	gs := validGenesis()
	gs.Infos = append(gs.Infos, gs.Infos[0])
	require.ErrorContains(t, gs.Validate(), "duplicate DID")
}

func TestGenesisState_Validate_DuplicateAddress(t *testing.T) {
	gs := validGenesis()
	gs.Infos = append(gs.Infos, DidInfo{
		Did:     "9999999999999",
		Address: genesisTestAddress,
		Pubkey:  genesisTestPubkey,
		Status:  DID_STATUS_ACTIVE,
	})
	require.ErrorContains(t, gs.Validate(), "duplicate address")
}

func TestGenesisState_Validate_InvalidDIDLength(t *testing.T) {
	gs := validGenesis()
	gs.Infos[0].Did = genesisTestShortValue
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

func TestGenesisState_Validate_InvalidDidStatus(t *testing.T) {
	gs := validGenesis()
	gs.Infos[0].Status = DidStatus(99)
	require.ErrorContains(t, gs.Validate(), "invalid DID status")
}

// --- Svcs ---

func TestGenesisState_Validate_DuplicateServiceSid(t *testing.T) {
	gs := validGenesis()
	gs.Svcs = append(gs.Svcs, gs.Svcs[0])
	require.ErrorContains(t, gs.Validate(), "duplicate service sid")
}

func TestGenesisState_Validate_InvalidServiceSidLength(t *testing.T) {
	gs := validGenesis()
	gs.Svcs[0].Sid = "a" // too short (< 2)
	require.ErrorContains(t, gs.Validate(), "sid length")
}

func TestGenesisState_Validate_InvalidServiceNameEmpty(t *testing.T) {
	gs := validGenesis()
	gs.Svcs[0].Name = ""
	require.ErrorContains(t, gs.Validate(), "name length")
}

func TestGenesisState_Validate_InvalidServiceNameTooLong(t *testing.T) {
	gs := validGenesis()
	gs.Svcs[0].Name = strings.Repeat("x", 21)
	require.ErrorContains(t, gs.Validate(), "name length")
}

func TestGenesisState_Validate_InvalidServiceDescriptionTooLong(t *testing.T) {
	gs := validGenesis()
	gs.Svcs[0].Description = strings.Repeat("x", 1025)
	require.ErrorContains(t, gs.Validate(), "description length")
}

func TestGenesisState_Validate_InvalidServiceIssuerDIDLength(t *testing.T) {
	gs := validGenesis()
	gs.Svcs[0].Issuers = []string{"tooshort"}
	require.ErrorContains(t, gs.Validate(), "issuer DID length")
}

func TestGenesisState_Validate_InvalidServiceStatus(t *testing.T) {
	gs := validGenesis()
	gs.Svcs[0].Status = ServiceStatus(99)
	require.ErrorContains(t, gs.Validate(), "invalid service status")
}

// --- Vcs ---

func TestGenesisState_Validate_DuplicateCredential(t *testing.T) {
	gs := validGenesis()
	gs.Vcs = append(gs.Vcs, gs.Vcs[0])
	require.ErrorContains(t, gs.Validate(), "duplicate credential")
}

func TestGenesisState_Validate_InvalidCredentialDIDLength(t *testing.T) {
	gs := validGenesis()
	gs.Vcs[0].Did = genesisTestShortValue
	require.ErrorContains(t, gs.Validate(), "DID length")
}

func TestGenesisState_Validate_InvalidCredentialSidLength(t *testing.T) {
	gs := validGenesis()
	gs.Vcs[0].Sid = "a" // too short
	require.ErrorContains(t, gs.Validate(), "sid length")
}

func TestGenesisState_Validate_InvalidCredentialHashEmpty(t *testing.T) {
	gs := validGenesis()
	gs.Vcs[0].Hash = ""
	require.ErrorContains(t, gs.Validate(), "hash length")
}

func TestGenesisState_Validate_InvalidCredentialHashTooLong(t *testing.T) {
	gs := validGenesis()
	gs.Vcs[0].Hash = strings.Repeat("x", 129)
	require.ErrorContains(t, gs.Validate(), "hash length")
}

func TestGenesisState_Validate_InvalidCredentialUriTooLong(t *testing.T) {
	gs := validGenesis()
	gs.Vcs[0].Uri = strings.Repeat("x", 1025)
	require.ErrorContains(t, gs.Validate(), "uri length")
}

func TestGenesisState_Validate_CredentialDataTooLarge(t *testing.T) {
	gs := validGenesis()
	gs.Vcs[0].Data = []byte(strings.Repeat("x", maxCredentialDataLength+1))
	require.ErrorContains(t, gs.Validate(), "data length exceeds")
}

// --- Flogs ---

func TestGenesisState_Validate_FlogOrphanCredential(t *testing.T) {
	gs := validGenesis()
	gs.Flogs[0].Did = "9999999999999" // no matching vc
	require.ErrorContains(t, gs.Validate(), "not found in vcs")
}

func TestGenesisState_Validate_FlogInvalidDIDLength(t *testing.T) {
	gs := validGenesis()
	gs.Flogs[0].Did = genesisTestShortValue
	require.ErrorContains(t, gs.Validate(), "DID length")
}

func TestGenesisState_Validate_FlogInvalidSidLength(t *testing.T) {
	gs := validGenesis()
	gs.Flogs[0].Sid = "a" // too short
	require.ErrorContains(t, gs.Validate(), "sid length")
}

func TestGenesisState_Validate_FlogFilterTooLong(t *testing.T) {
	gs := validGenesis()
	gs.Flogs[0].Filters = [][]byte{[]byte(strings.Repeat("x", 1025))}
	require.ErrorContains(t, gs.Validate(), "filter length exceeds")
}

func TestGenesisState_Validate_DuplicateFlog(t *testing.T) {
	gs := validGenesis()
	gs.Flogs = append(gs.Flogs, gs.Flogs[0])
	require.ErrorContains(t, gs.Validate(), "duplicate filter logger")
}
