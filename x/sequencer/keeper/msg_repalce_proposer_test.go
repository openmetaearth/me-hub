package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
)

func (suite *SequencerTestSuite) TestReplaceProposerSafeBlockHeight() {
	suite.SetupTest()
	goCtx := sdk.WrapSDKContext(suite.Ctx)

	rollappId := suite.CreateDefaultRollapp()
	// Create proposer (old proposer)
	oldProposerAddr := suite.CreateDefaultSequencer(suite.Ctx, rollappId)
	
	// Create another sequencer (new proposer)
	newProposerAddr := suite.CreateDefaultSequencer(suite.Ctx, rollappId)

	// Set state info end height to something high, greater than math.MaxInt64
	stateInfo := rollapptypes.StateInfo{
		StateInfoIndex: rollapptypes.StateInfoIndex{
			RollappId: rollappId,
			Index:     1,
		},
		StartHeight: 9223372036854775800, // close to MaxInt64
		NumBlocks:   100, // StartHeight + NumBlocks = 9223372036854775900 > math.MaxInt64
		Status:      1,
	}
	suite.App.RollappKeeper.SetStateInfo(suite.Ctx, stateInfo)
	suite.App.RollappKeeper.SetLatestStateInfoIndex(suite.Ctx, rollapptypes.StateInfoIndex{
		RollappId: rollappId,
		Index:     1,
	})

	// Test ReplaceProposer with BlockHeight <= stateInfo end height
	// We pass BlockHeight = 10. Since 10 <= end height, it must be rejected!
	msg := &types.MsgReplaceProposerRequest{
		Creator: alice, // rollapp creator
		ReplaceProposer: &types.MsgRepalceProposer{
			RollappId:   rollappId,
			OldProposer: oldProposerAddr,
			NewProposer: newProposerAddr,
			BlockHeight: 10,
		},
	}
	_, err := suite.msgServer.ReplaceProposer(goCtx, msg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "replace proposer block height")
}
