package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (k MsgServer) SendToModule(goCtx context.Context, msg *types.MsgSendToModule) (*types.MsgSendToModuleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.daoKeeper.IsGlobalDao(ctx, msg.Sender) {
		return nil, types.ErrCheckGlobalDao
	}
	if !types.IsAllowedSendToModuleTarget(msg.Receiver) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "module %q is not an allowed SendToModule target", msg.Receiver)
	}

	err := k.bankKeeper.Extend().SendCoinsFromAccountToModuleWithTag(
		ctx,
		sdk.MustAccAddressFromBech32(msg.Sender),
		msg.Receiver,
		msg.Amount,
		"SendToModule",
	)
	if err != nil {
		return nil, err
	}
	return &types.MsgSendToModuleResponse{}, nil
}
