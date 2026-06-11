package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

// MsgServer defines the Msg server
type MsgServer struct {
	Keeper
}

// NewMsgServer returns a new MsgServer instance
func NewMsgServer(keeper Keeper) MsgServer {
	return MsgServer{Keeper: keeper}
}

// VerifyMessage handles the VerifyMessage message
func (s MsgServer) VerifyMessage(goCtx context.Context, msg *types.Msg) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := s.Keeper.VerifyMessage(ctx, *msg); err != nil {
		return nil, err
	}
	return &types.MsgResponse{}, nil
}