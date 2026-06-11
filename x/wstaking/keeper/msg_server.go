package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (k msgServer) WithdrawDelegatorReward(goCtx context.Context, msg *types.MsgWithdrawDelegatorReward) (*types.MsgWithdrawDelegatorRewardResponse, error) {
	// ...
	
	// Calculate total gas cost
	totalGasCost := msg.Gas * k.stakingKeeper.GetParams(ctx).MinGasPrice

	// Calculate total rewards to claim
	totalRewardsToClaim := k.GetDelegatorWithdrawableRewards(ctx, msg.DelegatorAddress, msg.ValidatorAddress)

	// Check if total gas cost is greater than total rewards to claim
	if totalGasCost.GT(totalRewardsToClaim) {
		return nil, sdk.ErrInvalidRequest.Wrap("total gas cost is greater than total rewards to claim")
	}

	// ...
}