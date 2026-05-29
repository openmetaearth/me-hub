package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wdistri/types"
)

func (suite *KeeperTestSuite) TestAllocateBlockRewardDustLoss() {
	// Setup a context for a block height that triggers distribution
	ctx := suite.HelperNewContextWith(int64(types.OneDayTotalBlocks))

	// 3 regions with share 1 each (total share 3)
	regionShares := []int{1, 1, 1}
	addrs := suite.mockGetRegionI(ctx, regionShares...)

	// Total reward = 10 (not divisible by 3)
	// Without the fix:
	// - region 1 gets 3
	// - region 2 gets 3
	// - region 3 gets 3
	// Total distributed = 9, dust loss = 1
	// With the fix:
	// - region 1 gets 3
	// - region 2 gets 3
	// - region 3 gets 10 - 6 = 4 (remaining treasury balance)
	totalReward := int64(10)
	suite.SetMockGetBalance(ctx, sdk.NewInt(totalReward))

	// We expect the first two regions to get 3, and the last region to get 4
	wantReward := []coinAndAddr{
		{num: 3, addr: addrs[0]},
		{num: 3, addr: addrs[1]},
		{num: 4, addr: addrs[2]},
	}
	suite.setMockSendCoinsFromModuleToAccountExpect(ctx, wantReward...)

	// Trigger AllocateBlockRewardEveryday
	err := suite.App.DistrKeeper.AllocateBlockRewardEveryday(ctx, abci.RequestEndBlock{Height: ctx.BlockHeight()})
	suite.Require().NoError(err)
}
