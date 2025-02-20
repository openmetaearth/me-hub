package keeper

import (
	"context"
	"encoding/json"
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

func (k msgServer) UpdateDao(goCtx context.Context, msg *types.MsgUpdateDao) (*types.MsgUpdateDaoResponse, error) {
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

	oldByte, err := json.Marshal(oldAddresses)
	if err != nil {
		panic(err)
	}

	newByte, err := json.Marshal(msg.DaoAddresses)
	if err != nil {
		panic(err)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDaoUpdated,
			sdk.NewAttribute(types.AttributeKeyLastDaoAddresses, string(oldByte)),
			sdk.NewAttribute(types.AttributeKeyNewDaoAddresses, string(newByte)),
		),
	)

	return &types.MsgUpdateDaoResponse{}, nil
}
