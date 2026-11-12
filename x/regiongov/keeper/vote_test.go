package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestTreasuryElectorates() {
	usa := operatorAccount(s.usaValidator.OperatorAddress)
	chn := s.addCountryRegion("chn")
	dao := s.Dao.GlobalDao

	for _, tc := range []struct {
		proposer string
		regions  []string
	}{
		{usa, []string{wstakingtypes.MeEarthRegionId}},
		{usa, []string{"usa"}},
		{dao, []string{"usa", "chn"}},
		{dao, []string{wstakingtypes.MeEarthRegionId, "usa"}},
	} {
		var withdrawals []types.Withdrawal
		for _, id := range tc.regions {
			withdrawals = append(withdrawals, s.withdrawal(id))
		}
		id := s.submit(tc.proposer, withdrawals...)
		s.Require().ElementsMatch([]string{usa, chn}, s.electorAccounts(id))
		s.Require().Equal(types.ElectorateKindCountryRegionOperators, s.mustProposal(id).Electorate)
	}

	for _, ids := range [][]string{{"chn"}, {"usa", "chn"}, {wstakingtypes.MeEarthRegionId, "chn"}} {
		var withdrawals []types.Withdrawal
		for _, id := range ids {
			withdrawals = append(withdrawals, s.withdrawal(id))
		}
		_, err := s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
			Proposer: usa, Kind: types.ProposalKindTreasuryWithdraw, Title: "another country", VotingPeriod: 60,
			Treasury: &types.TreasuryWithdrawPayload{Withdrawals: withdrawals},
		})
		s.Require().ErrorIs(err, types.ErrUnauthorized)
	}

	_, err := s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer: dao, Kind: types.ProposalKindTreasuryWithdraw, Title: "unknown region", VotingPeriod: 60,
		Treasury: &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{s.withdrawal("missing")}},
	})
	s.Require().ErrorIs(err, types.ErrUnknownRegion)
}

func (s *KeeperTestSuite) TestEquityCanWithdrawOnlyItsOwnRegion() {
	s.addCountryRegion("chn")
	usaEquity, chnEquity := sample.AccAddress(), sample.AccAddress()
	s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: []string{usaEquity}})
	s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "chn", Addresses: []string{chnEquity}})

	own := s.submit(usaEquity, s.withdrawal("usa"))
	s.Require().ElementsMatch([]string{usaEquity}, s.electorAccounts(own))
	s.Require().Equal(types.ElectorateKindRegionEquity, s.mustProposal(own).Electorate)
	s.Require().Equal(operatorAccount(s.usaValidator.OperatorAddress), s.mustProposal(own).RequiredOperator)

	for _, ids := range [][]string{{"chn"}, {"usa", "chn"}, {wstakingtypes.MeEarthRegionId}} {
		var withdrawals []types.Withdrawal
		for _, id := range ids {
			withdrawals = append(withdrawals, s.withdrawal(id))
		}
		_, err := s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
			Proposer: usaEquity, Kind: types.ProposalKindTreasuryWithdraw, Title: "another region", VotingPeriod: 60,
			Treasury: &types.TreasuryWithdrawPayload{Withdrawals: withdrawals},
		})
		s.Require().ErrorIs(err, types.ErrUnauthorized)
	}

	_, err := s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer: usaEquity, Kind: types.ProposalKindReplaceRegionOperator, Title: "replace", VotingPeriod: 60,
		ReplaceOperator: &types.ReplaceRegionOperatorPayload{
			RegionId: "usa", NewOperator: sdk.ValAddress(sdk.MustAccAddressFromBech32(sample.AccAddress())).String(),
		},
	})
	s.Require().ErrorIs(err, types.ErrUnauthorized)

	_, err = s.msgServer.FreezeEquity(s.Ctx, &types.MsgFreezeEquity{
		Authority: s.Dao.GlobalDao, RegionId: "chn", Address: chnEquity,
	})
	s.Require().NoError(err)
	_, err = s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer: chnEquity, Kind: types.ProposalKindTreasuryWithdraw, Title: "frozen equity", VotingPeriod: 60,
		Treasury: &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{s.withdrawal("chn")}},
	})
	s.Require().ErrorIs(err, types.ErrEquityFrozen)

	_, err = s.msgServer.FreezeRegion(s.Ctx, &types.MsgFreezeRegion{
		Authority: s.Dao.GlobalDao, RegionIds: []string{"usa"},
	})
	s.Require().NoError(err)
	_, err = s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer: usaEquity, Kind: types.ProposalKindTreasuryWithdraw, Title: "frozen region", VotingPeriod: 60,
		Treasury: &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{s.withdrawal("usa")}},
	})
	s.Require().ErrorIs(err, types.ErrRegionFrozen)
}

