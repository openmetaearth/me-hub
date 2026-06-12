package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/rollapp/types"
)

func (k msgServer) UpdateRollapp(goCtx context.Context, msg *types.MsgUpdateRollapp) (*types.MsgUpdateRollappResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	rollapp, found := k.GetRollapp(ctx, msg.RollappId)
	if !found {
		return nil, types.ErrUnknownRollappID
	}
	// Update authorization is tied to the stored rollapp owner. The deployer
	// whitelist can govern creation, but it must not let a whitelisted account
	// mutate another creator's rollapp after ownership has been assigned.
	if msg.Creator != rollapp.Creator && !k.daoKeeper.IsDao(ctx, msg.Creator) {
		return nil, types.ErrUnauthorizedRollappCreator.Wrapf("only rollapp creator or DAO can update rollapp %s", msg.RollappId)
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
