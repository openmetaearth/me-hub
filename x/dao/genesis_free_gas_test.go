package dao_test

import (
	"testing"

	"github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/x/dao"
	daotypes "github.com/openmetaearth/me-hub/x/dao/types"
	"github.com/stretchr/testify/require"
)

func newTestAddress() string {
	privKey := secp256k1.GenPrivKey()
	return sdk.AccAddress(privKey.PubKey().Address()).String()
}

func TestExportImportMustPreserveFreeGasAccounts(t *testing.T) {
	sourceApp := apptesting.Setup(t, false)
	sourceCtx := sourceApp.GetBaseApp().NewContext(false, types.Header{})

	daoAddresses := daotypes.DaoAddresses{
		GlobalDao:      newTestAddress(),
		MeidDao:        newTestAddress(),
		DevOperator:    newTestAddress(),
		AirdropAddress: newTestAddress(),
	}
	freeGasAccounts := []string{
		newTestAddress(),
		newTestAddress(),
	}

	sourceApp.DaoKeeper.SetDaoAddresses(sourceCtx, daoAddresses)
	for _, address := range freeGasAccounts {
		sourceApp.DaoKeeper.SetFreeGasAccount(sourceCtx, address)
		require.True(t, sourceApp.DaoKeeper.CheckFreeGasAccount(sourceCtx, address))
	}

	exported := dao.ExportGenesis(sourceCtx, sourceApp.DaoKeeper)
	require.Equal(t, daoAddresses, exported.DaoAddresses)
	require.ElementsMatch(t, freeGasAccounts, exported.FreeGasAccounts)
	require.NoError(t, exported.Validate())

	exportedJSON := sourceApp.AppCodec().MustMarshalJSON(exported)
	require.Contains(t, string(exportedJSON), "free_gas_accounts")
	var decoded daotypes.GenesisState
	sourceApp.AppCodec().MustUnmarshalJSON(exportedJSON, &decoded)
	require.ElementsMatch(t, freeGasAccounts, decoded.FreeGasAccounts)

	importedApp := apptesting.Setup(t, false)
	importedCtx := importedApp.GetBaseApp().NewContext(false, types.Header{})
	dao.InitGenesis(importedCtx, importedApp.DaoKeeper, *exported)

	imported := dao.ExportGenesis(importedCtx, importedApp.DaoKeeper)
	require.Equal(t, daoAddresses, imported.DaoAddresses)
	require.ElementsMatch(t, freeGasAccounts, imported.FreeGasAccounts)
	for _, address := range freeGasAccounts {
		require.True(t, importedApp.DaoKeeper.CheckFreeGasAccount(importedCtx, address))
	}
}
