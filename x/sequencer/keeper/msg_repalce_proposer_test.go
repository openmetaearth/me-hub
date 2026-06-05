package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
)

func (suite *SequencerTestSuite) TestReplaceProposerRollappBinding() {
	suite.SetupTest()
	goCtx := sdk.WrapSDKContext(suite.Ctx)

	rollappId := suite.CreateDefaultRollapp()
	// Create proposer (old proposer)
	oldProposerAddr := suite.CreateDefaultSequencer(suite.Ctx, rollappId)
	_ = oldProposerAddr
	
	// Create another sequencer (new proposer)
	newProposerAddr := suite.CreateDefaultSequencer(suite.Ctx, rollappId)

	// Test case: Check rollapp binding
	// Let's create a sequencer for a different rollapp.
	otherRollappId := suite.CreateDefaultRollapp()
	otherSeqAddr := suite.CreateDefaultSequencer(suite.Ctx, otherRollappId)

	msgWrongOld := &types.MsgReplaceProposerRequest{
		Creator: alice,
		ReplaceProposer: &types.MsgRepalceProposer{
			RollappId:   rollappId,
			OldProposer: otherSeqAddr, // does not belong to rollappId
			NewProposer: newProposerAddr,
			BlockHeight: 100,
		},
	}
	_, err := suite.msgServer.ReplaceProposer(goCtx, msgWrongOld)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "does not belong to rollapp")
}
