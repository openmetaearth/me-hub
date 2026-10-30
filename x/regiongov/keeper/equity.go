package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/openmetaearth/me-hub/x/regiongov/types"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (k Keeper) submitEquityUpdate(ctx sdk.Context, msg *types.MsgSubmitEquityUpdate) (uint64, error) {
	if err := msg.ValidateBasic(); err != nil {
		return 0, err
	}
	region, proposer, err := k.requireRegionOperator(ctx, msg.RegionId, msg.Proposer)
	if err != nil {
		return 0, err
	}
	if k.isDao(ctx, proposer) {
		return 0, sdkerrors.Wrap(types.ErrUnauthorized, "dao cannot set equity addresses")
	}
	addresses, err := canonicalAddresses(msg.Addresses)
	if err != nil {
		return 0, err
	}
	id := k.queuePendingAction(ctx, types.PendingActionKindSetEquity, region.RegionId, proposer, addresses)
	k.emitEquity(ctx, types.EventTypeSubmitEquityUpdate, id, region.RegionId, proposer, "")
	return id, nil
}

func (k Keeper) submitFreezeEquity(ctx sdk.Context, msg *types.MsgSubmitFreezeEquity) (uint64, error) {
	if err := msg.ValidateBasic(); err != nil {
		return 0, err
	}
	return k.submitEquityFreezeChange(ctx, msg.RegionId, msg.Proposer, msg.Addresses, true)
}

func (k Keeper) submitUnfreezeEquity(ctx sdk.Context, msg *types.MsgSubmitUnfreezeEquity) (uint64, error) {
	if err := msg.ValidateBasic(); err != nil {
		return 0, err
	}
	return k.submitEquityFreezeChange(ctx, msg.RegionId, msg.Proposer, msg.Addresses, false)
}

func (k Keeper) submitEquityFreezeChange(ctx sdk.Context, regionID, proposer string, targets []string, frozen bool) (uint64, error) {
	region, proposer, addresses, err := k.prepareOperatorEquityChange(ctx, regionID, proposer, targets)
	if err != nil {
		return 0, err
	}
	for _, address := range addresses {
		if frozen && k.isEquityFrozen(ctx, region.RegionId, address) {
			return 0, sdkerrors.Wrapf(types.ErrEquityFrozen, "%s", address)
		}
		if !frozen && !k.isEquityFrozen(ctx, region.RegionId, address) {
			return 0, sdkerrors.Wrapf(types.ErrEquityNotFrozen, "%s", address)
		}
	}
	kind, eventType := types.PendingActionKindFreezeEquity, types.EventTypeSubmitFreezeEquity
	if !frozen {
		kind, eventType = types.PendingActionKindUnfreezeEquity, types.EventTypeSubmitUnfreezeEquity
	}
	id := k.queuePendingAction(ctx, kind, region.RegionId, proposer, addresses)
	for _, address := range addresses {
		k.emitEquity(ctx, eventType, id, region.RegionId, proposer, address)
	}
	return id, nil
}

func (k Keeper) approvePendingAction(ctx sdk.Context, msg *types.MsgApprovePendingAction) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	if err := k.requireDao(ctx, msg.Authority); err != nil {
		return err
	}
	action, err := k.openPendingAction(ctx, msg.ActionId)
	if err != nil {
		return err
	}
	if err := k.applyPendingAction(ctx, action); err != nil {
		return err
	}
	action.Status = types.PendingActionStatusApproved
	k.SetPendingAction(ctx, action)
	k.emitEquity(ctx, types.EventTypeApprovePendingAction, action.Id, action.RegionId, action.Proposer, "")
	return nil
}

func (k Keeper) rejectPendingAction(ctx sdk.Context, msg *types.MsgRejectPendingAction) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	if err := k.requireDao(ctx, msg.Authority); err != nil {
		return err
	}
	action, err := k.openPendingAction(ctx, msg.ActionId)
	if err != nil {
		return err
	}
	action.Status = types.PendingActionStatusRejected
	k.SetPendingAction(ctx, action)
	k.emitEquity(ctx, types.EventTypeRejectPendingAction, action.Id, action.RegionId, action.Proposer, "")
	return nil
}

