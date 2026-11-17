package keeper_test

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

func (s *KeeperTestSuite) TestEquityVoteRequiresCurrentMembership() {
	const membersReplaced = "members replaced"
	for _, state := range []string{membersReplaced, "membership cleared", "membership missing"} {
		s.Run(state, func() {
			s.SetupTest()
			operator := operatorAccount(s.usaValidator.OperatorAddress)
			equities := []string{sample.AccAddress(), sample.AccAddress(), sample.AccAddress()}
			newEquity := sample.AccAddress()
			k := s.App.RegionGovKeeper
			k.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: equities})
			alreadyVoted := s.submit(equities[0], s.withdrawal("usa"))
			notYetVoted := s.submit(equities[0], s.withdrawal("usa"))
			for _, voter := range equities[1:] {
				s.Require().NoError(s.vote(voter, alreadyVoted, types.VoteOptionYes))
			}

			switch state {
			case membersReplaced:
				// Exercise the real update path: it canonicalizes uppercase input
				// before DAO approval persists the replacement membership.
				update, err := s.msgServer.SubmitEquityUpdate(s.Ctx, &types.MsgSubmitEquityUpdate{
					Proposer: operator, RegionId: "usa", Addresses: []string{strings.ToUpper(equities[0]), newEquity},
				})
				s.Require().NoError(err)
				equity, found := k.GetEquity(s.Ctx, "usa")
				s.Require().True(found)
				s.Require().Equal(equities, equity.Addresses, "membership must not change before DAO approval")
				_, err = s.msgServer.ApprovePendingAction(s.Ctx, &types.MsgApprovePendingAction{
					Authority: s.Dao.GlobalDao, ActionId: update.ActionId,
				})
				s.Require().NoError(err)
				equity, found = k.GetEquity(s.Ctx, "usa")
				s.Require().True(found)
				s.Require().Equal([]string{equities[0], newEquity}, equity.Addresses, "approved addresses must be canonical")
			case "membership cleared":
				k.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa"})
			case "membership missing":
				s.Ctx.KVStore(s.App.GetKey(types.StoreKey)).Delete(types.EquityKey("usa"))
			}
			// Membership in another region cannot preserve this proposal's rights.
			s.addCountryRegion("chn")
			k.SetEquity(s.Ctx, types.RegionEquity{RegionId: "chn", Addresses: equities[1:]})
			for _, id := range []uint64{alreadyVoted, notYetVoted} {
				before := s.mustProposal(id)
				votes := k.GetVotes(s.Ctx, id)
				events := len(s.Ctx.EventManager().Events())
				for _, voter := range []string{equities[1], equities[2], newEquity} {
					for _, option := range []types.VoteOption{types.VoteOptionYes, types.VoteOptionNo} {
						s.Require().ErrorIs(s.vote(voter, id, option), types.ErrUnauthorized)
					}
				}
				s.Require().Equal(before, s.mustProposal(id))
				s.Require().Equal(votes, k.GetVotes(s.Ctx, id))
				s.Require().Len(s.Ctx.EventManager().Events(), events)
				s.Require().ElementsMatch(equities, s.equityProposalElectorAccounts(id))
				if state == membersReplaced {
					s.Require().NoError(s.vote(equities[0], id, types.VoteOptionYes))
				} else {
					s.Require().ErrorIs(s.vote(equities[0], id, types.VoteOptionYes), types.ErrUnauthorized)
				}
				// Neither independent approver is required to be an equity member.
				s.Require().NoError(s.vote(operator, id, types.VoteOptionYes))
				s.Require().NoError(s.vote(s.Dao.GlobalDao, id, types.VoteOptionYes))
			}
			s.endBlockAt(s.mustProposal(alreadyVoted).VotingEndTime)
			s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(alreadyVoted).Status, "existing ballots survive removal")
			s.Require().Equal(types.ProposalStatusRejected, s.mustProposal(notYetVoted).Status, "removal must not reduce the original denominator")
		})
	}
}

func (s *KeeperTestSuite) TestEquityVoteAllowedAfterMembershipRestored() {
	equity := sample.AccAddress()
	k := s.App.RegionGovKeeper
	k.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: []string{equity}})
	id := s.submit(equity, s.withdrawal("usa"))
	k.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa"})
	s.Require().ErrorIs(s.vote(equity, id, types.VoteOptionYes), types.ErrUnauthorized)
	k.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: []string{equity}})
	s.Require().NoError(s.vote(equity, id, types.VoteOptionYes))
	s.Require().ErrorIs(s.vote(equity, id, types.VoteOptionNo), types.ErrAlreadyVoted)
}

func (s *KeeperTestSuite) TestOperatorVoteRequiresCurrentRegionOwnership() {
	for _, kind := range []types.ProposalKind{types.ProposalKindTreasuryWithdraw, types.ProposalKindReplaceRegionOperator} {
		s.Run(kind.String(), func() {
			s.SetupTest()
			oldOperator := operatorAccount(s.usaValidator.OperatorAddress)
			chn := s.addCountryRegion("chn")
			jpn := s.addCountryRegion("jpn")
			msg := &types.MsgSubmitProposal{
				Proposer: s.Dao.GlobalDao, Kind: kind, Title: "current operator required", VotingPeriod: 60,
			}
			if kind == types.ProposalKindTreasuryWithdraw {
				msg.Treasury = &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{s.withdrawal("usa")}}
			} else {
				msg.ReplaceOperator = &types.ReplaceRegionOperatorPayload{
					RegionId: "jpn", NewOperator: sdk.ValAddress(sdk.MustAccAddressFromBech32(sample.AccAddress())).String(),
				}
			}
			alreadyVoted, err := s.msgServer.SubmitProposal(s.Ctx, msg)
			s.Require().NoError(err)
			notYetVoted, err := s.msgServer.SubmitProposal(s.Ctx, msg)
			s.Require().NoError(err)
			s.Require().NoError(s.vote(oldOperator, alreadyVoted.Id, types.VoteOptionYes))
			newOperator := sample.AccAddress()
			s.Require().NoError(s.App.StakingKeeper.SetRegionOperator(s.Ctx, "usa", sdk.ValAddress(sdk.MustAccAddressFromBech32(newOperator)).String()))

			for _, id := range []uint64{alreadyVoted.Id, notYetVoted.Id} {
				before := s.mustProposal(id)
				votes := s.App.RegionGovKeeper.GetVotes(s.Ctx, id)
				events := len(s.Ctx.EventManager().Events())
				for _, voter := range []string{oldOperator, newOperator} {
					for _, option := range []types.VoteOption{types.VoteOptionYes, types.VoteOptionNo} {
						s.Require().ErrorIs(s.vote(voter, id, option), types.ErrUnauthorized)
					}
				}
				s.Require().Equal(before, s.mustProposal(id))
				s.Require().Equal(votes, s.App.RegionGovKeeper.GetVotes(s.Ctx, id))
				s.Require().Len(s.Ctx.EventManager().Events(), events)
				s.Require().ElementsMatch([]string{oldOperator, chn, jpn}, s.electorAccounts(id))
				s.Require().NoError(s.vote(chn, id, types.VoteOptionYes))
				s.Require().NoError(s.vote(s.Dao.GlobalDao, id, types.VoteOptionYes))
			}
			s.endBlockAt(s.mustProposal(alreadyVoted.Id).VotingEndTime)
			s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(alreadyVoted.Id).Status, "earlier operator approval still counts")
			s.Require().Equal(types.ProposalStatusRejected, s.mustProposal(notYetVoted.Id).Status, "lost authority is an absent vote, not a smaller denominator")
		})
	}
}
