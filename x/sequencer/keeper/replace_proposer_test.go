package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
)

func (suite *SequencerTestSuite) TestReplaceProposerKeepsOldProposerBondUntilUnbondingTime() {
	suite.SetupTest()

	startTime := time.Unix(1_700_000_000, 0).UTC()
	suite.Ctx = suite.Ctx.WithBlockHeight(10).WithBlockTime(startTime)

	keeper := suite.App.SequencerKeeper
	rollappId := suite.CreateDefaultRollapp()
	oldProposer := suite.CreateDefaultSequencer(suite.Ctx, rollappId)
	newProposer := suite.CreateDefaultSequencer(suite.Ctx, rollappId)
	oldAddr := sdk.MustAccAddressFromBech32(oldProposer)
	balanceBeforeFinalize := suite.App.BankKeeper.GetBalance(suite.Ctx, oldAddr, bond.Denom)

	err := keeper.SetReplaceProposer(suite.Ctx, &types.MsgRepalceProposer{
		RollappId:   rollappId,
		OldProposer: oldProposer,
		NewProposer: newProposer,
		BlockHeight: 12,
	})
	suite.Require().NoError(err)

	stateInfo := &rollapptypes.StateInfo{
		StartHeight: 11,
		NumBlocks:   2,
	}

	err = keeper.ProcSequencerByPendingStates(suite.Ctx, rollappId, oldProposer, stateInfo)
	suite.Require().NoError(err)

	oldSeq, found := keeper.GetSequencer(suite.Ctx, oldProposer)
	suite.Require().True(found)
	suite.Equal(types.Unbonding, oldSeq.Status)
	suite.False(oldSeq.Proposer)
	suite.Equal(startTime.Add(keeper.UnbondingTime(suite.Ctx)), oldSeq.UnbondTime)

	newSeq, found := keeper.GetSequencer(suite.Ctx, newProposer)
	suite.Require().True(found)
	suite.Equal(types.Bonded, newSeq.Status)
	suite.True(newSeq.Proposer)

	err = keeper.RollappHooks().AfterStateFinalized(suite.Ctx, rollappId, stateInfo)
	suite.Require().NoError(err)

	balanceAfterFinalize := suite.App.BankKeeper.GetBalance(suite.Ctx, oldAddr, bond.Denom)
	suite.True(balanceBeforeFinalize.IsEqual(balanceAfterFinalize), "old proposer bond was refunded before unbonding time")

	oldSeq, found = keeper.GetSequencer(suite.Ctx, oldProposer)
	suite.Require().True(found)
	suite.Equal(types.Unbonding, oldSeq.Status)
	suite.False(oldSeq.Tokens.IsZero())

	replaceProposer, err := keeper.GetReplaceProposer(suite.Ctx, rollappId)
	suite.Require().NoError(err)
	suite.Nil(replaceProposer)

	matureBeforeUnbondTime := keeper.GetMatureUnbondingSequencers(suite.Ctx, oldSeq.UnbondTime.Add(-time.Second))
	suite.Len(matureBeforeUnbondTime, 0)

	keeper.UnbondAllMatureSequencers(suite.Ctx, oldSeq.UnbondTime.Add(time.Second))

	oldSeq, found = keeper.GetSequencer(suite.Ctx, oldProposer)
	suite.Require().True(found)
	suite.Equal(types.Unbonded, oldSeq.Status)
	suite.True(oldSeq.Tokens.IsZero())

	balanceAfterUnbondTime := suite.App.BankKeeper.GetBalance(suite.Ctx, oldAddr, bond.Denom)
	suite.True(balanceAfterFinalize.Add(bond).IsEqual(balanceAfterUnbondTime), "old proposer bond should refund after unbonding time")
}
