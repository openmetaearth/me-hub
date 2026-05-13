package wdistri

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wdistri/keeper"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	// record the proposer for when we payout on the next block
	consAddr := sdk.ConsAddress(ctx.BlockHeader().ProposerAddress)
	k.SetPreviousProposerConsAddr(ctx, consAddr)
}

func EndBlock(ctx sdk.Context, k keeper.Keeper) {
	if err := k.AllocateBlockRewardEveryday(ctx); err != nil {
		ctx.Logger().Error("AllocateBlockRewardEveryday", "err", err)
	}
}
