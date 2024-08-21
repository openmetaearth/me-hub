package keeper

import (
	"context"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

func (k MsgServer) CurrentDeposit(ctx context.Context, deposit *types.MsgCurrentDeposit) (*types.MsgCurrentDepositResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k MsgServer) CurrentWithdraw(ctx context.Context, withdraw *types.MsgCurrentWithdraw) (*types.MsgCurrentWithdrawResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k MsgServer) FixedDeposit(ctx context.Context, deposit *types.MsgFixedDeposit) (*types.MsgFixedDepositResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k MsgServer) FixedWithdraw(ctx context.Context, withdraw *types.MsgFixedWithdraw) (*types.MsgFixedWithdrawResponse, error) {
	//TODO implement me
	panic("implement me")
}
