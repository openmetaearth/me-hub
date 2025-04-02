package keeper

import (
	"context"
	sdkerrors "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	wbanktypes "github.com/st-chain/me-hub/x/wbank/types"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

func (k MsgServer) WithdrawFromTreasury(goCtx context.Context, msg *types.MsgWithdrawFromTreasury) (*types.MsgWithdrawFromTreasuryResp, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.daoKeeper.IsGlobalDao(ctx, msg.Withdrawer) {
		return nil, types.ErrCheckGlobalDao
	}

	err := k.bankKeeper.Extend().SendCoinsFromModuleToAccountWithTag(
		ctx,
		wbanktypes.TreasuryPoolName,
		sdk.MustAccAddressFromBech32(msg.Receiver),
		msg.Amount,
		fmt.Sprintf("WithdrawFromTreasury_SendCoinsFromTreasuryToAccount"),
	)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrWithdrawFromTreasury, "%v", err)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeWithdrawFromTreasury,
			sdk.NewAttribute(types.AttributeKeyWithdrawer, msg.Withdrawer),
			sdk.NewAttribute(sdk.AttributeKeySender, k.authKeeper.GetModuleAddress(wbanktypes.TreasuryPoolName).String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
	)
	return &types.MsgWithdrawFromTreasuryResp{}, nil
}
