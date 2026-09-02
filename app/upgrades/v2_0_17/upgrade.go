package v2_0_17 //nolint:revive

import (
	"fmt"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	appkeepers "github.com/openmetaearth/me-hub/app/keepers"
	"github.com/openmetaearth/me-hub/app/upgrades"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v2.0.17
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ upgrades.BaseAppParamManager,
	keepers *appkeepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)
		logger.Info("upgrade starting...")

		stakingParams := keepers.StakingKeeper.GetParams(ctx)
		stakingParams.UnbondingTime = time.Hour * 24 * 7 * 3 // 3 weeks
		if err := keepers.StakingKeeper.SetParams(ctx, stakingParams); err != nil {
			panic("failed to set Staking params: " + err.Error())
		}

		bscParams := keepers.BscKeeper.GetParams(ctx)
		bscParams.MaxDelegate = sdkmath.NewInt(200_000_000)
		if err := keepers.BscKeeper.SetParams(ctx, &bscParams); err != nil {
			return nil, fmt.Errorf("failed to set BSC max send to external amount: %w", err)
		}

		tronParams := keepers.TronKeeper.GetParams(ctx)
		tronParams.MaxDelegate = sdkmath.NewInt(200_000_000)
		if err := keepers.TronKeeper.SetParams(ctx, &tronParams); err != nil {
			return nil, fmt.Errorf("failed to set Tron max send to external amount: %w", err)
		}

		logger.Info("upgrade finished successfully.")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
