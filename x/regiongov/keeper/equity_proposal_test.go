package keeper_test

import (
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/regiongov"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

// Equity snapshots, like new country snapshots, store canonical ballot accounts.
func (s *KeeperTestSuite) equityProposalElectorAccounts(id uint64) []string {
	s.T().Helper()
	s.Require().Equal(types.ElectorateKindRegionEquity, s.mustProposal(id).Electorate)
	electors := s.App.RegionGovKeeper.GetElectors(s.Ctx, id)
	accounts := make([]string, 0, len(electors))
	for _, elector := range electors {
		account, err := sdk.AccAddressFromBech32(elector.Operator)
		s.Require().NoError(err)
		s.Require().Equal(account.String(), elector.Operator)
		s.Require().Equal("usa", elector.RegionId)
		accounts = append(accounts, account.String())
	}
	return accounts
}

func (s *KeeperTestSuite) TestEquityProposalApprovalAndTransfer() {
	operator := operatorAccount(s.usaValidator.OperatorAddress)
	otherOperator := s.addCountryRegion("chn")
	equities := []string{sample.AccAddress(), sample.AccAddress(), sample.AccAddress()}
	otherEquity := sample.AccAddress()
	// Imported/legacy lists can contain repeated or noncanonical addresses and
	// the separate approvers. Only distinct, counted equity accounts survive.
	s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{
		RegionId: "usa",
		Addresses: []string{
			equities[0], strings.ToUpper(equities[1]), equities[1], equities[2], equities[2],
			operator, strings.ToUpper(operator), s.Dao.GlobalDao, strings.ToUpper(s.Dao.GlobalDao),
		},
	})
	s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "chn", Addresses: []string{otherEquity}})
	s.Require().Equal(uint32(51), s.App.RegionGovKeeper.GetParams(s.Ctx).PassThresholdPercent)

	withdrawal := s.withdrawal("usa")
	withdrawal.Receiver = sample.AccAddress()
	withdrawal.Amount = sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 100))
	receiver := sdk.MustAccAddressFromBech32(withdrawal.Receiver)
	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, "usa")
	s.Require().True(found)
	treasury := sdk.MustAccAddressFromBech32(region.RegionTreasureAddr)
	before := s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom)

	id := s.submit(equities[0], withdrawal)
	s.Require().Equal(operatorAccount(s.usaValidator.OperatorAddress), s.mustProposal(id).RequiredOperator)
	s.Require().ElementsMatch(equities, s.equityProposalElectorAccounts(id))
	for _, approver := range []string{operator, s.Dao.GlobalDao} {
		_, found := s.App.RegionGovKeeper.GetElector(s.Ctx, id, approver)
		s.Require().False(found, "a separate approver must not be a counted equity")
	}
	for _, outsider := range []string{sample.AccAddress(), otherEquity, otherOperator, operatorAccount(s.meEarthValidator.OperatorAddress)} {
		s.Require().ErrorIs(s.vote(outsider, id, types.VoteOptionYes), types.ErrUnauthorized)
	}

	// The two other equities can pass without a ballot from the proposer.
	s.Require().NoError(s.vote(strings.ToUpper(equities[1]), id, types.VoteOptionYes))
	s.Require().ErrorIs(s.vote(equities[1], id, types.VoteOptionNo), types.ErrAlreadyVoted)
	s.Require().NoError(s.vote(equities[2], id, types.VoteOptionYes))
	s.Require().NoError(s.vote(operator, id, types.VoteOptionYes))
	s.Require().ErrorIs(s.vote(operator, id, types.VoteOptionNo), types.ErrAlreadyVoted)
	s.Require().NoError(s.vote(s.Dao.GlobalDao, id, types.VoteOptionYes))
	s.Require().Equal(types.ProposalStatusVoting, s.mustProposal(id).Status)
	s.Require().Equal(before, s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom))
	s.Require().True(s.App.BankKeeper.GetBalance(s.Ctx, receiver, params.BaseDenom).IsZero())

	end := s.mustProposal(id).VotingEndTime
	s.endBlockAt(end.Add(-time.Second))
	s.Require().Equal(types.ProposalStatusVoting, s.mustProposal(id).Status)
	s.Ctx = s.Ctx.WithBlockTime(end)
	s.Require().ErrorIs(s.vote(equities[0], id, types.VoteOptionYes), types.ErrVotingClosed)
	s.endBlockAt(end)
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(id).Status)
	s.Require().Equal(before.Amount.SubRaw(100), s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom).Amount)
	s.Require().Equal(sdk.NewInt(100), s.App.BankKeeper.GetBalance(s.Ctx, receiver, params.BaseDenom).Amount)

	s.Ctx = s.Ctx.WithEventManager(sdk.NewEventManager())
	s.endBlockAt(end.Add(time.Second))
	s.Require().Empty(s.Ctx.EventManager().Events())
	s.Require().Equal(before.Amount.SubRaw(100), s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom).Amount)
	s.Require().Equal(sdk.NewInt(100), s.App.BankKeeper.GetBalance(s.Ctx, receiver, params.BaseDenom).Amount)
}

