package keeper_test

import (
	"context"
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	wnftkeeper "github.com/openmetaearth/me-hub/x/wnft/keeper"
	"github.com/openmetaearth/me-hub/x/wnft/types"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/app/apptesting"
)

const wnftQueryCreator = "me139mq752delxv78jvtmwxhasyrycufsvr0mue6u"

func TestNftFilterRejectsInvalidOwnerForClassOwnerQuery(t *testing.T) {
	app := apptesting.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	goCtx := sdk.WrapSDKContext(ctx)
	msgServer := wnftkeeper.NewMsgServerImpl(app.WNFTKeeper, app.WNFTKeeper.Keeper)
	createWNFTClass(t, msgServer, goCtx, "owner-filter-class")

	_, err := app.WNFTKeeper.NftFilter(goCtx, &types.QueryNftFilterRequest{
		ClassId: "owner-filter-class",
		Owner:   "not-a-bech32-address",
	})

	require.Error(t, err)
	require.True(t, sdkerrors.ErrInvalidAddress.Is(err))
}

func TestNftFilterRejectsInvalidOwnerBeforeClassLookup(t *testing.T) {
	app := apptesting.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	goCtx := sdk.WrapSDKContext(ctx)

	_, err := app.WNFTKeeper.NftFilter(goCtx, &types.QueryNftFilterRequest{
		ClassId: "missing-owner-filter-class",
		Owner:   "not-a-bech32-address",
	})

	require.Error(t, err)
	require.True(t, sdkerrors.ErrInvalidAddress.Is(err))
}

func TestNftFilterRejectsInvalidOwnerForAllClassesOwnerQuery(t *testing.T) {
	app := apptesting.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	goCtx := sdk.WrapSDKContext(ctx)

	_, err := app.WNFTKeeper.NftFilter(goCtx, &types.QueryNftFilterRequest{
		Owner: "not-a-bech32-address",
	})

	require.Error(t, err)
	require.True(t, sdkerrors.ErrInvalidAddress.Is(err))
}

func TestNftFilterRejectsInvalidRequests(t *testing.T) {
	app := apptesting.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	goCtx := sdk.WrapSDKContext(ctx)

	_, err := app.WNFTKeeper.NftFilter(goCtx, nil)
	require.Error(t, err)
	require.True(t, sdkerrors.ErrInvalidRequest.Is(err))
}

func TestNftFilterRejectsUnsupportedFilterCombination(t *testing.T) {
	app := apptesting.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	goCtx := sdk.WrapSDKContext(ctx)

	_, err := app.WNFTKeeper.NftFilter(goCtx, &types.QueryNftFilterRequest{
		Owner:   wnftQueryCreator,
		TokenId: "1",
	})
	require.Error(t, err)
	require.True(t, sdkerrors.ErrInvalidRequest.Is(err))
}

func TestNftFilterReturnsClassOwnerResults(t *testing.T) {
	app := apptesting.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	goCtx := sdk.WrapSDKContext(ctx)
	msgServer := wnftkeeper.NewMsgServerImpl(app.WNFTKeeper, app.WNFTKeeper.Keeper)
	classID := "owner-filter-valid-class"

	createWNFTClass(t, msgServer, goCtx, classID)
	_, err := msgServer.MintNFT(goCtx, types.NewMsgMintNFT(
		classID,
		"1",
		"ipfs://owner-filter-token-1",
		"",
		wnftQueryCreator,
		wnftQueryCreator,
	))
	require.NoError(t, err)

	res, err := app.WNFTKeeper.NftFilter(goCtx, &types.QueryNftFilterRequest{
		ClassId: classID,
		Owner:   wnftQueryCreator,
	})

	require.NoError(t, err)
	require.Len(t, res.Nfts, 1)
	require.Equal(t, classID, res.Nfts[0].ClassId)
	require.Equal(t, "1", res.Nfts[0].TokenId)
	require.Equal(t, wnftQueryCreator, res.Nfts[0].Owner)
	require.Equal(t, "ipfs://owner-filter-token-1", res.Nfts[0].Uri)
}

func createWNFTClass(t *testing.T, msgServer types.MsgServer, goCtx context.Context, classID string) {
	t.Helper()

	_, err := msgServer.NewClass(goCtx, types.NewMsgNewClass(
		classID,
		wnftQueryCreator,
		"Owner Filter Class",
		"OWN",
		"",
		"",
		"",
		1,
	))
	require.NoError(t, err)
}
