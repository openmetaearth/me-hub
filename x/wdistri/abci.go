package wdistri

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wdistri/keeper"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	// record the proposer for when we payout on the next block
	consAddr := sdk.ConsAddress(req.Header.ProposerAddress)
	k.SetPreviousProposerConsAddr(ctx, consAddr)
}

func EndBlock(ctx sdk.Context, req abci.RequestEndBlock, k keeper.Keeper) {
	if err := k.AllocateBlockRewardEveryday(ctx, req); err != nil {
		ctx.Logger().Error("AllocateBlockRewardEveryday", "err", err)
	}
}
