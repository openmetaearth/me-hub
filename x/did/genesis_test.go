package did_test

import (
	"testing"

	"github.com/openmetaearth/me-hub/testutil/nullify"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/did"

	_ "github.com/evmos/ethermint/crypto/ethsecp256k1"
	keepertest "github.com/openmetaearth/me-hub/testutil/keeper"
	"github.com/openmetaearth/me-hub/x/did/types"
	"github.com/stretchr/testify/require"
)

func TestInitExportGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Infos: []types.DidInfo{
			{
				Did:     "1000000000000001",
				Address: sample.AccAddress(), // 使用 sample 生成有效地址
				Pubkey:  "{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"AjkBriaNQIyoihm/Op5a53ovjdThnbs8G3GhSdErW7Mt\"}",
				Status:  types.DID_STATUS_ACTIVE,
			},
		},
		Svcs: []types.Service{
			{
				Sid:         "kyc",
				Name:        "kyc",
				Description: "this is kyc test service.",
				Issuers:     []string{"00000000000001"},
				Status:      types.SERVICE_STATUS_ACTIVE,
			},
		},
		Vcs: []types.Credential{
			{
				Did:  "1000000000000001",
				Sid:  "kyc",
				Hash: "0000000000000000001",
				Uri:  "http://metaearth.com/files/0001.vc",
				Data: []byte("test"),
			},
		},
		Flogs: []types.FilterLogger{
			{
				Did: "1000000000000001",
				Sid: "kyc",
				Filters: [][]byte{
					[]byte("A0"),
				},
			},
		},
	}

	k, ctx := keepertest.DidKeeper(t)
	did.InitGenesis(ctx, k, genesisState)
	got := did.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(genesisState)
	nullify.Fill(*got)
	require.ElementsMatch(t, genesisState.Infos, got.Infos)
	require.ElementsMatch(t, genesisState.Svcs, got.Svcs)
	require.ElementsMatch(t, genesisState.Vcs, got.Vcs)
	require.ElementsMatch(t, genesisState.Flogs, got.Flogs)
}
