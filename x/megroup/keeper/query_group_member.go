package keeper

import (
	"context"
	"fmt"

	"github.com/openmetaearth/me-hub/x/megroup/types"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
)

func (k Keeper) GroupMemberAll(goCtx context.Context, req *types.QueryGroupAllMemberRequest) (*types.QueryGroupAllMemberResponse, error) {
	if req == nil {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, "")
	}

	var groupMembers []types.GroupMember
	ctx := sdk.UnwrapSDKContext(goCtx)

	grpMemberPrefix := fmt.Sprintf("%s%d/", types.GroupMemberKey, req.GroupID)
	groupMemberStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(grpMemberPrefix))

	pageRes, err := query.Paginate(groupMemberStore, req.Pagination, func(key []byte, value []byte) error {
		var groupMember types.GroupMember
		if err := k.cdc.Unmarshal(value, &groupMember); err != nil {
			return err
		}

		groupMembers = append(groupMembers, groupMember)
		return nil
	})

	if err != nil {
		return nil, errors.Wrapf(sdkerrors.ErrLogic, "query.Paginate error.err = %s. req.Pagination = %s", err.Error(), req.Pagination.String())
	}

	return &types.QueryGroupAllMemberResponse{GroupID: req.GroupID, GroupMember: groupMembers, Pagination: pageRes}, nil
}

func (k Keeper) GroupMember(goCtx context.Context, req *types.QueryGetGroupMemberRequest) (*types.QueryGetGroupMemberResponse, error) {
	if req == nil {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, "")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	joined, found := k.GetMemberJoined(ctx, req.Address)
	if !found {
		return nil, errors.Wrapf(sdkerrors.ErrKeyNotFound, "can not found group by member address")
	}
	grpMemberPrefix := fmt.Sprintf("%s%d/", types.GroupMemberKey, joined.GroupId)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(grpMemberPrefix))
	data := store.Get([]byte(req.Address))
	if nil == data {
		return nil, errors.Wrapf(sdkerrors.ErrKeyNotFound, "can not found groupMember by memberJoined info. join groupID = %d,memberAddress = %s", joined.GroupId, req.Address)
	}
	var groupMember types.GroupMember
	if err := k.cdc.Unmarshal(data, &groupMember); err != nil {
		return nil, errors.Wrapf(sdkerrors.ErrLogic, "cdc.Unmarshal groupMember data error.err = %s.", err.Error())
	}

	return &types.QueryGetGroupMemberResponse{GroupMember: groupMember}, nil
}
