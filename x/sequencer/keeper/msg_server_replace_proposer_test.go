package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	commontypes "github.com/openmetaearth/me-hub/x/common/types"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
)

func (suite *SequencerTestSuite) TestReplaceProposerAllowsHandoffAtLastStateInfoEndHeight() {
	suite.SetupTest()
	goCtx := sdk.WrapSDKContext(suite.Ctx)
	rollappId, oldProposer, newProposer := suite.setupReplaceProposerHandoffState(1, 10)

	msg, err := types.NewMsgReplaceProposerRequest(alice, rollappId, oldProposer, newProposer, 10)
	suite.Require().NoError(err)

	_, err = suite.msgServer.ReplaceProposer(goCtx, msg)

	suite.Require().NoError(err)
	pending, err := suite.App.SequencerKeeper.GetReplaceProposer(suite.Ctx, rollappId)
	suite.Require().NoError(err)
	suite.Require().NotNil(pending)
	suite.Require().Equal(int64(10), pending.ReplaceProposer.BlockHeight)
	suite.Require().NoError(suite.App.SequencerKeeper.IsExceedAuthoredBlockHeight(suite.Ctx, rollappId, newProposer, 11, 1))
}

func (suite *SequencerTestSuite) TestReplaceProposerRejectsHandoffBeforeLastStateInfoEndHeight() {
	suite.SetupTest()
	goCtx := sdk.WrapSDKContext(suite.Ctx)
	rollappId, oldProposer, newProposer := suite.setupReplaceProposerHandoffState(1, 10)

	msg, err := types.NewMsgReplaceProposerRequest(alice, rollappId, oldProposer, newProposer, 9)
	suite.Require().NoError(err)

	_, err = suite.msgServer.ReplaceProposer(goCtx, msg)

	suite.Require().Error(err)
	suite.Require().ErrorContains(err, "must be at least last state info end height 10")
	pending, pendingErr := suite.App.SequencerKeeper.GetReplaceProposer(suite.Ctx, rollappId)
	suite.Require().NoError(pendingErr)
	suite.Require().Nil(pending)
}

func (suite *SequencerTestSuite) setupReplaceProposerHandoffState(startHeight, numBlocks uint64) (string, string, string) {
	rollappId := suite.CreateDefaultRollapp()
	oldProposer := suite.CreateDefaultSequencer(suite.Ctx, rollappId)
	newProposer := suite.CreateDefaultSequencer(suite.Ctx, rollappId)
	stateInfoIndex := rollapptypes.StateInfoIndex{RollappId: rollappId, Index: 1}

	suite.App.RollappKeeper.SetStateInfo(suite.Ctx, rollapptypes.StateInfo{
		StateInfoIndex: stateInfoIndex,
		Sequencer:      oldProposer,
		StartHeight:    startHeight,
		NumBlocks:      numBlocks,
		Status:         commontypes.Status_FINALIZED,
	})
	suite.App.RollappKeeper.SetLatestStateInfoIndex(suite.Ctx, stateInfoIndex)

	return rollappId, oldProposer, newProposer
}
