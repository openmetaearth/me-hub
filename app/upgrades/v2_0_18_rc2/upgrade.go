package v2_0_18_rc2 //nolint:revive

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	appkeepers "github.com/openmetaearth/me-hub/app/keepers"
	"github.com/openmetaearth/me-hub/app/upgrades"
)

// CreateUpgradeHandler creates an incremental upgrade handler for v2.0.18.rc2.
// Existing module state and parameters are preserved by in-place migrations.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ upgrades.BaseAppParamManager,
	_ *appkeepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)
		logger.Info("upgrade starting...")

		updatedVM, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}

		logger.Info("upgrade finished successfully.")
		return updatedVM, nil
	}
}
