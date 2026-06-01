package keeper_test

import (
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
)

func (suite *SequencerTestSuite) TestReplaceProposerStaleRequestClearedAfterOldProposerUnbonds() {
	suite.SetupTest()

	rollappId := suite.CreateDefaultRollapp()
	oldProposer := suite.CreateDefaultSequencer(suite.Ctx, rollappId)
	newProposer := suite.CreateDefaultSequencer(suite.Ctx, rollappId)
	fallbackProposer := suite.CreateDefaultSequencer(suite.Ctx, rollappId)

	// MsgRepalceProposer is the existing generated type name in this module.
	err := suite.App.SequencerKeeper.SetReplaceProposer(suite.Ctx, &types.MsgRepalceProposer{
		RollappId:   rollappId,
		OldProposer: oldProposer,
		NewProposer: newProposer,
		BlockHeight: 10,
	})
	suite.Require().NoError(err)

	_, err = suite.msgServer.Unbond(suite.Ctx, &types.MsgUnbond{Creator: oldProposer})
	suite.Require().NoError(err)

	oldSequencer, found := suite.App.SequencerKeeper.GetSequencer(suite.Ctx, oldProposer)
	suite.Require().True(found)
	suite.Require().Equal(types.Unbonding, oldSequencer.Status)
	suite.Require().False(oldSequencer.Proposer)

	err = suite.App.SequencerKeeper.ProcSequencerByPendingStates(
		suite.Ctx,
		rollappId,
		oldProposer,
		&rollapptypes.StateInfo{
			StartHeight: 10,
			NumBlocks:   1,
		},
	)
	suite.Require().NoError(err)
	suite.Require().False(suite.App.SequencerKeeper.IsHasReplaceProposer(suite.Ctx, rollappId))

	currentProposer := ""
	candidateProposer := ""
	for _, sequencer := range suite.App.SequencerKeeper.GetSequencersByRollappByStatus(suite.Ctx, rollappId, types.Bonded) {
		if sequencer.Proposer {
			currentProposer = sequencer.SequencerAddress
		} else {
			candidateProposer = sequencer.SequencerAddress
		}
	}
	suite.Require().NotEmpty(currentProposer)
	suite.Require().NotEmpty(candidateProposer)
	suite.Require().Contains([]string{newProposer, fallbackProposer}, currentProposer)
	suite.Require().Contains([]string{newProposer, fallbackProposer}, candidateProposer)

	err = suite.App.SequencerKeeper.SetReplaceProposer(suite.Ctx, &types.MsgRepalceProposer{
		RollappId:   rollappId,
		OldProposer: currentProposer,
		NewProposer: candidateProposer,
		BlockHeight: 20,
	})
	suite.Require().NoError(err)
}

func (suite *SequencerTestSuite) TestReplaceProposerStaleRequestClearedAfterOldProposerRecordMissing() {
	suite.SetupTest()

	rollappId := suite.CreateDefaultRollapp()
	oldProposer := suite.CreateDefaultSequencer(suite.Ctx, rollappId)
	newProposer := suite.CreateDefaultSequencer(suite.Ctx, rollappId)

	// MsgRepalceProposer is the existing generated type name in this module.
	err := suite.App.SequencerKeeper.SetReplaceProposer(suite.Ctx, &types.MsgRepalceProposer{
		RollappId:   rollappId,
		OldProposer: oldProposer,
		NewProposer: newProposer,
		BlockHeight: 10,
	})
	suite.Require().NoError(err)

	store := suite.Ctx.KVStore(suite.App.GetKey(types.StoreKey))
	store.Delete(types.SequencerKey(oldProposer))

	_, found := suite.App.SequencerKeeper.GetSequencer(suite.Ctx, oldProposer)
	suite.Require().False(found)

	err = suite.App.SequencerKeeper.ProcSequencerByPendingStates(
		suite.Ctx,
		rollappId,
		oldProposer,
		&rollapptypes.StateInfo{
			StartHeight: 10,
			NumBlocks:   1,
		},
	)
	suite.Require().NoError(err)
	suite.Require().False(suite.App.SequencerKeeper.IsHasReplaceProposer(suite.Ctx, rollappId))
}
