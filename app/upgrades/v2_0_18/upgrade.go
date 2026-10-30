package v2_0_18 //nolint:revive

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	appkeepers "github.com/openmetaearth/me-hub/app/keepers"
	"github.com/openmetaearth/me-hub/app/upgrades"
	"github.com/openmetaearth/me-hub/x/regiongov"
	regiongovtypes "github.com/openmetaearth/me-hub/x/regiongov/types"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v2.0.18.
// The upgrade adds the regiongov store and writes its default params.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ upgrades.BaseAppParamManager,
	keepers *appkeepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)
		logger.Info("upgrade starting...")

		if mod, ok := mm.Modules[regiongovtypes.ModuleName].(module.HasConsensusVersion); ok {
			fromVM[regiongovtypes.ModuleName] = mod.ConsensusVersion()
		}
		regiongov.InitGenesis(ctx, keepers.RegionGovKeeper, *regiongovtypes.DefaultGenesis())

		logger.Info("upgrade finished successfully.")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
