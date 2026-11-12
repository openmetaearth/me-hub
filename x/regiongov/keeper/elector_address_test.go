package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

func (s *KeeperTestSuite) TestElectorsUseAccountAddresses() {
	usa := operatorAccount(s.usaValidator.OperatorAddress)
	chn := s.addCountryRegion("chn")
	equity := sample.AccAddress()
	s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: []string{equity}})
	id := s.submit(usa, s.withdrawal("usa"))

	response, err := s.App.RegionGovKeeper.Electors(sdk.WrapSDKContext(s.Ctx), &types.QueryElectorsRequest{ProposalId: id})
	s.Require().NoError(err)
	accounts := make([]string, 0, len(response.Electors))
	store := s.Ctx.KVStore(s.App.GetKey(types.StoreKey))
	for _, elector := range response.Electors {
		account, err := sdk.AccAddressFromBech32(elector.Operator)
		s.Require().NoError(err)
		s.Require().Equal(account.String(), elector.Operator)
		accounts = append(accounts, elector.Operator)
		var stored types.Elector
		s.App.AppCodec().MustUnmarshal(store.Get(types.ElectorKey(id, elector.Operator)), &stored)
		s.Require().Equal(elector, stored, "new snapshots must store the account, not just convert query output")
	}
	s.Require().ElementsMatch([]string{usa, chn}, accounts)
	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, "usa")
	s.Require().True(found)
	s.Require().Equal(s.usaValidator.OperatorAddress, region.OperatorAddress, "staking retains validator addresses")

	for _, outsider := range []string{sample.AccAddress(), equity, operatorAccount(s.meEarthValidator.OperatorAddress), operatorAccount(s.experienceValidator.OperatorAddress)} {
		s.Require().ErrorIs(s.vote(outsider, id, types.VoteOptionYes), types.ErrUnauthorized)
	}
	s.Require().ErrorIs(s.vote(s.usaValidator.OperatorAddress, id, types.VoteOptionYes), sdkerrors.ErrInvalidAddress)
	s.Require().Empty(s.App.RegionGovKeeper.GetVotes(s.Ctx, id))
	s.Require().NoError(s.vote(usa, id, types.VoteOptionYes))
	s.Require().ErrorIs(s.vote(usa, id, types.VoteOptionNo), types.ErrAlreadyVoted)
	s.Require().NoError(s.vote(chn, id, types.VoteOptionYes))
	s.Require().Equal(types.DaoVoteUnspecified, s.mustProposal(id).DaoVote)
	s.Require().NoError(s.vote(s.Dao.GlobalDao, id, types.VoteOptionYes))
	s.endBlockAt(s.mustProposal(id).VotingEndTime)
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(id).Status)
}

func (s *KeeperTestSuite) TestLegacyElectorsQueryPreservesVotingPermissions() {
	usa := operatorAccount(s.usaValidator.OperatorAddress)
	chn := s.addCountryRegion("chn")
	id := s.submit(usa, s.withdrawal("usa"))

	// Existing state uses account keys but validator-address values. Bypass the
	// setter to simulate an old snapshot without running a state migration.
	store := s.Ctx.KVStore(s.App.GetKey(types.StoreKey))
	before := make(map[string][]byte)
	for _, elector := range s.App.RegionGovKeeper.GetElectors(s.Ctx, id) {
		account, err := elector.AccountAddress()
		s.Require().NoError(err)
		elector.Operator = sdk.ValAddress(account).String()
		wire := s.App.AppCodec().MustMarshal(&elector)
		store.Set(types.ElectorKey(id, account.String()), wire)
		before[account.String()] = wire
	}

	response, err := s.App.RegionGovKeeper.Electors(sdk.WrapSDKContext(s.Ctx), &types.QueryElectorsRequest{ProposalId: id})
	s.Require().NoError(err)
	accounts := make([]string, 0, len(response.Electors))
	for _, elector := range response.Electors {
		accounts = append(accounts, elector.Operator)
	}
	s.Require().ElementsMatch([]string{usa, chn}, accounts)
	for account, wire := range before {
		s.Require().Equal(wire, store.Get(types.ElectorKey(id, account)), "queries must not rewrite old snapshots")
	}

	s.Require().ErrorIs(s.vote(sample.AccAddress(), id, types.VoteOptionYes), types.ErrUnauthorized)
	s.Require().NoError(s.vote(usa, id, types.VoteOptionYes))
	s.Require().ErrorIs(s.vote(usa, id, types.VoteOptionNo), types.ErrAlreadyVoted)
	_, err = s.msgServer.FreezeRegion(s.Ctx, &types.MsgFreezeRegion{Authority: s.Dao.GlobalDao, RegionIds: []string{"chn"}})
	s.Require().NoError(err)
	s.Require().ErrorIs(s.vote(chn, id, types.VoteOptionYes), types.ErrRegionFrozen)
	s.Require().NoError(s.vote(s.Dao.GlobalDao, id, types.VoteOptionYes))
	s.endBlockAt(s.mustProposal(id).VotingEndTime)
	s.Require().Equal(types.ProposalStatusRejected, s.mustProposal(id).Status, "freezing must not reduce the original denominator")
}

func (s *KeeperTestSuite) TestOperatorReplacementKeepsSnapshottedAccount() {
	oldAccount := operatorAccount(s.usaValidator.OperatorAddress)
	newAccount := sample.AccAddress()
	newOperator := sdk.ValAddress(sdk.MustAccAddressFromBech32(newAccount)).String()
	pending, err := s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer: oldAccount, Kind: types.ProposalKindTreasuryWithdraw, Title: "long vote", VotingPeriod: 120,
		Treasury: &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{s.withdrawal("usa")}},
	})
	s.Require().NoError(err)
	replacement, err := s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer: s.Dao.GlobalDao, Kind: types.ProposalKindReplaceRegionOperator, Title: "replace operator", VotingPeriod: 60,
		ReplaceOperator: &types.ReplaceRegionOperatorPayload{RegionId: "usa", NewOperator: newOperator},
	})
	s.Require().NoError(err)
	s.Require().NoError(s.vote(oldAccount, replacement.Id, types.VoteOptionYes))
	s.Require().NoError(s.vote(s.Dao.GlobalDao, replacement.Id, types.VoteOptionYes))
	s.endBlockAt(s.startTime.Add(time.Minute))
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(replacement.Id).Status)
	s.Require().Equal(types.ProposalStatusVoting, s.mustProposal(pending.Id).Status)

	s.Require().ErrorIs(s.vote(newAccount, pending.Id, types.VoteOptionYes), types.ErrUnauthorized)
	s.Require().NoError(s.vote(oldAccount, pending.Id, types.VoteOptionYes))
	s.Require().NoError(s.vote(s.Dao.GlobalDao, pending.Id, types.VoteOptionYes))
	fresh := s.submit(newAccount, s.withdrawal("usa"))
	s.Require().ErrorIs(s.vote(oldAccount, fresh, types.VoteOptionYes), types.ErrUnauthorized)
	s.Require().NoError(s.vote(newAccount, fresh, types.VoteOptionYes))
	s.Require().NoError(s.vote(s.Dao.GlobalDao, fresh, types.VoteOptionYes))
	s.endBlockAt(s.mustProposal(fresh).VotingEndTime)
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(pending.Id).Status)
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(fresh).Status)
}
