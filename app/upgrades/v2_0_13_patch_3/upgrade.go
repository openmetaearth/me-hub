package v2_0_13_patch_3

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	appkeepers "github.com/st-chain/me-hub/app/keepers"
	"github.com/st-chain/me-hub/app/upgrades"
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

		logger.Info("3. set batch timeout to 24 hours for bsc and tron.")
		params := keepers.TronKeeper.GetParams(ctx)
		params.AverageBlockTime = 6000
		params.ExternalBatchTimeout = 86400000
		err := keepers.TronKeeper.SetParams(ctx, &params)
		if err != nil {
			panic(fmt.Sprintf("failed to set bsc params during upgrade: %v", err))
		}

		params = keepers.BscKeeper.GetParams(ctx)
		params.AverageBlockTime = 6000
		params.ExternalBatchTimeout = 86400000
		err = keepers.BscKeeper.SetParams(ctx, &params)
		if err != nil {
			panic(fmt.Sprintf("failed to set bsc params during upgrade: %v", err))
		}

		logger.Info("upgrade finished successfully.")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
