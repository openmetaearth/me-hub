package rollapp_test

import (
	"testing"

	keepertest "github.com/openmetaearth/me-hub/testutil/keeper"
	"github.com/openmetaearth/me-hub/testutil/nullify"
	"github.com/openmetaearth/me-hub/x/rollapp"
	"github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/stretchr/testify/require"
)

func TestInitExportGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		RollappList: []types.Rollapp{
			{
				RollappId: "0",
			},
			{
				RollappId: "1",
			},
		},
		StateInfoList: []types.StateInfo{
			{
				StateInfoIndex: types.StateInfoIndex{
					RollappId: "0",
					Index:     0,
				},
			},
			{
				StateInfoIndex: types.StateInfoIndex{
					RollappId: "1",
					Index:     1,
				},
			},
		},
		LatestStateInfoIndexList: []types.StateInfoIndex{
			{
				RollappId: "0",
			},
			{
				RollappId: "1",
			},
		},
		BlockHeightToFinalizationQueueList: []types.BlockHeightToFinalizationQueue{
			{
				CreationHeight: 0,
			},
			{
				CreationHeight: 1,
			},
		},
		SkipDelayRollappList: []string{"1"},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.RollappKeeper(t)
	rollapp.InitGenesis(ctx, *k, genesisState)
	got := rollapp.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(genesisState)
	nullify.Fill(*got)

	require.ElementsMatch(t, genesisState.RollappList, got.RollappList)
	require.ElementsMatch(t, genesisState.StateInfoList, got.StateInfoList)
	require.ElementsMatch(t, genesisState.LatestStateInfoIndexList, got.LatestStateInfoIndexList)
	require.ElementsMatch(t, genesisState.BlockHeightToFinalizationQueueList, got.BlockHeightToFinalizationQueueList)
	require.ElementsMatch(t, genesisState.SkipDelayRollappList, got.SkipDelayRollappList)
	require.True(t, k.IsSkipDelayRollapp(ctx, "1"))
	require.False(t, k.IsSkipDelayRollapp(ctx, "0"))
	// this line is used by starport scaffolding # genesis/test/assert
}

func TestGenesisExportImportPreservesSkipDelayRollapps(t *testing.T) {
	rollappID := "rollapp_1234-1"
	genesisState := types.GenesisState{
		Params:               types.DefaultParams(),
		RollappList:          []types.Rollapp{{RollappId: rollappID}},
		SkipDelayRollappList: []string{rollappID},
	}

	k, ctx := keepertest.RollappKeeper(t)
	rollapp.InitGenesis(ctx, *k, genesisState)
	require.True(t, k.IsSkipDelayRollapp(ctx, rollappID))
	require.ElementsMatch(t, []string{rollappID}, k.GetSkipDelayRollapps(ctx))

	exported := rollapp.ExportGenesis(ctx, *k)

	importedKeeper, importedCtx := keepertest.RollappKeeper(t)
	rollapp.InitGenesis(importedCtx, *importedKeeper, *exported)

	require.True(t, importedKeeper.IsSkipDelayRollapp(importedCtx, rollappID), "skip-delay rollapp was dropped by genesis export/import")
	require.ElementsMatch(t, []string{rollappID}, importedKeeper.GetSkipDelayRollapps(importedCtx))
}
