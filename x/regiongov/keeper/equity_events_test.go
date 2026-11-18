package keeper_test

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

func (s *KeeperTestSuite) TestPendingActionEventsIncludeOriginalKind() {
	for _, kind := range []types.PendingActionKind{
		types.PendingActionKindSetEquity,
		types.PendingActionKindFreezeEquity,
		types.PendingActionKindUnfreezeEquity,
	} {
		for _, eventType := range []string{types.EventTypeApprovePendingAction, types.EventTypeRejectPendingAction} {
			s.Run(kind.String()+"/"+eventType, func() {
				s.SetupTest()
				operator := operatorAccount(s.usaValidator.OperatorAddress)
				addresses := []string{sample.AccAddress(), sample.AccAddress()}
				s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: addresses})
				s.Ctx = s.Ctx.WithEventManager(sdk.NewEventManager())
				id, submitEvent := s.submitEquityActionForEventTest(kind, operator, addresses...)

				// Every target has its own event; address and action kind are distinct.
				attrs := []sdk.Attribute{
					sdk.NewAttribute(types.AttributeKeyActionID, strconv.FormatUint(id, 10)),
					sdk.NewAttribute(types.AttributeKeyRegionID, "usa"),
					sdk.NewAttribute(types.AttributeKeyProposer, operator),
				}
				var submitEvents sdk.Events
				for _, address := range addresses {
					submitEvents = append(submitEvents, sdk.NewEvent(submitEvent, attrs...).AppendAttributes(
						sdk.NewAttribute(types.AttributeKeyAddress, address)))
				}
				s.Require().Equal(submitEvents, s.Ctx.EventManager().Events())
				s.Ctx = s.Ctx.WithEventManager(sdk.NewEventManager())

				decide := func(authority string) error {
					if eventType == types.EventTypeApprovePendingAction {
						_, err := s.msgServer.ApprovePendingAction(s.Ctx, &types.MsgApprovePendingAction{Authority: authority, ActionId: id})
						return err
					}
					_, err := s.msgServer.RejectPendingAction(s.Ctx, &types.MsgRejectPendingAction{Authority: authority, ActionId: id})
					return err
				}
				s.Require().ErrorIs(decide(operator), types.ErrUnauthorized)
				s.Require().Empty(s.Ctx.EventManager().Events())
				s.Require().NoError(decide(s.Dao.GlobalDao))
				attrs = append(attrs, sdk.NewAttribute(types.AttributeKeyActionKind, kind.String()))
				var decisionEvents sdk.Events
				for _, address := range addresses {
					decisionEvents = append(decisionEvents, sdk.NewEvent(eventType, attrs...).AppendAttributes(
						sdk.NewAttribute(types.AttributeKeyAddress, address)))
				}
				s.Require().Equal(decisionEvents, s.Ctx.EventManager().Events())

				// Repeating an approval/rejection must not emit another event.
				s.Ctx = s.Ctx.WithEventManager(sdk.NewEventManager())
				s.Require().ErrorIs(decide(s.Dao.GlobalDao), types.ErrPendingActionClosed)
				s.Require().Empty(s.Ctx.EventManager().Events())
			})
		}
	}
}

func (s *KeeperTestSuite) TestClearEquityEmitsEventsWithoutAddress() {
	for _, approve := range []bool{true, false} {
		s.SetupTest()
		operator := operatorAccount(s.usaValidator.OperatorAddress)
		original := []string{sample.AccAddress()}
		s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: original})
		s.Ctx = s.Ctx.WithEventManager(sdk.NewEventManager())
		id, _ := s.submitEquityActionForEventTest(types.PendingActionKindSetEquity, operator)
		attrs := []sdk.Attribute{
			sdk.NewAttribute(types.AttributeKeyActionID, strconv.FormatUint(id, 10)),
			sdk.NewAttribute(types.AttributeKeyRegionID, "usa"),
			sdk.NewAttribute(types.AttributeKeyProposer, operator),
		}
		s.Require().Equal(sdk.Events{sdk.NewEvent(types.EventTypeSubmitEquityUpdate, attrs...)}, s.Ctx.EventManager().Events())
		s.Ctx = s.Ctx.WithEventManager(sdk.NewEventManager())
		eventType := types.EventTypeRejectPendingAction
		if approve {
			_, err := s.msgServer.ApprovePendingAction(s.Ctx, &types.MsgApprovePendingAction{Authority: s.Dao.GlobalDao, ActionId: id})
			s.Require().NoError(err)
			eventType = types.EventTypeApprovePendingAction
		} else {
			_, err := s.msgServer.RejectPendingAction(s.Ctx, &types.MsgRejectPendingAction{Authority: s.Dao.GlobalDao, ActionId: id})
			s.Require().NoError(err)
		}
		attrs = append(attrs, sdk.NewAttribute(types.AttributeKeyActionKind, types.PendingActionKindSetEquity.String()))
		s.Require().Equal(sdk.Events{sdk.NewEvent(eventType, attrs...)}, s.Ctx.EventManager().Events())
		equity, found := s.App.RegionGovKeeper.GetEquity(s.Ctx, "usa")
		s.Require().True(found)
		if approve {
			s.Require().Empty(equity.Addresses)
		} else {
			s.Require().Equal(original, equity.Addresses)
		}
	}
}

func (s *KeeperTestSuite) submitEquityActionForEventTest(kind types.PendingActionKind, operator string, addresses ...string) (uint64, string) {
	s.T().Helper()
	switch kind {
	case types.PendingActionKindSetEquity:
		response, err := s.msgServer.SubmitEquityUpdate(s.Ctx, &types.MsgSubmitEquityUpdate{
			Proposer: operator, RegionId: "usa", Addresses: addresses,
		})
		s.Require().NoError(err)
		return response.ActionId, types.EventTypeSubmitEquityUpdate
	case types.PendingActionKindFreezeEquity:
		response, err := s.msgServer.SubmitFreezeEquity(s.Ctx, &types.MsgSubmitFreezeEquity{
			Proposer: operator, RegionId: "usa", Addresses: addresses,
		})
		s.Require().NoError(err)
		return response.ActionId, types.EventTypeSubmitFreezeEquity
	case types.PendingActionKindUnfreezeEquity:
		for _, address := range addresses {
			s.App.RegionGovKeeper.SetEquityFreeze(s.Ctx, types.EquityFreeze{RegionId: "usa", Address: address, Frozen: true})
		}
		response, err := s.msgServer.SubmitUnfreezeEquity(s.Ctx, &types.MsgSubmitUnfreezeEquity{
			Proposer: operator, RegionId: "usa", Addresses: addresses,
		})
		s.Require().NoError(err)
		return response.ActionId, types.EventTypeSubmitUnfreezeEquity
	default:
		s.T().Fatalf("unexpected action kind %s", kind)
		return 0, ""
	}
}
