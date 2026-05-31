package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
)

func (suite *SequencerTestSuite) TestReplaceProposerRejectsSelfReplacementBeforeStateLookup() {
	suite.SetupTest()

	proposer := sample.AccAddress()
	msg := &types.MsgReplaceProposerRequest{
		Creator: sample.AccAddress(),
		ReplaceProposer: &types.MsgRepalceProposer{
			RollappId:   "rollapp_1234-1",
			OldProposer: proposer,
			NewProposer: proposer,
			BlockHeight: 10,
		},
	}

	_, err := suite.msgServer.ReplaceProposer(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().ErrorIs(err, sdkerrors.ErrInvalidRequest)
	suite.Require().ErrorContains(err, "old proposer and new proposer cannot be the same address")
}
