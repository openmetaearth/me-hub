package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

func (s *KeeperTestSuite) TestEquityCrossRegionSubmissionRejected() {
	chnOperator := s.addCountryRegion("chn")
	equity := sample.AccAddress()
	s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "chn", Addresses: []string{equity}})
	withdrawal, assertNoTransfer := s.crossRegionWithdrawal("usa")

	seq := s.App.RegionGovKeeper.GetProposalSeq(s.Ctx)
	proposals := s.App.RegionGovKeeper.GetAllProposals(s.Ctx)
	electors := s.App.RegionGovKeeper.GetAllElectors(s.Ctx)
	events := s.Ctx.EventManager().Events()
	queued := func() []uint64 {
		// Inspect the queue without consuming entries in the real context.
		ctx, _ := s.Ctx.CacheContext()
		return s.App.RegionGovKeeper.DequeueMatureProposals(ctx, s.startTime.Add(time.Hour))
	}
	queue := queued()
	_, err := s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer: equity, Kind: types.ProposalKindTreasuryWithdraw, Title: "chn equity cannot withdraw usa treasury", VotingPeriod: 60,
		Treasury: &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{withdrawal}},
	})
	s.Require().ErrorIs(err, types.ErrUnauthorized)
	s.Require().ErrorContains(err, "an equity address can withdraw only from its own region, got usa")
	s.Require().Equal(seq, s.App.RegionGovKeeper.GetProposalSeq(s.Ctx))
	s.Require().Equal(proposals, s.App.RegionGovKeeper.GetAllProposals(s.Ctx))
	s.Require().Equal(electors, s.App.RegionGovKeeper.GetAllElectors(s.Ctx))
	s.Require().Equal(queue, queued())
	s.Require().Equal(events, s.Ctx.EventManager().Events())
	assertNoTransfer()

	// The same equity can still submit and execute a withdrawal from chn.
	ownWithdrawal := s.withdrawal("chn")
	ownWithdrawal.Receiver = sample.AccAddress()
	id := s.submit(equity, ownWithdrawal)
	s.Require().Equal(seq+1, id)
	s.Require().Equal(types.ElectorateKindRegionEquity, s.mustProposal(id).Electorate)
	s.Require().Equal([]string{equity}, s.electorAccounts(id))
	for _, voter := range []string{equity, chnOperator, s.Dao.GlobalDao} {
		s.Require().NoError(s.vote(voter, id, types.VoteOptionYes))
	}
	s.endBlockAt(s.mustProposal(id).VotingEndTime)
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(id).Status)
	s.Require().Equal(ownWithdrawal.Amount, s.App.BankKeeper.GetAllBalances(s.Ctx, sdk.MustAccAddressFromBech32(ownWithdrawal.Receiver)))
	assertNoTransfer()
}

func (s *KeeperTestSuite) TestEquityCrossRegionVotesRejectedOnEquityProposal() {
	chnOperator := s.addCountryRegion("chn")
	usaEquities := []string{sample.AccAddress(), sample.AccAddress(), sample.AccAddress()}
	chnEquities := []string{sample.AccAddress(), sample.AccAddress(), sample.AccAddress()}
	s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: usaEquities})
	s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "chn", Addresses: chnEquities})
	withdrawal, assertNoTransfer := s.crossRegionWithdrawal("usa")
	id := s.submit(usaEquities[0], withdrawal)
	s.Require().ElementsMatch(usaEquities, s.equityProposalElectorAccounts(id))
	s.Require().Equal(s.usaValidator.OperatorAddress, s.mustProposal(id).RequiredOperator)

	// Both separate approvals are present, but only one of three local equities
	// votes yes. Foreign equities and even the foreign operator cannot add votes.
	for _, voter := range []string{usaEquities[0], operatorAccount(s.usaValidator.OperatorAddress), s.Dao.GlobalDao} {
		s.Require().NoError(s.vote(voter, id, types.VoteOptionYes))
	}
	s.assertCrossRegionVotesUnauthorized(id, append(chnEquities, chnOperator)...)
	s.Require().Len(s.App.RegionGovKeeper.GetVotes(s.Ctx, id), 2)
	s.Require().Equal(types.DaoVoteYes, s.mustProposal(id).DaoVote)
	assertNoTransfer()
	s.endBlockAt(s.mustProposal(id).VotingEndTime)
	s.Require().Equal(types.ProposalStatusRejected, s.mustProposal(id).Status)
	assertNoTransfer()
}

