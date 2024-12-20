package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/megroup/types"
)

// GetGroupCount get the total number of group
func (k Keeper) GetGroupCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.GroupCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetGroupCount set the total number of group
func (k Keeper) SetGroupCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.GroupCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendGroup appends a group in the store with a new id and update the count
func (k Keeper) AppendGroup(
	ctx sdk.Context,
	group types.Group,
) uint64 {
	// Create the group
	count := k.GetGroupCount(ctx)

	// Set the ID of the appended value
	group.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GroupKey))
	appendedValue := k.cdc.MustMarshal(&group)
	store.Set(GetGroupIDBytes(group.Id), appendedValue)

	// Update group count
	k.SetGroupCount(ctx, count+1)

	return count
}

// SetGroup set a specific group in the store
func (k Keeper) SetGroup(ctx sdk.Context, group types.Group) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GroupKey))
	b := k.cdc.MustMarshal(&group)
	store.Set(GetGroupIDBytes(group.Id), b)
}

// GetGroup returns a group from its id
func (k Keeper) GetGroup(ctx sdk.Context, id uint64) (val types.Group, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GroupKey))
	b := store.Get(GetGroupIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveGroup removes a group from the store
func (k Keeper) RemoveGroup(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GroupKey))
	store.Delete(GetGroupIDBytes(id))
}

// GetAllGroup returns all group
func (k Keeper) GetAllGroup(ctx sdk.Context) (list []types.Group) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GroupKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Group
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetGroupIDBytes returns the byte representation of the ID
func GetGroupIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetGroupIDFromBytes returns ID in uint64 format from a byte array
func GetGroupIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
