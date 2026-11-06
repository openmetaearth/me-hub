package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestSubmitVoteAndEndBlock() {
	s.SetupTest()

	meEarth := operatorAccount(s.meEarthValidator.OperatorAddress)
	experience := operatorAccount(s.experienceValidator.OperatorAddress)
	usa := operatorAccount(s.usaValidator.OperatorAddress)
	dao := s.Dao.GlobalDao

	passed, err := s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer:     usa,
		Kind:         types.ProposalKindTreasuryWithdraw,
		Title:        "global treasury",
		VotingPeriod: 60,
		Treasury: &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{
			s.withdrawal(wstakingtypes.MeEarthRegionId),
		}},
	})
	s.Require().NoError(err)
	s.Require().ElementsMatch([]string{usa}, s.electorAccounts(passed.Id))
	s.Require().NotContains(s.electorAccounts(passed.Id), meEarth)
	s.Require().NotContains(s.electorAccounts(passed.Id), experience)
	stored, found := s.App.RegionGovKeeper.GetProposal(s.Ctx, passed.Id)
	s.Require().True(found)
	s.Require().Equal(types.ProposalStatusVoting, stored.Status)
	s.Require().Equal(types.DaoVoteUnspecified, stored.DaoVote)
	s.Require().Equal(types.ElectorateKindCountryRegionOperators, stored.Electorate)
	s.Require().Equal(s.startTime.Add(60*time.Second), stored.VotingEndTime)

	half := s.submit(usa, s.withdrawal(wstakingtypes.MeEarthRegionId))
	noDao := s.submit(dao, s.withdrawal(wstakingtypes.ExperienceRegionId))
	country := s.submit(usa, s.withdrawal("usa"))
	s.Require().ElementsMatch([]string{usa}, s.electorAccounts(country))
	s.Require().NotContains(s.electorAccounts(country), dao)

	s.Require().ErrorIs(s.vote(sample.AccAddress(), passed.Id, types.VoteOptionYes), types.ErrUnauthorized)
	s.Require().ErrorIs(s.vote(meEarth, passed.Id, types.VoteOptionYes), types.ErrUnauthorized)
	s.Require().ErrorIs(s.vote(experience, passed.Id, types.VoteOptionYes), types.ErrUnauthorized)
	s.Require().NoError(s.vote(usa, passed.Id, types.VoteOptionYes))
	s.Require().ErrorIs(s.vote(usa, passed.Id, types.VoteOptionNo), types.ErrAlreadyVoted)
	s.Require().NoError(s.vote(dao, passed.Id, types.VoteOptionYes))

	s.Require().NoError(s.vote(dao, half, types.VoteOptionYes))

	s.Require().NoError(s.vote(usa, noDao, types.VoteOptionYes))

	s.Require().NoError(s.vote(usa, country, types.VoteOptionYes))
	s.Require().NoError(s.vote(dao, country, types.VoteOptionYes))

	s.endBlockAt(s.startTime)
	for _, id := range []uint64{passed.Id, half, noDao, country} {
		proposal, found := s.App.RegionGovKeeper.GetProposal(s.Ctx, id)
		s.Require().True(found)
		s.Require().Equal(types.ProposalStatusVoting, proposal.Status)
	}

	s.Ctx = s.Ctx.WithBlockTime(stored.VotingEndTime)
	s.Require().ErrorIs(s.vote(sample.AccAddress(), passed.Id, types.VoteOptionYes), types.ErrVotingClosed)

	s.endBlockAt(stored.VotingEndTime)

	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(passed.Id).Status)
	s.Require().Equal(types.ProposalStatusRejected, s.mustProposal(half).Status)
	s.Require().Equal(types.ProposalStatusRejected, s.mustProposal(noDao).Status)
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(country).Status)

	s.endBlockAt(stored.VotingEndTime.Add(time.Second))
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(passed.Id).Status)
}

