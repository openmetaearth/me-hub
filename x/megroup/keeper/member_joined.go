package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/megroup/types"
)

// SetMemberJoined set a specific memberJoined in the store from its index
func (k Keeper) SetMemberJoined(ctx sdk.Context, memberJoined types.MemberJoined) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MemberJoinedKeyPrefix))
	b := k.cdc.MustMarshal(&memberJoined)
	store.Set(types.MemberJoinedKey(
		memberJoined.Address,
	), b)
}

// GetMemberJoined returns a memberJoined from its index
func (k Keeper) GetMemberJoined(
	ctx sdk.Context,
	address string,

) (val types.MemberJoined, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MemberJoinedKeyPrefix))

	b := store.Get(types.MemberJoinedKey(
		address,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveMemberJoined removes a memberJoined from the store
func (k Keeper) RemoveMemberJoined(
	ctx sdk.Context,
	address string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MemberJoinedKeyPrefix))
	store.Delete(types.MemberJoinedKey(
		address,
	))
}

// GetAllMemberJoined returns all memberJoined
func (k Keeper) GetAllMemberJoined(ctx sdk.Context) (list []types.MemberJoined) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MemberJoinedKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.MemberJoined
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
