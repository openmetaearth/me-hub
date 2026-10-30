package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

func (k Keeper) freezeRegion(ctx sdk.Context, msg *types.MsgFreezeRegion) error {
	return k.setRegionFrozen(ctx, msg, true, types.EventTypeFreezeRegion)
}

func (k Keeper) unfreezeRegion(ctx sdk.Context, msg *types.MsgUnfreezeRegion) error {
	return k.setRegionFrozen(ctx, &types.MsgFreezeRegion{
		Authority: msg.Authority,
		RegionIds: msg.RegionIds,
	}, false, types.EventTypeUnfreezeRegion)
}

func (k Keeper) setRegionFrozen(ctx sdk.Context, msg *types.MsgFreezeRegion, frozen bool, eventType string) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	if err := k.requireDao(ctx, msg.Authority); err != nil {
		return err
	}
	// Resolve the entire batch before writing state or emitting events.
	regionIDs := make([]string, 0, len(msg.RegionIds))
	for _, id := range msg.RegionIds {
		region, found := k.lookupRegion(ctx, id)
		if !found {
			return sdkerrors.Wrapf(types.ErrUnknownRegion, "%s", id)
		}
		regionIDs = append(regionIDs, region.RegionId)
	}
	for _, id := range regionIDs {
		k.SetRegionFreeze(ctx, types.RegionFreeze{
			RegionId: id,
			Frozen:   frozen,
		})
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			eventType,
			sdk.NewAttribute(types.AttributeKeyRegionID, id),
		))
	}
	return nil
}