func (s *KeeperTestSuite) TestDaoVetoBeforeEndBlock() {
	s.SetupTest()
	usa := operatorAccount(s.usaValidator.OperatorAddress)
	withdrawal := s.withdrawal("usa")
	withdrawal.Receiver = sample.AccAddress()
	id := s.submit(usa, withdrawal)
	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, "usa")
	s.Require().True(found)
	treasury := sdk.MustAccAddressFromBech32(region.RegionTreasureAddr)
	before := s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom)
	// Veto must win even after every operator and the DAO have approved.
	s.Require().NoError(s.vote(usa, id, types.VoteOptionYes))
	s.Require().NoError(s.vote(s.Dao.GlobalDao, id, types.VoteOptionYes))
	s.Require().Equal(types.DaoVoteYes, s.mustProposal(id).DaoVote)
	s.Require().Equal(types.ProposalStatusVoting, s.mustProposal(id).Status)

	s.Require().NoError(s.vote(s.Dao.GlobalDao, id, types.VoteOptionNo))
	s.Require().Equal(types.DaoVoteNo, s.mustProposal(id).DaoVote)
	s.Require().Equal(types.ProposalStatusRejected, s.mustProposal(id).Status)
	s.Require().ErrorIs(s.vote(s.Dao.GlobalDao, id, types.VoteOptionYes), types.ErrVotingClosed)

	s.Ctx = s.Ctx.WithEventManager(sdk.NewEventManager())
	s.endBlockAt(s.mustProposal(id).VotingEndTime)
	s.Require().Equal(types.ProposalStatusRejected, s.mustProposal(id).Status)
	s.Require().Empty(s.Ctx.EventManager().Events())
	s.Require().Equal(before, s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom))
	receiver := sdk.MustAccAddressFromBech32(withdrawal.Receiver)
	s.Require().True(s.App.BankKeeper.GetBalance(s.Ctx, receiver, params.BaseDenom).IsZero())
}

func (s *KeeperTestSuite) TestSubmitAuthorization() {
	s.SetupTest()
	usa := operatorAccount(s.usaValidator.OperatorAddress)
	meEarth := operatorAccount(s.meEarthValidator.OperatorAddress)

	_, err := s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer:     usa,
		Kind:         types.ProposalKindTreasuryWithdraw,
		Title:        "too short",
		VotingPeriod: 1,
		Treasury:     &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{s.withdrawal("usa")}},
	})
	s.Require().ErrorIs(err, types.ErrInvalidVotingPeriod)

	_, err = s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer:     sample.AccAddress(),
		Kind:         types.ProposalKindTreasuryWithdraw,
		Title:        "outsider",
		VotingPeriod: 60,
		Treasury:     &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{s.withdrawal("usa")}},
	})
	s.Require().ErrorIs(err, types.ErrUnauthorized)

	mixed, err := s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer:     s.Dao.GlobalDao,
		Kind:         types.ProposalKindTreasuryWithdraw,
		Title:        "mixed",
		VotingPeriod: 60,
		Treasury: &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{
			s.withdrawal(wstakingtypes.MeEarthRegionId),
			s.withdrawal("usa"),
		}},
	})
	s.Require().NoError(err)
	s.Require().Equal(types.ElectorateKindCountryRegionOperators, s.mustProposal(mixed.Id).Electorate)
	s.Require().ElementsMatch([]string{usa}, s.electorAccounts(mixed.Id))

	_, err = s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer:     usa,
		Kind:         types.ProposalKindTreasuryWithdraw,
		Title:        "own country with special regions",
		VotingPeriod: 60,
		Treasury: &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{
			s.withdrawal(wstakingtypes.MeEarthRegionId),
			s.withdrawal("usa"),
			s.withdrawal(wstakingtypes.ExperienceRegionId),
		}},
	})
	s.Require().NoError(err)

	res, err := s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer:     meEarth,
		Kind:         types.ProposalKindReplaceRegionOperator,
		Title:        "replace",
		VotingPeriod: 60,
		ReplaceOperator: &types.ReplaceRegionOperatorPayload{
			RegionId:    "usa",
			NewOperator: sdk.ValAddress(sdk.MustAccAddressFromBech32(sample.AccAddress())).String(),
		},
	})
	s.Require().NoError(err)
	s.Require().ElementsMatch([]string{usa}, s.electorAccounts(res.Id))

	s.App.RegionGovKeeper.SetRegionFreeze(s.Ctx, types.RegionFreeze{RegionId: "usa", Frozen: true})
	_, err = s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer:     usa,
		Kind:         types.ProposalKindTreasuryWithdraw,
		Title:        "frozen",
		VotingPeriod: 60,
		Treasury:     &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{s.withdrawal(wstakingtypes.MeEarthRegionId)}},
	})
	s.Require().ErrorIs(err, types.ErrRegionFrozen)
}

