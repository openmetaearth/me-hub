package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"mechain.wstaking.v1/types"
)

// MsgWithdrawDelegatorReward defines the MsgWithdrawDelegatorReward request type.
func (k msgServer) WithdrawDelegatorReward(goCtx context.Context, msg *types.MsgWithdrawDelegatorReward) (*types.MsgWithdrawDelegatorRewardResponse, error) {
	// ...

	// Calculate gas fee with a cap to prevent excessive fees
	gasFeeCap := sdk.NewInt(100) // 0.1% of the transaction value
	gasFee := sdk.NewInt(0)
	if msg.Reward.Amount.GT(gasFeeCap) {
		gasFee = msg.Reward.Amount.Quo(sdk.NewInt(1000)) // 0.1% of the transaction value
	} else {
		gasFee = msg.Reward.Amount
	}

	// ...
}