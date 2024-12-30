package v2_0_3

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	appkeepers "github.com/st-chain/me-hub/app/keepers"
	"github.com/st-chain/me-hub/app/upgrades"
	streamermoduletypes "github.com/st-chain/me-hub/x/streamer/types"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v4
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ upgrades.BaseAppParamManager,
	keepers *appkeepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)
		logger.Info("upgrade starting...")
		for n, m := range mm.Modules {
			if mod, ok := m.(module.HasConsensusVersion); ok {
				fromVM[n] = mod.ConsensusVersion()
			}
		}

		// create a new module account
		macc := types.NewEmptyModuleAccount(streamermoduletypes.ModuleName)
		maccI := (keepers.AccountKeeper.NewAccount(ctx, macc)).(types.ModuleAccountI) // set the account number
		keepers.AccountKeeper.SetModuleAccount(ctx, maccI)

		//upgrade for migrate did info
		logger.Info("upgrade finished.")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
