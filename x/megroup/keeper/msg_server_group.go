package keeper

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/st-chain/me-hub/x/megroup/types"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateGroup(goCtx context.Context, msg *types.MsgCreateGroup) (*types.MsgCreateGroupResponse, error) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), "created_group")
	ctx := sdk.UnwrapSDKContext(goCtx)

	//check permission
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}
	if !(k.daoKeeper.IsGlobalDao(ctx, msg.Creator) ||
		k.daoKeeper.GetMeidDao(ctx).Equals(creator)) {
		errLogBytes := types.EmitNewGroupError("only global admin can  create group", msg)
		return nil, errors.Wrapf(types.ErrCheckGlobalAdmin, string(errLogBytes))
	}

	if msg.GroupInfo == nil {
		return nil, errors.Wrapf(types.ErrCreate, "group info is nil")
	}

	//check region name
	_, err = k.stakingKeeper.CheckRegionName(strings.ToUpper(msg.GroupInfo.RegionID))
	if err != nil {
		errLogBytes := types.EmitNewGroupError(fmt.Sprintf("region id %s illegal", msg.GroupInfo.RegionID), msg)
		return nil, errors.Wrapf(types.ErrCreate, string(errLogBytes))
	}
	_, found := k.stakingKeeper.GetMeid(ctx, msg.GroupInfo.Admin)
	if !found {
		errLogBytes := types.EmitNewGroupError("only MEID user can create group", msg)
		return nil, errors.Wrapf(types.ErrMeidNotExists, string(errLogBytes))
	}

	joined, found := k.GetMemberJoined(ctx, msg.GroupInfo.Admin)
	if found {
		errLogBytes := types.EmitNewGroupError(fmt.Sprintf("admin has joined a group (groupID:%d)", joined.GroupId), msg)
		return nil, errors.Wrapf(types.ErrCreate, string(errLogBytes))
	}

	newGroupID := k.GetGroupCount(ctx)
	var group = types.Group{
		Creator: msg.Creator,
		GroupInfo: &types.GroupInfo{
			Id:          newGroupID,
			Admin:       msg.GroupInfo.Admin,
			Metadata:    msg.GroupInfo.Metadata,
			Version:     1,
			TotalWeight: math.NewInt(0).String(),
			CreatedAt:   ctx.BlockTime(),
			RegionID:    msg.GroupInfo.RegionID,
		},
	}

	id := k.AppendGroup(
		ctx,
		group,
	)
	//Mark admin has joined the group
	k.SetMemberJoined(ctx, types.MemberJoined{
		Address:      msg.GroupInfo.Admin,
		GroupId:      newGroupID,
		MemberListId: 0,
	})
	//TODO: add admin to group member

	ctx.EventManager().EmitEvent(sdk.NewEvent("EventGroupCreated",
		sdk.NewAttribute("group_id", fmt.Sprintf("%d", newGroupID)),
		sdk.NewAttribute("admin", msg.GroupInfo.Admin),
		sdk.NewAttribute("region_id", msg.GroupInfo.RegionID),
		sdk.NewAttribute("metadata", msg.GroupInfo.Metadata),
	))
	return &types.MsgCreateGroupResponse{
		Id: id,
	}, nil
}

func (k msgServer) UpdateGroup(goCtx context.Context, msg *types.MsgUpdateGroup) (*types.MsgUpdateGroupResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Checks that the element exists
	val, found := k.GetGroup(ctx, msg.Id)
	if !found {
		return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}
	if msg.GroupInfo == nil {
		return nil, errors.Wrapf(types.ErrCreate, "group info is nil")
	}
	// if new admin
	if msg.GroupInfo.Admin != val.GroupInfo.Admin {
		//TODO: join group
		//TODO: abort if admin has joined other group
	}

	var group = types.Group{
		Creator:   msg.Creator,
		Id:        msg.Id,
		GroupInfo: nil,
	}
	group.GroupInfo = &types.GroupInfo{
		// not mut
		Id:          val.GroupInfo.Id,
		TotalWeight: val.GroupInfo.TotalWeight,
		CreatedAt:   val.GroupInfo.CreatedAt,
		// mut
		RegionID: msg.GroupInfo.RegionID,
		Admin:    msg.GroupInfo.Admin,
		Metadata: msg.GroupInfo.Metadata,
		Version:  val.GroupInfo.Version + 1,
	}

	k.SetGroup(ctx, group)
	//TODO: emit event
	return &types.MsgUpdateGroupResponse{}, nil
}

func (k msgServer) DeleteGroup(goCtx context.Context, msg *types.MsgDeleteGroup) (*types.MsgDeleteGroupResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Checks that the element exists
	val, found := k.GetGroup(ctx, msg.Id)
	if !found {
		return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	//check has existed for 365 days
	if !val.GroupInfo.CreatedAt.Before(ctx.BlockTime().Add(-time.Hour * 24 * 365)) {
		return nil, errors.Wrapf(types.ErrDeleteGroup, "group has existed for 365 days")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}
	//remove all member in group
	k.LoadMemberStoreByGroupID(ctx, msg.Id).DestroyThisGroup()

	//remove group
	k.RemoveGroup(ctx, msg.Id)
	//TODO: emit event
	return &types.MsgDeleteGroupResponse{}, nil
}
