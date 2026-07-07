package v3_0_0

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/openmetaearth/me-hub/app/upgrades"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v3.0.0
// This upgrade migrates the chain from Cosmos SDK v0.47 to v0.50.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *upgrades.UpgradeKeepers,
) upgradetypes.UpgradeHandler {
	return func(goCtx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(goCtx)
		logger := ctx.Logger().With("upgrade", UpgradeName)

		// Run all module migrations first.
		// NOTE: DO NOT PUT ANY STATE CHANGES BEFORE RunMigrations().
		migrations, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}

		// Migrate wstaking validators from V1 (v47-era) protobuf format to V2 (v50-era).
		// The Validator protobuf schema changed field numbers between the two SDK versions;
		// this step rewrites every on-chain validator record to the new layout.
		if err := keepers.StakingKeeper.MigrateValidatorsFromV1(ctx); err != nil {
			return nil, fmt.Errorf("failed to migrate wstaking validators from V1: %w", err)
		}
		logger.Info("successfully migrated wstaking validators to V2 format")

		logger.Info("upgrade finished successfully.")
		return migrations, nil
	}
}
