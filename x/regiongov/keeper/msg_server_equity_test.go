package keeper_test

import (
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

func (s *KeeperTestSuite) TestBatchEquityFreezeApproval() {
	s.SetupTest()
	operator := operatorAccount(s.usaValidator.OperatorAddress)
	dao := s.Dao.GlobalDao
	addresses := []string{sample.AccAddress(), sample.AccAddress()}
	s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: addresses})
	for _, freeze := range []bool{true, false} {
		submit := func(targets []string) (uint64, error) {
			if freeze {
				response, err := s.msgServer.SubmitFreezeEquity(s.Ctx, &types.MsgSubmitFreezeEquity{Proposer: operator, RegionId: "usa", Addresses: targets})
				if err != nil {
					return 0, err
				}
				return response.ActionId, nil
			}
			response, err := s.msgServer.SubmitUnfreezeEquity(s.Ctx, &types.MsgSubmitUnfreezeEquity{Proposer: operator, RegionId: "usa", Addresses: targets})
			if err != nil {
				return 0, err
			}
			return response.ActionId, nil
		}
		for _, invalid := range [][]string{nil, {addresses[0], addresses[0]}, {addresses[0], "invalid"}, {addresses[0], sample.AccAddress()}} {
			seq := s.App.RegionGovKeeper.GetPendingActionSeq(s.Ctx)
			events := len(s.Ctx.EventManager().Events())
			_, err := submit(invalid)
			s.Require().Error(err)
			s.Require().Equal(seq, s.App.RegionGovKeeper.GetPendingActionSeq(s.Ctx))
			s.Require().Len(s.Ctx.EventManager().Events(), events)
		}
		// A later address with an incompatible state must not queue a partial batch.
		s.App.RegionGovKeeper.SetEquityFreeze(s.Ctx, types.EquityFreeze{RegionId: "usa", Address: addresses[1], Frozen: freeze})
		seq := s.App.RegionGovKeeper.GetPendingActionSeq(s.Ctx)
		_, err := submit(addresses)
		s.Require().Error(err)
		s.Require().Equal(seq, s.App.RegionGovKeeper.GetPendingActionSeq(s.Ctx))
		s.App.RegionGovKeeper.SetEquityFreeze(s.Ctx, types.EquityFreeze{RegionId: "usa", Address: addresses[1], Frozen: !freeze})
		id, err := submit(addresses)
		s.Require().NoError(err)
		action, found := s.App.RegionGovKeeper.GetPendingAction(s.Ctx, id)
		s.Require().True(found)
		s.Require().Equal(addresses, action.Addresses)
		for _, address := range addresses {
			state, _ := s.App.RegionGovKeeper.GetEquityFreeze(s.Ctx, "usa", address)
			s.Require().Equal(!freeze, state.Frozen)
		}
		_, err = s.msgServer.ApprovePendingAction(s.Ctx, &types.MsgApprovePendingAction{Authority: dao, ActionId: id})
		s.Require().NoError(err)
		for _, address := range addresses {
			state, found := s.App.RegionGovKeeper.GetEquityFreeze(s.Ctx, "usa", address)
			s.Require().True(found)
			s.Require().Equal(freeze, state.Frozen)
		}
	}
}

