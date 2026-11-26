package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

func (s *KeeperTestSuite) TestDaoCanVoteOnlyOncePerProposal() {
	const replaceOperator = "replace operator"
	for _, electorate := range []string{"treasury", replaceOperator, "equity treasury"} {
		for _, tc := range []struct {
			name     string
			option   types.VoteOption
			ballot   types.DaoVote
			yesCount int
			want     types.ProposalStatus
		}{
			{"dao veto", types.VoteOptionNo, types.DaoVoteNo, 2, types.ProposalStatusRejected},
			{"dao approval", types.VoteOptionYes, types.DaoVoteYes, 2, types.ProposalStatusPassed},
			{"dao approval below threshold", types.VoteOptionYes, types.DaoVoteYes, 1, types.ProposalStatusRejected},
		} {
			s.Run(electorate+"/"+tc.name, func() {
				s.SetupTest()
				operator := operatorAccount(s.usaValidator.OperatorAddress)
				voters := []string{operator, s.addCountryRegion("chn"), s.addCountryRegion("jpn")}
				proposer := operator
				if electorate == "equity treasury" {
					voters = []string{sample.AccAddress(), sample.AccAddress(), sample.AccAddress()}
					proposer = voters[0]
					s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: voters})
				}
				withdrawal := s.withdrawal("usa")
				withdrawal.Receiver = sample.AccAddress()
				newOperator := sdk.ValAddress(sdk.MustAccAddressFromBech32(sample.AccAddress()))
				msg := &types.MsgSubmitProposal{
					Proposer: proposer, Kind: types.ProposalKindTreasuryWithdraw, Title: "single dao ballot", VotingPeriod: 60,
					Treasury: &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{withdrawal}},
				}
				if electorate == replaceOperator {
					msg.Kind = types.ProposalKindReplaceRegionOperator
					msg.Treasury = nil
					msg.ReplaceOperator = &types.ReplaceRegionOperatorPayload{RegionId: "usa", NewOperator: newOperator.String()}
				}
				response, err := s.msgServer.SubmitProposal(s.Ctx, msg)
				s.Require().NoError(err)
				id := response.Id
				proposal := s.mustProposal(id)
				s.Require().ElementsMatch(voters, s.electorAccounts(id))
				s.Require().Equal(uint32(51), s.App.RegionGovKeeper.GetParams(s.Ctx).PassThresholdPercent)
				queue := types.VotingQueueKey(proposal.VotingEndTime, id)
				store := s.Ctx.KVStore(s.App.GetKey(types.StoreKey))
				region, found := s.App.StakingKeeper.GetRegion(s.Ctx, "usa")
				s.Require().True(found)
				treasury := sdk.MustAccAddressFromBech32(region.RegionTreasureAddr)
				receiver := sdk.MustAccAddressFromBech32(withdrawal.Receiver)
				balance := s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom)

				// The first DAO ballot does not close voting for the electors.
				s.Ctx = s.Ctx.WithEventManager(sdk.NewEventManager())
				s.Require().NoError(s.vote(s.Dao.GlobalDao, id, tc.option))
				expected := proposal
				expected.DaoVote = tc.ballot
				s.Require().Equal(expected, s.mustProposal(id))
				s.Require().True(store.Has(queue))
				s.Require().Empty(s.App.RegionGovKeeper.GetVotes(s.Ctx, id))
				s.Require().Len(s.Ctx.EventManager().Events(), 1)
				s.Require().Equal(types.EventTypeVote, s.Ctx.EventManager().Events()[0].Type)
				for _, voter := range voters[:tc.yesCount] {
					s.Require().NoError(s.vote(voter, id, types.VoteOptionYes))
				}
				if electorate == "equity treasury" {
					s.Require().NoError(s.vote(operator, id, types.VoteOptionYes))
				}
				votes := s.App.RegionGovKeeper.GetVotes(s.Ctx, id)
				// Neither repeating the first vote nor changing it is allowed, even
				// in a later block when the proposal is still open.
				s.Ctx = s.Ctx.WithBlockTime(proposal.VotingEndTime.Add(-time.Nanosecond))
				s.Ctx = s.Ctx.WithEventManager(sdk.NewEventManager())
				for _, option := range []types.VoteOption{types.VoteOptionYes, types.VoteOptionNo} {
					s.Require().ErrorIs(s.vote(s.Dao.GlobalDao, id, option), types.ErrAlreadyVoted)
				}
				s.Require().Equal(expected, s.mustProposal(id))
				s.Require().Equal(votes, s.App.RegionGovKeeper.GetVotes(s.Ctx, id))
				s.Require().True(store.Has(queue))
				s.Require().Empty(s.Ctx.EventManager().Events())
				s.endBlockAt(s.Ctx.BlockTime())
				s.Require().Equal(types.ProposalStatusVoting, s.mustProposal(id).Status)
				currentRegion, found := s.App.StakingKeeper.GetRegion(s.Ctx, "usa")
				s.Require().True(found)
				s.Require().Equal(region, currentRegion)
				s.Require().Equal(balance, s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom))
				s.Require().ErrorIs(s.vote(voters[0], id, types.VoteOptionNo), types.ErrAlreadyVoted)
				before := s.mustProposal(id)
				s.Ctx = s.Ctx.WithBlockTime(proposal.VotingEndTime)
				for _, option := range []types.VoteOption{types.VoteOptionYes, types.VoteOptionNo} {
					s.Require().ErrorIs(s.vote(s.Dao.GlobalDao, id, option), types.ErrVotingClosed)
				}
				s.Require().Equal(before, s.mustProposal(id), "deadline must reject changes even before EndBlocker runs")
				s.Ctx = s.Ctx.WithEventManager(sdk.NewEventManager())
				s.endBlockAt(proposal.VotingEndTime)
				s.Require().Equal(tc.want, s.mustProposal(id).Status)
				s.Require().Equal(before.DaoVote, s.mustProposal(id).DaoVote)
				s.Require().False(store.Has(queue))
				currentRegion, found = s.App.StakingKeeper.GetRegion(s.Ctx, "usa")
				s.Require().True(found)
				if tc.want == types.ProposalStatusPassed && electorate == replaceOperator {
					s.Require().Equal(newOperator.String(), currentRegion.OperatorAddress)
					validator, err := s.App.StakingKeeper.GetValidator(s.Ctx, newOperator)
					s.Require().NoError(err)
					s.Require().Equal(s.usaValidator.ConsensusPubkey, validator.ConsensusPubkey)
				} else {
					s.Require().Equal(region, currentRegion)
				}
				if tc.want == types.ProposalStatusPassed && electorate != replaceOperator {
					s.Require().Equal(balance.Amount.SubRaw(1), s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom).Amount)
					s.Require().Equal(sdkmath.OneInt(), s.App.BankKeeper.GetBalance(s.Ctx, receiver, params.BaseDenom).Amount)
				} else {
					s.Require().Equal(balance, s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom))
					s.Require().True(s.App.BankKeeper.GetBalance(s.Ctx, receiver, params.BaseDenom).IsZero())
				}
				final := s.mustProposal(id)
				s.Require().ErrorIs(s.vote(s.Dao.GlobalDao, id, types.VoteOptionYes), types.ErrVotingClosed)
				s.Ctx = s.Ctx.WithEventManager(sdk.NewEventManager())
				s.endBlockAt(proposal.VotingEndTime.Add(time.Second))
				s.Require().Equal(final, s.mustProposal(id))
				s.Require().Empty(s.Ctx.EventManager().Events())
			})
		}
	}
}
