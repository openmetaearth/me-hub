package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/rollapp/types"
	wstakingtypes "github.com/st-chain/me-hub/x/wstaking/types"
)

func (k MsgServer) SendToModule(goCtx context.Context, msg *wstakingtypes.MsgSendToModule) (*wstakingtypes.MsgSendToModuleResponse, error) {
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

	return &wstakingtypes.MsgSendToModuleResponse{}, nil
}
