package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (ah AnteHandler) CheckTx(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
	// ...

	for _, msg := range tx.GetMsgs() {
		switch msg := msg.(type) {
		case *types.MsgWithdrawDelegatorReward:
			// Calculate total gas cost
			totalGasCost := tx.GetGas() * ctx.MinGasPrices().AmountOf("umec")

			// Calculate total rewards to claim
			totalRewardsToClaim := ctx.KVStore(keeper.StoreKey).Get(types.GetDelegatorWithdrawableRewardsKey(msg.DelegatorAddress, msg.ValidatorAddress))

			// Check if total gas cost is greater than total rewards to claim
			if totalGasCost.GT(totalRewardsToClaim) {
				return ctx, sdk.ErrInvalidRequest.Wrap("total gas cost is greater than total rewards to claim")
			}
		}
	}

	// ...
}