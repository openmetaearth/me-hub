package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/dymensionxyz/dymension/v3/x/rollup/types"

	//"github.com/dymensionxyz/dymension/v3/x/rollup/types"

	//"github.com/dymensionxyz/dymension/v3/x/rollup/types"
	"me-hub/x/rollup/types"
)

// NewHandler creates an sdk.Handler for all the staking type messages
func NewHandler(k Keeper, bk bankKeeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgStake:
			return handleMsgStake(ctx, k, bk, msg)
		case types.MsgUnstake:
			return handleMsgUnstake(ctx, k, bk, msg)
		default:
			return nil, fmt.Errorf("unrecognized staking message type: %T", msg)

		}
	}
}

func handleMsgStake(ctx sdk.Context, k Keeper, bk bankKeeper.Keeper, msg types.MsgStake) (*sdk.Result, error) {
	err := bk.SendCoins(ctx, msg.Delegator, k.GetModuleAddress(), sdk.NewCoins(msg.Amount))
	if err != nil {
		return nil, err
	}

	k.StakeTokens(ctx, msg.Delegator, msg.Amount)

	return &sdk.Result{}, nil
}

func handleMsgUnstake(ctx sdk.Context, k Keeper, bk bankKeeper.Keeper, msg types.MsgUnstake) (*sdk.Result, error) {
	err := k.UnstakeTokens(ctx, msg.Delegator, msg.Amount)
	if err != nil {
		return nil, err
	}

	err = bk.SendCoins(ctx, k.GetModuleAddress(), msg.Delegator, sdk.NewCoins(msg.Amount))
	if err != nil {
		return nil, err
	}

	return &sdk.Result{}, nil
}
