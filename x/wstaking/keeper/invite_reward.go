package keeper

import (
	"fmt"
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

	if err := k.bankKeeper.Extend().SendCoinsWithTag(ctx,
		sdk.MustAccAddressFromBech32(region.RegionTreasureAddr),
		sdk.MustAccAddressFromBech32(inviter),
		sdk.NewCoins(sdk.NewCoin(params.BaseDenom, types.InviteReward)),
		fmt.Sprintf("SendInviteReward_%s", region.RegionId),
	); err != nil {
		return fmt.Errorf("send invite reward to inviter, %v", err)
	}

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