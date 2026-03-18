package v2_0_13_patch_4

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	appkeepers "github.com/st-chain/me-hub/app/keepers"
	"github.com/st-chain/me-hub/app/upgrades"
	bsctypes "github.com/st-chain/me-hub/x/bsc/types"
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

		logger.Info("1. upgrade for x/gravity module, set params")
		keepers.TronKeeper.ClearGenesis(ctx)
		keepers.BscKeeper.ClearGenesis(ctx)

		accounts := keepers.AccountKeeper.GetAllAccounts(ctx)
		for _, acc := range accounts {
			balances := keepers.BankKeeper.GetAllBalances(ctx, acc.GetAddress())
			for _, balance := range balances {
				if balance.Denom == "uusdc" || balance.Denom == "uusdt" {
					err := keepers.BankKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), bsctypes.ModuleName, sdk.NewCoins(balance))
					if err != nil {
						panic(fmt.Errorf("failed to delegate coins to gravity module account: %s", err.Error()))
					}
					err = keepers.BankKeeper.BurnCoins(ctx, bsctypes.ModuleName, sdk.NewCoins(balance))
					if err != nil {
						panic(fmt.Errorf("failed to burn coins from gravity module account: %s", err.Error()))
					}
				}
			}
		}
		logger.Info("upgrade finished successfully.")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
