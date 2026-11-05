package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/regiongov"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

func (s *KeeperTestSuite) TestPendingActionPruningBoundaryAndStatus() {
	operator := operatorAccount(s.usaValidator.OperatorAddress)
	statuses := []types.PendingActionStatus{
		types.PendingActionStatusApproved,
		types.PendingActionStatusRejected,
		types.PendingActionStatusPending,
	}
	for i, state := range statuses {
		s.App.RegionGovKeeper.SetPendingAction(s.Ctx, types.PendingAction{
			Id: uint64(i + 1), Kind: types.PendingActionKindSetEquity, Status: state,
			RegionId: "usa", Proposer: operator, SubmitTime: s.startTime,
		})
	}
	s.App.RegionGovKeeper.SetPendingActionSeq(s.Ctx, 3)
	// The current retention policy is one day; test both sides of its strict boundary.
	boundary := s.startTime.Add(24 * time.Hour)
	for _, now := range []time.Time{boundary.Add(-time.Nanosecond), boundary} {
		s.endBlockAt(now)
		s.Require().Len(s.App.RegionGovKeeper.GetAllPendingActions(s.Ctx), 3)
	}
	s.Ctx = s.Ctx.WithEventManager(sdk.NewEventManager())
	s.endBlockAt(boundary.Add(time.Nanosecond))
	remaining := s.App.RegionGovKeeper.GetAllPendingActions(s.Ctx)
	s.Require().Len(remaining, 1)
	s.Require().Equal(uint64(3), remaining[0].Id)
	s.Require().Equal(types.PendingActionStatusPending, remaining[0].Status)
	s.Require().Equal(uint64(3), s.App.RegionGovKeeper.GetPendingActionSeq(s.Ctx))
	s.Require().Empty(s.Ctx.EventManager().Events())

	for _, id := range []uint64{1, 2} {
		_, found := s.App.RegionGovKeeper.GetPendingAction(s.Ctx, id)
		s.Require().False(found)
		_, err := s.App.RegionGovKeeper.PendingAction(s.Ctx, &types.QueryPendingActionRequest{Id: id})
		s.Require().Equal(codes.NotFound, status.Code(err))
	}
	exported := regiongov.ExportGenesis(s.Ctx, s.App.RegionGovKeeper)
	s.Require().NoError(exported.Validate())
	s.Require().Equal(remaining, exported.PendingActions)
	s.Require().Equal(uint64(3), exported.PendingActionSeq)

	// Repeated end blocks are harmless; even very old pending actions survive.
	s.endBlockAt(boundary.Add(30 * 24 * time.Hour))
	s.Require().Equal(remaining, s.App.RegionGovKeeper.GetAllPendingActions(s.Ctx))
	response, err := s.msgServer.SubmitEquityUpdate(s.Ctx, &types.MsgSubmitEquityUpdate{
		Proposer: operator, RegionId: "usa", Addresses: []string{sample.AccAddress()},
	})
	s.Require().NoError(err)
	s.Require().Equal(uint64(4), response.ActionId)
}

func (s *KeeperTestSuite) TestOldPendingActionPrunedAfterDaoDecision() {
	for _, kind := range []types.PendingActionKind{
		types.PendingActionKindSetEquity,
		types.PendingActionKindFreezeEquity,
		types.PendingActionKindUnfreezeEquity,
	} {
		for _, approve := range []bool{true, false} {
			s.SetupTest()
			operator := operatorAccount(s.usaValidator.OperatorAddress)
			addresses := []string{sample.AccAddress(), sample.AccAddress()}
			s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: addresses})
			id, _ := s.submitEquityActionForEventTest(kind, operator, addresses...)
			s.endBlockAt(s.startTime.Add(2 * 24 * time.Hour))
			_, found := s.App.RegionGovKeeper.GetPendingAction(s.Ctx, id)
			s.Require().True(found)

			if approve {
				_, err := s.msgServer.ApprovePendingAction(s.Ctx, &types.MsgApprovePendingAction{Authority: s.Dao.GlobalDao, ActionId: id})
				s.Require().NoError(err)
			} else {
				_, err := s.msgServer.RejectPendingAction(s.Ctx, &types.MsgRejectPendingAction{Authority: s.Dao.GlobalDao, ActionId: id})
				s.Require().NoError(err)
			}
			action, found := s.App.RegionGovKeeper.GetPendingAction(s.Ctx, id)
			s.Require().True(found) // Retained until EndBlocker, not deleted by the Msg.
			s.Require().Equal(s.startTime, action.SubmitTime)
			events := append(sdk.Events(nil), s.Ctx.EventManager().Events()...)
			equity, _ := s.App.RegionGovKeeper.GetEquity(s.Ctx, "usa")
			freeze, freezeFound := s.App.RegionGovKeeper.GetEquityFreeze(s.Ctx, "usa", addresses[0])

			// Approval time does not start a new retention window.
			s.endBlockAt(s.Ctx.BlockTime())
			_, found = s.App.RegionGovKeeper.GetPendingAction(s.Ctx, id)
			s.Require().False(found)
			afterEquity, _ := s.App.RegionGovKeeper.GetEquity(s.Ctx, "usa")
			afterFreeze, afterFreezeFound := s.App.RegionGovKeeper.GetEquityFreeze(s.Ctx, "usa", addresses[0])
			s.Require().Equal(equity, afterEquity)
			s.Require().Equal(freeze, afterFreeze)
			s.Require().Equal(freezeFound, afterFreezeFound)
			s.Require().Equal(events, s.Ctx.EventManager().Events())
			s.Require().Equal(id, s.App.RegionGovKeeper.GetPendingActionSeq(s.Ctx))
		}
	}
}
