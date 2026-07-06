package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/rollapp/types"
)

func (k msgServer) UpdateRollapp(goCtx context.Context, msg *types.MsgUpdateRollapp) (*types.MsgUpdateRollappResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.RollappsEnabled(ctx) {
		return nil, types.ErrRollappsDisabled
	}

	rollapp, found := k.GetRollapp(ctx, msg.RollappId)
	if !found {
		return nil, types.ErrUnknownRollappID
	}

	if rollapp.Frozen {
		return nil, types.ErrRollappJailed
	}

	if msg.Creator != rollapp.Creator {
		return nil, types.ErrUnauthorizedRollappCreator
	}

	if msg.MaxSequencers != 0 {
		rollapp.MaxSequencers = msg.MaxSequencers
	}
	if msg.ChannelId != "" {
		rollapp.ChannelId = msg.ChannelId
	}
	if len(msg.PermissionedAddresses) != 0 {
		rollapp.PermissionedAddresses = msg.PermissionedAddresses
	}

	k.SetRollapp(ctx, rollapp)
	return &types.MsgUpdateRollappResponse{}, nil
}
