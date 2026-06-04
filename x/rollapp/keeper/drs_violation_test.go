package keeper_test

import (
	common "github.com/openmetaearth/me-hub/x/common/types"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	sequencertypes "github.com/openmetaearth/me-hub/x/sequencer/types"
)

func (suite *RollappTestSuite) TestHandleDRSViolationFreezesRollappAndRevertsPendingState() {
	suite.SetupTest()
	ctx := &suite.Ctx
	keeper := suite.App.RollappKeeper

	rollappID := suite.CreateDefaultRollapp()
	proposer := suite.CreateDefaultSequencer(*ctx, rollappID)
	suite.CreateDefaultSequencer(*ctx, rollappID)

	_, err := suite.PostStateUpdate(*ctx, rollappID, proposer, 1, 10)
	suite.Require().NoError(err)

	stateInfo, err := keeper.FindStateInfoByHeight(*ctx, rollappID, 1)
	suite.Require().NoError(err)
	suite.Require().Equal(common.Status_PENDING, stateInfo.Status)

	err = keeper.HandleDRSViolation(*ctx, rollappID)
	suite.Require().NoError(err)

	rollapp, found := keeper.GetRollapp(*ctx, rollappID)
	suite.Require().True(found)
	suite.Require().True(rollapp.Frozen)

	stateInfo, err = keeper.FindStateInfoByHeight(*ctx, rollappID, 1)
	suite.Require().NoError(err)
	suite.Require().Equal(common.Status_REVERTED, stateInfo.Status)

	sequencers := suite.App.SequencerKeeper.GetSequencersByRollapp(*ctx, rollappID)
	for _, sequencer := range sequencers {
		suite.Require().Equal(sequencertypes.Unbonded, sequencer.Status)
	}
}

func (suite *RollappTestSuite) TestUpdateStateRejectsFrozenRollapp() {
	suite.SetupTest()
	ctx := &suite.Ctx
	keeper := suite.App.RollappKeeper

	rollappID := suite.CreateDefaultRollapp()
	proposer := suite.CreateDefaultSequencer(*ctx, rollappID)

	rollapp, found := keeper.GetRollapp(*ctx, rollappID)
	suite.Require().True(found)
	rollapp.Frozen = true
	keeper.SetRollapp(*ctx, rollapp)

	_, err := suite.PostStateUpdate(*ctx, rollappID, proposer, 1, 10)
	suite.Require().ErrorIs(err, rollapptypes.ErrRollappFrozen)
}
