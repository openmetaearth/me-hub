package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/rollapp/types"
)

func (suite *RollappTestSuite) TestUpdateRollappRejectsNonCreator() {
	suite.SetupTest()
	goCtx := sdk.WrapSDKContext(suite.Ctx)

	rollapp := types.Rollapp{
		RollappId:             "rollapp1",
		Creator:               alice,
		MaxSequencers:         1,
		PermissionedAddresses: []string{alice},
		ChannelId:             "channel-original",
	}
	suite.App.RollappKeeper.SetRollapp(suite.Ctx, rollapp)

	_, err := suite.msgServer.UpdateRollapp(
		goCtx,
		types.NewMsgUpdateRollapp(bob, rollapp.RollappId, "channel-evil", 2, []string{bob}),
	)

	suite.Require().ErrorIs(err, types.ErrUnauthorizedRollappCreator)
	suite.Require().ErrorContains(err, "only rollapp creator or DAO can update rollapp")
	stored := suite.App.RollappKeeper.MustGetRollapp(suite.Ctx, rollapp.RollappId)
	suite.Require().Equal(rollapp.ChannelId, stored.ChannelId)
	suite.Require().Equal(rollapp.MaxSequencers, stored.MaxSequencers)
	suite.Require().Equal(rollapp.PermissionedAddresses, stored.PermissionedAddresses)
}

func (suite *RollappTestSuite) TestUpdateRollappAllowsCreator() {
	suite.SetupTest()
	goCtx := sdk.WrapSDKContext(suite.Ctx)

	rollapp := types.Rollapp{
		RollappId:     "rollapp1",
		Creator:       alice,
		MaxSequencers: 1,
		ChannelId:     "channel-original",
	}
	suite.App.RollappKeeper.SetRollapp(suite.Ctx, rollapp)

	_, err := suite.msgServer.UpdateRollapp(
		goCtx,
		types.NewMsgUpdateRollapp(alice, rollapp.RollappId, "channel-updated", 2, []string{bob}),
	)

	suite.Require().NoError(err)
	stored := suite.App.RollappKeeper.MustGetRollapp(suite.Ctx, rollapp.RollappId)
	suite.Require().Equal("channel-updated", stored.ChannelId)
	suite.Require().Equal(uint64(2), stored.MaxSequencers)
	suite.Require().Equal([]string{bob}, stored.PermissionedAddresses)
}

func (suite *RollappTestSuite) TestUpdateRollappAllowsDao() {
	suite.SetupTest()
	suite.InitializeDao()
	goCtx := sdk.WrapSDKContext(suite.Ctx)

	rollapp := types.Rollapp{
		RollappId:     "rollapp1",
		Creator:       alice,
		MaxSequencers: 1,
		ChannelId:     "channel-original",
	}
	suite.App.RollappKeeper.SetRollapp(suite.Ctx, rollapp)

	_, err := suite.msgServer.UpdateRollapp(
		goCtx,
		types.NewMsgUpdateRollapp(suite.Dao.GlobalDao, rollapp.RollappId, "channel-dao", 2, []string{bob}),
	)

	suite.Require().NoError(err)
	stored := suite.App.RollappKeeper.MustGetRollapp(suite.Ctx, rollapp.RollappId)
	suite.Require().Equal("channel-dao", stored.ChannelId)
	suite.Require().Equal(uint64(2), stored.MaxSequencers)
	suite.Require().Equal([]string{bob}, stored.PermissionedAddresses)
}
