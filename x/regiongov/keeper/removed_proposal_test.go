package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/protobuf/encoding/protowire"

	"github.com/openmetaearth/me-hub/x/regiongov"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

func (s *KeeperTestSuite) TestRemovedProposalSubmissionRejected() {
	before := regiongov.ExportGenesis(s.Ctx, s.App.RegionGovKeeper)
	_, err := s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer: s.Dao.GlobalDao, Kind: types.ProposalKind(3),
		Title: "removed proposal", VotingPeriod: 60,
	})
	s.Require().ErrorContains(err, "unknown proposal kind")
	s.Require().Equal(before, regiongov.ExportGenesis(s.Ctx, s.App.RegionGovKeeper))
	s.Require().Empty(s.App.RegionGovKeeper.DequeueMatureProposals(s.Ctx, s.startTime.Add(time.Hour)))
}

func (s *KeeperTestSuite) TestRemovedProposalHistoricalFinalization() {
	for _, status := range []types.ProposalStatus{types.ProposalStatusFailed, types.ProposalStatusRejected} {
		s.Run(status.String(), func() {
			s.SetupTest()
			operator := operatorAccount(s.usaValidator.OperatorAddress)
			id := s.submit(operator, s.withdrawal("usa"))
			if status == types.ProposalStatusFailed {
				// Reproduce ballots already recorded before this proposal kind was removed.
				s.Require().NoError(s.vote(operator, id, types.VoteOptionYes))
				s.Require().NoError(s.vote(s.Dao.GlobalDao, id, types.VoteOptionYes))
			}
			proposal := s.mustProposal(id)
			proposal.Kind = types.ProposalKind(3)
			proposal.Treasury = nil
			wire := s.App.AppCodec().MustMarshal(&proposal)
			// The old field 15 requested a minimum of 120 seconds. New code must
			// skip this payload rather than applying it or interpreting another action.
			wire = protowire.AppendTag(wire, 15, protowire.BytesType)
			wire = protowire.AppendBytes(wire, []byte{0x08, 0x78})
			store := s.Ctx.KVStore(s.App.GetKey(types.StoreKey))
			store.Set(types.ProposalKey(id), wire)
			s.Require().Equal(proposal, s.mustProposal(id))
			response, err := s.App.RegionGovKeeper.Proposal(sdk.WrapSDKContext(s.Ctx), &types.QueryProposalRequest{Id: id})
			s.Require().NoError(err)
			s.Require().Equal(proposal, response.Proposal)

			before := regiongov.ExportGenesis(s.Ctx, s.App.RegionGovKeeper)
			for _, voter := range []string{operator, s.Dao.GlobalDao} {
				for _, option := range []types.VoteOption{types.VoteOptionYes, types.VoteOptionNo} {
					s.Require().ErrorIs(s.vote(voter, id, option), types.ErrInvalidProposal)
					s.Require().Equal(before, regiongov.ExportGenesis(s.Ctx, s.App.RegionGovKeeper))
				}
			}
			s.Require().Equal(wire, store.Get(types.ProposalKey(id)), "rejected votes must not rewrite legacy storage")
			s.Require().NoError(before.Validate())
			// Export/import must retain historical ballots and restore the voting queue.
			s.SetupTest()
			regiongov.InitGenesis(s.Ctx, s.App.RegionGovKeeper, *before)
			s.Require().Equal(before, regiongov.ExportGenesis(s.Ctx, s.App.RegionGovKeeper))
			region, found := s.App.StakingKeeper.GetRegion(s.Ctx, "usa")
			s.Require().True(found)
			treasury := sdk.MustAccAddressFromBech32(region.RegionTreasureAddr)
			balance := s.App.BankKeeper.GetAllBalances(s.Ctx, treasury)
			s.endBlockAt(proposal.VotingEndTime.Add(-time.Nanosecond))
			s.Require().Equal(types.ProposalStatusVoting, s.mustProposal(id).Status)
			s.endBlockAt(proposal.VotingEndTime)
			s.Require().Equal(status, s.mustProposal(id).Status)
			s.Require().Equal(before.Params, s.App.RegionGovKeeper.GetParams(s.Ctx))
			s.Require().Equal(balance, s.App.BankKeeper.GetAllBalances(s.Ctx, treasury))
			s.Require().False(s.Ctx.KVStore(s.App.GetKey(types.StoreKey)).Has(types.VotingQueueKey(proposal.VotingEndTime, id)))
			s.Require().NoError(regiongov.ExportGenesis(s.Ctx, s.App.RegionGovKeeper).Validate())
			s.endBlockAt(proposal.VotingEndTime.Add(time.Second))
			s.Require().Equal(status, s.mustProposal(id).Status)
		})
	}
}
