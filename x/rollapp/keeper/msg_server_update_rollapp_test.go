package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/rollapp/types"
)

func (suite *RollappTestSuite) TestUpdateRollappRejectsWhenRollappsDisabled() {
	suite.SetupTest()
	goCtx := sdk.WrapSDKContext(suite.Ctx)

	msg := types.MsgCreateRollapp{
		Creator:               alice,
		RollappId:             "rollapp_1234-1",
		MaxSequencers:         2,
		PermissionedAddresses: []string{alice},
	}
	_, err := suite.msgServer.CreateRollapp(goCtx, &msg)
	suite.Require().NoError(err)

	params := suite.App.RollappKeeper.GetParams(suite.Ctx)
	params.RollappsEnabled = false
	suite.App.RollappKeeper.SetParams(suite.Ctx, params)

	_, err = suite.msgServer.UpdateRollapp(goCtx, &types.MsgUpdateRollapp{
		Creator:       alice,
		RollappId:     msg.RollappId,
		MaxSequencers: 3,
	})
	suite.Require().ErrorIs(err, types.ErrRollappsDisabled)
}

func (suite *RollappTestSuite) TestUpdateRollappValidatesMutatedStateBeforePersisting() {
	suite.SetupTest()
	goCtx := sdk.WrapSDKContext(suite.Ctx)

	msg := types.MsgCreateRollapp{
		Creator:               alice,
		RollappId:             "rollapp_1234-1",
		MaxSequencers:         2,
		PermissionedAddresses: []string{alice},
	}
	_, err := suite.msgServer.CreateRollapp(goCtx, &msg)
	suite.Require().NoError(err)

	_, err = suite.msgServer.UpdateRollapp(goCtx, &types.MsgUpdateRollapp{
		Creator:               alice,
		RollappId:             msg.RollappId,
		MaxSequencers:         1,
		PermissionedAddresses: []string{alice, bob},
	})
	suite.Require().ErrorIs(err, types.ErrTooManyPermissionedAddresses)

	rollapp, found := suite.App.RollappKeeper.GetRollapp(suite.Ctx, msg.RollappId)
	suite.Require().True(found)
	suite.Require().Equal(uint64(2), rollapp.MaxSequencers)
	suite.Require().Equal([]string{alice}, rollapp.PermissionedAddresses)
}
