package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/utils/uevent"

	"github.com/openmetaearth/me-hub/x/sequencer/types"
)

// UpdateWhitelistedRelayers defines a method for updating the sequencer's whitelisted relater list.
func (k msgServer) UpdateWhitelistedRelayers(goCtx context.Context, msg *types.MsgUpdateWhitelistedRelayers) (*types.MsgUpdateWhitelistedRelayersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	seq, err := k.RealSequencer(ctx, msg.Creator)
	if err != nil {
		return nil, err
	}
	defer func() {
		k.SetSequencer(ctx, seq)
	}()

	seq.SetWhitelistedRelayers(msg.Relayers)

	err = uevent.EmitTypedEvent(ctx, &types.EventUpdateWhitelistedRelayers{
		Creator:  msg.Creator,
		Relayers: msg.Relayers,
	})
	if err != nil {
		return nil, fmt.Errorf("emit event: %w", err)
	}

	return &types.MsgUpdateWhitelistedRelayersResponse{}, nil
}
