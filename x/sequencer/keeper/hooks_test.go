package keeper_test

import (
	"time"

	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
)

func (suite *SequencerTestSuite) TestFraudSubmittedHook() {
	suite.SetupTest()
	suite.Ctx = suite.Ctx.WithBlockHeight(10)
	suite.Ctx = suite.Ctx.WithBlockTime(time.Now())

	keeper := suite.App.SequencerKeeper

	rollappId := suite.CreateDefaultRollapp()

	numOfSequencers := 5

	// create 5 sequencers for rollapp1
	seqAddrs := make([]string, numOfSequencers)
	for i := 0; i < numOfSequencers; i++ {
		seqAddrs[i] = suite.CreateDefaultSequencer(suite.Ctx, rollappId)
	}
	proposer := seqAddrs[0]

	err := keeper.RollappHooks().FraudSubmitted(suite.Ctx, rollappId, 0, proposer)
	suite.Require().NoError(err)

	// check if proposer is slashed
	sequencer, found := keeper.GetSequencer(suite.Ctx, proposer)
	suite.Require().True(found)
	suite.Require().True(sequencer.Jailed)
	suite.Require().Equal(sequencer.Status, types.Unbonded)

	// check if other sequencers are unbonded
	for i := 1; i < numOfSequencers; i++ {
		sequencer, found := keeper.GetSequencer(suite.Ctx, seqAddrs[i])
		suite.Require().True(found)
		suite.Require().False(sequencer.Proposer)
		suite.Require().Equal(sequencer.Status, types.Unbonded)
	}
}

func (suite *SequencerTestSuite) TestReplaceProposerAfterStateFinalizedRollappBinding() {
	suite.SetupTest()
	keeper := suite.App.SequencerKeeper

	rollappId := suite.CreateDefaultRollapp()
	oldProposerAddr := suite.CreateDefaultSequencer(suite.Ctx, rollappId)
	
	// Set the old proposer as unbonding status so forceRemoveUnbondingSequencer can be run (otherwise it returns ErrInvalidSequencerStatus)
	oldSeq, found := keeper.GetSequencer(suite.Ctx, oldProposerAddr)
	suite.Require().True(found)
	oldSeq.Status = types.Unbonding
	keeper.SetSequencer(suite.Ctx, oldSeq)

	// Create sequencer on another rollapp
	otherRollappId := suite.CreateDefaultRollapp()
	otherProposerAddr := suite.CreateDefaultSequencer(suite.Ctx, otherRollappId)
	
	// 1. Valid case matching rollappId
	replaceProposerMsg := types.MsgRepalceProposer{
		RollappId:   rollappId,
		OldProposer: oldProposerAddr,
		NewProposer: otherProposerAddr,
		BlockHeight: 10,
	}
	err := keeper.SetReplaceProposer(suite.Ctx, &replaceProposerMsg)
	suite.Require().NoError(err)

	stateInfo := rollapptypes.StateInfo{
		StateInfoIndex: rollapptypes.StateInfoIndex{RollappId: rollappId, Index: 1},
		StartHeight:    1,
		NumBlocks:      15, // BlockHeight 10 <= StartHeight + NumBlocks - 1 (15)
	}

	err = keeper.RollappHooks().AfterStateFinalized(suite.Ctx, rollappId, &stateInfo)
	suite.Require().NoError(err)

	// 2. Corrupted ReplaceProposer state: oldProposer belongs to wrong rollapp
	// Cleanup and set a new one
	keeper.DeleteReplaceProposer(suite.Ctx, rollappId)
	
	// Create another sequencer on the other rollapp
	otherSeq, found := keeper.GetSequencer(suite.Ctx, otherProposerAddr)
	suite.Require().True(found)
	otherSeq.Status = types.Unbonding
	keeper.SetSequencer(suite.Ctx, otherSeq)

	badReplaceProposerMsg := types.MsgRepalceProposer{
		RollappId:   rollappId,
		OldProposer: otherProposerAddr, // Belongs to otherRollappId instead of rollappId
		NewProposer: oldProposerAddr,
		BlockHeight: 10,
	}
	// Bypass SetReplaceProposer check (which would block it at Tx submission time) to test the hook safety net
	store := suite.Ctx.KVStore(suite.App.GetKey(types.StoreKey))
	bz, err := suite.App.AppCodec().Marshal(&types.MsgStoreReplaceProposer{
		ReplaceProposer: badReplaceProposerMsg,
		HubBlockHeight:  suite.Ctx.BlockHeight(),
	})
	suite.Require().NoError(err)
	store.Set(types.RepalceRollappProposerKey(rollappId), bz)

	// Trigger hook with the wrong binding — should recover (no error) and delete the bad entry
	err = keeper.RollappHooks().AfterStateFinalized(suite.Ctx, rollappId, &stateInfo)
	suite.Require().NoError(err)
	// Verify the corrupted entry was cleaned up
	suite.Require().False(keeper.IsHasReplaceProposer(suite.Ctx, rollappId))

	// 3. Unknown old sequencer — should also recover and delete
	unknownBadMsg := types.MsgRepalceProposer{
		RollappId:   rollappId,
		OldProposer: "cosmos1nonexistent",
		NewProposer: oldProposerAddr,
		BlockHeight: 10,
	}
	bz, err = suite.App.AppCodec().Marshal(&types.MsgStoreReplaceProposer{
		ReplaceProposer: unknownBadMsg,
		HubBlockHeight:  suite.Ctx.BlockHeight(),
	})
	suite.Require().NoError(err)
	store.Set(types.RepalceRollappProposerKey(rollappId), bz)

	err = keeper.RollappHooks().AfterStateFinalized(suite.Ctx, rollappId, &stateInfo)
	suite.Require().NoError(err)
	suite.Require().False(keeper.IsHasReplaceProposer(suite.Ctx, rollappId))
}