func (s *KeeperTestSuite) TestExecuteTreasuryAndReplaceOperator() {
	s.SetupTest()
	usa := operatorAccount(s.usaValidator.OperatorAddress)
	dao := s.Dao.GlobalDao
	receiver := sample.AccAddress()

	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, "usa")
	s.Require().True(found)
	before := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(region.RegionTreasureAddr), params.BaseDenom)
	amount := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 100))

	id := s.submit(usa, types.Withdrawal{
		RegionId: "usa",
		Receiver: receiver,
		Amount:   amount,
	})
	s.Require().NoError(s.vote(usa, id, types.VoteOptionYes))
	s.Require().NoError(s.vote(dao, id, types.VoteOptionYes))
	s.Require().Equal(types.ProposalStatusVoting, s.mustProposal(id).Status)

	s.endBlockAt(s.mustProposal(id).VotingEndTime)
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(id).Status)

	after := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(region.RegionTreasureAddr), params.BaseDenom)
	got := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(receiver), params.BaseDenom)
	s.Require().Equal(before.Amount.Sub(amount[0].Amount), after.Amount)
	s.Require().Equal(amount[0].Amount, got.Amount)

	newOwner := sdk.MustAccAddressFromBech32(sample.AccAddress())
	newOperator := sdk.ValAddress(newOwner)
	replace, err := s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer:     dao,
		Kind:         types.ProposalKindReplaceRegionOperator,
		Title:        "replace usa",
		VotingPeriod: 60,
		ReplaceOperator: &types.ReplaceRegionOperatorPayload{
			RegionId:    "usa",
			NewOperator: newOperator.String(),
		},
	})
	s.Require().NoError(err)
	s.Require().NoError(s.vote(usa, replace.Id, types.VoteOptionYes))
	s.Require().NoError(s.vote(dao, replace.Id, types.VoteOptionYes))
	s.Require().Equal(types.ProposalStatusVoting, s.mustProposal(replace.Id).Status)

	s.endBlockAt(s.mustProposal(replace.Id).VotingEndTime)
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(replace.Id).Status)

	updated, found := s.App.StakingKeeper.GetRegion(s.Ctx, "usa")
	s.Require().True(found)
	s.Require().Equal(newOperator.String(), updated.OperatorAddress)
	oldAddr, err := sdk.ValAddressFromBech32(s.usaValidator.OperatorAddress)
	s.Require().NoError(err)
	_, found = s.App.StakingKeeper.GetValidator(s.Ctx, oldAddr)
	s.Require().False(found)
	validator := s.App.StakingKeeper.MustGetValidator(s.Ctx, newOperator)
	s.Require().Equal("usa", validator.Description.RegionID)
	s.Require().Equal(newOwner.String(), validator.OwnerAddress)
	s.Require().Equal(s.usaValidator.ConsensusPubkey, validator.ConsensusPubkey)
	s.Require().Equal(s.usaValidator.Tokens, validator.Tokens)
	resetEvents := 0
	for _, event := range s.Ctx.EventManager().Events() {
		if event.Type == wstakingtypes.EventTypeResetValidator {
			resetEvents++
		}
	}
	s.Require().Equal(1, resetEvents)
	globalRegion, found := s.App.StakingKeeper.GetRegion(s.Ctx, wstakingtypes.MeEarthRegionId)
	s.Require().True(found)
	s.Require().Equal(s.meEarthValidator.OperatorAddress, globalRegion.OperatorAddress)
	updates, err := s.App.StakingKeeper.ApplyAndReturnValidatorSetUpdates(s.Ctx)
	s.Require().NoError(err)
	s.Require().Empty(updates)
}

func (s *KeeperTestSuite) TestReplaceOperatorRejectsExistingValidatorAtExecution() {
	s.SetupTest()
	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, "usa")
	s.Require().True(found)
	validators := s.App.StakingKeeper.GetAllValidators(s.Ctx)
	stakes := s.App.StakingKeeper.GetAllStakes(s.Ctx)
	replace, err := s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer: s.Dao.GlobalDao, Kind: types.ProposalKindReplaceRegionOperator,
		Title: "cannot steal another region's validator", VotingPeriod: 60,
		ReplaceOperator: &types.ReplaceRegionOperatorPayload{RegionId: "usa", NewOperator: s.meEarthValidator.OperatorAddress},
	})
	s.Require().NoError(err)
	s.Require().NoError(s.vote(operatorAccount(region.OperatorAddress), replace.Id, types.VoteOptionYes))
	s.Require().NoError(s.vote(s.Dao.GlobalDao, replace.Id, types.VoteOptionYes))
	s.endBlockAt(s.mustProposal(replace.Id).VotingEndTime)
	s.Require().Equal(types.ProposalStatusFailed, s.mustProposal(replace.Id).Status)
	unchanged, found := s.App.StakingKeeper.GetRegion(s.Ctx, "usa")
	s.Require().True(found)
	s.Require().Equal(region, unchanged)
	s.Require().Equal(validators, s.App.StakingKeeper.GetAllValidators(s.Ctx))
	s.Require().Equal(stakes, s.App.StakingKeeper.GetAllStakes(s.Ctx))
}

func (s *KeeperTestSuite) mustProposal(id uint64) types.Proposal {
	proposal, found := s.App.RegionGovKeeper.GetProposal(s.Ctx, id)
	s.Require().True(found)
	return proposal
}
