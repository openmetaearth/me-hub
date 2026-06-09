package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"mechain.wstaking.v1/types"
)

// Keeper defines the wstaking module keeper
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      codec.BinaryCodec
}

// NewKeeper returns a new wstaking Keeper instance
func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey) *Keeper {
	return &Keeper{
		storeKey: key,
		cdc:      cdc,
	}
}

// WithdrawDelegatorReward handles the MsgWithdrawDelegatorReward request
func (k Keeper) WithdrawDelegatorReward(ctx sdk.Context, msg *types.MsgWithdrawDelegatorReward) (*types.MsgWithdrawDelegatorRewardResponse, error) {
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