func (s *KeeperTestSuite) TestEquityProposalRequiresAllThreeApprovals() {
	operator := operatorAccount(s.usaValidator.OperatorAddress)
	equities := []string{sample.AccAddress(), sample.AccAddress(), sample.AccAddress()}
	s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{
		RegionId: "usa", Addresses: append(append([]string{}, equities...), operator, s.Dao.GlobalDao),
	})
	withdrawal := s.withdrawal("usa")
	withdrawal.Receiver = sample.AccAddress()
	var rejected []uint64
	for _, tc := range []struct {
		name      string
		equityYes int
		operator  types.VoteOption
		daoYes    bool
		daoVeto   bool
	}{
		{name: "missing operator", equityYes: 3, daoYes: true},
		{name: "missing dao", equityYes: 3, operator: types.VoteOptionYes},
		{name: "below equity threshold", equityYes: 1, operator: types.VoteOptionYes, daoYes: true},
		{name: "approvers are not equity votes", operator: types.VoteOptionYes, daoYes: true},
		{name: "operator no is not an immediate veto", equityYes: 3, operator: types.VoteOptionNo, daoYes: true},
		{name: "dao yes can become an immediate veto", equityYes: 3, operator: types.VoteOptionYes, daoYes: true, daoVeto: true},
	} {
		s.Run(tc.name, func() {
			id := s.submit(equities[0], withdrawal)
			rejected = append(rejected, id)
			s.Require().ElementsMatch(equities, s.equityProposalElectorAccounts(id))
			for _, equity := range equities[:tc.equityYes] {
				s.Require().NoError(s.vote(equity, id, types.VoteOptionYes))
			}
			if tc.operator != types.VoteOptionUnspecified {
				s.Require().NoError(s.vote(operator, id, tc.operator))
				s.Require().ErrorIs(s.vote(operator, id, types.VoteOptionYes), types.ErrAlreadyVoted)
			}
			if tc.daoYes {
				s.Require().NoError(s.vote(s.Dao.GlobalDao, id, types.VoteOptionYes))
			}
			s.Require().Equal(types.ProposalStatusVoting, s.mustProposal(id).Status)
			if tc.daoVeto {
				s.Require().NoError(s.vote(s.Dao.GlobalDao, id, types.VoteOptionNo))
				s.Require().Equal(types.DaoVoteNo, s.mustProposal(id).DaoVote)
				s.Require().Equal(types.ProposalStatusRejected, s.mustProposal(id).Status)
				s.Require().ErrorIs(s.vote(s.Dao.GlobalDao, id, types.VoteOptionYes), types.ErrVotingClosed)
			}
		})
	}
	s.endBlockAt(s.mustProposal(rejected[0]).VotingEndTime)
	for _, id := range rejected {
		s.Require().Equal(types.ProposalStatusRejected, s.mustProposal(id).Status)
	}
	s.Require().True(s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(withdrawal.Receiver), params.BaseDenom).IsZero())
}

