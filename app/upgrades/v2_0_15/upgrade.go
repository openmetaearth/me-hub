package v2_0_15 //nolint:revive

import (
	"context"
	"time"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

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
	return func(goCtx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(goCtx)
		logger := ctx.Logger().With("upgrade", UpgradeName)
		logger.Info("upgrade starting...")

		stakingParams, err := keepers.StakingKeeper.GetParams(ctx)
		if err != nil {
			return nil, err
		}
		stakingParams.UnbondingTime = time.Hour * 24 * 7 * 3 // 3 weeks
		if err := keepers.StakingKeeper.SetParams(ctx, stakingParams); err != nil {
			panic("failed to set Staking params: " + err.Error())
		}

		logger.Info("upgrade finished successfully.")
		return mm.RunMigrations(goCtx, configurator, fromVM)
	}
}
