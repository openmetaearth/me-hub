package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (k Keeper) HasKycRewardPaid(ctx sdk.Context, account sdk.AccAddress, regionId string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.KycRewardPaidKeyPrefix))
	return store.Has(types.KycRewardPaidKey(account, regionId))
}

func (k Keeper) SetKycRewardPaid(ctx sdk.Context, account sdk.AccAddress, regionId string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.KycRewardPaidKeyPrefix))
	store.Set(types.KycRewardPaidKey(account, regionId), []byte{1})
}
