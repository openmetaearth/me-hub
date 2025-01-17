package wdistri_test

import (
	"testing"

	keepertest "github.com/st-chain/me-hub/testutil/keeper"
	"github.com/st-chain/me-hub/testutil/nullify"
	"github.com/st-chain/me-hub/x/wdistri"
	"github.com/st-chain/me-hub/x/wdistri/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.WdistriKeeper(t)
	wdistri.InitGenesis(ctx, *k, genesisState)
	got := wdistri.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
