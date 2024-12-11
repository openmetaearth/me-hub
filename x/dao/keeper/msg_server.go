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

	isGlobalDao := k.IsGlobalDao(ctx, msg.Creator)
	if !isGlobalDao {
		return nil, types.ErrCreatorNotDao
	}

	oldAddresses, found := k.GetDaoAddresses(ctx)
	if !found {
		return nil, types.ErrNotFound
	}

	k.SetDaoAddresses(ctx, msg.DaoAddresses)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDaoUpdated,
			sdk.NewAttribute(types.AttributeKeyLastGlobalDao, oldAddresses.GlobalDao),
			sdk.NewAttribute(types.AttributeKeyCurrentGlobalDao, msg.DaoAddresses.GlobalDao),

			sdk.NewAttribute(types.AttributeKeyLastMeidDao, oldAddresses.MeidDao),
			sdk.NewAttribute(types.AttributeKeyCurrentMeidDao, msg.DaoAddresses.MeidDao),

			sdk.NewAttribute(types.AttributeKeyLastDevOperator, oldAddresses.DevOperator),
			sdk.NewAttribute(types.AttributeKeyCurrentDevOperator, msg.DaoAddresses.DevOperator),

			sdk.NewAttribute(types.AttributeKeyLastAirdrop, oldAddresses.AirdropAddress),
			sdk.NewAttribute(types.AttributeKeyCurrentAirdrop, msg.DaoAddresses.AirdropAddress),
		),
	)

	return &types.MsgUpdateGlobalDaoResponse{}, nil
}
