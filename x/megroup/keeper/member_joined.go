package keeper

import (
	"cosmossdk.io/errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/megroup/types"
)

// SetMemberJoined set a specific memberJoined in the store from its index
func (k Keeper) SetMemberJoined(ctx sdk.Context, memberJoined types.MemberJoined) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MemberJoinedKeyPrefix))
	b := k.cdc.MustMarshal(&memberJoined)
	store.Set(types.MemberJoinedKey(
		memberJoined.Address,
	), b)
}

func (k *Keeper) AddGroupMember(ctx sdk.Context, grpMember *types.GroupMember) error {
	grpMemberPrefix := fmt.Sprintf("%s%d/", types.GroupMemberKey, grpMember.GroupId)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(grpMemberPrefix))
	if nil != store.Get([]byte(grpMember.Member.Address)) {
		return errors.Wrapf(types.ErrGroupMemberRepeated, "member has been joined this group store.groupID = %d",
			grpMember.GroupId)
	}
	val := k.cdc.MustMarshal(grpMember)
	store.Set([]byte(grpMember.Member.Address), val)
	return nil
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

func (k *Keeper) deleteMemberFormGroup(ctx sdk.Context, groupID uint64, address string) error {
	grpMemberPrefix := fmt.Sprintf("%s%d/", types.GroupMemberKey, groupID)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(grpMemberPrefix))
	addrBytes := []byte(address)
	if nil == store.Get(addrBytes) {
		return errors.Wrapf(types.ErrGroupMemberNotExist, "can not found member in group.addr: %s,groupID: %d", address, groupID)
	}
	store.Delete(addrBytes)
	return nil
}
