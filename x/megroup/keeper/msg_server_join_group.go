package keeper

import (
	"context"

	"github.com/st-chain/me-hub/x/megroup/types"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
			return nil, types.ErrPermissionDenied
		}
	}

	//TODO: check applicant is MEID user

	//check group is exist
	if _, found := k.GetGroup(ctx, msg.GroupId); !found {
		return nil, errors.Wrapf(types.ErrGroupNotExist, "group id %d", msg.GroupId)
	}

	//add to group_member
	groupMenberStore := k.LoadMemberStoreByGroupID(ctx, msg.GroupId)
	groupMenberStore.AppendGroupMember(ctx, types.GroupMember{
		GroupID: msg.GroupId,
		Member: &types.Member{
			Address: msg.ApplicantAddress,
			AddedAt: ctx.BlockTime(),
		},
	})

	//TODO: set applicant was joined
	//TODO: emit event
	return &types.MsgJoinGroupResponse{}, nil
}
