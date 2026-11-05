package keeper_test

import (
	"fmt"
	"time"

	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/regiongov"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

func minVotingPeriodProposal(proposer string, value, duration uint64) *types.MsgSubmitProposal {
	return &types.MsgSubmitProposal{
		Proposer: proposer, Kind: types.ProposalKindUpdateMinVotingPeriod,
		Title: "update minimum voting period", VotingPeriod: duration,
		UpdateMinVotingPeriod: &types.UpdateMinVotingPeriodPayload{MinVotingPeriod: value},
	}
}

func (s *KeeperTestSuite) TestUpdateMinVotingPeriodAuthorizationAndElectorate() {
	usa := operatorAccount(s.usaValidator.OperatorAddress)
	china := s.addCountryRegion("chn")
	global := operatorAccount(s.meEarthValidator.OperatorAddress)
	experience := operatorAccount(s.experienceValidator.OperatorAddress)
	equity := sample.AccAddress()
	s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: []string{equity}})
	for _, proposer := range []string{usa, global, experience, s.Dao.GlobalDao} {
		response, err := s.msgServer.SubmitProposal(s.Ctx, minVotingPeriodProposal(proposer, 120, 60))
		s.Require().NoError(err)
		proposal := s.mustProposal(response.Id)
		s.Require().Equal(types.ProposalKindUpdateMinVotingPeriod, proposal.Kind)
		s.Require().Equal(uint64(120), proposal.UpdateMinVotingPeriod.MinVotingPeriod)
		s.Require().Equal(types.ElectorateKindCountryRegionOperators, proposal.Electorate)
		s.Require().ElementsMatch([]string{usa, china}, s.electorAccounts(response.Id))
		s.Require().ErrorIs(s.vote(global, response.Id, types.VoteOptionYes), types.ErrUnauthorized)
		s.Require().ErrorIs(s.vote(experience, response.Id, types.VoteOptionYes), types.ErrUnauthorized)
		s.Require().ErrorIs(s.vote(equity, response.Id, types.VoteOptionYes), types.ErrUnauthorized)
	}
	for _, proposer := range []string{equity, sample.AccAddress()} {
		seq := s.App.RegionGovKeeper.GetProposalSeq(s.Ctx)
		_, err := s.msgServer.SubmitProposal(s.Ctx, minVotingPeriodProposal(proposer, 120, 60))
		s.Require().ErrorIs(err, types.ErrUnauthorized)
		s.Require().Equal(seq, s.App.RegionGovKeeper.GetProposalSeq(s.Ctx))
	}
	s.App.RegionGovKeeper.SetRegionFreeze(s.Ctx, types.RegionFreeze{RegionId: "chn", Frozen: true})
	response, err := s.msgServer.SubmitProposal(s.Ctx, minVotingPeriodProposal(global, 120, 60))
	s.Require().NoError(err)
	s.Require().ElementsMatch([]string{usa}, s.electorAccounts(response.Id))
	s.App.RegionGovKeeper.SetRegionFreeze(s.Ctx, types.RegionFreeze{RegionId: "usa", Frozen: true})
	_, err = s.msgServer.SubmitProposal(s.Ctx, minVotingPeriodProposal(usa, 120, 60))
	s.Require().ErrorIs(err, types.ErrRegionFrozen)
	_, err = s.msgServer.SubmitProposal(s.Ctx, minVotingPeriodProposal(s.Dao.GlobalDao, 120, 60))
	s.Require().ErrorIs(err, types.ErrEmptyElectorate)
	s.Require().Equal(uint64(60), s.App.RegionGovKeeper.GetParams(s.Ctx).MinVotingPeriod)
}

func (s *KeeperTestSuite) TestUpdateMinVotingPeriodTallyAndVeto() {
	for _, test := range []struct {
		name     string
		yesVotes int
		daoYes   bool
		veto     bool
		passed   bool
	}{
		{name: "half does not meet threshold", yesVotes: 1, daoYes: true},
		{name: "no DAO approval", yesVotes: 2},
		{name: "DAO veto overrides approvals", yesVotes: 2, daoYes: true, veto: true},
		{name: "threshold and DAO approval", yesVotes: 2, daoYes: true, passed: true},
	} {
		s.Run(test.name, func() {
			s.SetupTest()
			usa := operatorAccount(s.usaValidator.OperatorAddress)
			china := s.addCountryRegion("chn")
			before := s.App.RegionGovKeeper.GetParams(s.Ctx)
			response, err := s.msgServer.SubmitProposal(s.Ctx, minVotingPeriodProposal(usa, 120, 60))
			s.Require().NoError(err)
			for _, voter := range []string{usa, china}[:test.yesVotes] {
				s.Require().NoError(s.vote(voter, response.Id, types.VoteOptionYes))
			}
			if test.daoYes {
				s.Require().NoError(s.vote(s.Dao.GlobalDao, response.Id, types.VoteOptionYes))
			}
			if test.veto {
				s.Require().NoError(s.vote(s.Dao.GlobalDao, response.Id, types.VoteOptionNo))
				s.Require().Equal(types.ProposalStatusRejected, s.mustProposal(response.Id).Status)
			}
			s.Require().Equal(before, s.App.RegionGovKeeper.GetParams(s.Ctx))
			s.endBlockAt(s.mustProposal(response.Id).VotingEndTime)
			expected := before
			status := types.ProposalStatusRejected
			if test.passed {
				expected.MinVotingPeriod = 120
				status = types.ProposalStatusPassed
			}
			s.Require().Equal(status, s.mustProposal(response.Id).Status)
			s.Require().Equal(expected, s.App.RegionGovKeeper.GetParams(s.Ctx))
		})
	}
}

