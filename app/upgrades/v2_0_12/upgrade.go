package v2_0_12

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	appkeepers "github.com/st-chain/me-hub/app/keepers"
	"github.com/st-chain/me-hub/app/upgrades"
	kyckeeper "github.com/st-chain/me-hub/x/kyc/keeper"
	wstakingkeeper "github.com/st-chain/me-hub/x/wstaking/keeper"
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

		migrateFixedDeposit(ctx, keepers.StakingKeeper, keepers.KycKeeper)
		logger.Info("upgrade finished.")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func migrateFixedDeposit(ctx sdk.Context, sk *wstakingkeeper.Keeper, kk *kyckeeper.Keeper) {
	fixedDeposits := sk.GetAllFixedDeposit(ctx)
	for _, fixedDeposit := range fixedDeposits {
		if fixedDeposit.Account == "" {
			continue
		}

		did, didFound := kk.GetDID(ctx, sdk.MustAccAddressFromBech32(fixedDeposit.Account))
		if !didFound {
			panic(fmt.Errorf("fixed deposit account: %s, did not found", fixedDeposit.Account))
		}
		kycData, ok := kk.GetKYC(ctx, did)
		if !ok {
			panic(fmt.Errorf("kyc data not found: %s", did))
		}
		regionId := string(kycData.Data)
		region, found := sk.GetRegion(ctx, regionId)
		if !found {
			panic(fmt.Errorf("region not found: %s", regionId))
		}
		region.FixedDepositAmount = region.FixedDepositAmount.Add(fixedDeposit.Principal.Amount)
		sk.SetRegion(ctx, region)
	}
}
