package keeper_test

import (
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

func (s *KeeperTestSuite) TestLegacyEquityProposerKeepsCountryElectorate() {
	operator := operatorAccount(s.usaValidator.OperatorAddress)
	equity := sample.AccAddress()
	s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: []string{equity}})
	id := s.submit(s.Dao.GlobalDao, s.withdrawal("usa"))
	// Before equity voting was enabled, an equity proposer also stored a country
	// electorate. Existing proposals must use their stored snapshot, not reroute
	// based on the proposer's current role during voting or tally.
	proposal := s.mustProposal(id)
	proposal.Proposer = equity
	s.App.RegionGovKeeper.SetProposal(s.Ctx, proposal)
	s.Require().Equal(types.ElectorateKindCountryRegionOperators, proposal.Electorate)
	s.Require().Empty(proposal.RequiredOperator)
	s.Require().ErrorIs(s.vote(equity, id, types.VoteOptionYes), types.ErrUnauthorized)
	s.Require().NoError(s.vote(operator, id, types.VoteOptionYes))
	s.Require().NoError(s.vote(s.Dao.GlobalDao, id, types.VoteOptionYes))
	s.endBlockAt(proposal.VotingEndTime)
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(id).Status)
}

func (s *KeeperTestSuite) TestEquityProposalRequiresValidRegionOperator() {
	equity := sample.AccAddress()
	s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: []string{equity}})
	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, "usa")
	s.Require().True(found)
	// Defensive imported-state case; normal region creation binds a validator.
	region.OperatorAddress = ""
	s.App.StakingKeeper.SetRegion(s.Ctx, region)
	seq := s.App.RegionGovKeeper.GetProposalSeq(s.Ctx)
	_, err := s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer: equity, Kind: types.ProposalKindTreasuryWithdraw, Title: "missing required operator", VotingPeriod: 60,
		Treasury: &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{s.withdrawal("usa")}},
	})
	s.Require().ErrorIs(err, types.ErrInvalidProposal)
	s.Require().Equal(seq, s.App.RegionGovKeeper.GetProposalSeq(s.Ctx))
	s.Require().Empty(s.App.RegionGovKeeper.GetAllProposals(s.Ctx))
	s.Require().Empty(s.App.RegionGovKeeper.GetAllElectors(s.Ctx))
}
