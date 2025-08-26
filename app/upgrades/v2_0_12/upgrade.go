package v2_0_12

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	appkeepers "github.com/st-chain/me-hub/app/keepers"
	"github.com/st-chain/me-hub/app/upgrades"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v2.0.12
// This upgrade focuses on migrating wnft module class data structures
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

		logger.Info("1. Starting WNFT class data migration...")
		if err := migrateWNFTClassData(ctx, keepers); err != nil {
			return nil, fmt.Errorf("failed to migrate WNFT class data: %w", err)
		}

		logger.Info("upgrade finished successfully.")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

// migrateWNFTClassData migrates existing WNFT class data structures.
// This function updates class metadata and related fields for specific WNFT classes.
func migrateWNFTClassData(ctx sdk.Context, keepers *appkeepers.AppKeepers) error {
	logger := ctx.Logger().With("migration", "wnft_class_data")

	// Define the mapping of class IDs to their new metadata fields
	classUpdates := []struct {
		Id          string
		Name        string
		Symbol      string
		Description string
	}{
		{
			Id:          "495393167",
			Name:        "ME_ExplorerA1000",
			Symbol:      "ME_ExplorerA1000",
			Description: "The Explorer Gold Hunter wields advanced technology, symbolizing wealth and prosperity. It drives the economic growth of new territories, bringing wealth and opportunities to Meta Earth.",
		},
		{
			Id:          "506661488",
			Name:        "ME_ExplorerB1000",
			Symbol:      "ME_ExplorerB1000",
			Description: "The Explorer Freedom moves freely like the wind, symbolizing vast vision and limitless possibilities. It leads the new territories toward wisdom and prosperity, continually driving the development of the future.",
		},
		{
			Id:          "697811991",
			Name:        "ME_PioneerB1000",
			Symbol:      "ME_PioneerB1000",
			Description: "The Pioneer Serenity rests lazily on the clouds, exuding calmness on the outside but with inner strength. It represents peace and wisdom after battle, safeguarding the balance of Meta Earth.",
		},
		{
			Id:          "767391917",
			Name:        "ME_PioneerA1000",
			Symbol:      "ME_PioneerA1000",
			Description: "The \"Pioneer·Might\" wields weapons that symbolize speed and power. Fearlessly charging forward, it paves the way in Meta Earth, becoming the vanguard of world expansion with its boundless fighting spirit.",
		},
	}

	migratedCount := 0

	// Iterate through all classes and update if ID matches
	for _, classUpdate := range classUpdates {
		class, f := keepers.WNFTKeeper.GetClass(ctx, classUpdate.Id)
		if !f {
			panic(fmt.Sprintf("NFT class with ID %s not found", classUpdate.Id))
		}
		// Update the fields as required
		class.Name = classUpdate.Name
		class.Symbol = classUpdate.Symbol
		class.Description = classUpdate.Description

		// Save the updated class back to the keeper
		if err := keepers.WNFTKeeper.UpdateClass(ctx, class); err != nil {
			panic(fmt.Sprintf("Failed to update WNFT class", "id", class.Id, "error", err))
		}
		migratedCount++
	}
	logger.Info(fmt.Sprintf("Migration completed. Successfully migrated: %d", migratedCount))
	return nil
}
