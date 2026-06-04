package rollapp_test

import (
	"testing"

	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/x/rollapp"
	"github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisExportImportPreservesSkipDelayRollapps(t *testing.T) {
	sourceApp := apptesting.Setup(t, false)
	sourceCtx := sourceApp.GetBaseApp().NewContext(false, cometbftproto.Header{})
	rollappID := "rollapp_1234-1"
	disabledRollappID := "rollapp_5678-1"

	sourceApp.RollappKeeper.SetRollapp(sourceCtx, types.Rollapp{RollappId: rollappID})
	sourceApp.RollappKeeper.SetRollapp(sourceCtx, types.Rollapp{RollappId: disabledRollappID})
	sourceApp.RollappKeeper.SetSkipDelayRollapp(sourceCtx, rollappID, true)
	sourceApp.RollappKeeper.SetSkipDelayRollapp(sourceCtx, disabledRollappID, false)
	require.True(t, sourceApp.RollappKeeper.IsSkipDelayRollapp(sourceCtx, rollappID))
	require.False(t, sourceApp.RollappKeeper.IsSkipDelayRollapp(sourceCtx, disabledRollappID))

	exported := rollapp.ExportGenesis(sourceCtx, *sourceApp.RollappKeeper)
	require.ElementsMatch(t, []string{rollappID}, exported.SkipDelayRollappList)
	require.NoError(t, exported.Validate())

	exportedJSON := sourceApp.AppCodec().MustMarshalJSON(exported)
	require.Contains(t, string(exportedJSON), "skip_delay_rollapp_list")
	var decoded types.GenesisState
	sourceApp.AppCodec().MustUnmarshalJSON(exportedJSON, &decoded)
	require.ElementsMatch(t, []string{rollappID}, decoded.SkipDelayRollappList)

	importedApp := apptesting.Setup(t, false)
	importedCtx := importedApp.GetBaseApp().NewContext(false, cometbftproto.Header{})
	rollapp.InitGenesis(importedCtx, *importedApp.RollappKeeper, *exported)

	require.True(t, importedApp.RollappKeeper.IsSkipDelayRollapp(importedCtx, rollappID))
	require.False(t, importedApp.RollappKeeper.IsSkipDelayRollapp(importedCtx, disabledRollappID))
	require.ElementsMatch(t, []string{rollappID}, importedApp.RollappKeeper.GetSkipDelayRollapps(importedCtx))
}
