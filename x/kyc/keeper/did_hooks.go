package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/kyc/types"
)

func (k Keeper) AfterDidStatusUpdated(ctx sdk.Context, info didtypes.DidInfo) error {
	if info.Status != didtypes.DID_STATUS_INACTIVE {
		return nil
	}
	return k.revokeKycForDid(ctx, info)
}

func (k Keeper) revokeKycForDid(ctx sdk.Context, didInfo didtypes.DidInfo) error {
	kyc, found := k.GetKYC(ctx, didInfo.Did)
	regionID := didInfo.RegionId
	if found && len(kyc.Data) > 0 {
		regionID = string(kyc.Data)
	}

	didInfo.RegionId = ""
	didInfo.KycLevel = didtypes.KYC_LEVEL_NONE
	didInfo.Status = didtypes.DID_STATUS_INACTIVE
	k.SetDidInfo(ctx, didInfo.Did, didInfo)

	if !found {
		return nil
	}

	k.DeleteKYC(ctx, didInfo.Did)

	filters, _ := k.GetFilters(ctx, didInfo.Did)
	if len(filters) > 0 {
		k.DeleteFilters(ctx, didInfo.Did, filters)
	}

	if didInfo.Address != "" && regionID != "" {
		if err := k.DeleteApproveReward(ctx, didInfo.Address, regionID); err != nil {
			return errors.Wrap(err, "delete reward failed")
		}
	}

	event := sdk.NewEvent(
		types.EventTypeRemove,
		sdk.NewAttribute(types.AttributeKeyAddress, didInfo.Address),
		sdk.NewAttribute(types.AttributeKeyRegionId, regionID),
	)
	ctx.EventManager().EmitEvent(event)

	return k.handlerReg.HandleEvent(sdk.WrapSDKContext(ctx), types.EventTypeRemove, event)
}

var _ interface {
	AfterDidStatusUpdated(sdk.Context, didtypes.DidInfo) error
} = Keeper{}