func (s *KeeperTestSuite) TestEquityCrossRegionVotesRejectedOnDAOProposal() {
	chnOperator := s.addCountryRegion("chn")
	chnEquities := []string{sample.AccAddress(), sample.AccAddress(), sample.AccAddress()}
	s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "chn", Addresses: chnEquities})
	withdrawal, assertNoTransfer := s.crossRegionWithdrawal("usa")
	id := s.submit(s.Dao.GlobalDao, withdrawal)
	s.Require().Equal(types.ElectorateKindCountryRegionOperators, s.mustProposal(id).Electorate)
	s.Require().Empty(s.mustProposal(id).RequiredOperator)
	s.Require().ElementsMatch([]string{operatorAccount(s.usaValidator.OperatorAddress), chnOperator}, s.electorAccounts(id))
	s.Require().NoError(s.vote(s.Dao.GlobalDao, id, types.VoteOptionYes))

	s.assertCrossRegionVotesUnauthorized(id, chnEquities...)
	s.Require().Empty(s.App.RegionGovKeeper.GetVotes(s.Ctx, id))
	s.Require().Equal(types.DaoVoteYes, s.mustProposal(id).DaoVote)
	assertNoTransfer()
	s.endBlockAt(s.mustProposal(id).VotingEndTime)
	s.Require().Equal(types.ProposalStatusRejected, s.mustProposal(id).Status)
	assertNoTransfer()
}

func (s *KeeperTestSuite) assertCrossRegionVotesUnauthorized(id uint64, voters ...string) {
	s.T().Helper()
	proposal := s.mustProposal(id)
	electors := s.App.RegionGovKeeper.GetElectors(s.Ctx, id)
	votes := s.App.RegionGovKeeper.GetVotes(s.Ctx, id)
	events := s.Ctx.EventManager().Events()
	for _, voter := range voters {
		s.Require().ErrorIs(s.vote(voter, id, types.VoteOptionYes), types.ErrUnauthorized)
		_, found := s.App.RegionGovKeeper.GetVote(s.Ctx, id, voter)
		s.Require().False(found, "foreign voter %s must not have a stored ballot", voter)
	}
	s.Require().Equal(proposal, s.mustProposal(id))
	s.Require().Equal(electors, s.App.RegionGovKeeper.GetElectors(s.Ctx, id))
	s.Require().Equal(votes, s.App.RegionGovKeeper.GetVotes(s.Ctx, id))
	s.Require().Equal(events, s.Ctx.EventManager().Events())
}

// Use a funded, existing treasury and a fresh receiver so authorization failures
// cannot be mistaken for an unknown region, insufficient funds, or a self-send.
func (s *KeeperTestSuite) crossRegionWithdrawal(regionID string) (types.Withdrawal, func()) {
	s.T().Helper()
	withdrawal := s.withdrawal(regionID)
	withdrawal.Receiver = sample.AccAddress()
	withdrawal.Amount = sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 100))
	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, regionID)
	s.Require().True(found)
	treasury := sdk.MustAccAddressFromBech32(region.RegionTreasureAddr)
	receiver := sdk.MustAccAddressFromBech32(withdrawal.Receiver)
	treasuryBefore := s.App.BankKeeper.GetAllBalances(s.Ctx, treasury)
	receiverBefore := s.App.BankKeeper.GetAllBalances(s.Ctx, receiver)
	s.Require().True(treasuryBefore.IsAllGTE(withdrawal.Amount))
	return withdrawal, func() {
		s.T().Helper()
		s.Require().Equal(treasuryBefore, s.App.BankKeeper.GetAllBalances(s.Ctx, treasury))
		s.Require().Equal(receiverBefore, s.App.BankKeeper.GetAllBalances(s.Ctx, receiver))
	}
}