func (k Keeper) freezeEquity(ctx sdk.Context, msg *types.MsgFreezeEquity) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	region, address, err := k.prepareDaoEquityChange(ctx, msg.Authority, msg.RegionId, msg.Address)
	if err != nil {
		return err
	}
	k.setEquityFrozen(ctx, region.RegionId, address, true)
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeFreezeEquity,
		sdk.NewAttribute(types.AttributeKeyRegionID, region.RegionId),
		sdk.NewAttribute(types.AttributeKeyAddress, address),
	))
	return nil
}

func (k Keeper) unfreezeEquity(ctx sdk.Context, msg *types.MsgUnfreezeEquity) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	region, address, err := k.prepareDaoEquityChange(ctx, msg.Authority, msg.RegionId, msg.Address)
	if err != nil {
		return err
	}
	k.setEquityFrozen(ctx, region.RegionId, address, false)
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeUnfreezeEquity,
		sdk.NewAttribute(types.AttributeKeyRegionID, region.RegionId),
		sdk.NewAttribute(types.AttributeKeyAddress, address),
	))
	return nil
}

func (k Keeper) prepareOperatorEquityChange(ctx sdk.Context, regionID, proposer string, addresses []string) (wstakingtypes.Region, string, []string, error) {
	region, canonicalProposer, err := k.requireRegionOperator(ctx, regionID, proposer)
	if err != nil {
		return wstakingtypes.Region{}, "", nil, err
	}
	if k.isDao(ctx, canonicalProposer) {
		return wstakingtypes.Region{}, "", nil, sdkerrors.Wrap(types.ErrUnauthorized, "dao freezes equity addresses directly")
	}
	canonicalAddresses := make([]string, 0, len(addresses))
	for _, address := range addresses {
		canonical, err := k.requireEquityAddress(ctx, region.RegionId, address)
		if err != nil {
			return wstakingtypes.Region{}, "", nil, err
		}
		canonicalAddresses = append(canonicalAddresses, canonical)
	}
	return region, canonicalProposer, canonicalAddresses, nil
}

func (k Keeper) prepareDaoEquityChange(ctx sdk.Context, authority, regionID, address string) (wstakingtypes.Region, string, error) {
	if err := k.requireDao(ctx, authority); err != nil {
		return wstakingtypes.Region{}, "", err
	}
	region, found := k.lookupRegion(ctx, regionID)
	if !found {
		return wstakingtypes.Region{}, "", sdkerrors.Wrapf(types.ErrUnknownRegion, "%s", regionID)
	}
	canonicalAddress, err := k.requireEquityAddress(ctx, region.RegionId, address)
	if err != nil {
		return wstakingtypes.Region{}, "", err
	}
	return region, canonicalAddress, nil
}

func (k Keeper) requireRegionOperator(ctx sdk.Context, regionID, proposer string) (wstakingtypes.Region, string, error) {
	region, found := k.lookupRegion(ctx, regionID)
	if !found {
		return wstakingtypes.Region{}, "", sdkerrors.Wrapf(types.ErrUnknownRegion, "%s", regionID)
	}
	if k.isRegionFrozen(ctx, region.RegionId) {
		return wstakingtypes.Region{}, "", sdkerrors.Wrapf(types.ErrRegionFrozen, "%s", region.RegionId)
	}
	canonical, err := canonicalAccAddress(proposer)
	if err != nil {
		return wstakingtypes.Region{}, "", err
	}
	operator, ok := operatorAccount(region.OperatorAddress)
	if !ok || operator.String() != canonical {
		return wstakingtypes.Region{}, "", sdkerrors.Wrap(types.ErrUnauthorized, "proposer is not the region operator")
	}
	return region, canonical, nil
}

func (k Keeper) requireDao(ctx sdk.Context, authority string) error {
	canonical, err := canonicalAccAddress(authority)
	if err != nil {
		return err
	}
	if !k.isDao(ctx, canonical) {
		return sdkerrors.Wrap(types.ErrUnauthorized, "authority is not the dao")
	}
	return nil
}

