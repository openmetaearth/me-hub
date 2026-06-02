package sequencer_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/openmetaearth/me-hub/testutil/keeper"
	"github.com/openmetaearth/me-hub/testutil/nullify"
	"github.com/openmetaearth/me-hub/x/sequencer"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
	"github.com/stretchr/testify/require"
)

func TestInitGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		SequencerList: []types.Sequencer{
			{
				SequencerAddress: "0",
				Status:           types.Bonded,
				Proposer:         true,
			},
			{
				SequencerAddress: "1",
				Status:           types.Bonded,
			},
		},
		ReplaceProposerList: []types.MsgStoreReplaceProposer{
			{
				ReplaceProposer: types.MsgRepalceProposer{
					RollappId:   "rollapp_1234-1",
					OldProposer: "0",
					NewProposer: "1",
					BlockHeight: 42,
				},
				HubBlockHeight: 7,
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.SequencerKeeper(t)
	sequencer.InitGenesis(ctx, *k, genesisState)
	got := sequencer.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.SequencerList, got.SequencerList)
	require.ElementsMatch(t, genesisState.ReplaceProposerList, got.ReplaceProposerList)
	// this line is used by starport scaffolding # genesis/test/assert
}

func TestExportGenesis(t *testing.T) {
	params := types.Params{
		MinBond:       sdk.NewCoin("dym", sdk.NewInt(100)),
		UnbondingTime: 100,
	}
	sequencerList := []types.Sequencer{
		{
			SequencerAddress: "0",
			Status:           types.Bonded,
			Proposer:         true,
		},
		{
			SequencerAddress: "1",
			Status:           types.Bonded,
		},
	}
	k, ctx := keepertest.SequencerKeeper(t)
	k.SetParams(ctx, params)
	for _, sequencer := range sequencerList {
		k.SetSequencer(ctx, sequencer)
	}
	got := sequencer.ExportGenesis(ctx, *k)
	require.NotNil(t, got)
	require.Equal(t, params, got.Params)
	require.ElementsMatch(t, sequencerList, got.SequencerList)
}

func TestGenesisExportImportPreservesPendingReplaceProposer(t *testing.T) {
	pending := types.MsgRepalceProposer{
		RollappId:   "rollapp_1234-1",
		OldProposer: "old-proposer",
		NewProposer: "new-proposer",
		BlockHeight: 42,
	}

	k, ctx := keepertest.SequencerKeeper(t)
	ctx = ctx.WithBlockHeight(7)
	require.NoError(t, k.SetReplaceProposer(ctx, &pending))

	beforeExport, err := k.GetReplaceProposer(ctx, pending.RollappId)
	require.NoError(t, err)
	require.NotNil(t, beforeExport)

	exported := sequencer.ExportGenesis(ctx, *k)

	importedKeeper, importedCtx := keepertest.SequencerKeeper(t)
	sequencer.InitGenesis(importedCtx, *importedKeeper, *exported)

	afterImport, err := importedKeeper.GetReplaceProposer(importedCtx, pending.RollappId)
	require.NoError(t, err)
	require.NotNil(t, afterImport, "pending replace proposer was dropped by genesis export/import")
	require.Equal(t, *beforeExport, *afterImport)
}
