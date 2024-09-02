package keeper

import (
	"context"
	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

func (k Keeper) GetGlobalAdminFeePoolAddr(ctx sdk.Context) sdk.AccAddress {
	addr := sdk.AccAddress(crypto.AddressHash([]byte(types.GlobalAdminFeePool)))
	account := k.AuthKeeper.GetAccount(ctx, addr)
	if account == nil {
		k.AuthKeeper.SetAccount(ctx, k.AuthKeeper.NewAccountWithAddress(ctx, addr))
		return addr
	}
	return account.GetAddress()
}

func (k MsgServer) WithdrawFromGlobalDaoFeePool(goCtx context.Context, msg *types.MsgWithdrawFromGlobalDaoFeePool) (*types.MsgWithdrawFromGlobalDaoFeePoolResp, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.DaoKeeper.IsGlobalDao(ctx, msg.Withdrawer) {
		return nil, types.ErrCheckGlobalDao
	}

	toAddr, err := sdk.AccAddressFromBech32(msg.Withdrawer)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrUnknownAccount, "receiver account %s format error %s", msg.Withdrawer, err)
	}

	fromAddr := k.GetGlobalAdminFeePoolAddr(ctx)
	err = k.BankKeeper.SendCoins(
		ctx,
		fromAddr,
		toAddr,
		msg.Amount)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "retrieve fee from global fee pool error: from(%s), to (%s)", fromAddr, toAddr.String())
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