func (k Keeper) requireEquityAddress(ctx sdk.Context, regionID, address string) (string, error) {
	canonical, err := canonicalAccAddress(address)
	if err != nil {
		return "", err
	}
	equity, found := k.GetEquity(ctx, regionID)
	if !found || !containsAddress(equity.Addresses, canonical) {
		return "", sdkerrors.Wrapf(types.ErrNotEquityAddress, "%s", canonical)
	}
	return canonical, nil
}

func (k Keeper) isEquityFrozen(ctx sdk.Context, regionID, address string) bool {
	freeze, found := k.GetEquityFreeze(ctx, regionID, address)
	return found && freeze.Frozen
}

func (k Keeper) setEquityFrozen(ctx sdk.Context, regionID, address string, frozen bool) {
	k.SetEquityFreeze(ctx, types.EquityFreeze{
		RegionId: regionID,
		Address:  address,
		Frozen:   frozen,
	})
}

func (k Keeper) queuePendingAction(ctx sdk.Context, kind types.PendingActionKind, regionID, proposer string, addresses []string) uint64 {
	id := k.GetPendingActionSeq(ctx) + 1
	k.SetPendingActionSeq(ctx, id)
	k.SetPendingAction(ctx, types.PendingAction{
		Id:         id,
		Kind:       kind,
		Status:     types.PendingActionStatusPending,
		RegionId:   regionID,
		Proposer:   proposer,
		Addresses:  addresses,
		SubmitTime: ctx.BlockTime(),
	})
	return id
}

func (k Keeper) openPendingAction(ctx sdk.Context, id uint64) (types.PendingAction, error) {
	action, found := k.GetPendingAction(ctx, id)
	if !found {
		return types.PendingAction{}, sdkerrors.Wrapf(types.ErrPendingActionNotFound, "%d", id)
	}
	if action.Status != types.PendingActionStatusPending {
		return types.PendingAction{}, sdkerrors.Wrapf(types.ErrPendingActionClosed, "%d", id)
	}
	return action, nil
}

func (k Keeper) applyPendingAction(ctx sdk.Context, action types.PendingAction) error {
	switch action.Kind {
	case types.PendingActionKindSetEquity:
		k.SetEquity(ctx, types.RegionEquity{
			RegionId:  action.RegionId,
			Addresses: action.Addresses,
		})
	case types.PendingActionKindFreezeEquity, types.PendingActionKindUnfreezeEquity:
		msg := types.MsgSubmitFreezeEquity{Proposer: action.Proposer, RegionId: action.RegionId, Addresses: action.Addresses}
		if err := msg.ValidateBasic(); err != nil {
			return err
		}
		for _, address := range action.Addresses {
			k.setEquityFrozen(ctx, action.RegionId, address, action.Kind == types.PendingActionKindFreezeEquity)
		}
	default:
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "unknown pending action kind %s", action.Kind.String())
	}
	return nil
}

func (k Keeper) emitEquity(ctx sdk.Context, eventType string, actionID uint64, regionID, proposer, address string) {
	attrs := []sdk.Attribute{
		sdk.NewAttribute(types.AttributeKeyActionID, strconv.FormatUint(actionID, 10)),
		sdk.NewAttribute(types.AttributeKeyRegionID, regionID),
		sdk.NewAttribute(types.AttributeKeyProposer, proposer),
	}
	if address != "" {
		attrs = append(attrs, sdk.NewAttribute(types.AttributeKeyAddress, address))
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(eventType, attrs...))
}

func canonicalAddresses(addresses []string) ([]string, error) {
	out := make([]string, 0, len(addresses))
	for _, address := range addresses {
		canonical, err := canonicalAccAddress(address)
		if err != nil {
			return nil, err
		}
		out = append(out, canonical)
	}
	return out, nil
}

func containsAddress(addresses []string, address string) bool {
	for _, item := range addresses {
		if item == address {
			return true
		}
	}
	return false
}
