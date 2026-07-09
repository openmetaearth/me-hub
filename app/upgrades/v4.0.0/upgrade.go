package v4_0_0

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/openmetaearth/me-hub/app/upgrades"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v4.0.0.
//
// This upgrade adds the lightclient module store key introduced by the
// Dymension Hub v4 settlement stack. Full params→collections migrations for
// settlement modules are intentionally deferred (TODO) and can be layered on
// later once mainnet state shape is confirmed.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ *upgrades.UpgradeKeepers,
) upgradetypes.UpgradeHandler {
	return func(goCtx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(goCtx)
		logger := ctx.Logger().With("upgrade", UpgradeName)

		// NOTE: DO NOT PUT ANY STATE CHANGES BEFORE RunMigrations().
		migrations, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}

		// TODO: migrate settlement module params from x/params subspaces into
		// collections-backed stores where needed (rollapp/sequencer/delayedack/eibc/lightclient).

		logger.Info("added lightclient store key; deferred full collections migrations")
		logger.Info("upgrade finished successfully.")
		return migrations, nil
	}
}