func (s *KeeperTestSuite) TestEquityProposalUsesInclusiveCurrentThreshold() {
	operator := operatorAccount(s.usaValidator.OperatorAddress)
	equities := make([]string, 100)
	for i := range equities {
		equities[i] = sample.AccAddress()
	}
	for _, tc := range []struct {
		name            string
		eligible, yes   int
		submitThreshold uint32
		finalThreshold  uint32
		want            types.ProposalStatus
	}{
		{"exactly default 51 percent", 100, 51, 51, 51, types.ProposalStatusPassed},
		{"below default 51 percent", 100, 50, 51, 51, types.ProposalStatusRejected},
		{"inclusive custom threshold", 2, 1, 50, 50, types.ProposalStatusPassed},
		{"raised threshold applies at tally", 3, 2, 51, 100, types.ProposalStatusRejected},
		{"lowered threshold applies at tally", 2, 1, 100, 50, types.ProposalStatusPassed},
	} {
		s.Run(tc.name, func() {
			govParams := s.App.RegionGovKeeper.GetParams(s.Ctx)
			govParams.PassThresholdPercent = tc.submitThreshold
			s.Require().NoError(s.App.RegionGovKeeper.SetParams(s.Ctx, govParams))
			s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: equities[:tc.eligible]})
			id := s.submit(equities[0], s.withdrawal("usa"))
			for _, equity := range equities[:tc.yes] {
				s.Require().NoError(s.vote(equity, id, types.VoteOptionYes))
			}
			s.Require().NoError(s.vote(operator, id, types.VoteOptionYes))
			s.Require().NoError(s.vote(s.Dao.GlobalDao, id, types.VoteOptionYes))
			govParams.PassThresholdPercent = tc.finalThreshold
			s.Require().NoError(s.App.RegionGovKeeper.SetParams(s.Ctx, govParams))
			s.endBlockAt(s.mustProposal(id).VotingEndTime)
			s.Require().Equal(tc.want, s.mustProposal(id).Status)
		})
	}
}

func (s *KeeperTestSuite) TestEquityProposalAuthorizationAndOperatorRouting() {
	usa := operatorAccount(s.usaValidator.OperatorAddress)
	chn := s.addCountryRegion("chn")
	equity := sample.AccAddress()
	outsider := sample.AccAddress()
	// Operators and the DAO retain their existing proposal route even when
	// listed as equities, including an operator of a different active region.
	s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: []string{equity, usa, chn, s.Dao.GlobalDao}})
	s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "chn", Addresses: []string{equity}})
	var countryProposals []uint64
	for _, tc := range []struct {
		name     string
		proposer string
		regions  []string
	}{
		{"dao also equity", s.Dao.GlobalDao, []string{"usa", "chn"}},
		{"own operator also equity", usa, []string{"usa"}},
		{"operator mixed treasuries unchanged", usa, []string{"usa", wstakingtypes.MeEarthRegionId, wstakingtypes.ExperienceRegionId}},
		{"other active operator also equity", chn, []string{"usa"}},
	} {
		s.Run(tc.name, func() {
			var withdrawals []types.Withdrawal
			for _, regionID := range tc.regions {
				withdrawals = append(withdrawals, s.withdrawal(regionID))
			}
			id := s.submit(tc.proposer, withdrawals...)
			countryProposals = append(countryProposals, id)
			s.Require().Equal(types.ElectorateKindCountryRegionOperators, s.mustProposal(id).Electorate)
			s.Require().Empty(s.mustProposal(id).RequiredOperator)
			s.Require().ElementsMatch([]string{usa, chn}, s.electorAccounts(id))
			s.Require().ErrorIs(s.vote(equity, id, types.VoteOptionYes), types.ErrUnauthorized)
			s.Require().NoError(s.vote(usa, id, types.VoteOptionYes))
			s.Require().NoError(s.vote(chn, id, types.VoteOptionYes))
			s.Require().NoError(s.vote(s.Dao.GlobalDao, id, types.VoteOptionYes))
		})
	}

	for _, tc := range []struct {
		name     string
		proposer string
		regions  []string
	}{
		{"outsider", outsider, []string{"usa"}},
		{"foreign special region", equity, []string{wstakingtypes.MeEarthRegionId}},
		{"foreign experience region", equity, []string{wstakingtypes.ExperienceRegionId}},
		{"two equity memberships still cannot combine regions", equity, []string{"usa", "chn"}},
		{"duplicate withdrawal region", equity, []string{"usa", "usa"}},
	} {
		s.Run(tc.name, func() {
			seq := s.App.RegionGovKeeper.GetProposalSeq(s.Ctx)
			events := len(s.Ctx.EventManager().Events())
			var withdrawals []types.Withdrawal
			for _, regionID := range tc.regions {
				withdrawals = append(withdrawals, s.withdrawal(regionID))
			}
			_, err := s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
				Proposer: tc.proposer, Kind: types.ProposalKindTreasuryWithdraw,
				Title: tc.name, VotingPeriod: 60,
				Treasury: &types.TreasuryWithdrawPayload{Withdrawals: withdrawals},
			})
			s.Require().Error(err)
			s.Require().Equal(seq, s.App.RegionGovKeeper.GetProposalSeq(s.Ctx))
			s.Require().Len(s.Ctx.EventManager().Events(), events)
		})
	}
	for _, msg := range []*types.MsgSubmitProposal{
		{
			Proposer: equity, Kind: types.ProposalKindReplaceRegionOperator, Title: "equity cannot replace operator", VotingPeriod: 60,
			ReplaceOperator: &types.ReplaceRegionOperatorPayload{RegionId: "usa", NewOperator: sdk.ValAddress(sdk.MustAccAddressFromBech32(outsider)).String()},
		},
	} {
		_, err := s.msgServer.SubmitProposal(s.Ctx, msg)
		s.Require().ErrorIs(err, types.ErrUnauthorized)
	}

	s.endBlockAt(s.mustProposal(countryProposals[0]).VotingEndTime)
	for _, id := range countryProposals {
		s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(id).Status)
	}
}

