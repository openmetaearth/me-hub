package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	wnftkeeper "github.com/openmetaearth/me-hub/x/wnft/keeper"
	"github.com/openmetaearth/me-hub/x/wnft/types"
	"github.com/stretchr/testify/require"
)

const wnftKeeperSender = "me139mq752delxv78jvtmwxhasyrycufsvr0mue6u"

func TestNewClassRejectsReservedKycClassID(t *testing.T) {
	app := apptesting.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	goCtx := sdk.WrapSDKContext(ctx)
	msgServer := wnftkeeper.NewMsgServerImpl(app.WNFTKeeper, app.WNFTKeeper.Keeper)

	beforeClass, beforeFound := app.WNFTKeeper.GetClass(ctx, types.ReservedKycClassID)

	_, err := msgServer.NewClass(goCtx, types.NewMsgNewClass(
		types.ReservedKycClassID,
		wnftKeeperSender,
		"KYC",
		"KYC",
		"",
		"",
		"",
		1,
	))

	require.Error(t, err)
	require.True(t, types.ErrReservedClassId.Is(err))

	afterClass, afterFound := app.WNFTKeeper.GetClass(ctx, types.ReservedKycClassID)
	require.Equal(t, beforeFound, afterFound)
	if beforeFound {
		require.Equal(t, beforeClass, afterClass)
	}
}
