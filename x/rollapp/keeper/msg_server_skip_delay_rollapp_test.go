package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/stretchr/testify/suite"
)

func TestSkipDelayRollappTestSuite(t *testing.T) {
	suite.Run(t, new(SkipDelayRollappTestSuite))
}

type SkipDelayRollappTestSuite struct {
	RollappTestSuite
}

// TestSkipDelayRollapp_GlobalDaoAllowed verifies that GlobalDao can toggle skip-delay
func (suite *SkipDelayRollappTestSuite) TestSkipDelayRollapp_GlobalDaoAllowed() {
	suite.SetupTest()
	suite.InitializeDao()
	goCtx := sdk.WrapSDKContext(suite.Ctx)

	// Create a rollapp first
	rollappId := "rollapp1"
	rollapp := types.Rollapp{
		RollappId:     rollappId,
		Creator:       alice,
		MaxSequencers: 1,
	}
	suite.App.RollappKeeper.SetRollapp(suite.Ctx, rollapp)

	// GlobalDao should be allowed to toggle skip delay
	msg := types.MsgSkipDelayRollapp{
		Creator:   suite.Dao.GlobalDao,
		RollappId: rollappId,
		Skip:      true,
	}

	_, err := suite.msgServer.SkipDelayRollapp(goCtx, &msg)
	suite.Require().NoError(err)

	// Verify skip-delay is enabled
	isSkip := suite.App.RollappKeeper.IsSkipDelayRollapp(suite.Ctx, rollappId)
	suite.Require().True(isSkip, "skip-delay should be enabled for rollapp after GlobalDao sets it")

	// Disable skip
	msg.Skip = false
	_, err = suite.msgServer.SkipDelayRollapp(goCtx, &msg)
	suite.Require().NoError(err)

	isSkip = suite.App.RollappKeeper.IsSkipDelayRollapp(suite.Ctx, rollappId)
	suite.Require().False(isSkip, "skip-delay should be disabled after GlobalDao unsets it")
}

// TestSkipDelayRollapp_MeidDaoRejected verifies that MeidDao cannot toggle skip-delay
// This is the core fix for issue #105: MeidDao should NOT be able to bypass fraud dispute delays
func (suite *SkipDelayRollappTestSuite) TestSkipDelayRollapp_MeidDaoRejected() {
	suite.SetupTest()
	suite.InitializeDao()
	goCtx := sdk.WrapSDKContext(suite.Ctx)

	// Create a rollapp first
	rollappId := "rollapp1"
	rollapp := types.Rollapp{
		RollappId:     rollappId,
		Creator:       alice,
		MaxSequencers: 1,
	}
	suite.App.RollappKeeper.SetRollapp(suite.Ctx, rollapp)

	// MeidDao should NOT be allowed to toggle skip delay
	msg := types.MsgSkipDelayRollapp{
		Creator:   suite.Dao.MeidDao,
		RollappId: rollappId,
		Skip:      true,
	}

	_, err := suite.msgServer.SkipDelayRollapp(goCtx, &msg)
	suite.Require().Error(err, "MeidDao should not be allowed to toggle skip-delay")
	suite.Require().ErrorIs(err, types.ErrCheckGlobalDao)

	// Verify skip-delay is NOT enabled
	isSkip := suite.App.RollappKeeper.IsSkipDelayRollapp(suite.Ctx, rollappId)
	suite.Require().False(isSkip, "skip-delay should not be enabled when MeidDao attempts to set it")
}

// TestSkipDelayRollapp_UnauthorizedRejected verifies that non-DAO accounts cannot toggle skip-delay
func (suite *SkipDelayRollappTestSuite) TestSkipDelayRollapp_UnauthorizedRejected() {
	suite.SetupTest()
	suite.InitializeDao()
	goCtx := sdk.WrapSDKContext(suite.Ctx)

	// Create a rollapp
	rollappId := "rollapp1"
	rollapp := types.Rollapp{
		RollappId:     rollappId,
		Creator:       alice,
		MaxSequencers: 1,
	}
	suite.App.RollappKeeper.SetRollapp(suite.Ctx, rollapp)

	// Random account should NOT be allowed
	msg := types.MsgSkipDelayRollapp{
		Creator:   alice,
		RollappId: rollappId,
		Skip:      true,
	}

	_, err := suite.msgServer.SkipDelayRollapp(goCtx, &msg)
	suite.Require().Error(err, "non-DAO account should not be allowed to toggle skip-delay")
	suite.Require().ErrorIs(err, types.ErrCheckGlobalDao)
}

// TestSkipDelayRollapp_UnknownRollapp verifies that skip-delay cannot be set for non-existent rollapps
func (suite *SkipDelayRollappTestSuite) TestSkipDelayRollapp_UnknownRollapp() {
	suite.SetupTest()
	suite.InitializeDao()
	goCtx := sdk.WrapSDKContext(suite.Ctx)

	msg := types.MsgSkipDelayRollapp{
		Creator:   suite.Dao.GlobalDao,
		RollappId: "nonexistent-rollapp",
		Skip:      true,
	}

	_, err := suite.msgServer.SkipDelayRollapp(goCtx, &msg)
	suite.Require().Error(err, "should not be able to set skip-delay for non-existent rollapp")
	suite.Require().ErrorIs(err, types.ErrUnknownRollappID)
}

// TestSkipDelayRollapp_TogglePreventsFraudBypass verifies the security fix:
// After enabling and then disabling skip-delay, IsSkipDelayRollapp returns false
// This prevents the selective fraud attack described in issue #105
func (suite *SkipDelayRollappTestSuite) TestSkipDelayRollapp_TogglePreventsFraudBypass() {
	suite.SetupTest()
	suite.InitializeDao()
	goCtx := sdk.WrapSDKContext(suite.Ctx)

	rollappId := "rollapp1"
	rollapp := types.Rollapp{
		RollappId:     rollappId,
		Creator:       alice,
		MaxSequencers: 1,
	}
	suite.App.RollappKeeper.SetRollapp(suite.Ctx, rollapp)

	// Initially, skip should be false (default)
	isSkip := suite.App.RollappKeeper.IsSkipDelayRollapp(suite.Ctx, rollappId)
	suite.Require().False(isSkip, "skip-delay should be false by default")

	// Enable skip-delay
	msg := types.MsgSkipDelayRollapp{
		Creator:   suite.Dao.GlobalDao,
		RollappId: rollappId,
		Skip:      true,
	}
	_, err := suite.msgServer.SkipDelayRollapp(goCtx, &msg)
	suite.Require().NoError(err)

	isSkip = suite.App.RollappKeeper.IsSkipDelayRollapp(suite.Ctx, rollappId)
	suite.Require().True(isSkip, "skip-delay should be enabled")

	// Disable skip-delay
	msg.Skip = false
	_, err = suite.msgServer.SkipDelayRollapp(goCtx, &msg)
	suite.Require().NoError(err)

	isSkip = suite.App.RollappKeeper.IsSkipDelayRollapp(suite.Ctx, rollappId)
	suite.Require().False(isSkip, "skip-delay should be disabled after toggle")
}

// TestIsSkipDelayRollapp_NonExistentRollapp verifies IsSkipDelayRollapp returns false for unknown rollapps
func (suite *SkipDelayRollappTestSuite) TestIsSkipDelayRollapp_NonExistentRollapp() {
	suite.SetupTest()

	isSkip := suite.App.RollappKeeper.IsSkipDelayRollapp(suite.Ctx, "nonexistent")
	suite.Require().False(isSkip, "IsSkipDelayRollapp should return false for non-existent rollapp")
}