func (s *KeeperTestSuite) TestFrozenRegionIsOutOfTheSnapshot() {
	chn := s.addCountryRegion("chn")
	_, err := s.msgServer.FreezeRegion(s.Ctx, &types.MsgFreezeRegion{
		Authority: s.Dao.GlobalDao, RegionIds: []string{"usa"},
	})
	s.Require().NoError(err)

	id := s.submit(chn, s.withdrawal(wstakingtypes.MeEarthRegionId))
	s.Require().Equal([]string{chn}, s.electorAccounts(id))

	newOperator := sdk.ValAddress(sdk.MustAccAddressFromBech32(sample.AccAddress()))
	replaced, err := s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer: s.Dao.GlobalDao, Kind: types.ProposalKindReplaceRegionOperator, Title: "replace frozen operator", VotingPeriod: 60,
		ReplaceOperator: &types.ReplaceRegionOperatorPayload{RegionId: "usa", NewOperator: newOperator.String()},
	})
	s.Require().NoError(err)
	s.Require().Equal([]string{chn}, s.electorAccounts(replaced.Id))
	s.Require().NoError(s.vote(chn, replaced.Id, types.VoteOptionYes))
	s.Require().NoError(s.vote(s.Dao.GlobalDao, replaced.Id, types.VoteOptionYes))
	s.endBlockAt(s.mustProposal(replaced.Id).VotingEndTime)
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(replaced.Id).Status)
	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, "usa")
	s.Require().True(found)
	s.Require().Equal(newOperator.String(), region.OperatorAddress)
	_, found = s.App.StakingKeeper.GetValidator(s.Ctx, s.usaValidator.GetOperator())
	s.Require().False(found)
	validator, found := s.App.StakingKeeper.GetValidator(s.Ctx, newOperator)
	s.Require().True(found)
	s.Require().Equal(s.usaValidator.ConsensusPubkey, validator.ConsensusPubkey)
	s.Require().Equal("usa", validator.Description.RegionID)
}

func (s *KeeperTestSuite) TestEmptyElectorateDoesNotOpen() {
	for _, name := range []string{"all country regions frozen", "country region without operator"} {
		s.Run(name, func() {
			s.SetupTest()
			if name == "all country regions frozen" {
				_, err := s.msgServer.FreezeRegion(s.Ctx, &types.MsgFreezeRegion{
					Authority: s.Dao.GlobalDao, RegionIds: []string{"usa"},
				})
				s.Require().NoError(err)
			} else {
				// Defensive legacy-state fixture: normal NewRegion requires a validator.
				region, found := s.App.StakingKeeper.GetRegion(s.Ctx, "usa")
				s.Require().True(found)
				region.OperatorAddress = ""
				s.App.StakingKeeper.SetRegion(s.Ctx, region)
			}
			before := s.App.ExportState(s.Ctx)
			events := len(s.Ctx.EventManager().Events())
			_, err := s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
				Proposer: s.Dao.GlobalDao, Kind: types.ProposalKindTreasuryWithdraw, Title: "empty electorate", VotingPeriod: 60,
				Treasury: &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{s.withdrawal("usa")}},
			})
			s.Require().ErrorIs(err, types.ErrEmptyElectorate)
			s.Require().Empty(s.App.RegionGovKeeper.GetAllProposals(s.Ctx))
			s.Require().Equal(before, s.App.ExportState(s.Ctx))
			s.Require().Len(s.Ctx.EventManager().Events(), events)
		})
	}
}

func (s *KeeperTestSuite) TestDuplicateOperatorHasOneVote() {
	s.addCountryRegion("chn")
	// NewRegion forbids duplicate validator bindings. Retain coverage for imported
	// legacy/malformed state by changing only this fixture through the real keeper.
	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, "chn")
	s.Require().True(found)
	region.OperatorAddress = s.usaValidator.OperatorAddress
	s.App.StakingKeeper.SetRegion(s.Ctx, region)

	id := s.submit(s.Dao.GlobalDao, s.withdrawal("usa"), s.withdrawal("chn"))
	usa := operatorAccount(s.usaValidator.OperatorAddress)
	s.Require().Equal([]string{usa}, s.electorAccounts(id))
	s.Require().NoError(s.vote(usa, id, types.VoteOptionYes))
	s.Require().ErrorIs(s.vote(usa, id, types.VoteOptionNo), types.ErrAlreadyVoted)
}

