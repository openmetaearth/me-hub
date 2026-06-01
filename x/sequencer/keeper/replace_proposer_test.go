package keeper_test

import (
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
)

func (suite *SequencerTestSuite) TestProcSequencerByPendingStatesInvalidNewProposerStatusReturnsSdkError() {
	suite.SetupTest()

	rollappId := suite.CreateDefaultRollapp()
	oldProposer := suite.CreateDefaultSequencer(suite.Ctx, rollappId)
	newProposer := suite.CreateDefaultSequencer(suite.Ctx, rollappId)

	newSequencer, found := suite.App.SequencerKeeper.GetSequencer(suite.Ctx, newProposer)
	suite.Require().True(found)
	newSequencer.Status = types.Unbonding
	suite.App.SequencerKeeper.UpdateSequencer(suite.Ctx, newSequencer, types.Bonded)

	err := suite.App.SequencerKeeper.SetReplaceProposer(suite.Ctx, &types.MsgRepalceProposer{
		RollappId:   rollappId,
		OldProposer: oldProposer,
		NewProposer: newProposer,
		BlockHeight: 10,
	})
	suite.Require().NoError(err)

	err = suite.App.SequencerKeeper.ProcSequencerByPendingStates(
		suite.Ctx,
		rollappId,
		oldProposer,
		&rollapptypes.StateInfo{
			StartHeight: 10,
			NumBlocks:   1,
		},
	)

	suite.Require().ErrorIs(err, types.ErrInvalidSequencerStatus)
	suite.Require().ErrorContains(err, "is not bonded")
}
