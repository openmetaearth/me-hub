package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func (k *Keeper) SetApproveReward(ctx sdk.Context, address, inviter, issuer, originId string) error {
	return k.stkKeeper.KycReward(ctx, sdk.MustAccAddressFromBech32(address), inviter, originId, issuer)
}

func (k *Keeper) DeleteApproveReward(ctx sdk.Context, address string) error {
	return k.stkKeeper.RemoveKycReward(ctx, sdk.MustAccAddressFromBech32(address))

}
