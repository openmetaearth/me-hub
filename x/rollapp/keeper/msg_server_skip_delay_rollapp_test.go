package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/rollapp/types"
)

func (suite *RollappTestSuite) TestSkipDelayRollappRequiresGlobalDao() {
	suite.SetupTest()
	suite.InitializeDao()
	goCtx := sdk.WrapSDKContext(suite.Ctx)

	rollapp := types.Rollapp{
		RollappId:     "rollapp1",
		Creator:       alice,
		Version:       1,
		MaxSequencers: 1,
	}
	suite.App.RollappKeeper.SetRollapp(suite.Ctx, rollapp)

	_, err := suite.msgServer.SkipDelayRollapp(goCtx, types.NewMsgSkipDelayRollapp(suite.Dao.MeidDao, rollapp.RollappId, true))
	suite.Require().ErrorIs(err, types.ErrCheckGlobalDao)
	suite.Require().False(suite.App.RollappKeeper.IsSkipDelayRollapp(suite.Ctx, rollapp.RollappId))

	_, err = suite.msgServer.SkipDelayRollapp(goCtx, types.NewMsgSkipDelayRollapp(suite.Dao.GlobalDao, rollapp.RollappId, true))
	suite.Require().NoError(err)
	suite.Require().True(suite.App.RollappKeeper.IsSkipDelayRollapp(suite.Ctx, rollapp.RollappId))
}
