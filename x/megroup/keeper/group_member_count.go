package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/megroup/types"
)

// SetGroupMemberCount set a specific groupMemberCount in the store from its index
func (k Keeper) SetGroupMemberCount(ctx sdk.Context, groupID, number uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GroupMemberCountKeyPrefix))
	store.Set(types.GetBytesFromUint64(groupID), types.GetBytesFromUint64(number))
}

// GetGroupMemberCount returns a groupMemberCount from its index
func (k Keeper) GetGroupMemberCount(ctx sdk.Context, groupId uint64) (uint64, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GroupMemberCountKeyPrefix))

	val := store.Get(types.GetBytesFromUint64(groupId))
	if nil == val {
		return 0, false
	}
	return types.GetUint64FromBytes(val), true
}

// RemoveGroupMemberCount removes a groupMemberCount from the store
func (k Keeper) RemoveGroupMemberCount(
	ctx sdk.Context,
	groupId uint64,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GroupMemberCountKeyPrefix))
	store.Delete(types.GetBytesFromUint64(groupId))
}

// GetAllGroupMemberCount returns all groupMemberCount
/*
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

*/