func (s *KeeperTestSuite) TestUpdateMinVotingPeriodUsesCurrentBoundsAndChangesOnlyMinimum() {
	for _, value := range []uint64{1, 30, 60, 120, types.MaxVotingPeriod} {
		s.Run(fmt.Sprint(value), func() {
			s.SetupTest()
			usa := operatorAccount(s.usaValidator.OperatorAddress)
			// A proposal to lower the minimum may not use its proposed new
			// minimum as its own voting period before that change executes.
			_, err := s.msgServer.SubmitProposal(s.Ctx, minVotingPeriodProposal(usa, value, 59))
			s.Require().ErrorIs(err, types.ErrInvalidVotingPeriod)
			response, err := s.msgServer.SubmitProposal(s.Ctx, minVotingPeriodProposal(usa, value, 60))
			s.Require().NoError(err)
			s.Require().Equal(uint64(60), s.App.RegionGovKeeper.GetParams(s.Ctx).MinVotingPeriod)

			// Change unrelated parameters while voting is open: execution
			// must read current params rather than restore a submission snapshot.
			current := s.App.RegionGovKeeper.GetParams(s.Ctx)
			current.PassThresholdPercent = 75
			current.EquityWithdrawPercent = 7
			s.Require().NoError(s.App.RegionGovKeeper.SetParams(s.Ctx, current))
			s.Require().NoError(s.vote(usa, response.Id, types.VoteOptionYes))
			s.Require().NoError(s.vote(s.Dao.GlobalDao, response.Id, types.VoteOptionYes))
			s.endBlockAt(s.mustProposal(response.Id).VotingEndTime)
			s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(response.Id).Status)
			current.MinVotingPeriod = value
			s.Require().Equal(current, s.App.RegionGovKeeper.GetParams(s.Ctx))

			// Both the new proposal kind and existing treasury submissions
			// must use the updated minimum for subsequent proposals.
			if value > 1 {
				_, err = s.msgServer.SubmitProposal(s.Ctx, minVotingPeriodProposal(usa, 60, value-1))
				s.Require().ErrorIs(err, types.ErrInvalidVotingPeriod)
			}
			next, err := s.msgServer.SubmitProposal(s.Ctx, minVotingPeriodProposal(usa, 60, value))
			s.Require().NoError(err)
			s.Require().Equal(s.Ctx.BlockTime().Add(time.Duration(value)*time.Second), s.mustProposal(next.Id).VotingEndTime)
			treasury := &types.MsgSubmitProposal{
				Proposer: usa, Kind: types.ProposalKindTreasuryWithdraw, Title: "existing kind",
				VotingPeriod: value, Treasury: &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{s.withdrawal("usa")}},
			}
			if value > 1 {
				treasury.VotingPeriod = value - 1
				_, err = s.msgServer.SubmitProposal(s.Ctx, treasury)
				s.Require().ErrorIs(err, types.ErrInvalidVotingPeriod)
			}
			treasury.VotingPeriod = value
			_, err = s.msgServer.SubmitProposal(s.Ctx, treasury)
			s.Require().NoError(err)
		})
	}
}

func (s *KeeperTestSuite) TestUpdateMinVotingPeriodKeepsOpenDeadlinesAndGenesis() {
	usa := operatorAccount(s.usaValidator.OperatorAddress)
	change, err := s.msgServer.SubmitProposal(s.Ctx, minVotingPeriodProposal(usa, 120, 60))
	s.Require().NoError(err)
	// Submit another formerly-valid 60-second proposal ten seconds later so
	// it is still open when the minimum increases to 120 seconds.
	s.Ctx = s.Ctx.WithBlockTime(s.startTime.Add(10 * time.Second))
	open := s.submit(usa, s.withdrawal("usa"))
	s.Require().NoError(s.vote(usa, change.Id, types.VoteOptionYes))
	s.Require().NoError(s.vote(s.Dao.GlobalDao, change.Id, types.VoteOptionYes))
	s.Require().NoError(s.vote(usa, open, types.VoteOptionYes))
	s.Require().NoError(s.vote(s.Dao.GlobalDao, open, types.VoteOptionYes))
	openBefore := s.mustProposal(open)
	s.endBlockAt(s.mustProposal(change.Id).VotingEndTime)
	s.Require().Equal(uint64(120), s.App.RegionGovKeeper.GetParams(s.Ctx).MinVotingPeriod)
	s.Require().Equal(openBefore, s.mustProposal(open))

	exported := regiongov.ExportGenesis(s.Ctx, s.App.RegionGovKeeper)
	s.Require().NoError(exported.Validate())
	s.Require().Less(exported.Proposals[0].VotingPeriod, exported.Params.MinVotingPeriod)

	// The already-open queue entry still matures at its original deadline.
	s.endBlockAt(openBefore.VotingEndTime.Add(-time.Nanosecond))
	s.Require().Equal(types.ProposalStatusVoting, s.mustProposal(open).Status)
	s.endBlockAt(openBefore.VotingEndTime)
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(open).Status)

	// Import the earlier snapshot into a fresh real app and ensure the open
	// proposal/deadline, passed update payload and changed params survive.
	s.SetupTest()
	regiongov.InitGenesis(s.Ctx, s.App.RegionGovKeeper, *exported)
	s.Require().Equal(exported, regiongov.ExportGenesis(s.Ctx, s.App.RegionGovKeeper))
	s.endBlockAt(openBefore.VotingEndTime)
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(open).Status)
	s.Require().Equal(openBefore.VotingEndTime, s.mustProposal(open).VotingEndTime)
	s.Require().Equal(uint64(120), s.App.RegionGovKeeper.GetParams(s.Ctx).MinVotingPeriod)
}
