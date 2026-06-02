package dao_test

import (
	"bytes"
	"testing"

	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	dao "github.com/openmetaearth/me-hub/x/dao"
	daotypes "github.com/openmetaearth/me-hub/x/dao/types"
	"github.com/stretchr/testify/require"
)

func TestExportImportPreservesFreeGasAccounts(t *testing.T) {
	sourceApp := apptesting.Setup(t, false)
	sourceCtx := sourceApp.NewContext(false, cometbftproto.Header{Height: 1})
	daoAddresses := testDaoAddresses()
	freeGasAccounts := []string{testDaoAddress(5), testDaoAddress(6)}

	dao.InitGenesis(sourceCtx, sourceApp.DaoKeeper, daotypes.GenesisState{
		DaoAddresses:    daoAddresses,
		FreeGasAccounts: freeGasAccounts,
	})
	for _, address := range freeGasAccounts {
		require.True(t, sourceApp.DaoKeeper.CheckFreeGasAccount(sourceCtx, address))
	}

	exported := dao.ExportGenesis(sourceCtx, sourceApp.DaoKeeper)
	require.Equal(t, daoAddresses, exported.DaoAddresses)
	require.ElementsMatch(t, freeGasAccounts, exported.FreeGasAccounts)

	importApp := apptesting.Setup(t, false)
	importCtx := importApp.NewContext(false, cometbftproto.Header{Height: 1})
	dao.InitGenesis(importCtx, importApp.DaoKeeper, *exported)

	importedDaoAddresses, found := importApp.DaoKeeper.GetDaoAddresses(importCtx)
	require.True(t, found)
	require.Equal(t, daoAddresses, importedDaoAddresses)
	for _, address := range freeGasAccounts {
		require.True(t, importApp.DaoKeeper.CheckFreeGasAccount(importCtx, address))
	}
}

func testDaoAddresses() daotypes.DaoAddresses {
	return daotypes.DaoAddresses{
		GlobalDao:      testDaoAddress(1),
		MeidDao:        testDaoAddress(2),
		DevOperator:    testDaoAddress(3),
		AirdropAddress: testDaoAddress(4),
	}
}

func testDaoAddress(seed byte) string {
	return types.AccAddress(bytes.Repeat([]byte{seed}, 20)).String()
}
