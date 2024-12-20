package keeper

import (
	"context"

	"github.com/st-chain/me-hub/x/megroup/types"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GroupMemberAll(goCtx context.Context, req *types.QueryAllGroupMemberRequest) (*types.QueryAllGroupMemberResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var groupMembers []types.GroupMember
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	groupMemberStore := prefix.NewStore(store, types.KeyPrefix(types.GroupMemberKey))

	pageRes, err := query.Paginate(groupMemberStore, req.Pagination, func(key []byte, value []byte) error {
		var groupMember types.GroupMember
		if err := k.cdc.Unmarshal(value, &groupMember); err != nil {
			return err
		}

		groupMembers = append(groupMembers, groupMember)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllGroupMemberResponse{GroupMember: groupMembers, Pagination: pageRes}, nil
}

func (k Keeper) GroupMember(goCtx context.Context, req *types.QueryGetGroupMemberRequest) (*types.QueryGetGroupMemberResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	groupMember, found := k.GetMemberJoined(ctx, req.Address)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}
	gm, found := k.LoadMemberStoreByGroupID(ctx, groupMember.GroupId).GetGroupMember(groupMember.MemberListId)
	if !found {
		//THIS SHOULD NOT HAPPEN
		return nil, errors.Wrapf(sdkerrors.ErrKeyNotFound, "THIS SHOULD NEVER HAPPEN :group member not found %s", groupMember.Address)
	}
	return &types.QueryGetGroupMemberResponse{GroupMember: gm}, nil
}
