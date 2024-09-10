package kyc_test

import (
	"github.com/st-chain/me-hub/testutil/nullify"
	didtypes "github.com/st-chain/me-hub/x/did/types"
	"github.com/st-chain/me-hub/x/kyc"
	"testing"

	_ "github.com/evmos/ethermint/crypto/ethsecp256k1"
	keepertest "github.com/st-chain/me-hub/testutil/keeper"
	"github.com/st-chain/me-hub/x/kyc/types"
	"github.com/stretchr/testify/require"
)

func TestInitExportGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Issuer: didtypes.DidInfo{
			Did:    "1000000000000001",
			Pubkey: "{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"AjkBriaNQIyoihm/Op5a53ovjdThnbs8G3GhSdErW7Mt\"}",
			Status: didtypes.DID_STATUS_ACTIVE,
		},
	}

	k, ctx := keepertest.KycKeeper(t)
	kyc.InitGenesis(ctx, *k, genesisState)
	got := kyc.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(genesisState)
	nullify.Fill(*got)
	require.ElementsMatch(t, genesisState.Issuer, got.Issuer)
}
