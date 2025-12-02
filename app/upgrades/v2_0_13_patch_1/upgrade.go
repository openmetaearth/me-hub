package v2_0_13

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	appkeepers "github.com/st-chain/me-hub/app/keepers"
	"github.com/st-chain/me-hub/app/upgrades"
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

		logger.Info("1. migrate denoms for existing bridge tokens...")
		//denomMap := make(map[string]struct{})
		//bscBridgeTokens := []types.BridgeToken{}
		//keepers.BscKeeper.IterateBridgeTokenByDenom(ctx, func(bt *types.BridgeToken) bool {
		//	bscBridgeTokens = append(bscBridgeTokens, *bt)
		//	denomMap[bt.Denom] = struct{}{}
		//	return false
		//})
		//
		//tronBridgeTokens := []types.BridgeToken{}
		//keepers.TronKeeper.IterateBridgeTokenByDenom(ctx, func(bt *types.BridgeToken) bool {
		//	tronBridgeTokens = append(tronBridgeTokens, *bt)
		//	denomMap[bt.Denom] = struct{}{}
		//	return false
		//})
		//
		//bridgeDenoms := make([]string, 0, len(denomMap))
		//for denom := range denomMap {
		//	bridgeDenoms = append(bridgeDenoms, denom)
		//}
		//
		//for _, bt := range bscBridgeTokens {
		//	bt.Denom = types.BridgeTokenPrefix + bt.Denom
		//	keepers.BscKeeper.DelBridgeToken(ctx, &bt)
		//	keepers.BscKeeper.SetBridgeToken(ctx, &bt)
		//}
		//
		//for _, bt := range tronBridgeTokens {
		//	bt.Denom = types.BridgeTokenPrefix + bt.Denom
		//	keepers.TronKeeper.DelBridgeToken(ctx, &bt)
		//	keepers.TronKeeper.SetBridgeToken(ctx, &bt)
		//}
		//
		//logger.Info("2. migrate specific denom metadata...")
		//// Migrate denom metadata for each denom in the map
		//for _, oldDenom := range bridgeDenoms {
		//	denomMetaData, found := keepers.BankKeeper.GetDenomMetaData(ctx, oldDenom)
		//	if !found {
		//		logger.Info("denom metadata not found, skipping", "denom", oldDenom)
		//		continue
		//	}
		//	newDenom := types.BridgeTokenPrefix + oldDenom
		//
		//	logger.Info("migrating denom metadata", "old", oldDenom, "new", newDenom)
		//
		//	// Update all denom units that reference the old denom
		//	for i, d := range denomMetaData.DenomUnits {
		//		if d.Denom == oldDenom {
		//			denomMetaData.DenomUnits[i].Denom = newDenom
		//		}
		//	}
		//
		//	// Update the base denom
		//	denomMetaData.Base = newDenom
		//
		//	// Save the updated metadata with new base denom
		//	keepers.BankKeeper.SetDenomMetaData(ctx, denomMetaData)
		//
		//	logger.Info("denom metadata migrated successfully", "old", oldDenom, "new", newDenom)
		//}
		//
		//logger.Info("3. migrate account balances...")
		//// Collect all accounts that need balance migration
		//var accountsToMigrate []sdk.AccAddress
		//accountBalances := make(map[string]sdk.Coins) // Map account address to coins that need migration
		//
		//// Iterate through all accounts in the auth keeper
		//keepers.AccountKeeper.IterateAccounts(ctx, func(account authtypes.AccountI) bool {
		//	addr := account.GetAddress()
		//	balances := keepers.BankKeeper.GetAllBalances(ctx, addr)
		//
		//	// Check if this account has any balances with old denoms
		//	var coinsToMigrate sdk.Coins
		//	for _, coin := range balances {
		//		if newDenom, needsMigration := oldToNewDenomMap[coin.Denom]; needsMigration {
		//			coinsToMigrate = coinsToMigrate.Add(sdk.NewCoin(newDenom, coin.Amount))
		//		}
		//	}
		//
		//	if !coinsToMigrate.IsZero() {
		//		accountsToMigrate = append(accountsToMigrate, addr)
		//		accountBalances[addr.String()] = coinsToMigrate
		////	}
		//
		//	return false // continue iteration
		//})
		//
		//logger.Info("found accounts to migrate", "count", len(accountsToMigrate))
		//
		//// Migrate balances for each account
		//for _, addr := range accountsToMigrate {
		//	balances := keepers.BankKeeper.GetAllBalances(ctx, addr)
		//	coinsToMigrate := accountBalances[addr.String()]
		//
		//	// Burn old denom coins
		//	for _, coin := range balances {
		//		if _, needsMigration := oldToNewDenomMap[coin.Denom]; needsMigration {
		//			// Burn the old denom from the account
		//			if err := keepers.BankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		//				logger.Error("failed to send coins to module", "address", addr.String(), "coin", coin.String(), "error", err)
		//				continue
		//			}
		//			if err := keepers.BankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		//				logger.Error("failed to burn old denom", "address", addr.String(), "coin", coin.String(), "error", err)
		//				continue
		//			}
		//		}
		//	}
		//
		//	// Mint and send new denom coins
		//	if err := keepers.BankKeeper.MintCoins(ctx, types.ModuleName, coinsToMigrate); err != nil {
		//		logger.Error("failed to mint new denom", "address", addr.String(), "coins", coinsToMigrate.String(), "error", err)
		//		continue
		//	}
		//	if err := keepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coinsToMigrate); err != nil {
		//		logger.Error("failed to send new denom to account", "address", addr.String(), "coins", coinsToMigrate.String(), "error", err)
		//		continue
		//	}
		//}

		logger.Info("upgrade finished successfully.")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
