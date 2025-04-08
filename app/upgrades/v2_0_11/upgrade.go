package v2_0_11

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	appkeepers "github.com/st-chain/me-hub/app/keepers"
	"github.com/st-chain/me-hub/app/upgrades"
	kyckeeper "github.com/st-chain/me-hub/x/kyc/keeper"
	wstakingkeeper "github.com/st-chain/me-hub/x/wstaking/keeper"
	wstakingtypes "github.com/st-chain/me-hub/x/wstaking/types"
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
		MigrateDelegation(ctx, keepers.StakingKeeper, keepers.KycKeeper)
		logger.Info("upgrade finished.")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func MigrateDelegation(ctx sdk.Context, stakingKeeper *wstakingkeeper.Keeper, kk *kyckeeper.Keeper) {
	// experience region has no validator address
	expRegion, isFound := stakingKeeper.GetRegion(ctx, wstakingtypes.ExperienceRegionId)
	if !isFound {
		panic(fmt.Errorf("should have experience region"))
	}
	stakingKeeper.IterateAllDelegation(ctx, func(_ int64, del stakingtypes.Delegation) (stop bool) {
		did, didFound := kk.GetDID(ctx, sdk.MustAccAddressFromBech32(del.DelegatorAddress))
		if didFound {
			kyc, kycFound := kk.GetKYC(ctx, did)
			if kycFound {
				region, regionFound := stakingKeeper.GetRegion(ctx, string(kyc.Data))
				if regionFound {
					del.ValidatorAddress = region.OperatorAddress
				}
			}
		} else {
			del.ValidatorAddress = expRegion.OperatorAddress
		}
		stakingKeeper.SetDelegation(ctx, del)
		return false
	})
}
