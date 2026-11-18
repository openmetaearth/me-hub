package keeper_test

import (
	"fmt"
	"time"

	"github.com/openmetaearth/me-hub/x/regiongov"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

func (s *KeeperTestSuite) TestProposalsUseCurrentMinimumVotingPeriod() {
	usa := operatorAccount(s.usaValidator.OperatorAddress)
	for _, value := range []uint64{1, 30, 60, 120, types.MaxVotingPeriod} {
		s.Run(fmt.Sprint(value), func() {
			params := s.App.RegionGovKeeper.GetParams(s.Ctx)
			params.MinVotingPeriod = value
			s.Require().NoError(s.App.RegionGovKeeper.SetParams(s.Ctx, params))
			for _, kind := range []types.ProposalKind{types.ProposalKindTreasuryWithdraw, types.ProposalKindReplaceRegionOperator} {
				msg := &types.MsgSubmitProposal{
					Proposer: usa, Kind: kind, Title: "current minimum", VotingPeriod: value,
				}
				if kind == types.ProposalKindTreasuryWithdraw {
					msg.Treasury = &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{s.withdrawal("usa")}}
				} else {
					msg.ReplaceOperator = &types.ReplaceRegionOperatorPayload{RegionId: "usa", NewOperator: s.usaValidator.OperatorAddress}
				}
				if value > 1 {
					msg.VotingPeriod = value - 1
					seq := s.App.RegionGovKeeper.GetProposalSeq(s.Ctx)
					_, err := s.msgServer.SubmitProposal(s.Ctx, msg)
					s.Require().ErrorIs(err, types.ErrInvalidVotingPeriod)
					s.Require().Equal(seq, s.App.RegionGovKeeper.GetProposalSeq(s.Ctx))
				}
				msg.VotingPeriod = value
				response, err := s.msgServer.SubmitProposal(s.Ctx, msg)
				s.Require().NoError(err)
				s.Require().Equal(s.Ctx.BlockTime().Add(time.Duration(value)*time.Second), s.mustProposal(response.Id).VotingEndTime)
			}
			s.Require().Equal(params, s.App.RegionGovKeeper.GetParams(s.Ctx))
		})
	}
}

func (s *KeeperTestSuite) TestMinimumChangeKeepsOpenDeadlinesAndGenesis() {
	usa := operatorAccount(s.usaValidator.OperatorAddress)
	open := s.submit(usa, s.withdrawal("usa"))
	s.Require().NoError(s.vote(usa, open, types.VoteOptionYes))
	s.Require().NoError(s.vote(s.Dao.GlobalDao, open, types.VoteOptionYes))
	openBefore := s.mustProposal(open)
	params := s.App.RegionGovKeeper.GetParams(s.Ctx)
	params.MinVotingPeriod = 120
	s.Require().NoError(s.App.RegionGovKeeper.SetParams(s.Ctx, params))
	s.Require().Equal(openBefore, s.mustProposal(open))
	exported := regiongov.ExportGenesis(s.Ctx, s.App.RegionGovKeeper)
	s.Require().NoError(exported.Validate())
	s.Require().Less(exported.Proposals[0].VotingPeriod, exported.Params.MinVotingPeriod)
	s.endBlockAt(openBefore.VotingEndTime.Add(-time.Nanosecond))
	s.Require().Equal(types.ProposalStatusVoting, s.mustProposal(open).Status)
	s.endBlockAt(openBefore.VotingEndTime)
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(open).Status)
	s.Require().Equal(params, s.App.RegionGovKeeper.GetParams(s.Ctx))

	s.SetupTest()
	regiongov.InitGenesis(s.Ctx, s.App.RegionGovKeeper, *exported)
	s.Require().Equal(exported, regiongov.ExportGenesis(s.Ctx, s.App.RegionGovKeeper))
	s.endBlockAt(openBefore.VotingEndTime)
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(open).Status)
	s.Require().Equal(openBefore.VotingEndTime, s.mustProposal(open).VotingEndTime)
	s.Require().Equal(params, s.App.RegionGovKeeper.GetParams(s.Ctx))
}