func (s *KeeperTestSuite) TestEquityUpdateNeedsDaoApproval() {
	s.SetupTest()

	usa := operatorAccount(s.usaValidator.OperatorAddress)
	meEarth := operatorAccount(s.meEarthValidator.OperatorAddress)
	dao := s.Dao.GlobalDao
	first := sample.AccAddress()
	second := sample.AccAddress()
	third := sample.AccAddress()

	_, err := s.msgServer.SubmitEquityUpdate(s.Ctx, &types.MsgSubmitEquityUpdate{
		Proposer: dao, RegionId: "usa", Addresses: []string{first},
	})
	s.Require().ErrorIs(err, types.ErrUnauthorized)
	_, err = s.msgServer.SubmitEquityUpdate(s.Ctx, &types.MsgSubmitEquityUpdate{
		Proposer: meEarth, RegionId: "usa", Addresses: []string{first},
	})
	s.Require().ErrorIs(err, types.ErrUnauthorized)
	_, err = s.msgServer.SubmitEquityUpdate(s.Ctx, &types.MsgSubmitEquityUpdate{
		Proposer: usa, RegionId: "missing", Addresses: []string{first},
	})
	s.Require().ErrorIs(err, types.ErrUnknownRegion)

	pending, err := s.msgServer.SubmitEquityUpdate(s.Ctx, &types.MsgSubmitEquityUpdate{
		Proposer: usa, RegionId: "usa", Addresses: []string{first, second},
	})
	s.Require().NoError(err)
	_, found := s.App.RegionGovKeeper.GetEquity(s.Ctx, "usa")
	s.Require().False(found)

	_, err = s.msgServer.ApprovePendingAction(s.Ctx, &types.MsgApprovePendingAction{
		Authority: usa, ActionId: pending.ActionId,
	})
	s.Require().ErrorIs(err, types.ErrUnauthorized)
	_, err = s.msgServer.ApprovePendingAction(s.Ctx, &types.MsgApprovePendingAction{
		Authority: dao, ActionId: 99,
	})
	s.Require().ErrorIs(err, types.ErrPendingActionNotFound)

	_, err = s.msgServer.ApprovePendingAction(s.Ctx, &types.MsgApprovePendingAction{
		Authority: dao, ActionId: pending.ActionId,
	})
	s.Require().NoError(err)
	equity, found := s.App.RegionGovKeeper.GetEquity(s.Ctx, "usa")
	s.Require().True(found)
	s.Require().Equal([]string{first, second}, equity.Addresses)
	_, err = s.msgServer.ApprovePendingAction(s.Ctx, &types.MsgApprovePendingAction{
		Authority: dao, ActionId: pending.ActionId,
	})
	s.Require().ErrorIs(err, types.ErrPendingActionClosed)

	replacement, err := s.msgServer.SubmitEquityUpdate(s.Ctx, &types.MsgSubmitEquityUpdate{
		Proposer: usa, RegionId: "usa", Addresses: []string{third},
	})
	s.Require().NoError(err)
	_, err = s.msgServer.RejectPendingAction(s.Ctx, &types.MsgRejectPendingAction{
		Authority: usa, ActionId: replacement.ActionId,
	})
	s.Require().ErrorIs(err, types.ErrUnauthorized)
	_, err = s.msgServer.RejectPendingAction(s.Ctx, &types.MsgRejectPendingAction{
		Authority: dao, ActionId: replacement.ActionId,
	})
	s.Require().NoError(err)
	equity, found = s.App.RegionGovKeeper.GetEquity(s.Ctx, "usa")
	s.Require().True(found)
	s.Require().Equal([]string{first, second}, equity.Addresses)

	cleared, err := s.msgServer.SubmitEquityUpdate(s.Ctx, &types.MsgSubmitEquityUpdate{
		Proposer: usa, RegionId: "usa",
	})
	s.Require().NoError(err)
	_, err = s.msgServer.ApprovePendingAction(s.Ctx, &types.MsgApprovePendingAction{
		Authority: dao, ActionId: cleared.ActionId,
	})
	s.Require().NoError(err)
	equity, found = s.App.RegionGovKeeper.GetEquity(s.Ctx, "usa")
	s.Require().True(found)
	s.Require().Empty(equity.Addresses)
}

