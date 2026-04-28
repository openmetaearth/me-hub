package types_test

import (
	"testing"

	_ "github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/openmetaearth/me-hub/x/did/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	test := types.GenesisState{
		Infos: []types.DidInfo{
			types.DidInfo{
				Did:    "1000000000001",
				Pubkey: "{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"AjkBriaNQIyoihm/Op5a53ovjdThnbs8G3GhSdErW7Mt\"}",
				Status: types.DID_STATUS_ACTIVE,
			},
		},
		Svcs: []types.Service{
			types.Service{
				Sid:         "kyc",
				Name:        "kyc",
				Description: "this is kyc test service.",
				Issuers:     []string{"1000000000001"},
				Status:      types.SERVICE_STATUS_ACTIVE,
			},
		},
		Vcs: []types.Credential{
			types.Credential{
				Did:  "1000000000001",
				Sid:  "kyc",
				Hash: "0000000000000000001",
				Uri:  "http://metaearth.com/files/0001.vc",
				Data: []byte("test"),
			},
		},
		Flogs: []types.FilterLogger{
			types.FilterLogger{
				Did: "000000000001",
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
