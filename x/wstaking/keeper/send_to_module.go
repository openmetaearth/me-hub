package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (k MsgServer) SendToModule(goCtx context.Context, msg *types.MsgSendToModule) (*types.MsgSendToModuleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.daoKeeper.IsGlobalDao(ctx, msg.Sender) {
		return nil, types.ErrCheckGlobalDao
	}

	err := k.bankKeeper.Extend().SendCoinsFromAccountToModuleWithTag(
		ctx,
		sdk.MustAccAddressFromBech32(msg.Sender),
		msg.Receiver,
		msg.Amount,
		fmt.Sprintf("SendToModule"),
	)
	if err != nil {
		return nil, err
	}
	return &types.MsgSendToModuleResponse{}, nil
}
