package keeper

import (
	"context"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

func (k MsgServer) FixedDeposit(ctx context.Context, deposit *types.MsgFixedDeposit) (*types.MsgFixedDepositResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k MsgServer) FixedWithdraw(ctx context.Context, withdraw *types.MsgFixedWithdraw) (*types.MsgFixedWithdrawResponse, error) {
	//TODO implement me
	panic("implement me")
}
