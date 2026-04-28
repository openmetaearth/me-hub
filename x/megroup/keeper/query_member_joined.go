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

func (k Keeper) MemberJoinedAll(goCtx context.Context, req *types.QueryAllMemberJoinedRequest) (*types.QueryAllMemberJoinedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var memberJoineds []types.MemberJoined
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	memberJoinedStore := prefix.NewStore(store, types.KeyPrefix(types.MemberJoinedKeyPrefix))

	pageRes, err := query.Paginate(memberJoinedStore, req.Pagination, func(key []byte, value []byte) error {
		var memberJoined types.MemberJoined
		if err := k.cdc.Unmarshal(value, &memberJoined); err != nil {
			return err
		}

		memberJoineds = append(memberJoineds, memberJoined)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllMemberJoinedResponse{MemberJoined: memberJoineds, Pagination: pageRes}, nil
}

func (k Keeper) MemberJoined(goCtx context.Context, req *types.QueryGetMemberJoinedRequest) (*types.QueryGetMemberJoinedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetMemberJoined(
		ctx,
		req.Address,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetMemberJoinedResponse{MemberJoined: val}, nil
}
