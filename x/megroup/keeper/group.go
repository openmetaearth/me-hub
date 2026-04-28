package keeper

import (
	"cosmossdk.io/errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/megroup/types"
)

// GetGroupCount get the total number of group
func (k Keeper) GetLastGroupID(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.GroupLastIDKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return types.GetUint64FromBytes(bz)
}

// SetGroupCount set the total number of group
func (k Keeper) SetLastGroupID(ctx sdk.Context, groupID uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.GroupLastIDKey)
	store.Set(byteKey, types.GetBytesFromUint64(groupID))
}

// AppendGroup appends a group in the store with a new id and update the count
func (k Keeper) AppendGroup(ctx sdk.Context, group *types.GroupInfo) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GroupKey))
	if nil != store.Get(types.GetBytesFromUint64(group.Id)) {
		return errors.Wrapf(types.ErrGroupCreateRepeated, "group id has bee existed.groupID = %d", group.Id)

	}
	appendedValue := k.cdc.MustMarshal(group)

	store.Set(types.GetBytesFromUint64(group.Id), appendedValue)

	// Update group count
	k.SetLastGroupID(ctx, group.Id)

	return nil
}

func (k Keeper) SetGroupInfo(ctx sdk.Context, group types.GroupInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GroupKey))
	b := k.cdc.MustMarshal(&group)
	store.Set(types.GetBytesFromUint64(group.Id), b)
}

// GetGroup returns a group from its id
func (k Keeper) GetGroupInfo(ctx sdk.Context, id uint64) (val types.GroupInfo, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GroupKey))
	b := store.Get(types.GetBytesFromUint64(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveGroup removes a group from the store
func (k Keeper) RemoveGroup(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GroupKey))
	store.Delete(types.GetBytesFromUint64(id))
}

// GetAllGroup returns all group
func (k Keeper) GetAllGroup(ctx sdk.Context) (list []types.GroupInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GroupKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.GroupInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) DeleteGroupAssociateWithRegion(ctx sdk.Context, regionID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GroupRegionKey))
	store.Delete([]byte(regionID))
}

func (k Keeper) SetGroupToRegion(ctx sdk.Context, regionID string, groupID uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GroupRegionKey))
	store.Set([]byte(regionID), types.GetBytesFromUint64(groupID))
}

func (k Keeper) GetGroupIdByRegion(ctx sdk.Context, regionID string) (uint64, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GroupRegionKey))
	data := store.Get([]byte(regionID))
	if nil == data {
		return 0, false
	}
	return types.GetUint64FromBytes(data), true
}

func (k Keeper) UpdateGroupAdmin(ctx sdk.Context, regionID string, admin string) {
	groupId, found := k.GetGroupIdByRegion(ctx, regionID)
	if !found {
		return
	}
	group, found := k.GetGroupInfo(ctx, groupId)
	if !found {
		return
	}
	group.Admin = admin
	k.SetGroupInfo(ctx, group)
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EvtUpdateGroupAdmin,
		sdk.NewAttribute("group_id", fmt.Sprintf("%d", groupId)),
		sdk.NewAttribute("admin", admin),
		sdk.NewAttribute("region_id", regionID),
	))
}