func (s *KeeperTestSuite) TestEquityFreezeByOperatorAndDao() {
	s.SetupTest()

	usa := operatorAccount(s.usaValidator.OperatorAddress)
	dao := s.Dao.GlobalDao
	operatorTarget := sample.AccAddress()
	daoTarget := sample.AccAddress()
	outsider := sample.AccAddress()

	listed, err := s.msgServer.SubmitEquityUpdate(s.Ctx, &types.MsgSubmitEquityUpdate{
		Proposer: usa, RegionId: "usa", Addresses: []string{operatorTarget, daoTarget},
	})
	s.Require().NoError(err)
	_, err = s.msgServer.ApprovePendingAction(s.Ctx, &types.MsgApprovePendingAction{
		Authority: dao, ActionId: listed.ActionId,
	})
	s.Require().NoError(err)

	_, err = s.msgServer.SubmitFreezeEquity(s.Ctx, &types.MsgSubmitFreezeEquity{
		Proposer: usa, RegionId: "usa", Addresses: []string{outsider},
	})
	s.Require().ErrorIs(err, types.ErrNotEquityAddress)
	_, err = s.msgServer.FreezeEquity(s.Ctx, &types.MsgFreezeEquity{
		Authority: usa, RegionId: "usa", Address: daoTarget,
	})
	s.Require().ErrorIs(err, types.ErrUnauthorized)
	_, err = s.msgServer.SubmitUnfreezeEquity(s.Ctx, &types.MsgSubmitUnfreezeEquity{
		Proposer: usa, RegionId: "usa", Addresses: []string{operatorTarget},
	})
	s.Require().ErrorIs(err, types.ErrEquityNotFrozen)

	pending, err := s.msgServer.SubmitFreezeEquity(s.Ctx, &types.MsgSubmitFreezeEquity{
		Proposer: usa, RegionId: "usa", Addresses: []string{operatorTarget},
	})
	s.Require().NoError(err)
	notFrozen, found := s.App.RegionGovKeeper.GetEquityFreeze(s.Ctx, "usa", operatorTarget)
	s.Require().False(found)
	s.Require().False(notFrozen.Frozen)

	_, err = s.msgServer.FreezeEquity(s.Ctx, &types.MsgFreezeEquity{
		Authority: dao, RegionId: "usa", Address: daoTarget,
	})
	s.Require().NoError(err)
	frozen, found := s.App.RegionGovKeeper.GetEquityFreeze(s.Ctx, "usa", daoTarget)
	s.Require().True(found)
	s.Require().True(frozen.Frozen)

	_, err = s.msgServer.ApprovePendingAction(s.Ctx, &types.MsgApprovePendingAction{
		Authority: dao, ActionId: pending.ActionId,
	})
	s.Require().NoError(err)
	frozen, found = s.App.RegionGovKeeper.GetEquityFreeze(s.Ctx, "usa", operatorTarget)
	s.Require().True(found)
	s.Require().True(frozen.Frozen)
	_, err = s.msgServer.SubmitFreezeEquity(s.Ctx, &types.MsgSubmitFreezeEquity{
		Proposer: usa, RegionId: "usa", Addresses: []string{operatorTarget},
	})
	s.Require().ErrorIs(err, types.ErrEquityFrozen)

	unfreeze, err := s.msgServer.SubmitUnfreezeEquity(s.Ctx, &types.MsgSubmitUnfreezeEquity{
		Proposer: usa, RegionId: "usa", Addresses: []string{operatorTarget},
	})
	s.Require().NoError(err)
	_, err = s.msgServer.ApprovePendingAction(s.Ctx, &types.MsgApprovePendingAction{
		Authority: dao, ActionId: unfreeze.ActionId,
	})
	s.Require().NoError(err)
	thawed, found := s.App.RegionGovKeeper.GetEquityFreeze(s.Ctx, "usa", operatorTarget)
	s.Require().True(found)
	s.Require().False(thawed.Frozen)

	_, err = s.msgServer.UnfreezeEquity(s.Ctx, &types.MsgUnfreezeEquity{
		Authority: dao, RegionId: "usa", Address: daoTarget,
	})
	s.Require().NoError(err)
	thawed, found = s.App.RegionGovKeeper.GetEquityFreeze(s.Ctx, "usa", daoTarget)
	s.Require().True(found)
	s.Require().False(thawed.Frozen)

	s.App.RegionGovKeeper.SetRegionFreeze(s.Ctx, types.RegionFreeze{RegionId: "usa", Frozen: true})
	_, err = s.msgServer.SubmitEquityUpdate(s.Ctx, &types.MsgSubmitEquityUpdate{
		Proposer: usa, RegionId: "usa", Addresses: []string{operatorTarget},
	})
	s.Require().ErrorIs(err, types.ErrRegionFrozen)
	_, err = s.msgServer.SubmitFreezeEquity(s.Ctx, &types.MsgSubmitFreezeEquity{
		Proposer: usa, RegionId: "usa", Addresses: []string{daoTarget},
	})
	s.Require().ErrorIs(err, types.ErrRegionFrozen)
}
