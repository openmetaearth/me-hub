package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
)

func (suite *SequencerTestSuite) TestReplaceProposerUnbondingFlow() {
	suite.SetupTest()
	rollappId := suite.CreateDefaultRollapp()

	// 1. Create old proposer and new proposer sequencers
	oldProposerAddr := suite.CreateDefaultSequencer(suite.Ctx, rollappId)
	newProposerAddr := suite.CreateDefaultSequencer(suite.Ctx, rollappId)

	// Verify old proposer is indeed the proposer and status is Bonded
	oldProposer, found := suite.App.SequencerKeeper.GetSequencer(suite.Ctx, oldProposerAddr)
	suite.Require().True(found)
	suite.Require().Equal(types.Bonded, oldProposer.Status)
	suite.Require().True(oldProposer.Proposer)

	newProposer, found := suite.App.SequencerKeeper.GetSequencer(suite.Ctx, newProposerAddr)
	suite.Require().True(found)
	suite.Require().Equal(types.Bonded, newProposer.Status)
	suite.Require().False(newProposer.Proposer)

	// 2. Set replace proposer request
	replaceMsg := &types.MsgRepalceProposer{
		RollappId:   rollappId,
		OldProposer: oldProposerAddr,
		NewProposer: newProposerAddr,
		BlockHeight: 100,
	}
	err := suite.App.SequencerKeeper.SetReplaceProposer(suite.Ctx, replaceMsg)
	suite.Require().NoError(err)

	// 3. ProcSequencerByPendingStates at height lower than replacement height
	stateInfoLow := &rollapptypes.StateInfo{
		StartHeight: 1,
		NumBlocks:   50,
	}
	err = suite.App.SequencerKeeper.ProcSequencerByPendingStates(suite.Ctx, rollappId, oldProposerAddr, stateInfoLow)
	suite.Require().NoError(err)

	// Old proposer should still be bonded and proposer
	oldProposer, _ = suite.App.SequencerKeeper.GetSequencer(suite.Ctx, oldProposerAddr)
	suite.Require().Equal(types.Bonded, oldProposer.Status)
	suite.Require().True(oldProposer.Proposer)

	// 4. ProcSequencerByPendingStates at/above replacement height
	stateInfoHigh := &rollapptypes.StateInfo{
		StartHeight: 51,
		NumBlocks:   50, // 51 + 50 - 1 = 100
	}
	suite.Ctx = suite.Ctx.WithBlockTime(time.Now())
	err = suite.App.SequencerKeeper.ProcSequencerByPendingStates(suite.Ctx, rollappId, oldProposerAddr, stateInfoHigh)
	suite.Require().NoError(err)

	// 5. Verify status transitions
	oldProposer, _ = suite.App.SequencerKeeper.GetSequencer(suite.Ctx, oldProposerAddr)
	suite.Require().Equal(types.Unbonding, oldProposer.Status)
	suite.Require().False(oldProposer.Proposer)
	suite.Require().False(oldProposer.UnbondTime.IsZero())

	newProposer, _ = suite.App.SequencerKeeper.GetSequencer(suite.Ctx, newProposerAddr)
	suite.Require().Equal(types.Bonded, newProposer.Status)
	suite.Require().True(newProposer.Proposer)

	// Check if old sequencer is in the unbonding queue
	unbondingSeqs := suite.App.SequencerKeeper.GetMatureUnbondingSequencers(suite.Ctx, oldProposer.UnbondTime.Add(1*time.Second))
	suite.Require().Len(unbondingSeqs, 1)
	suite.Require().Equal(oldProposerAddr, unbondingSeqs[0].SequencerAddress)

	// 6. Test maturation path (UnbondAllMatureSequencers)
	suite.App.SequencerKeeper.UnbondAllMatureSequencers(suite.Ctx, oldProposer.UnbondTime.Add(1*time.Second))
	oldProposer, _ = suite.App.SequencerKeeper.GetSequencer(suite.Ctx, oldProposerAddr)
	suite.Require().Equal(types.Unbonded, oldProposer.Status)

	// Verify it was removed from the queue
	unbondingSeqsAfter := suite.App.SequencerKeeper.GetMatureUnbondingSequencers(suite.Ctx, oldProposer.UnbondTime.Add(1*time.Second))
	suite.Require().Len(unbondingSeqsAfter, 0)
}

func (suite *SequencerTestSuite) TestReplaceProposerForceRemoveQueueCleanup() {
	suite.SetupTest()
	rollappId := suite.CreateDefaultRollapp()

	// 1. Create old proposer and new proposer sequencers
	oldProposerAddr := suite.CreateDefaultSequencer(suite.Ctx, rollappId)
	newProposerAddr := suite.CreateDefaultSequencer(suite.Ctx, rollappId)

	// 2. Set replace proposer request
	replaceMsg := &types.MsgRepalceProposer{
		RollappId:   rollappId,
		OldProposer: oldProposerAddr,
		NewProposer: newProposerAddr,
		BlockHeight: 100,
	}
	err := suite.App.SequencerKeeper.SetReplaceProposer(suite.Ctx, replaceMsg)
	suite.Require().NoError(err)

	// 3. ProcSequencerByPendingStates to transition old proposer to Unbonding and enqueue it
	stateInfo := &rollapptypes.StateInfo{
		StartHeight: 1,
		NumBlocks:   100,
	}
	suite.Ctx = suite.Ctx.WithBlockTime(time.Now())
	err = suite.App.SequencerKeeper.ProcSequencerByPendingStates(suite.Ctx, rollappId, oldProposerAddr, stateInfo)
	suite.Require().NoError(err)

	oldProposer, _ := suite.App.SequencerKeeper.GetSequencer(suite.Ctx, oldProposerAddr)
	suite.Require().Equal(types.Unbonding, oldProposer.Status)

	// Verify old proposer is in the unbonding queue
	unbondingSeqs := suite.App.SequencerKeeper.GetMatureUnbondingSequencers(suite.Ctx, oldProposer.UnbondTime.Add(1*time.Second))
	suite.Require().Len(unbondingSeqs, 1)

	// 4. Trigger force remove unbonding sequencer (dispute period ended / state finalized)
	err = suite.App.SequencerKeeper.RollappHooks().AfterStateFinalized(suite.Ctx, rollappId, stateInfo)
	suite.Require().NoError(err)

	// 5. Verify sequencer is deleted from store
	_, found := suite.App.SequencerKeeper.GetSequencer(suite.Ctx, oldProposerAddr)
	suite.Require().False(found)

	// 6. Verify sequencer is removed from the unbonding queue as well!
	unbondingSeqsAfter := suite.App.SequencerKeeper.GetMatureUnbondingSequencers(suite.Ctx, oldProposer.UnbondTime.Add(1*time.Second))
	suite.Require().Len(unbondingSeqsAfter, 0)
}
