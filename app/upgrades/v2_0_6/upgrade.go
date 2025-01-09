package v2_0_6

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/nft"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	appkeepers "github.com/st-chain/me-hub/app/keepers"
	"github.com/st-chain/me-hub/app/upgrades"
	kyctypes "github.com/st-chain/me-hub/x/kyc/types"
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

		_, classExist := keepers.WNFTKeeper.GetClass(ctx, kyctypes.ModuleName)
		if !classExist {
			err := keepers.WNFTKeeper.SaveClass(ctx, nft.Class{
				Id:          kyctypes.ModuleName,
				Name:        kyctypes.ModuleName,
				Symbol:      "SBT",
				Description: "",
				Uri:         "",
				UriHash:     "",
				Data:        nil,
				TotalSupply: 0,
			})
			if err != nil {
				panic(err)
			}
		}

		//upgrade for migrate did info
		logger.Info("upgrade finished.")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
