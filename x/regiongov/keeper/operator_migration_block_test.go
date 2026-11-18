package keeper_test

import (
	"time"

	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

// Exercise the full application hook order, not just regiongov.EndBlocker:
// operator migration finishes in one EndBlock and the following block still runs.
func (s *KeeperTestSuite) TestReplaceOperatorAcrossApplicationBlocks() {
	s.SetupTest()
	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, "usa")
	s.Require().True(found)
	newAddr := sdk.ValAddress(s.NewAccounts(1)[0])
	response, err := s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer: s.Dao.GlobalDao, Kind: types.ProposalKindReplaceRegionOperator,
		Title: "replace operator atomically", VotingPeriod: 60,
		ReplaceOperator: &types.ReplaceRegionOperatorPayload{RegionId: region.RegionId, NewOperator: newAddr.String()},
	})
	s.Require().NoError(err)
	s.Require().NoError(s.vote(operatorAccount(region.OperatorAddress), response.Id, types.VoteOptionYes))
	s.Require().NoError(s.vote(s.Dao.GlobalDao, response.Id, types.VoteOptionYes))
	deadline := s.mustProposal(response.Id).VotingEndTime
	consAddr, err := s.usaValidator.GetConsAddr()
	s.Require().NoError(err)
	for i, blockTime := range []time.Time{deadline.Add(-time.Second), deadline, deadline.Add(time.Second)} {
		header := cometbftproto.Header{ChainID: apptesting.TestChainID, Height: int64(i + 1), Time: blockTime, ProposerAddress: consAddr}
		s.Ctx = s.Ctx.WithBlockHeader(header).WithChainID(header.ChainID).WithEventManager(sdk.NewEventManager())
		_, err = s.App.BeginBlocker(s.Ctx)
		s.Require().NoError(err)
		_, err = s.App.EndBlocker(s.Ctx)
		s.Require().NoError(err)
		updated, found := s.App.StakingKeeper.GetRegion(s.Ctx, region.RegionId)
		s.Require().True(found)
		if i == 0 {
			s.Require().Equal(types.ProposalStatusVoting, s.mustProposal(response.Id).Status)
			s.Require().Equal(region.OperatorAddress, updated.OperatorAddress)
		} else {
			s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(response.Id).Status)
			s.Require().Equal(newAddr.String(), updated.OperatorAddress)
			validator, err := s.App.StakingKeeper.GetValidatorByConsAddr(s.Ctx, consAddr)
			s.Require().NoError(err)
			s.Require().Equal(newAddr.String(), validator.OperatorAddress)
			s.Require().Equal(s.usaValidator.ConsensusPubkey, validator.ConsensusPubkey)
		}
	}
}