func (s *KeeperTestSuite) TestBallotRulesAndTally() {
	usa := operatorAccount(s.usaValidator.OperatorAddress)
	chn := s.addCountryRegion("chn")
	jpn := s.addCountryRegion("jpn")
	dao := s.Dao.GlobalDao
	withdrawal := s.withdrawal("usa")
	withdrawal.Receiver = sample.AccAddress()
	receiver := sdk.MustAccAddressFromBech32(withdrawal.Receiver)
	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, "usa")
	s.Require().True(found)
	treasury := sdk.MustAccAddressFromBech32(region.RegionTreasureAddr)
	before := s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom)

	oneThird := s.submit(usa, withdrawal)
	s.Require().ElementsMatch([]string{usa, chn, jpn}, s.electorAccounts(oneThird))
	s.Require().ErrorIs(s.vote(sample.AccAddress(), oneThird, types.VoteOptionYes), types.ErrUnauthorized)
	s.Require().NoError(s.vote(chn, oneThird, types.VoteOptionYes))
	s.Require().ErrorIs(s.vote(chn, oneThird, types.VoteOptionNo), types.ErrAlreadyVoted)
	s.Require().NoError(s.vote(usa, oneThird, types.VoteOptionNo))
	s.Require().NoError(s.vote(dao, oneThird, types.VoteOptionYes))

	twoThirds := s.submit(usa, withdrawal)
	s.Require().NoError(s.vote(chn, twoThirds, types.VoteOptionYes))
	s.Require().NoError(s.vote(jpn, twoThirds, types.VoteOptionYes))
	s.Require().NoError(s.vote(dao, twoThirds, types.VoteOptionYes))

	noDao := s.submit(usa, withdrawal)
	for _, voter := range []string{usa, chn, jpn} {
		s.Require().NoError(s.vote(voter, noDao, types.VoteOptionYes))
	}

	s.endBlockAt(s.startTime)
	for _, id := range []uint64{oneThird, twoThirds, noDao} {
		s.Require().Equal(types.ProposalStatusVoting, s.mustProposal(id).Status)
	}
	end := s.mustProposal(oneThird).VotingEndTime
	s.Ctx = s.Ctx.WithBlockTime(end.Add(-time.Second))
	s.Require().NoError(s.vote(jpn, oneThird, types.VoteOptionNo))
	s.Ctx = s.Ctx.WithBlockTime(end)
	// This operator has not voted on twoThirds: the deadline, not a duplicate
	// ballot or lack of authorization, must prevent the vote.
	s.Require().ErrorIs(s.vote(usa, twoThirds, types.VoteOptionYes), types.ErrVotingClosed)
	s.endBlockAt(end)
	s.Require().Equal(types.ProposalStatusRejected, s.mustProposal(oneThird).Status)
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(twoThirds).Status)
	s.Require().Equal(types.ProposalStatusRejected, s.mustProposal(noDao).Status)
	// Only the passing proposal pays, through wstaking and the actual bank keeper.
	s.Require().Equal(before.Amount.SubRaw(1), s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom).Amount)
	s.Require().Equal(sdk.OneInt(), s.App.BankKeeper.GetBalance(s.Ctx, receiver, params.BaseDenom).Amount)

	s.Ctx = s.Ctx.WithEventManager(sdk.NewEventManager())
	s.endBlockAt(end.Add(time.Second))
	s.Require().Empty(s.Ctx.EventManager().Events())
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(twoThirds).Status)
	s.Require().Equal(sdk.OneInt(), s.App.BankKeeper.GetBalance(s.Ctx, receiver, params.BaseDenom).Amount)
	s.Require().Equal(before.Amount.SubRaw(1), s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom).Amount)
}

func (s *KeeperTestSuite) TestFreezeAfterSnapshotBlocksThatVote() {
	usa := operatorAccount(s.usaValidator.OperatorAddress)
	chn := s.addCountryRegion("chn")
	id := s.submit(s.Dao.GlobalDao, s.withdrawal("usa"), s.withdrawal("chn"))
	s.Require().ElementsMatch([]string{usa, chn}, s.electorAccounts(id))
	_, err := s.msgServer.FreezeRegion(s.Ctx, &types.MsgFreezeRegion{
		Authority: s.Dao.GlobalDao, RegionIds: []string{"usa"},
	})
	s.Require().NoError(err)
	s.Require().ErrorIs(s.vote(usa, id, types.VoteOptionYes), types.ErrRegionFrozen)
	s.Require().NoError(s.vote(chn, id, types.VoteOptionYes))
	s.Require().NoError(s.vote(s.Dao.GlobalDao, id, types.VoteOptionYes))
	// Freezing after submission does not reduce the denominator to one.
	s.Require().ElementsMatch([]string{usa, chn}, s.electorAccounts(id))
	s.endBlockAt(s.mustProposal(id).VotingEndTime)
	s.Require().Equal(types.ProposalStatusRejected, s.mustProposal(id).Status)
}

func (s *KeeperTestSuite) TestOneOfTwoFailsAndTwoOfTwoPasses() {
	usa := operatorAccount(s.usaValidator.OperatorAddress)
	chn := s.addCountryRegion("chn")
	withdrawal := s.withdrawal("usa")
	withdrawal.Receiver = sample.AccAddress()

	failing := s.submit(usa, withdrawal)
	s.Require().NoError(s.vote(chn, failing, types.VoteOptionYes))
	s.Require().NoError(s.vote(s.Dao.GlobalDao, failing, types.VoteOptionYes))

	passing := s.submit(usa, withdrawal)
	s.Require().NoError(s.vote(usa, passing, types.VoteOptionYes))
	s.Require().NoError(s.vote(chn, passing, types.VoteOptionYes))
	s.Require().NoError(s.vote(s.Dao.GlobalDao, passing, types.VoteOptionYes))

	s.endBlockAt(s.mustProposal(failing).VotingEndTime)
	s.Require().Equal(types.ProposalStatusRejected, s.mustProposal(failing).Status)
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(passing).Status)
	receiver := sdk.MustAccAddressFromBech32(withdrawal.Receiver)
	s.Require().Equal(sdk.OneInt(), s.App.BankKeeper.GetBalance(s.Ctx, receiver, params.BaseDenom).Amount)
}
