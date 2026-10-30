package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestBatchRegionFreeze() {
	s.SetupTest()
	dao := s.Dao.GlobalDao
	ids := []string{"usa", wstakingtypes.MeEarthRegionId}
	for _, frozen := range []bool{true, false} {
		apply := func(authority string, regions []string) error {
			if frozen {
				_, err := s.msgServer.FreezeRegion(s.Ctx, &types.MsgFreezeRegion{Authority: authority, RegionIds: regions})
				return err
			}
			_, err := s.msgServer.UnfreezeRegion(s.Ctx, &types.MsgUnfreezeRegion{Authority: authority, RegionIds: regions})
			return err
		}
		assertState := func(expected bool) {
			for _, id := range ids {
				state, _ := s.App.RegionGovKeeper.GetRegionFreeze(s.Ctx, id)
				s.Require().Equal(expected, state.Frozen)
			}
		}
		for _, invalid := range [][]string{nil, {""}, {"usa", "USA"}, {"usa", "missing"}} {
			before := len(s.Ctx.EventManager().Events())
			s.Require().Error(apply(dao, invalid))
			assertState(!frozen)
			s.Require().Len(s.Ctx.EventManager().Events(), before)
		}
		s.Require().ErrorIs(apply(operatorAccount(s.usaValidator.OperatorAddress), ids), types.ErrUnauthorized)
		assertState(!frozen)
		before := len(s.Ctx.EventManager().Events())
		s.Require().NoError(apply(dao, []string{"USA", wstakingtypes.MeEarthRegionId}))
		assertState(frozen)
		events := s.Ctx.EventManager().Events()[before:]
		s.Require().Len(events, len(ids))
		eventType := types.EventTypeFreezeRegion
		if !frozen {
			eventType = types.EventTypeUnfreezeRegion
		}
		for i, event := range events {
			s.Require().Equal(eventType, event.Type)
			s.Require().Equal(types.AttributeKeyRegionID, event.Attributes[0].Key)
			s.Require().Equal(ids[i], event.Attributes[0].Value)
		}
	}
}

func (s *KeeperTestSuite) TestDaoFreezesAndUnfreezesRegion() {
	s.SetupTest()

	usa := operatorAccount(s.usaValidator.OperatorAddress)
	meEarth := operatorAccount(s.meEarthValidator.OperatorAddress)
	dao := s.Dao.GlobalDao

	open := s.submit(usa, 60, s.withdrawal(wstakingtypes.MeEarthRegionId))
	s.Require().ElementsMatch([]string{usa}, s.electorAccounts(open))

	_, err := s.msgServer.FreezeRegion(s.Ctx, &types.MsgFreezeRegion{
		Authority: usa, RegionIds: []string{"usa"},
	})
	s.Require().ErrorIs(err, types.ErrUnauthorized)
	_, err = s.msgServer.FreezeRegion(s.Ctx, &types.MsgFreezeRegion{
		Authority: meEarth, RegionIds: []string{wstakingtypes.MeEarthRegionId},
	})
	s.Require().ErrorIs(err, types.ErrUnauthorized)
	_, err = s.msgServer.FreezeRegion(s.Ctx, &types.MsgFreezeRegion{
		Authority: dao, RegionIds: []string{"missing"},
	})
	s.Require().ErrorIs(err, types.ErrUnknownRegion)

	_, err = s.msgServer.FreezeRegion(s.Ctx, &types.MsgFreezeRegion{
		Authority: dao, RegionIds: []string{"usa"},
	})
	s.Require().NoError(err)
	frozen, found := s.App.RegionGovKeeper.GetRegionFreeze(s.Ctx, "usa")
	s.Require().True(found)
	s.Require().True(frozen.Frozen)

	s.Require().ErrorIs(s.vote(usa, open, types.VoteOptionYes), types.ErrRegionFrozen)
	_, err = s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer:     usa,
		Kind:         types.ProposalKindTreasuryWithdraw,
		Title:        "frozen operator",
		VotingPeriod: 60,
		Treasury:     &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{s.withdrawal("usa")}},
	})
	s.Require().ErrorIs(err, types.ErrRegionFrozen)
	_, err = s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer:     usa,
		Kind:         types.ProposalKindReplaceRegionOperator,
		Title:        "frozen replace",
		VotingPeriod: 60,
		ReplaceOperator: &types.ReplaceRegionOperatorPayload{
			RegionId:    wstakingtypes.MeEarthRegionId,
			NewOperator: sdk.ValAddress(sdk.MustAccAddressFromBech32(sample.AccAddress())).String(),
		},
	})
	s.Require().ErrorIs(err, types.ErrRegionFrozen)

	_, err = s.msgServer.UnfreezeRegion(s.Ctx, &types.MsgUnfreezeRegion{
		Authority: usa, RegionIds: []string{"usa"},
	})
	s.Require().ErrorIs(err, types.ErrUnauthorized)
	_, err = s.msgServer.UnfreezeRegion(s.Ctx, &types.MsgUnfreezeRegion{
		Authority: dao, RegionIds: []string{"usa"},
	})
	s.Require().NoError(err)
	frozen, found = s.App.RegionGovKeeper.GetRegionFreeze(s.Ctx, "usa")
	s.Require().True(found)
	s.Require().False(frozen.Frozen)

	s.Require().NoError(s.vote(usa, open, types.VoteOptionYes))
	restored := s.submit(usa, 60, s.withdrawal("usa"))
	s.Require().ElementsMatch([]string{usa}, s.electorAccounts(restored))
}
