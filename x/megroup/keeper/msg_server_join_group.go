package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/x/megroup/types"
)

func (k msgServer) JoinGroup(goCtx context.Context, msg *types.MsgJoinGroup) (*types.MsgJoinGroupResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	// allow join group when creator is applicant or global DAO or Meid DAO

	if msg.ApplicantAddress != msg.Creator {
		if !(k.daoKeeper.IsGlobalDao(ctx, msg.Creator) || k.daoKeeper.GetMeidDao(ctx).Equals(creator)) {
			return nil, errors.Wrapf(types.ErrPermissionDenied, "msg.ApplicantAddress != msg.Creator. but Creator is not adminAddress")
		}
	}

	groupInfo, found := k.GetGroup(ctx, msg.GroupId)
	if !found {
		return nil, errors.Wrapf(types.ErrGroupNotExist, fmt.Sprintf("msg's groupID = %d", msg.GroupId))
	}

	//
	_, found = k.stakingKeeper.GetMeid(ctx, msg.ApplicantAddress)
	if !found {
		errLogBytes := fmt.Sprintf("only MEID user can join group. user's address = %s", msg.Creator)
		return nil, errors.Wrapf(types.ErrPermissionDenied, string(errLogBytes))
	}

	joined, JoinGroupFound := k.GetMemberJoined(ctx, msg.ApplicantAddress)
	if JoinGroupFound && joined.GroupId > 0 {
		errLogBytes := fmt.Sprintf("user has joined a group (groupID:%d)", joined.GroupId)
		return nil, errors.Wrapf(types.ErrPermissionDenied, errLogBytes)
	}

	//set member's join group info
	k.SetMemberJoined(ctx, types.MemberJoined{
		Address: msg.ApplicantAddress,
		GroupId: msg.GroupId,
	})
	//add to group_member
	err = k.addGroupMember(ctx, &types.GroupMember{
		GroupID: msg.GroupId,
		Member: &types.Member{
			Address: msg.ApplicantAddress,
			AddedAt: ctx.BlockTime()}})
	if err != nil {
		return nil, err
	}
	if !JoinGroupFound { //send rewards if user has not joined group
		//get RegionTreasureAddr
		region, found := k.stakingKeeper.GetRegion(ctx, groupInfo.RegionID)
		if !found {
			return nil, errors.Wrapf(types.ErrRegionNotExist, fmt.Sprintf("group's region = %d", groupInfo.RegionID))
		}
		rewardsCoin := sdk.NewCoin(params.BaseDenom, math.NewInt(1000000))
		err = k.bankKeeper.SendCoins(ctx, sdk.MustAccAddressFromBech32(region.GetRegionTreasureAddr()),
			sdk.MustAccAddressFromBech32(msg.ApplicantAddress), sdk.NewCoins(rewardsCoin))
		if err != nil {
			return nil, errors.Wrapf(types.ErrProcData, fmt.Sprintf("transfer rewards coins error. err = %s,fromAddr = %s,toAddr = %s",
				err.Error(), region.GetRegionTreasureAddr(), msg.ApplicantAddress))
		}
		err = k.bankKeeper.SendCoins(ctx, sdk.MustAccAddressFromBech32(region.GetRegionTreasureAddr()),
			sdk.MustAccAddressFromBech32(groupInfo.Admin), sdk.NewCoins(rewardsCoin))
		if err != nil {
			return nil, errors.Wrapf(types.ErrProcData, fmt.Sprintf("transfer rewards coins error. err = %s,fromAddr = %s,toAddr = %s",
				err.Error(), region.GetRegionTreasureAddr(), groupInfo.Admin))
		}

	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EvtJoinGroup,
		sdk.NewAttribute("group_id", fmt.Sprintf("%d", msg.GroupId)),
		sdk.NewAttribute("creator", msg.Creator),
		sdk.NewAttribute("user", msg.ApplicantAddress),
		//1sdk.NewAttribute("metadata", msg.),
	))
	return &types.MsgJoinGroupResponse{}, nil
}

func (k msgServer) LeaveGroup(goCtx context.Context, req *types.MsgLeaveGroupRequest) (*types.MsgLeaveGroupResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	joined, found := k.GetMemberJoined(ctx, req.Creator)
	if !found {
		return nil, errors.Wrapf(types.ErrExcute, "can not found join group")
	}
	if joined.GroupId != req.GroupId {
		return nil, errors.Wrapf(types.ErrExcute, fmt.Sprint("group info dismatch.input group's id = %d,join gropp's id = %d",
			req.GroupId, joined.GroupId))
	}

	if err := k.deleteMemberFormGroup(ctx, req.GroupId, req.Creator); err != nil {
		return nil, err
	}

	joined.GroupId = 0
	k.SetMemberJoined(ctx, joined)
	return &types.MsgLeaveGroupResponse{}, nil
}
