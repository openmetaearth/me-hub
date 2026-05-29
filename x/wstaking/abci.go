package wstaking

import (
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wstaking/keeper"
)

// BeginBlocker will persist the current header and validator set as a historical entry
// and prune the oldest entry based on the HistoricalEntries parameter
func BeginBlock(ctx sdk.Context, k *keeper.Keeper) {
	totalRewardsPerBlockTemp := k.GetPerBlockMintCoinAmount(ctx)
	totalRewardsPerBlock := sdk.NewIntFromBigInt(&totalRewardsPerBlockTemp)
	regions := k.GetAllRegion(ctx)

	for _, region := range regions {
		rewards, err := k.Calculate(ctx, sdk.NewDecFromInt(totalRewardsPerBlock), region.DelegateAmount)
		if err != nil {
			k.Logger(ctx).Error("failed to calculate rewards for region", "region", region.RegionId, "error", err)
			continue
		}
		region.DelegateInterest = region.DelegateInterest.Add(rewards)
		k.SetRegion(ctx, region)
	}

	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	k.TrackHistoricalInfo(ctx)
}

func EndBlock(ctx sdk.Context, k *keeper.Keeper) []abci.ValidatorUpdate {
	k.ChangeDelegationValidator(ctx)
	return k.BlockValidatorUpdates(ctx)
}
