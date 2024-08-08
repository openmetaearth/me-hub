package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/dao/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) UpdateGlobalDao(goCtx context.Context, msg *types.MsgUpdateGlobalDao) (*types.MsgUpdateGlobalDaoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	lastAddress := k.GetGlobalDao(ctx)
	newAddress, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, err
	}

	if lastAddress.Equals(newAddress) {
		return nil, types.ErrLastAddressEqualNewAddress
	}

	isGlobalDao := k.IsGlobalDao(ctx, msg.Creator)
	if !isGlobalDao {
		return nil, types.ErrCreatorNotDao
	}

	k.SetGlobalDao(ctx, newAddress)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAdminUpdated,
			sdk.NewAttribute(types.AttributeKeyLastAdmin, lastAddress.String()),
			sdk.NewAttribute(types.AttributeKeyCurrentAdmin, msg.Address),
		),
	)

	return &types.MsgUpdateGlobalDaoResponse{}, nil
}
