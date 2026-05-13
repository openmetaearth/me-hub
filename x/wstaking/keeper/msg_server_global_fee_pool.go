package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (k MsgServer) WithdrawFromGlobalDaoFeePool(goCtx context.Context, msg *types.MsgWithdrawFromGlobalDaoFeePool) (*types.MsgWithdrawFromGlobalDaoFeePoolResp, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.daoKeeper.IsGlobalDao(ctx, msg.Withdrawer) {
		return nil, types.ErrCheckGlobalDao
	}

	toAddr, err := sdk.AccAddressFromBech32(msg.Withdrawer)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrUnknownAccount, "receiver account %s format error %s", msg.Withdrawer, err)
	}

	fromAddr := k.daoKeeper.GetGlobalDaoFeePoolAddr(ctx)
	err = k.bankKeeper.Extend().SendCoinsWithTag(
		ctx,
		fromAddr,
		toAddr,
		msg.Amount,
		"WithdrawFromGlobalDaoFeePool_SendCoinsFromGlobalDaoFeePoolToUserAccount",
	)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "retrieve fee from global fee pool error: from(%s), to (%s)", fromAddr, toAddr.String())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeWithdrawFromGlobalDaoFeePool,
			sdk.NewAttribute(sdk.AttributeKeySender, fromAddr.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, toAddr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
	)

	return &types.MsgWithdrawFromGlobalDaoFeePoolResp{}, nil
}
