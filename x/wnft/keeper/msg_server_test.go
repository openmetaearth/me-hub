package keeper_test

import (
	"testing"

	"cosmossdk.io/x/nft"
	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wnft/keeper"
	"github.com/openmetaearth/me-hub/x/wnft/types"
)

func setupMsgServer(t *testing.T) (*apptesting.KeeperTestHelper, types.MsgServer) {
	t.Helper()

	app := apptesting.Setup(t, false)
	ctx := app.GetBaseApp().NewContext(false, cometbftproto.Header{})

	err := app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	require.NoError(t, err)

	err = app.BankKeeper.SetParams(ctx, banktypes.DefaultParams())
	require.NoError(t, err)

	stakingParams := stakingtypes.DefaultParams()
	stakingParams.BondDenom = params.BaseDenom
	err = app.StakingKeeper.SetParams(ctx, stakingParams)
	require.NoError(t, err)

	helper := &apptesting.KeeperTestHelper{App: app, Ctx: ctx}
	return helper, keeper.NewMsgServerImpl(app.WNFTKeeper, app.WNFTKeeper.Keeper)
}

func TestMintNFTEnforcesClassTotalSupplyByMintCount(t *testing.T) {
	helper, msgServer := setupMsgServer(t)
	creator, _ := helper.NewAccount()
	receiver, _ := helper.NewAccount()

	_, err := msgServer.NewClass(helper.Ctx, &types.MsgNewClass{
		ClassId:     "class-total-supply",
		Sender:      creator.String(),
		Name:        "Total Supply",
		Symbol:      "TS",
		Description: "test",
		Uri:         "ipfs://class",
		UriHash:     "classhash",
		TotalSupply: 2,
	})
	require.NoError(t, err)

	for _, tokenID := range []string{"1", "2"} {
		_, err = msgServer.MintNFT(helper.Ctx, &types.MsgMintNFT{
			ClassId:  "class-total-supply",
			TokenId:  tokenID,
			Uri:      "ipfs://token-" + tokenID,
			UriHash:  "hash-" + tokenID,
			Creator:  creator.String(),
			Receiver: receiver.String(),
		})
		require.NoError(t, err)
	}

	require.Equal(t, uint64(2), helper.App.WNFTKeeper.GetTotalSupply(helper.Ctx, "class-total-supply"))

	_, err = msgServer.MintNFT(helper.Ctx, &types.MsgMintNFT{
		ClassId:  "class-total-supply",
		TokenId:  "3",
		Uri:      "ipfs://token-3",
		UriHash:  "hash-3",
		Creator:  creator.String(),
		Receiver: receiver.String(),
	})
	require.Error(t, err)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
	require.ErrorContains(t, err, "invalid token id")

	err = helper.App.WNFTKeeper.Keeper.Mint(helper.Ctx, nft.NFT{
		ClassId: "class-total-supply",
		Id:      "100",
		Uri:     "ipfs://token-100",
		UriHash: "hash-100",
	}, receiver)
	require.NoError(t, err)
	require.Equal(t, uint64(3), helper.App.WNFTKeeper.GetTotalSupply(helper.Ctx, "class-total-supply"))

	_, err = msgServer.MintNFT(helper.Ctx, &types.MsgMintNFT{
		ClassId:  "class-total-supply",
		TokenId:  "1",
		Uri:      "ipfs://token-1b",
		UriHash:  "hash-1b",
		Creator:  creator.String(),
		Receiver: receiver.String(),
	})
	require.Error(t, err)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
	require.ErrorContains(t, err, "total supply exceeded")
}
