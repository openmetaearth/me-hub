package v2_0_11

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	appkeepers "github.com/st-chain/me-hub/app/keepers"
	"github.com/st-chain/me-hub/app/upgrades"
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

		oldDid1 := "0000000000000001"
		newDid1 := "0000000000001"
		info1, found := keepers.DidKeeper.GetDidInfo(ctx, oldDid1)
		if !found {
			panic("did: info not found")
		}
		info1.Did = newDid1
		keepers.DidKeeper.SetDidInfo(ctx, newDid1, info1)
		keepers.DidKeeper.DeleteDidInfo(ctx, oldDid1)
		keepers.DidKeeper.SetDID(ctx, sdk.MustAccAddressFromBech32(info1.Address), newDid1)

		oldDid2 := "0000000000000002"
		newDid2 := "0000000000002"
		info2, found := keepers.DidKeeper.GetDidInfo(ctx, oldDid2)
		if !found {
			panic("did: info not found")
		}
		info2.Did = newDid2
		keepers.DidKeeper.SetDidInfo(ctx, newDid2, info2)
		keepers.DidKeeper.DeleteDidInfo(ctx, oldDid2)
		keepers.DidKeeper.DeleteDID(ctx, sdk.MustAccAddressFromBech32(info2.Address))
		keepers.DidKeeper.SetDID(ctx, sdk.MustAccAddressFromBech32(info2.Address), newDid2)

		service, found := keepers.KycKeeper.GetService(ctx)
		if !found {
			panic("kyc: service not found")
		}
		if len(service.Issuers) != 2 {
			panic("kyc: issuer count not match")
		}
		service.Issuers[0] = newDid1
		service.Issuers[1] = newDid2
		keepers.KycKeeper.SetService(ctx, service)

		didInfos := keepers.DidKeeper.GetDidInfos(ctx)
		for _, didInfo := range didInfos {
			vc, f := keepers.DidKeeper.GetCredential(ctx, didInfo.Did, service.Sid)
			if !f {
				if didInfo.Did != newDid1 && didInfo.Did != newDid2 {
					logger.Error("credential not found", "did", didInfo.String())
				}
			}
			keepers.DidKeeper.AddFilters(ctx, didInfo.Did, service.Sid, [][]byte{[]byte(didInfo.RegionId)}, vc)
		}
		logger.Info("upgrade finished.")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