func (s *KeeperTestSuite) TestEquityProposalFreezePreservesSnapshotAndVotes() {
	operator := operatorAccount(s.usaValidator.OperatorAddress)
	equities := []string{sample.AccAddress(), sample.AccAddress(), sample.AccAddress(), sample.AccAddress()}
	s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: equities})
	_, err := s.msgServer.FreezeEquity(s.Ctx, &types.MsgFreezeEquity{Authority: s.Dao.GlobalDao, RegionId: "usa", Address: equities[3]})
	s.Require().NoError(err)
	_, err = s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer: equities[3], Kind: types.ProposalKindTreasuryWithdraw, Title: "frozen proposer", VotingPeriod: 60,
		Treasury: &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{s.withdrawal("usa")}},
	})
	s.Require().ErrorIs(err, types.ErrEquityFrozen)

	alreadyVoted := s.submit(equities[0], s.withdrawal("usa"))
	notYetVoted := s.submit(equities[0], s.withdrawal("usa"))
	for _, id := range []uint64{alreadyVoted, notYetVoted} {
		s.Require().ElementsMatch(equities[:3], s.equityProposalElectorAccounts(id))
		s.Require().NoError(s.vote(operator, id, types.VoteOptionYes))
		s.Require().NoError(s.vote(s.Dao.GlobalDao, id, types.VoteOptionYes))
	}
	for _, equity := range equities[1:3] {
		s.Require().NoError(s.vote(equity, alreadyVoted, types.VoteOptionYes))
		_, err = s.msgServer.FreezeEquity(s.Ctx, &types.MsgFreezeEquity{Authority: s.Dao.GlobalDao, RegionId: "usa", Address: equity})
		s.Require().NoError(err)
		s.Require().ErrorIs(s.vote(equity, notYetVoted, types.VoteOptionYes), types.ErrEquityFrozen)
	}
	_, err = s.msgServer.UnfreezeEquity(s.Ctx, &types.MsgUnfreezeEquity{Authority: s.Dao.GlobalDao, RegionId: "usa", Address: equities[3]})
	s.Require().NoError(err)
	s.Require().ErrorIs(s.vote(equities[3], notYetVoted, types.VoteOptionYes), types.ErrUnauthorized)
	s.Require().NoError(s.vote(equities[0], notYetVoted, types.VoteOptionYes))
	for _, id := range []uint64{alreadyVoted, notYetVoted} {
		s.Require().ElementsMatch(equities[:3], s.equityProposalElectorAccounts(id))
	}
	s.endBlockAt(s.mustProposal(alreadyVoted).VotingEndTime)
	// Earlier yes ballots survive a freeze; later freezes cannot shrink the
	// denominator and make one of the original three equities a majority.
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(alreadyVoted).Status)
	s.Require().Equal(types.ProposalStatusRejected, s.mustProposal(notYetVoted).Status)
}

