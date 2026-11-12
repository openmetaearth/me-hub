package v2_0_18_rc1 //nolint:revive

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	appkeepers "github.com/openmetaearth/me-hub/app/keepers"
	"github.com/openmetaearth/me-hub/app/upgrades"
)

// CreateUpgradeHandler upgrades v2.0.18 without reinitializing regiongov.
// Only proposal required-operator encodings change; staking addresses and elector
// snapshots remain untouched.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ upgrades.BaseAppParamManager,
	keepers *appkeepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)
		logger.Info("upgrade starting...")

		versionMap, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}

		if err := keepers.RegionGovKeeper.MigrateRequiredOperatorAccounts(ctx); err != nil {
			return nil, err
		}

		logger.Info("upgrade finished successfully.")
		return versionMap, nil
	}
}
