package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/rollapp/types"
)

func (suite *RollappTestSuite) TestSkipDelayRollappRequiresGlobalDao() {
	suite.SetupTest()
	suite.InitializeDao()
	goCtx := sdk.WrapSDKContext(suite.Ctx)
	rollappID := suite.CreateDefaultRollapp()

	_, err := suite.msgServer.SkipDelayRollapp(goCtx, &types.MsgSkipDelayRollapp{
		Creator:   suite.Dao.MeidDao,
		RollappId: rollappID,
		Skip:      true,
	})
	suite.Require().ErrorIs(err, types.ErrCheckGlobalDao)
	suite.Require().False(suite.App.RollappKeeper.IsSkipDelayRollapp(suite.Ctx, rollappID))

	_, err = suite.msgServer.SkipDelayRollapp(goCtx, &types.MsgSkipDelayRollapp{
		Creator:   suite.Dao.GlobalDao,
		RollappId: rollappID,
		Skip:      true,
	})
	suite.Require().NoError(err)
	suite.Require().True(suite.App.RollappKeeper.IsSkipDelayRollapp(suite.Ctx, rollappID))
}
