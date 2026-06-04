package app_test

import (
	"testing"

	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/stretchr/testify/require"
)

func TestExportGenesisPreservesCommittedValidatorUpdates(t *testing.T) {
	sourceApp := apptesting.Setup(t, false)
	sourceCtx := sourceApp.NewContext(false, cometbftproto.Header{})
	sourceStore := sourceCtx.KVStore(sourceApp.GetKey(types.StoreKey))

	sourceValidatorUpdates := sourceStore.Get(types.ValidatorUpdatesKey)
	require.NotEmpty(t, sourceValidatorUpdates)

	exportedGenesis := sourceApp.StakingKeeper.ExportGenesis(sourceCtx)

	importApp := apptesting.Setup(t, false)
	importCtx := importApp.NewContext(false, cometbftproto.Header{})
	importStore := importCtx.KVStore(importApp.GetKey(types.StoreKey))
	importStore.Delete(types.ValidatorUpdatesKey)
	require.Empty(t, importStore.Get(types.ValidatorUpdatesKey))

	importApp.StakingKeeper.InitGenesis(importCtx, exportedGenesis)

	require.Equal(t, sourceValidatorUpdates, importStore.Get(types.ValidatorUpdatesKey))
}
