package keeper

import (
	"encoding/binary"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (k Keeper) SendInviteReward(ctx sdk.Context, inviter, invitee, regionId string) error {
	if inviter == "" {
		return nil
	}

	region, found := k.GetRegion(ctx, regionId)
	if !found {
		return types.ErrRegionNotExist
	}
	hasInviterReward := k.HasInviterReward(ctx, invitee)
	if hasInviterReward {
		return nil
	}
	if err := k.bankKeeper.Extend().SendCoinsWithTag(ctx,
		sdk.MustAccAddressFromBech32(region.RegionTreasureAddr),
		sdk.MustAccAddressFromBech32(inviter),
		sdk.NewCoins(sdk.NewCoin(params.BaseDenom, types.InviteReward)),
		fmt.Sprintf("SendInviteReward_%s", region.RegionId),
	); err != nil {
		return fmt.Errorf("send kyc reward to inviter, %v", err)
	}
	k.SetInviterReward(ctx, invitee)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventInviteReward,
			sdk.NewAttribute(types.AttributeKeyKycInviterRewardSender, region.RegionTreasureAddr),
			sdk.NewAttribute(types.AttributeKeyKycInviter, inviter),
			sdk.NewAttribute(types.AttributeKeyKycInvitee, invitee),
			sdk.NewAttribute(types.AttributeKeyKycInviterReward, types.InviteReward.String()),
		),
	)
	return nil
}

func (k Keeper) SetInviterReward(ctx sdk.Context, address string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.InviteKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, 1)
	store.Set([]byte(address), bz)
}

func (k Keeper) HasInviterReward(ctx sdk.Context, address string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.InviteKey)
	bz := store.Get([]byte(address))
	if bz == nil {
		return false
	}
	return true
}
