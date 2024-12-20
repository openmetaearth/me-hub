package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/st-chain/me-hub/x/megroup/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GroupAll(goCtx context.Context, req *types.QueryAllGroupRequest) (*types.QueryAllGroupResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var groups []types.Group
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	groupStore := prefix.NewStore(store, types.KeyPrefix(types.GroupKey))

	pageRes, err := query.Paginate(groupStore, req.Pagination, func(key []byte, value []byte) error {
		var group types.Group
		if err := k.cdc.Unmarshal(value, &group); err != nil {
			return err
		}

		groups = append(groups, group)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllGroupResponse{Group: groups, Pagination: pageRes}, nil
}

func (k Keeper) Group(goCtx context.Context, req *types.QueryGetGroupRequest) (*types.QueryGetGroupResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	group, found := k.GetGroup(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetGroupResponse{Group: group}, nil
}