func (s *KeeperTestSuite) TestEquityProposalRegionFreezeBlocksBothRoles() {
	operator := operatorAccount(s.usaValidator.OperatorAddress)
	equities := []string{sample.AccAddress(), sample.AccAddress()}
	s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: equities})
	id := s.submit(equities[0], s.withdrawal("usa"))
	_, err := s.msgServer.FreezeRegion(s.Ctx, &types.MsgFreezeRegion{Authority: s.Dao.GlobalDao, RegionIds: []string{"usa"}})
	s.Require().NoError(err)
	for _, voter := range []string{equities[0], equities[1], operator} {
		s.Require().ErrorIs(s.vote(voter, id, types.VoteOptionYes), types.ErrRegionFrozen)
	}
	_, err = s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer: equities[0], Kind: types.ProposalKindTreasuryWithdraw, Title: "frozen region", VotingPeriod: 60,
		Treasury: &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{s.withdrawal("usa")}},
	})
	s.Require().ErrorIs(err, types.ErrRegionFrozen)
	s.Require().NoError(s.vote(s.Dao.GlobalDao, id, types.VoteOptionYes))
	s.Require().ElementsMatch(equities, s.equityProposalElectorAccounts(id))
	_, err = s.msgServer.UnfreezeRegion(s.Ctx, &types.MsgUnfreezeRegion{Authority: s.Dao.GlobalDao, RegionIds: []string{"usa"}})
	s.Require().NoError(err)
	for _, voter := range []string{equities[0], equities[1], operator} {
		s.Require().NoError(s.vote(voter, id, types.VoteOptionYes))
	}
	s.endBlockAt(s.mustProposal(id).VotingEndTime)
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(id).Status)
}

func (s *KeeperTestSuite) TestEquityProposalSnapshotSurvivesEquityAndOperatorReplacement() {
	oldOperator := operatorAccount(s.usaValidator.OperatorAddress)
	oldEquities := []string{sample.AccAddress(), sample.AccAddress(), sample.AccAddress()}
	newEquities := []string{sample.AccAddress(), sample.AccAddress(), sample.AccAddress()}
	s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: oldEquities})
	approvedByOldOperator := s.submit(oldEquities[0], s.withdrawal("usa"))
	missingOldOperator := s.submit(oldEquities[0], s.withdrawal("usa"))

	update, err := s.msgServer.SubmitEquityUpdate(s.Ctx, &types.MsgSubmitEquityUpdate{Proposer: oldOperator, RegionId: "usa", Addresses: newEquities})
	s.Require().NoError(err)
	_, err = s.msgServer.ApprovePendingAction(s.Ctx, &types.MsgApprovePendingAction{Authority: s.Dao.GlobalDao, ActionId: update.ActionId})
	s.Require().NoError(err)
	newOperator := sdk.MustAccAddressFromBech32(sample.AccAddress())
	newValoper := sdk.ValAddress(newOperator).String()
	s.Require().NoError(s.App.StakingKeeper.SetRegionOperator(s.Ctx, "usa", newValoper))

	for _, id := range []uint64{approvedByOldOperator, missingOldOperator} {
		s.Require().Equal(operatorAccount(s.usaValidator.OperatorAddress), s.mustProposal(id).RequiredOperator)
		s.Require().ElementsMatch(oldEquities, s.equityProposalElectorAccounts(id))
		for _, newEquity := range newEquities {
			s.Require().ErrorIs(s.vote(newEquity, id, types.VoteOptionYes), types.ErrUnauthorized)
		}
		s.Require().ErrorIs(s.vote(newOperator.String(), id, types.VoteOptionYes), types.ErrUnauthorized)
		// Removed equity accounts retain their snapshotted voting rights.
		for _, oldEquity := range oldEquities[1:] {
			s.Require().NoError(s.vote(oldEquity, id, types.VoteOptionYes))
		}
		s.Require().NoError(s.vote(s.Dao.GlobalDao, id, types.VoteOptionYes))
	}
	s.Require().NoError(s.vote(oldOperator, approvedByOldOperator, types.VoteOptionYes))

	fresh := s.submit(newEquities[0], s.withdrawal("usa"))
	s.Require().Equal(newOperator.String(), s.mustProposal(fresh).RequiredOperator)
	s.Require().ElementsMatch(newEquities, s.equityProposalElectorAccounts(fresh))
	s.Require().ErrorIs(s.vote(oldOperator, fresh, types.VoteOptionYes), types.ErrUnauthorized)
	for _, equity := range newEquities[1:] {
		s.Require().NoError(s.vote(equity, fresh, types.VoteOptionYes))
	}
	s.Require().NoError(s.vote(newOperator.String(), fresh, types.VoteOptionYes))
	s.Require().NoError(s.vote(s.Dao.GlobalDao, fresh, types.VoteOptionYes))
	s.endBlockAt(s.mustProposal(fresh).VotingEndTime)
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(approvedByOldOperator).Status)
	s.Require().Equal(types.ProposalStatusRejected, s.mustProposal(missingOldOperator).Status)
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(fresh).Status)
}

