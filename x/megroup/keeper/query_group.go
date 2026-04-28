package keeper

import (
	"context"
	"cosmossdk.io/errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/openmetaearth/me-hub/x/megroup/types"
)

func (k Keeper) GroupAll(goCtx context.Context, req *types.QueryAllGroupRequest) (*types.QueryAllGroupResponse, error) {
	if req == nil {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, "")
	}

	var groups []types.GroupInfo
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	groupStore := prefix.NewStore(store, types.KeyPrefix(types.GroupKey))

	pageRes, err := query.Paginate(groupStore, req.Pagination, func(key []byte, value []byte) error {
		var group types.GroupInfo
		if err := k.cdc.Unmarshal(value, &group); err != nil {
			return err
		}

		groups = append(groups, group)
		return nil
	})

	if err != nil {
		return nil, errors.Wrap(sdkerrors.ErrLogic, fmt.Sprintf(" query.Paginate error.err = %s,req.Pagination = %s",
			err.Error(), req.Pagination.String()))
	}

	return &types.QueryAllGroupResponse{Group: groups, Pagination: pageRes}, nil
}

func (k Keeper) Group(goCtx context.Context, req *types.QueryGetGroupRequest) (*types.QueryGetGroupResponse, error) {
	if req == nil {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, "")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	group, found := k.GetGroupInfo(ctx, req.Id)
	if !found {
		return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("can not found group by groupID.groupID = %d", req.Id))
	}

	return &types.QueryGetGroupResponse{Group: group}, nil
}

func (k Keeper) GroupByMember(goCtx context.Context, req *types.QueryGroupByMemberRequest) (*types.QueryGetGroupResponse, error) {
	if req == nil {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, "")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	joined, found := k.GetMemberJoined(ctx, req.Address)
	if !found {
		return nil, errors.Wrapf(sdkerrors.ErrKeyNotFound, "can not found memberJoin info by address")
	}

	group, found := k.GetGroupInfo(ctx, joined.GroupId)
	if !found {
		return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("can not found group by memberJoin's groupID.groupID = %d", joined.GroupId))
	}
	return &types.QueryGetGroupResponse{Group: group}, nil
}
