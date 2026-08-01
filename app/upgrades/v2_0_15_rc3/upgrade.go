package v2_0_15_rc3 //nolint:revive

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
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

		maxSendToExternalUSDAmount := sdkmath.NewInt(10_000_000_000)
		bscParams := keepers.BscKeeper.GetParams(ctx)
		bscParams.MaxSendToExternalUSDAmount = maxSendToExternalUSDAmount
		if err := keepers.BscKeeper.SetParams(ctx, &bscParams); err != nil {
			return nil, fmt.Errorf("failed to set BSC max send to external amount: %w", err)
		}

		tronParams := keepers.TronKeeper.GetParams(ctx)
		tronParams.MaxSendToExternalUSDAmount = maxSendToExternalUSDAmount
		if err := keepers.TronKeeper.SetParams(ctx, &tronParams); err != nil {
			return nil, fmt.Errorf("failed to set Tron max send to external amount: %w", err)
		}

		logger.Info("upgrade finished successfully.")

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
