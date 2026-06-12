package v2_0_14_patch_2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	appkeepers "github.com/openmetaearth/me-hub/app/keepers"
	"github.com/openmetaearth/me-hub/app/upgrades"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v2.0.13
// This upgrade initializes the Gravity bridge module for BSC and Tron
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ upgrades.BaseAppParamManager,
	keepers *appkeepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)
		logger.Info("upgrade starting...")

		// Initialize consensus versions for all modules
		for n, m := range mm.Modules {
			if mod, ok := m.(module.HasConsensusVersion); ok {
				fromVM[n] = mod.ConsensusVersion()
			}
		}

		evmParams := keepers.EvmKeeper.GetParams(ctx)
		evmParams.EvmDenom = "umec"
		if err := keepers.EvmKeeper.SetParams(ctx, evmParams); err != nil {
			panic("failed to set EVM params: " + err.Error())
		}

		feemarketParams := keepers.FeeMarketKeeper.GetParams(ctx)
		feemarketParams.BaseFee = sdk.NewInt(4700000000)
		feemarketParams.MinGasPrice = sdk.NewDecFromInt(sdk.NewInt(4700000000))
		if err := keepers.FeeMarketKeeper.SetParams(ctx, feemarketParams); err != nil {
			panic("failed to set FeeMarket params: " + err.Error())
		}

		logger.Info("upgrade finished successfully.")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
