package keeper

import (
	"context"
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/megroup/types"
)

func (k msgServer) JoinGroup(goCtx context.Context, msg *types.MsgJoinGroup) (*types.MsgJoinGroupResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	//creator, err := sdk.AccAddressFromBech32(msg.Creator)
	//if err != nil {
	//	return nil, err
	//}

	// allow join group when creator is applicant or global DAO or Meid DAO

	if msg.ApplicantAddress != msg.Creator {
		if !(k.daoKeeper.IsGlobalDao(ctx, msg.Creator) || k.daoKeeper.GetMeidDao(ctx) == msg.Creator) {
			return nil, errors.Wrapf(types.ErrPermissionDenied, "msg.ApplicantAddress != msg.Creator. but Creator is not adminAddress")
		}
	}

	userAccAddr, err := sdk.AccAddressFromBech32(msg.ApplicantAddress)
	if err != nil {
		return nil, errors.Wrapf(types.ErrProcData, "sdk.AccAddressFromBech32 error.err = %s,addr = %s",
			err.Error(), msg.ApplicantAddress)
	}

	groupInfo, found := k.GetGroupInfo(ctx, msg.GroupId)
	if !found {
		return nil, errors.Wrapf(types.ErrGroupNotExist, "msg's groupID = %d", msg.GroupId)
	}

	_, isKycActive := k.GetDidAndKycActive(ctx, userAccAddr, groupInfo.RegionID)
	if !isKycActive {
		return nil, errors.Wrapf(types.ErrPermissionDenied, "can not found hight kyc level user's did active status in group's region"+
			"address = %s, group's regionID = %s", msg.ApplicantAddress, groupInfo.RegionID)
	}

	grpNumber, found := k.GetGroupMemberCount(ctx, msg.GroupId)
	if !found {
		return nil, errors.Wrap(types.ErrProcData, fmt.Sprintf("can not found group number count in JoinGroup"))
	}

	joined, JoinGroupFound := k.GetMemberJoined(ctx, msg.ApplicantAddress)
	if JoinGroupFound && joined.GroupId > 0 {
		errLogBytes := fmt.Sprintf("user has joined a group (groupID:%d)", joined.GroupId)
		return nil, errors.Wrap(types.ErrPermissionDenied, errLogBytes)
	}

	//set member's join group info
	k.SetMemberJoined(ctx, types.MemberJoined{
		Address: msg.ApplicantAddress,
		GroupId: msg.GroupId,
	})
	//add to group_member

	err = k.AddGroupMember(ctx, &types.GroupMember{
		GroupId: msg.GroupId,
		Member: &types.Member{
			Address: msg.ApplicantAddress,
			AddedAt: ctx.BlockTime()}})
	if err != nil {
		return nil, err
	}
	k.SetGroupMemberCount(ctx, msg.GroupId, grpNumber+1)

	if !JoinGroupFound { //send rewards if user has not joined group
		//get RegionTreasureAddr
		region, found := k.stakingKeeper.GetRegion(ctx, groupInfo.RegionID)
		if !found {
			return nil, errors.Wrap(types.ErrRegionNotExist, fmt.Sprintf("group's region: %s", groupInfo.RegionID))
		}
		rewardsCoin := sdk.NewCoin(params.BaseDenom, math.NewInt(1000000))
		err = k.bankKeeper.Extend().SendCoinsWithTag(ctx, sdk.MustAccAddressFromBech32(region.GetRegionTreasureAddr()),
			sdk.MustAccAddressFromBech32(msg.ApplicantAddress), sdk.NewCoins(rewardsCoin), fmt.Sprintf("JoinGroup_SendApplicantRewards_%s", region.RegionId))
		if err != nil {
			return nil, errors.Wrap(types.ErrProcData, fmt.Sprintf("transfer rewards coins error. err = %s,fromAddr = %s,toAddr = %s",
				err.Error(), region.GetRegionTreasureAddr(), msg.ApplicantAddress))
		}
		err = k.bankKeeper.Extend().SendCoinsWithTag(ctx, sdk.MustAccAddressFromBech32(region.GetRegionTreasureAddr()),
			sdk.MustAccAddressFromBech32(groupInfo.Admin), sdk.NewCoins(rewardsCoin), fmt.Sprintf("JoinGroup_SendAdminRewards_%s", region.RegionId))
		if err != nil {
			return nil, errors.Wrap(types.ErrProcData, fmt.Sprintf("transfer rewards coins error. err = %s,fromAddr = %s,toAddr = %s",
				err.Error(), region.GetRegionTreasureAddr(), groupInfo.Admin))
		}
		ctx.EventManager().EmitEvent(sdk.NewEvent(types.EvtJoinGroupReward,
			sdk.NewAttribute("applicant", msg.ApplicantAddress),
			sdk.NewAttribute("admin", groupInfo.Admin),
			sdk.NewAttribute("regionTreasureAddress", region.GetRegionTreasureAddr()),
			sdk.NewAttribute("rewards", rewardsCoin.String()),
		))
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EvtJoinGroup,
		sdk.NewAttribute("group_id", fmt.Sprintf("%d", msg.GroupId)),
		sdk.NewAttribute("creator", msg.Creator),
		sdk.NewAttribute("applicant", msg.ApplicantAddress),
	))
	return &types.MsgJoinGroupResponse{}, nil
}

func (k msgServer) LeaveGroup(goCtx context.Context, req *types.MsgLeaveGroupRequest) (*types.MsgLeaveGroupResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	groupInfo, found := k.GetGroupInfo(ctx, req.GroupId)
	if !found {
		return nil, errors.Wrap(types.ErrGroupNotExist, fmt.Sprintf("can not found gourp.groupID = %d", req.GroupId))
	}

	if req.Creator == groupInfo.Admin { //admin can not leave group
		return nil, errors.Wrapf(types.ErrExecute, "admin of group can not leave")
	}

	joined, found := k.GetMemberJoined(ctx, req.Creator)
	if !found {
		return nil, errors.Wrapf(types.ErrExecute, "can not found join group")
	}
	if joined.GroupId != req.GroupId {
		return nil, errors.Wrap(types.ErrExecute, fmt.Sprintf("group info dismatch.input group's id = %d,join gropp's id = %d", req.GroupId, joined.GroupId))
	}

	grpNumber, found := k.GetGroupMemberCount(ctx, req.GroupId)
	if !found {
		return nil, errors.Wrapf(types.ErrProcData, "can not found group number count in LeaveGroup")
	}
	if grpNumber == 0 {
		return nil, errors.Wrapf(types.ErrProcData, "group number is 0 in LeaveGroup")
	}

	if err := k.deleteMemberFormGroup(ctx, req.GroupId, req.Creator); err != nil {
		return nil, err
	}

	joined.GroupId = 0
	k.SetMemberJoined(ctx, joined)
	k.SetGroupMemberCount(ctx, req.GroupId, grpNumber-1)

	return &types.MsgLeaveGroupResponse{}, nil
}
