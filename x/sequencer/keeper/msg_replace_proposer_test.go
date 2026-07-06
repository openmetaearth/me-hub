package keeper_test

import (
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
)

func (suite *SequencerTestSuite) setLatestStateInfo(rollappID, sequencerAddr string) {
	stateInfo := rollapptypes.NewStateInfo(
		rollappID,
		1,
		sequencerAddr,
		1,
		10,
		"",
		0,
		0,
		rollapptypes.BlockDescriptors{},
	)
	suite.App.RollappKeeper.SetStateInfo(suite.Ctx, *stateInfo)
	suite.App.RollappKeeper.SetLatestStateInfoIndex(suite.Ctx, stateInfo.GetIndex())
}

func (suite *SequencerTestSuite) TestReplaceProposerRejectsNewProposerFromDifferentRollapp() {
	suite.SetupTest()

	rollappID := suite.CreateDefaultRollapp()
	otherRollappID := suite.CreateDefaultRollapp()
	oldProposer := suite.CreateDefaultSequencer(suite.Ctx, rollappID)
	newProposer := suite.CreateDefaultSequencer(suite.Ctx, otherRollappID)
	suite.setLatestStateInfo(rollappID, oldProposer)

	_, err := suite.msgServer.ReplaceProposer(suite.Ctx, &types.MsgReplaceProposerRequest{
		Creator: alice,
		ReplaceProposer: &types.MsgRepalceProposer{
			RollappId:   rollappID,
			OldProposer: oldProposer,
			NewProposer: newProposer,
			BlockHeight: 20,
		},
	})

	suite.Require().ErrorIs(err, types.ErrInvalidRequest)
	suite.Require().False(suite.App.SequencerKeeper.IsHasReplaceProposer(suite.Ctx, rollappID))
}

func (suite *SequencerTestSuite) TestReplaceProposerRejectsOldProposerFromDifferentRollapp() {
	suite.SetupTest()

	rollappID := suite.CreateDefaultRollapp()
	otherRollappID := suite.CreateDefaultRollapp()
	newProposer := suite.CreateDefaultSequencer(suite.Ctx, rollappID)
	oldProposer := suite.CreateDefaultSequencer(suite.Ctx, otherRollappID)
	suite.setLatestStateInfo(rollappID, newProposer)

	_, err := suite.msgServer.ReplaceProposer(suite.Ctx, &types.MsgReplaceProposerRequest{
		Creator: alice,
		ReplaceProposer: &types.MsgRepalceProposer{
			RollappId:   rollappID,
			OldProposer: oldProposer,
			NewProposer: newProposer,
			BlockHeight: 20,
		},
	})

	suite.Require().ErrorIs(err, types.ErrInvalidRequest)
	suite.Require().False(suite.App.SequencerKeeper.IsHasReplaceProposer(suite.Ctx, rollappID))
}

func (suite *SequencerTestSuite) TestReplaceProposerAcceptsSameRollappSequencers() {
	suite.SetupTest()

	rollappID := suite.CreateDefaultRollapp()
	oldProposer := suite.CreateDefaultSequencer(suite.Ctx, rollappID)
	newProposer := suite.CreateDefaultSequencer(suite.Ctx, rollappID)
	suite.setLatestStateInfo(rollappID, oldProposer)

	_, err := suite.msgServer.ReplaceProposer(suite.Ctx, &types.MsgReplaceProposerRequest{
		Creator: alice,
		ReplaceProposer: &types.MsgRepalceProposer{
			RollappId:   rollappID,
			OldProposer: oldProposer,
			NewProposer: newProposer,
			BlockHeight: 20,
		},
	})

	suite.Require().NoError(err)
	suite.Require().True(suite.App.SequencerKeeper.IsHasReplaceProposer(suite.Ctx, rollappID))
}
