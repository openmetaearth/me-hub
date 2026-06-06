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

	updatedRollapp := rollapp
	if msg.MaxSequencers != 0 {
		updatedRollapp.MaxSequencers = msg.MaxSequencers
	}
	if msg.ChannelId != "" {
		updatedRollapp.ChannelId = msg.ChannelId
	}
	if len(msg.PermissionedAddresses) != 0 {
		updatedRollapp.PermissionedAddresses = msg.PermissionedAddresses
	}

	if err := updatedRollapp.ValidateBasic(); err != nil {
		return nil, err
	}

	k.SetRollapp(ctx, updatedRollapp)
	return &types.MsgUpdateRollappResponse{}, nil
}
