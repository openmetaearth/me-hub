package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	nftkeeper "github.com/openmetaearth/me-hub/x/wnft/keeper"
	"github.com/openmetaearth/me-hub/x/wnft/types"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/app/apptesting"
)

const wnftCreator = "me139mq752delxv78jvtmwxhasyrycufsvr0mue6u"

func TestMintNFTRejectsLeadingZeroDuplicateTokenID(t *testing.T) {
	app := apptesting.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	goCtx := sdk.WrapSDKContext(ctx)
	msgServer := nftkeeper.NewMsgServerImpl(app.WNFTKeeper, app.WNFTKeeper.Keeper)

	classID := "limited-class"
	_, err := msgServer.NewClass(goCtx, types.NewMsgNewClass(
		classID,
		wnftCreator,
		"Limited Class",
		"LIMITED",
		"",
		"",
		"",
		1,
	))
	require.NoError(t, err)

	_, err = msgServer.MintNFT(goCtx, types.NewMsgMintNFT(
		classID,
		"1",
		"ipfs://token-1",
		"",
		wnftCreator,
		wnftCreator,
	))
	require.NoError(t, err)

	_, err = msgServer.MintNFT(goCtx, types.NewMsgMintNFT(
		classID,
		"01",
		"ipfs://token-01",
		"",
		wnftCreator,
		wnftCreator,
	))
	require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)

	_, found := app.WNFTKeeper.GetNFT(ctx, classID, "01")
	require.False(t, found)
}

func TestMintNFTRejectsNonCanonicalDecimalTokenID(t *testing.T) {
	app := apptesting.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	goCtx := sdk.WrapSDKContext(ctx)
	msgServer := nftkeeper.NewMsgServerImpl(app.WNFTKeeper, app.WNFTKeeper.Keeper)

	classID := "canonical-class"
	_, err := msgServer.NewClass(goCtx, types.NewMsgNewClass(
		classID,
		wnftCreator,
		"Canonical Class",
		"CANON",
		"",
		"",
		"",
		10,
	))
	require.NoError(t, err)

	for _, tokenID := range []string{"01", "010"} {
		_, err = msgServer.MintNFT(goCtx, types.NewMsgMintNFT(
			classID,
			tokenID,
			"ipfs://token",
			"",
			wnftCreator,
			wnftCreator,
		))
		require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)

		_, found := app.WNFTKeeper.GetNFT(ctx, classID, tokenID)
		require.False(t, found)
	}
}
