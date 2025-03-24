package v2_0_12

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/nft"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	appkeepers "github.com/st-chain/me-hub/app/keepers"
	"github.com/st-chain/me-hub/app/upgrades"
	didkeeper "github.com/st-chain/me-hub/x/did/keeper"
	didtypes "github.com/st-chain/me-hub/x/did/types"
	kyckeeper "github.com/st-chain/me-hub/x/kyc/keeper"
	kyctypes "github.com/st-chain/me-hub/x/kyc/types"
	wnftkeeper "github.com/st-chain/me-hub/x/wnft/keeper"
	wstakingkeeper "github.com/st-chain/me-hub/x/wstaking/keeper"
	"github.com/st-chain/me-hub/x/wstaking/types"
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

		migrateFixedDeposit(ctx, keepers.StakingKeeper, keepers.KycKeeper, keepers.DidKeeper, keepers.WNFTKeeper)
		logger.Info("upgrade finished.")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func migrateFixedDeposit(ctx sdk.Context, sk *wstakingkeeper.Keeper, kk *kyckeeper.Keeper, didKeeper *didkeeper.Keeper,
	nftKeeper *wnftkeeper.Keeper,
) {
	fixedDeposits := sk.GetAllFixedDeposit(ctx)
	didNumber := 9998887776660
	for _, fixedDeposit := range fixedDeposits {
		if fixedDeposit.Account == "" {
			continue
		}

		regionId := ""
		did, didFound := kk.GetDID(ctx, sdk.MustAccAddressFromBech32(fixedDeposit.Account))
		if !didFound {
			didStr := fmt.Sprintf("%d", didNumber)
			meid, found := sk.GetMeid(ctx, fixedDeposit.Account)
			if found {
				didInfo := didtypes.DidInfo{
					Did:      didStr,
					Address:  meid.Account,
					Pubkey:   "",
					RegionId: meid.RegionId,
					KycLevel: didtypes.KYC_LEVEL_ONE,
					Status:   didtypes.DID_STATUS_ACTIVE,
				}
				vc := didtypes.Credential{
					Did:  didStr,
					Sid:  "kyc",
					Uri:  "",
					Hash: "",
					Data: []byte(meid.RegionId),
				}
				sk.SetInviterReward(ctx, meid.Account)
				// write new data to the new module s storage
				didKeeper.SetDID(ctx, sdk.MustAccAddressFromBech32(meid.Account), didStr)
				didKeeper.SetDidInfo(ctx, didInfo.Did, didInfo)
				didKeeper.SetCredential(
					ctx,
					didInfo.Did,
					"kyc",
					vc,
				)
				didKeeper.AddFilters(ctx, didStr, "kyc", [][]byte{[]byte(meid.RegionId)}, vc)
				migrateNFTtoSBT(ctx, sk, meid, nftKeeper, kk, didStr)
			}
			didNumber++
			regionId = meid.RegionId
		} else {
			kycData, ok := kk.GetKYC(ctx, did)
			if !ok {
				panic(fmt.Errorf("kyc data not found: %s", did))
			}
			regionId = string(kycData.Data)
		}
		region, found := sk.GetRegion(ctx, regionId)
		if !found {
			panic(fmt.Errorf("region not found: %s", regionId))
		}
		region.FixedDepositAmount = region.FixedDepositAmount.Add(fixedDeposit.Principal.Amount)
		sk.SetRegion(ctx, region)
	}
}

func migrateNFTtoSBT(ctx sdk.Context,
	stakingKeeper *wstakingkeeper.Keeper,
	oldRecord types.Meid,
	nftKeeper *wnftkeeper.Keeper,
	kycKeeper *kyckeeper.Keeper,
	didStr string,
) {
	_, found := stakingKeeper.GetRegion(ctx, oldRecord.RegionId)
	if !found {
		panic(fmt.Sprintf("kyc: region %s not found", oldRecord.RegionId))
	}

	if err := kycKeeper.SetSBT(
		ctx,
		nft.NFT{
			ClassId: kyctypes.ModuleName,
			Id:      didStr,
			Uri:     "aa",
			UriHash: "aa",
			Data:    nil,
		},
		sdk.MustAccAddressFromBech32(oldRecord.Account),
	); err != nil {
		panic(fmt.Sprintf("account: %s, did: %s, error: %v", oldRecord.Account, didStr, err))
	}
	stakingKeeper.RemoveMeidNFT(ctx, oldRecord.Account, oldRecord.RegionId)
}
