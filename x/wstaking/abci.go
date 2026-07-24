package wstaking

import (
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/openmetaearth/me-hub/x/wstaking/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wstaking/keeper"
)

// BeginBlocker will persist the current header and validator set as a historical entry
// and prune the oldest entry based on the HistoricalEntries parameter
func BeginBlock(ctx sdk.Context, k *keeper.Keeper) {
	totalRewardsPerBlockTemp := k.GetPerBlockMintCoinAmount(ctx)
	totalRewardsPerBlock := sdkmath.NewIntFromBigInt(&totalRewardsPerBlockTemp)
	regions := k.GetAllRegion(ctx)

	for _, region := range regions {
		rewards, _ := k.Calculate(ctx, sdkmath.LegacyNewDecFromInt(totalRewardsPerBlock), region.DelegateAmount) // rate.MulInt(totalRewardsPerBlock.Mul(region.DelegateAmount)).Mul(sdk.NewDecWithPrec(1, sdk.MEExponent))
		region.DelegateInterest = region.DelegateInterest.Add(rewards)
		k.SetRegion(ctx, region)
	}

	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	k.TrackHistoricalInfo(ctx)
	// Initialize region cache once on first block; subsequent calls are no-ops due to sync.Once.
	k.InitRegionCache(ctx)
}

func EndBlock(ctx sdk.Context, k *keeper.Keeper) []abci.ValidatorUpdate {
	k.ChangeDelegationValidator(ctx)
	updates := k.BlockValidatorUpdates(ctx)
	// Refresh region cache at end of block so the next block's TXs see up-to-date data.
	k.SetRegionsCache(ctx, k.GetAllRegion(ctx))
	return updates
}
