package v2_0_2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/st-chain/me-hub/app/keepers"
	"github.com/st-chain/me-hub/app/upgrades"
	wstakingkeeper "github.com/st-chain/me-hub/x/wstaking/keeper"
	wstakingtypes "github.com/st-chain/me-hub/x/wstaking/types"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v4
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ upgrades.BaseAppParamManager,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)
		logger.Info("upgrade starting...")
		for n, m := range mm.Modules {
			if mod, ok := m.(module.HasConsensusVersion); ok {
				fromVM[n] = mod.ConsensusVersion()
			}
		}

		migrateRegions(ctx, keepers.StakingKeeper)
		logger.Info("upgrade finished.")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func migrateRegions(ctx sdk.Context, k *wstakingkeeper.Keeper) {
	validators := k.GetAllValidators(ctx)
	regions := k.GetAllRegion(ctx)

	regionIds := []string{}
	regionsMap := make(map[string]wstakingtypes.Region)
	for _, region := range regions {
		regionIds = append(regionIds, region.RegionId)

		region.RegionTreasureAddr = k.CreateRegionAccount(ctx, wstakingtypes.RegionAccountTypeBase, region.RegionId).String()
		region.DepositInterestAddr = k.CreateRegionAccount(ctx, wstakingtypes.RegionAccountTypeDepositInterest, region.RegionId).String()
		region.RegionShare = sdk.ZeroInt()

		regionsMap[region.RegionId] = region
	}

	for index, validator := range validators {

		validator.Description.RegionID = regionIds[index]
		k.SetValidator(ctx, validator)

		if region, ok := regionsMap[validator.Description.RegionID]; ok {
			region.OperatorAddress = validator.OperatorAddress
			region.RegionShare = validator.Tokens
			regionsMap[validator.Description.RegionID] = region
		}
	}

	for _, region := range regionsMap {
		k.SetRegion(ctx, region)
	}
}