func (s *KeeperTestSuite) TestEquityProposalGenesisRoundTrip() {
	operator := operatorAccount(s.usaValidator.OperatorAddress)
	equities := []string{sample.AccAddress(), sample.AccAddress(), sample.AccAddress()}
	s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: equities})
	equityProposal := s.submit(equities[0], s.withdrawal("usa"))
	countryProposal := s.submit(s.Dao.GlobalDao, s.withdrawal("usa"))
	for _, equity := range equities[1:] {
		s.Require().NoError(s.vote(equity, equityProposal, types.VoteOptionYes))
	}
	for _, id := range []uint64{equityProposal, countryProposal} {
		s.Require().NoError(s.vote(operator, id, types.VoteOptionYes))
		s.Require().NoError(s.vote(s.Dao.GlobalDao, id, types.VoteOptionYes))
	}
	genesis := regiongov.ExportGenesis(s.Ctx, s.App.RegionGovKeeper)
	s.Require().NoError(genesis.Validate())
	s.Require().ElementsMatch(equities, s.equityProposalElectorAccounts(equityProposal))
	s.Require().Equal([]string{operator}, s.electorAccounts(countryProposal))

	// Reuse the funded app and its real staking/bank state, but clear this
	// module so import must restore both electorates and the voting queue.
	store := s.Ctx.KVStore(s.App.GetKey(types.StoreKey))
	iterator := store.Iterator(nil, nil)
	var keys [][]byte
	for ; iterator.Valid(); iterator.Next() {
		keys = append(keys, append([]byte(nil), iterator.Key()...))
	}
	s.Require().NoError(iterator.Close())
	for _, key := range keys {
		store.Delete(key)
	}
	s.Require().Empty(s.App.RegionGovKeeper.GetAllProposals(s.Ctx))
	regiongov.InitGenesis(s.Ctx, s.App.RegionGovKeeper, *genesis)
	s.Require().Equal(genesis, regiongov.ExportGenesis(s.Ctx, s.App.RegionGovKeeper))
	s.Require().Equal(operatorAccount(s.usaValidator.OperatorAddress), s.mustProposal(equityProposal).RequiredOperator)
	s.Require().Equal(types.DaoVoteYes, s.mustProposal(equityProposal).DaoVote)
	s.Require().ElementsMatch(equities, s.equityProposalElectorAccounts(equityProposal))
	s.Require().Equal([]string{operator}, s.electorAccounts(countryProposal))
	vote, found := s.App.RegionGovKeeper.GetVote(s.Ctx, equityProposal, operator)
	s.Require().True(found)
	s.Require().Equal(types.VoteOptionYes, vote.Option)
	s.Require().ErrorIs(s.vote(operator, equityProposal, types.VoteOptionNo), types.ErrAlreadyVoted)
	s.Require().ErrorIs(s.vote(equities[1], equityProposal, types.VoteOptionNo), types.ErrAlreadyVoted)
	s.endBlockAt(s.mustProposal(equityProposal).VotingEndTime)
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(equityProposal).Status)
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(countryProposal).Status)
	s.Require().NoError(regiongov.ExportGenesis(s.Ctx, s.App.RegionGovKeeper).Validate())
}
