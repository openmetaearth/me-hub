package kyc_test

import (
	"testing"

	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	did "github.com/openmetaearth/me-hub/x/did"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	kyc "github.com/openmetaearth/me-hub/x/kyc"
	kyctypes "github.com/openmetaearth/me-hub/x/kyc/types"
	"github.com/stretchr/testify/require"
)

func TestKycGenesisPreservesInactiveServiceStatus(t *testing.T) {
	app := apptesting.Setup(t, false)
	ctx := app.GetBaseApp().NewContext(false, cometbftproto.Header{})
	inactiveService := didtypes.Service{
		Sid:         kyctypes.ModuleName,
		Name:        kyctypes.ModuleName,
		Description: "inactive service restored by DID genesis",
		Status:      didtypes.SERVICE_STATUS_INACTIVE,
	}

	did.InitGenesis(ctx, app.DidKeeper, didtypes.GenesisState{Svcs: []didtypes.Service{inactiveService}})
	service, found := app.KycKeeper.GetService(ctx)
	require.True(t, found)
	require.Equal(t, didtypes.SERVICE_STATUS_INACTIVE, service.Status)

	kyc.InitGenesis(ctx, *app.KycKeeper, kyctypes.GenesisState{})

	service, found = app.KycKeeper.GetService(ctx)
	require.True(t, found)
	require.Equal(t, didtypes.SERVICE_STATUS_INACTIVE, service.Status)
}

func TestKycGenesisDefaultsServiceActiveWhenMissing(t *testing.T) {
	app := apptesting.Setup(t, false)
	ctx := app.GetBaseApp().NewContext(false, cometbftproto.Header{})

	app.DidKeeper.DeleteService(ctx, kyctypes.ModuleName)
	_, found := app.KycKeeper.GetService(ctx)
	require.False(t, found)

	kyc.InitGenesis(ctx, *app.KycKeeper, kyctypes.GenesisState{})

	service, found := app.KycKeeper.GetService(ctx)
	require.True(t, found)
	require.Equal(t, didtypes.SERVICE_STATUS_ACTIVE, service.Status)
}
