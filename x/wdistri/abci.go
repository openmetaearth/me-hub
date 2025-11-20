package wdistri

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/wdistri/keeper"
)

func EndBlock(ctx sdk.Context, req abci.RequestEndBlock, k keeper.Keeper) {
	if err := k.AllocateBlockRewardEveryday(ctx, req); err != nil {
		ctx.Logger().Error("AllocateBlockRewardEveryday", "err", err)
	}
}
