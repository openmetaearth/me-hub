package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/openmetaearth/me-hub/x/megroup/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GroupMemberCountAll(goCtx context.Context, req *types.QueryAllGroupMemberCountRequest) (*types.QueryAllGroupMemberCountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var groupMemberCounts []types.GroupMemberCount
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	groupMemberCountStore := prefix.NewStore(store, types.KeyPrefix(types.GroupMemberCountKeyPrefix))

	pageRes, err := query.Paginate(groupMemberCountStore, req.Pagination, func(key []byte, value []byte) error {
		var groupMemberCount types.GroupMemberCount
		if err := k.cdc.Unmarshal(value, &groupMemberCount); err != nil {
			return err
		}

		groupMemberCounts = append(groupMemberCounts, groupMemberCount)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllGroupMemberCountResponse{GroupMemberCount: groupMemberCounts, Pagination: pageRes}, nil
}

func (k Keeper) GroupMemberCount(goCtx context.Context, req *types.QueryGetGroupMemberCountRequest) (*types.QueryGetGroupMemberCountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetGroupMemberCount(
		ctx,
		req.GroupId,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetGroupMemberCountResponse{GroupMemberCount: val}, nil
}
