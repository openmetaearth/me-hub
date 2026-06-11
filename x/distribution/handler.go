package handler

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/distribution/types"
)

// NewHandler returns a handler for distribution messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx context.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case *types.MsgClaimRewards:
			return handleMsgClaimRewards(ctx, k, msg)
		default:
			return nil, fmt.Errorf("unrecognized distribution message type: %T", msg)
		}
	}
}

func handleMsgClaimRewards(ctx context.Context, k Keeper, msg *types.MsgClaimRewards) (*sdk.Result, error) {
	// Claim the rewards
	rewards, err := k.ClaimRewards(ctx, msg.ValidatorAddr)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}