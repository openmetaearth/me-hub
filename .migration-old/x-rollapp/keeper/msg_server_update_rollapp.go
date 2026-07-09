package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/rollapp/types"
)

func (k msgServer) UpdateRollapp(goCtx context.Context, msg *types.MsgUpdateRollapp) (*types.MsgUpdateRollappResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check to see if there is an active whitelist
	if whitelist := k.DeployerWhitelist(ctx); len(whitelist) > 0 {
		if !k.IsAddressInDeployerWhiteList(ctx, msg.Creator) {
			return nil, types.ErrUnauthorizedRollappCreator
		}
	}

	rollapp, found := k.GetRollapp(ctx, msg.RollappId)
	if !found {
		return nil, types.ErrUnknownRollappID
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
