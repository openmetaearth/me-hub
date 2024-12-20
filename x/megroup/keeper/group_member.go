package keeper

import (
	"encoding/binary"
	"fmt"

	"github.com/st-chain/me-hub/x/megroup/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MemberStoreByGroupID struct {
	GroupID     uint64
	cdc         codec.BinaryCodec
	memberStore prefix.Store
	counts      prefix.Store
	k           Keeper
	ctx         sdk.Context
}

func (k Keeper) LoadMemberStoreByGroupID(ctx sdk.Context, groupID uint64) MemberStoreByGroupID {
	return MemberStoreByGroupID{
		GroupID:     groupID,
		cdc:         k.cdc,
		memberStore: prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GroupMemberKey+fmt.Sprintf("%d", groupID))),
		counts:      prefix.NewStore(ctx.KVStore(k.storeKey), []byte{}),
		ctx:         ctx,
	}
}

// GetGroupMemberCount get the total number of groupMember in a groupID
func (ms MemberStoreByGroupID) GetGroupMemberCount() uint64 {
	val, found := ms.k.GetGroupMemberCount(ms.ctx, ms.GroupID)
	if !found {
		return 0
	}
	// Parse bytes
	return val.Num
}

// SetGroupMemberCount set the total number of groupMember
func (ms MemberStoreByGroupID) SetGroupMemberCount(count uint64) {
	ms.k.SetGroupMemberCount(ms.ctx, types.GroupMemberCount{
		GroupId: ms.GroupID,
		Num:     count,
	})
}

// AppendGroupMember appends a groupMember in the store with a new id and update the count
func (ms MemberStoreByGroupID) AppendGroupMember(
	ctx sdk.Context,
	groupMember types.GroupMember,
) uint64 {
	// Create the groupMember
	count := ms.GetGroupMemberCount()

	// Set the ID of the appended value
	groupMember.Id = count

	appendedValue := ms.cdc.MustMarshal(&groupMember)
	ms.memberStore.Set(GetGroupMemberIDBytes(groupMember.Id), appendedValue)

	// Update groupMember count
	ms.SetGroupMemberCount(count + 1)

	return count
}

// SetGroupMember set a specific groupMember in the store
func (ms MemberStoreByGroupID) SetGroupMember(groupMember types.GroupMember) {
	b := ms.cdc.MustMarshal(&groupMember)
	ms.memberStore.Set(GetGroupMemberIDBytes(groupMember.Id), b)
}

// GetGroupMember returns a groupMember from its id
func (k MemberStoreByGroupID) GetGroupMember(id uint64) (val types.GroupMember, found bool) {
	b := k.memberStore.Get(GetGroupMemberIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveGroupMember removes a groupMember from the store
func (ms MemberStoreByGroupID) RemoveGroupMember(id uint64) {
	ms.memberStore.Delete(GetGroupMemberIDBytes(id))
}

// GetAllGroupMember returns all groupMember
func (ms MemberStoreByGroupID) GetAllGroupMember() (list []types.GroupMember) {
	iterator := sdk.KVStorePrefixIterator(ms.memberStore, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.GroupMember
		ms.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
func (ms MemberStoreByGroupID) DestroyThisGroup() {
	members := ms.GetAllGroupMember()
	for _, member := range members {
		ms.RemoveGroupMember(member.Id)
	}
	ms.k.RemoveGroupMemberCount(ms.ctx, ms.GroupID)
}

// GetGroupMemberIDBytes returns the byte representation of the ID
func GetGroupMemberIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetGroupMemberIDFromBytes returns ID in uint64 format from a byte array
func GetGroupMemberIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) GetAllGroupMember(ctx sdk.Context) (list []types.GroupMember) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GroupMemberKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.GroupMember
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetGroupMemberCount get the total number of groupMember
func (k Keeper) GetGroupTotalMemberCount(ctx sdk.Context) uint64 {
	groupCountList := k.GetAllGroupMemberCount(ctx)
	var totalCount uint64
	for _, num := range groupCountList {
		totalCount += num.Num
	}
	return totalCount
}
