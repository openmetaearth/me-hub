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

func (k MsgServer) NewFixedDepositCfg(ctx context.Context, cfg *types.MsgFixedDepositCfg) (*types.MsgFixedDepositCfgResp, error) {
	//TODO implement me
	panic("implement me")
}

func (k MsgServer) RemoveFixedDepositCfg(ctx context.Context, cfg *types.MsgRemoveFixedDepositCfg) (*types.MsgRemoveFixedDepositCfgResp, error) {
	//TODO implement me
	panic("implement me")
}

func (k MsgServer) SetFixedDepositCfgStatus(ctx context.Context, status *types.MsgSetFixedDepositCfgStatus) (*types.MsgSetFixedDepositCfgStatusResp, error) {
	//TODO implement me
	panic("implement me")
}

func (k MsgServer) SetFixedDepositCfgRate(ctx context.Context, rate *types.MsgSetFixedDepositCfgRate) (*types.MsgSetFixedDepositCfgRateResp, error) {
	//TODO implement me
	panic("implement me")
}
