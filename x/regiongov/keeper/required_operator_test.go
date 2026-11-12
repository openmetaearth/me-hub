package keeper_test

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/regiongov"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

func (s *KeeperTestSuite) TestMigratedRequiredOperatorQueryAndVoting() {
	k := s.App.RegionGovKeeper
	operator := operatorAccount(s.usaValidator.OperatorAddress)
	equity := sample.AccAddress()
	k.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: []string{equity}})
	id := s.submit(equity, s.withdrawal("usa"))
	proposal := s.mustProposal(id)
	s.Require().Equal(operator, proposal.RequiredOperator)
	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, "usa")
	s.Require().True(found)
	s.Require().Equal(s.usaValidator.OperatorAddress, region.OperatorAddress)

	// Bypass normalized writes to reproduce state persisted by v2.0.18.
	proposal.RequiredOperator = s.usaValidator.OperatorAddress
	store := s.Ctx.KVStore(s.App.GetKey(types.StoreKey))
	legacy := s.App.AppCodec().MustMarshal(&proposal)
	store.Set(types.ProposalKey(id), legacy)
	// The upgrade normalizes storage before normal query and voting traffic.
	s.Require().NoError(k.MigrateRequiredOperatorAccounts(s.Ctx))
	migrated := append([]byte(nil), store.Get(types.ProposalKey(id))...)
	s.Require().NotEqual(legacy, migrated)
	response, err := k.Proposal(sdk.WrapSDKContext(s.Ctx), &types.QueryProposalRequest{Id: id})
	s.Require().NoError(err)
	s.Require().Equal(operator, response.Proposal.RequiredOperator)
	list, err := k.Proposals(sdk.WrapSDKContext(s.Ctx), &types.QueryProposalsRequest{})
	s.Require().NoError(err)
	s.Require().Len(list.Proposals, 1)
	s.Require().Equal(operator, list.Proposals[0].RequiredOperator)
	s.Require().Equal(s.mustProposal(id), response.Proposal, "query returns the stored proposal")
	s.Require().Equal(response.Proposal, list.Proposals[0])
	s.Require().Equal(migrated, store.Get(types.ProposalKey(id)), "queries must not rewrite state")

	// A replacement must not inherit approval authority on the legacy proposal.
	newOperator := sdk.MustAccAddressFromBech32(sample.AccAddress())
	s.Require().NoError(s.App.StakingKeeper.SetRegionOperator(s.Ctx, "usa", sdk.ValAddress(newOperator).String()))
	s.Require().ErrorIs(s.vote(newOperator.String(), id, types.VoteOptionYes), types.ErrUnauthorized)
	s.Require().ErrorIs(s.vote(sample.AccAddress(), id, types.VoteOptionYes), types.ErrUnauthorized)
	s.Require().ErrorIs(s.vote(proposal.RequiredOperator, id, types.VoteOptionYes), sdkerrors.ErrInvalidAddress)
	s.Require().NoError(s.vote(operator, id, types.VoteOptionYes))
	s.Require().ErrorIs(s.vote(operator, id, types.VoteOptionYes), types.ErrAlreadyVoted)
	s.Require().NoError(s.vote(equity, id, types.VoteOptionYes))
	s.Require().NoError(s.vote(s.Dao.GlobalDao, id, types.VoteOptionYes))
	s.endBlockAt(proposal.VotingEndTime)
	s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(id).Status)
}

func (s *KeeperTestSuite) TestRequiredOperatorMigrationAndLegacyGenesis() {
	k := s.App.RegionGovKeeper
	operator := operatorAccount(s.usaValidator.OperatorAddress)
	equity := sample.AccAddress()
	k.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: []string{equity}})
	id := s.submit(equity, s.withdrawal("usa"))
	s.Require().NoError(s.vote(operator, id, types.VoteOptionYes))
	proposal := s.mustProposal(id)
	proposal.RequiredOperator = s.usaValidator.OperatorAddress
	store := s.Ctx.KVStore(s.App.GetKey(types.StoreKey))
	store.Set(types.ProposalKey(id), s.App.AppCodec().MustMarshal(&proposal))
	before := regiongov.ExportGenesis(s.Ctx, k)
	s.Require().NoError(before.Validate())
	// Deliberately change the live operator: migration must use the snapshot.
	newOperator := sdk.MustAccAddressFromBech32(sample.AccAddress())
	s.Require().NoError(s.App.StakingKeeper.SetRegionOperator(s.Ctx, "usa", sdk.ValAddress(newOperator).String()))
	regionsBefore := s.App.StakingKeeper.GetAllRegion(s.Ctx)
	s.Require().NoError(k.MigrateRequiredOperatorAccounts(s.Ctx))
	before.Proposals[0].RequiredOperator = operator
	s.Require().Equal(before, regiongov.ExportGenesis(s.Ctx, k))
	s.Require().Equal(regionsBefore, s.App.StakingKeeper.GetAllRegion(s.Ctx))
	s.Require().NoError(k.MigrateRequiredOperatorAccounts(s.Ctx))
	s.Require().Equal(before, regiongov.ExportGenesis(s.Ctx, k), "migration is idempotent")
	s.Require().ErrorIs(s.vote(operator, id, types.VoteOptionYes), types.ErrAlreadyVoted)
	s.Require().ErrorIs(s.vote(newOperator.String(), id, types.VoteOptionYes), types.ErrUnauthorized)

	// Import accepts old encodings but persists only canonical account addresses.
	for _, address := range []string{proposal.RequiredOperator, strings.ToUpper(operator)} {
		before.Proposals[0].RequiredOperator = address
		s.Require().NoError(before.Validate())
		regiongov.InitGenesis(s.Ctx, k, *before)
		s.Require().Equal(operator, s.mustProposal(id).RequiredOperator)
		s.Require().NoError(regiongov.ExportGenesis(s.Ctx, k).Validate())
	}
}

func (s *KeeperTestSuite) TestRequiredOperatorMigrationRejectsInvalidWithoutPartialWrites() {
	k := s.App.RegionGovKeeper
	equity := sample.AccAddress()
	k.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: []string{equity}})
	first := s.submit(equity, s.withdrawal("usa"))
	second := s.submit(equity, s.withdrawal("usa"))
	store := s.Ctx.KVStore(s.App.GetKey(types.StoreKey))
	for _, id := range []uint64{first, second} {
		proposal := s.mustProposal(id)
		proposal.RequiredOperator = s.usaValidator.OperatorAddress
		if id == second {
			proposal.RequiredOperator = "invalid"
		}
		store.Set(types.ProposalKey(id), s.App.AppCodec().MustMarshal(&proposal))
	}
	before := regiongov.ExportGenesis(s.Ctx, k)
	s.Require().ErrorContains(k.MigrateRequiredOperatorAccounts(s.Ctx), "required operator")
	s.Require().Equal(before, regiongov.ExportGenesis(s.Ctx, k))
}
