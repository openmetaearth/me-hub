package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/rollapp/types"
)

func (suite *RollappTestSuite) TestUpdateRollapp() {
	suite.SetupTest()
	goCtx := sdk.WrapSDKContext(suite.Ctx)

	original := types.Rollapp{
		RollappId:             "rollapp_1234-1",
		Creator:               alice,
		MaxSequencers:         2,
		PermissionedAddresses: []string{bob, carol},
		ChannelId:             "channel-0",
	}
	suite.App.RollappKeeper.SetRollapp(suite.Ctx, original)

	updateMsg := types.MsgUpdateRollapp{
		Creator:               alice,
		RollappId:             original.RollappId,
		MaxSequencers:         3,
		PermissionedAddresses: []string{bob},
	}

	updateResp, err := suite.msgServer.UpdateRollapp(goCtx, &updateMsg)
	suite.Require().NoError(err)
	suite.Require().EqualValues(types.MsgUpdateRollappResponse{}, *updateResp)

	updated, found := suite.App.RollappKeeper.GetRollapp(suite.Ctx, original.RollappId)
	suite.Require().True(found)
	suite.Require().EqualValues(original.RollappId, updated.RollappId)
	suite.Require().EqualValues(original.Creator, updated.Creator)
	suite.Require().EqualValues(uint64(3), updated.MaxSequencers)
	suite.Require().EqualValues([]string{bob}, updated.PermissionedAddresses)
	suite.Require().EqualValues(original.ChannelId, updated.ChannelId)
}

func (suite *RollappTestSuite) TestUpdateRollappUnknownRollappID() {
	suite.SetupTest()
	goCtx := sdk.WrapSDKContext(suite.Ctx)

	updateMsg := types.MsgUpdateRollapp{
		Creator:   alice,
		RollappId: "unknown-rollapp",
	}

	_, err := suite.msgServer.UpdateRollapp(goCtx, &updateMsg)
	suite.EqualError(err, types.ErrUnknownRollappID.Error())
}

func (suite *RollappTestSuite) TestUpdateRollappWhenDisabled() {
	suite.SetupTest()
	goCtx := sdk.WrapSDKContext(suite.Ctx)

	original := types.Rollapp{
		RollappId:     "rollapp_1234-1",
		Creator:       alice,
		MaxSequencers: 1,
	}
	suite.App.RollappKeeper.SetRollapp(suite.Ctx, original)

	params := suite.App.RollappKeeper.GetParams(suite.Ctx)
	params.RollappsEnabled = false
	suite.App.RollappKeeper.SetParams(suite.Ctx, params)

	updateMsg := types.MsgUpdateRollapp{
		Creator:       alice,
		RollappId:     original.RollappId,
		MaxSequencers: 2,
	}

	_, err := suite.msgServer.UpdateRollapp(goCtx, &updateMsg)
	suite.EqualError(err, types.ErrRollappsDisabled.Error())

	stored, found := suite.App.RollappKeeper.GetRollapp(suite.Ctx, original.RollappId)
	suite.Require().True(found)
	suite.Require().EqualValues(original, stored)
}

func (suite *RollappTestSuite) TestUpdateRollappFrozen() {
	suite.SetupTest()
	goCtx := sdk.WrapSDKContext(suite.Ctx)

	original := types.Rollapp{
		RollappId:     "rollapp_1234-1",
		Creator:       alice,
		MaxSequencers: 1,
		Frozen:        true,
	}
	suite.App.RollappKeeper.SetRollapp(suite.Ctx, original)

	updateMsg := types.MsgUpdateRollapp{
		Creator:       alice,
		RollappId:     original.RollappId,
		MaxSequencers: 2,
	}

	_, err := suite.msgServer.UpdateRollapp(goCtx, &updateMsg)
	suite.EqualError(err, types.ErrRollappJailed.Error())

	stored, found := suite.App.RollappKeeper.GetRollapp(suite.Ctx, original.RollappId)
	suite.Require().True(found)
	suite.Require().EqualValues(original, stored)
}

func (suite *RollappTestSuite) TestUpdateRollappUnauthorizedCreator() {
	suite.SetupTest()
	goCtx := sdk.WrapSDKContext(suite.Ctx)

	original := types.Rollapp{
		RollappId:     "rollapp_1234-1",
		Creator:       alice,
		MaxSequencers: 1,
	}
	suite.App.RollappKeeper.SetRollapp(suite.Ctx, original)

	updateMsg := types.MsgUpdateRollapp{
		Creator:       bob,
		RollappId:     original.RollappId,
		MaxSequencers: 2,
	}

	_, err := suite.msgServer.UpdateRollapp(goCtx, &updateMsg)
	suite.EqualError(err, types.ErrUnauthorizedRollappCreator.Error())

	stored, found := suite.App.RollappKeeper.GetRollapp(suite.Ctx, original.RollappId)
	suite.Require().True(found)
	suite.Require().EqualValues(original, stored)
}
