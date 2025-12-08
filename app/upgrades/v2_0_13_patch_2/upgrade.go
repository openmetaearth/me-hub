package v2_0_13_patch_2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	appkeepers "github.com/st-chain/me-hub/app/keepers"
	"github.com/st-chain/me-hub/app/upgrades"
	gravitykeeper "github.com/st-chain/me-hub/x/gravity/keeper"
	"github.com/st-chain/me-hub/x/gravity/types"
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

		logger.Info("1. upgrade for setting umec metadata.")
		keepers.BankKeeper.SetDenomMetaData(ctx, banktypes.Metadata{
			Description: "Denom metadata for MEC (umec)",
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    "umec",
					Exponent: uint32(0),
				},
				{
					Denom:    "MEC",
					Exponent: uint32(8),
				},
			},
			Base:    "umec",
			Display: "MEC",
			Name:    "MEC",
			Symbol:  "MEC",
			URI:     "",
			URIHash: "",
		})

		logger.Info("2. clear tron gengesis")
		ClearGenesis(ctx, keepers.TronKeeper)

		logger.Info("upgrade finished successfully.")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func ClearGenesis(ctx sdk.Context, k gravitykeeper.Keeper) {
	//genesis := gravitykeeper.ExportGenesis(ctx, k)
	k.IterateOutgoingTxBatches(ctx, func(batch *types.OutgoingTxBatch) bool {
		k.DeleteBatch(ctx, batch)
		return false
	})
	k.SetLastObservedEventNonce(ctx, 0)
	k.SetLastObservedBlockHeight(ctx, 0, 0)
	k.PruneAttestations(ctx)
	k.IterateUnbatchedTransactions(ctx, "", func(tx *types.OutgoingTransferTx) bool {
		err := k.DelUnbatchedTx(ctx, tx.Fee, tx.Id)
		if err != nil {
			panic(err)
		}
		return false
	})

	relayerSets := []types.RelayerSet{}
	k.IterateRelayerSets(ctx, false, func(relayerSet *types.RelayerSet) bool {
		relayerSets = append(relayerSets, *relayerSet)
		return false
	})

	for _, rs := range relayerSets {
		k.DeleteRelayerSetConfirm(ctx, rs.Nonce)
	}
	nextID := k.AutoIncrementID(ctx, types.KeyLastOutgoingBatchID)
	k.IterateBridgeTokenByDenom(ctx, func(token *types.BridgeToken) bool {
		k.DelBridgeToken(ctx, token)
		for i := uint64(0); i < nextID; i++ {
			k.DeleteBatchConfirm(ctx, 0, token.ContractAddress)
		}
		return false
	})
	if lastObserved := k.GetLastObservedRelayerSet(ctx); lastObserved != nil {
		k.DelLastObservedRelayerSet(ctx)
	}
	return
}
