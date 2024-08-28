package keeper

import (
	"context"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

// MsgServer is wrapper staking customParamsKeeper message server.
type MsgServer struct {
	stakingtypes.MsgServer
	*Keeper
}

var _ types.MsgServer = MsgServer{}

// NewMsgServerImpl returns an implementation of the staking wrapped MsgServer.
func NewMsgServerImpl(
	keeper *Keeper,
	stakingMsgSrv stakingtypes.MsgServer,
) MsgServer {
	return MsgServer{
		Keeper:    keeper,
		MsgServer: stakingMsgSrv,
	}
}

func (k MsgServer) WithdrawDelegatorReward(goCtx context.Context, msg *types.MsgWithdrawDelegatorReward) (*types.MsgWithdrawDelegatorRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("%v,address=%s", err, msg.DelegatorAddress)
	}

	amount, err := k.WithdrawDelegationRewards(ctx, delegatorAddress, sdk.ValAddress{})
	if err != nil {
		return nil, err
	}

	defer func() {
		for _, a := range amount {
			if a.Amount.IsInt64() {
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", "withdraw_reward"},
					float32(a.Amount.Int64()),
					[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
				)
			}
		}
	}()

	return &types.MsgWithdrawDelegatorRewardResponse{Amount: amount}, nil
}

//func (k MsgServer) UnmeidWithdrawDelegatorReward(goCtx context.Context, msg *types.MsgUnmeidWithdrawDelegatorReward) (*types.MsgUnmeidWithdrawDelegatorRewardResponse, error) {
//	ctx := sdk.UnwrapSDKContext(goCtx)
//
//	msgErr := msg.ValidateBasic()
//	if msgErr != nil {
//		return nil, msgErr
//	}
//
//	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
//	if err != nil {
//		return nil, err
//	}
//	amount, err := k.UnmeidWithdrawDelegationRewards(ctx, delegatorAddress, sdk.ValAddress{})
//	if err != nil {
//		return nil, err
//	}
//
//	defer func() {
//		for _, a := range amount {
//			if a.Amount.IsInt64() {
//				telemetry.SetGaugeWithLabels(
//					[]string{"tx", "msg", "withdraw_reward"},
//					float32(a.Amount.Int64()),
//					[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
//				)
//			}
//		}
//	}()
//
//	return &types.MsgUnmeidWithdrawDelegatorRewardResponse{Amount: amount}, nil
//}
