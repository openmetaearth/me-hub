package types_test

import (
	"testing"

	didtypes "github.com/openmetaearth/me-hub/x/did/types"

	"github.com/openmetaearth/me-hub/x/kyc/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is invalid",
			genState: types.DefaultGenesis(),
			valid:    false,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				Issuers: []didtypes.DidInfo{
					didtypes.DidInfo{
						Did:    "1000000000000001",
						Pubkey: "{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"AyfZ/7fojbKMioe5Oaw378EH4F8w2CGvZ7SwOCRvlCH8\"}",
						Status: didtypes.DID_STATUS_ACTIVE,
					},
				},
			},
			valid: true,
		},
		{
			desc: "did is invalid",
			genState: &types.GenesisState{
				Issuers: []didtypes.DidInfo{
					didtypes.DidInfo{
						Did:    "0x1000000000000001",
						Pubkey: "{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"AyfZ/7fojbKMioe5Oaw378EH4F8w2CGvZ7SwOCRvlCH8\"}",
						Status: didtypes.DID_STATUS_ACTIVE,
					},
				},
			},
			valid: false,
		},
		{
			desc: "pubkey is invalid",
			genState: &types.GenesisState{
				Issuers: []didtypes.DidInfo{
					didtypes.DidInfo{
						Did:    "0x1000000000000001",
						Pubkey: "{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"AyfZ/7fojbKMioe5Oaw378EH4F8w2CGvZ7SwOCRvlCH8\"}",
						Status: didtypes.DID_STATUS_ACTIVE,
					},
				},
			},
			valid: false,
		},
		{
			desc: "status is invalid",
			genState: &types.GenesisState{
				Issuers: []didtypes.DidInfo{
					didtypes.DidInfo{
						Did:    "0x1000000000000001",
						Pubkey: "",
						Status: didtypes.DID_STATUS_INACTIVE,
					},
				},
			},
			valid: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
