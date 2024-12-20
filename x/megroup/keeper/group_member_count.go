package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/megroup/types"
)

// SetGroupMemberCount set a specific groupMemberCount in the store from its index
func (k Keeper) SetGroupMemberCount(ctx sdk.Context, groupMemberCount types.GroupMemberCount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GroupMemberCountKeyPrefix))
	b := k.cdc.MustMarshal(&groupMemberCount)
	store.Set(types.GroupMemberCountKey(
		groupMemberCount.GroupId,
	), b)
}

// GetGroupMemberCount returns a groupMemberCount from its index
func (k Keeper) GetGroupMemberCount(
	ctx sdk.Context,
	groupId uint64,

) (val types.GroupMemberCount, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GroupMemberCountKeyPrefix))

	b := store.Get(types.GroupMemberCountKey(
		groupId,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveGroupMemberCount removes a groupMemberCount from the store
func (k Keeper) RemoveGroupMemberCount(
	ctx sdk.Context,
	groupId uint64,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GroupMemberCountKeyPrefix))
	store.Delete(types.GroupMemberCountKey(
		groupId,
	))
}

// GetAllGroupMemberCount returns all groupMemberCount
func (k Keeper) GetAllGroupMemberCount(ctx sdk.Context) (list []types.GroupMemberCount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GroupMemberCountKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.GroupMemberCount
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
