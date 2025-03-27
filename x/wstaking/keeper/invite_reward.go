package keeper

import (
	"encoding/binary"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

func (k Keeper) SendInviteReward(ctx sdk.Context, inviter, invitee, regionId string) error {
	if len(inviter) > 0 {
		region, found := k.GetRegion(ctx, regionId)
		if !found {
			return types.ErrRegionNotExist
		}
		hasInviterReward := k.HasInviterReward(ctx, invitee)
		if hasInviterReward {
			return nil
		}
		err := k.bankKeeper.Extend().SendCoinsWithTag(ctx,
			sdk.MustAccAddressFromBech32(region.RegionTreasureAddr),
			sdk.MustAccAddressFromBech32(inviter),
			sdk.NewCoins(sdk.NewCoin(params.BaseDenom, types.InviteReward)),
			fmt.Sprintf("SendInviteReward_InviteReward_%s", region.RegionId),
		)
		if err != nil {
			return fmt.Errorf("send kyc reward to inviter, %v", err)
		}
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(types.EventInviteReward,
				sdk.NewAttribute(types.AttributeKeyKycInviterRewardSender, region.RegionTreasureAddr),
				sdk.NewAttribute(types.AttributeKeyKycInviter, inviter),
				sdk.NewAttribute(types.AttributeKeyKycInvitee, invitee),
				sdk.NewAttribute(types.AttributeKeyKycInviterReward, types.InviteReward.String()),
			),
		)
	}
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
