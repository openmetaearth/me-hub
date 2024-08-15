package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

// GetMeid returns a meid from its index
func (k Keeper) GetMeid(ctx sdk.Context, account string) (val types.Meid, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MeidKeyPrefix))
	b := store.Get(types.MeidKey(account))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
