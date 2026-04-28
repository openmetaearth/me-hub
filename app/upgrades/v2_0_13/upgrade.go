package v2_0_13

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	appkeepers "github.com/openmetaearth/me-hub/app/keepers"
	"github.com/openmetaearth/me-hub/app/upgrades"
	"github.com/openmetaearth/me-hub/utils"
	bsctypes "github.com/openmetaearth/me-hub/x/bsc/types"
	gravitykeeper "github.com/openmetaearth/me-hub/x/gravity/keeper"
	"github.com/openmetaearth/me-hub/x/gravity/types"
	gravitytypes "github.com/openmetaearth/me-hub/x/gravity/types"
	trontypes "github.com/openmetaearth/me-hub/x/tron/types"
	"time"
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

		params := keepers.GovKeeper.GetParams(ctx)
		maxDepositPeriod := 30 * time.Minute
		params.MaxDepositPeriod = &maxDepositPeriod
		votingPeriod := 30 * time.Minute
		params.VotingPeriod = &votingPeriod
		if err := keepers.GovKeeper.SetParams(ctx, params); err != nil {
			panic(fmt.Sprintf("failed to set gov module params during upgrade: %s", err.Error()))
		}

		logger.Info("1. upgrade for x/gravity module, set params")
		// !important: proposalRelayers should be sorted to ensure consistency with solidity contracts.
		proposalRelayers := []string{
			"me1frjhlw9slyy7mrhmk0r4vytkyldxqtkf326amv",
			"me1c5zp26c0gq2klk87nrpff3y52u34zn4ydug2yd",
			"me1hrxxjeqae2y5wx3kxcljzns9f2lguygu9qngxh",
			"me14jazxhme3ptv00k52fza5rravx4xn27qs0slz2",
			"me1qdhu5h5g0qwhdpl4q553v7gcmltdr4w3lnqnjg",
		}

		// delegate total amount to module account
		delegateAmount := sdk.NewInt(1 * 1e8)
		//for _, relayerAddr := range proposalRelayers {
		//if err := keepers.BankKeeper.SendCoinsFromAccountToModule(ctx, sdk.MustAccAddressFromBech32(relayerAddr), bsctypes.ModuleName,
		//	sdk.NewCoins(sdk.NewCoin(params.BaseDenom, delegateAmount))); err != nil {
		//	panic(fmt.Sprintf("failed to delegate coins to relayer %s: %s", relayerAddr, err.Error()))
		//}
		//if err := keepers.BankKeeper.SendCoinsFromAccountToModule(ctx, sdk.MustAccAddressFromBech32(relayerAddr), trontypes.ModuleName,
		//	sdk.NewCoins(sdk.NewCoin(params.BaseDenom, delegateAmount))); err != nil {
		//	panic(fmt.Sprintf("failed to delegate coins to relayer %s: %s", relayerAddr, err.Error()))
		//}
		//}

		bscGenState := GenGravityGenesis(ctx.BlockHeight(), proposalRelayers, bsctypes.DefaultGenesisState(), delegateAmount, bsctypes.ModuleName)
		gravitykeeper.InitGenesis(ctx, keepers.BscKeeper, bscGenState)

		tronGenstate := GenGravityGenesis(ctx.BlockHeight(), proposalRelayers, trontypes.DefaultGenesisState(), delegateAmount, trontypes.ModuleName)
		gravitykeeper.InitGenesis(ctx, keepers.TronKeeper, tronGenstate)

		logger.Info("2. upgrade for setting umec metadata.")
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
		logger.Info("upgrade finished successfully.")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func GenGravityGenesis(height int64, proposalRelayers []string, defaultGenesis *gravitytypes.GenesisState, delegateAmount sdk.Int, moduleName string) *gravitytypes.GenesisState {
	// 1. set proposal relayers
	defaultGenesis.ProposalRelayer = types.ProposalRelayer{
		Relayers: proposalRelayers,
	}

	// 2. set relayers
	var err error
	for _, relayerAddr := range proposalRelayers {
		externalAddress := ""
		if moduleName == bsctypes.ModuleName {
			externalAddress, err = utils.MeBech32ToEth(relayerAddr)
			if err != nil {
				panic(err)
			}
		} else if moduleName == trontypes.ModuleName {
			externalAddress, err = utils.MeBech32ToTron(relayerAddr)
			if err != nil {
				panic(err)
			}
		}

		relayer := types.Relayer{
			RelayerAddress:  relayerAddr,
			ExternalAddress: externalAddress,
			DelegateAmount:  delegateAmount,
			StartHeight:     height,
			Online:          true,
			SlashTimes:      0,
		}
		defaultGenesis.Relayers = append(defaultGenesis.Relayers, relayer)
	}

	// 3.relayer set
	var totalPower uint64
	relayerSet := types.RelayerSet{
		Nonce:   0,
		Height:  uint64(height),
		Members: []types.BridgeValidator{},
	}
	for _, relayer := range defaultGenesis.Relayers {
		power := relayer.GetPower()
		if power.LTE(sdkmath.ZeroInt()) {
			continue
		}
		totalPower += power.Uint64()
		bridgeVal := types.BridgeValidator{
			Power:           power.Uint64(),
			ExternalAddress: relayer.ExternalAddress,
		}
		relayerSet.Members = append(relayerSet.Members, bridgeVal)
	}
	for i := range relayerSet.Members {
		relayerSet.Members[i].Power = sdkmath.NewUint(relayerSet.Members[i].Power).MulUint64(gravitytypes.PowerBase).QuoUint64(totalPower).Uint64()
	}
	defaultGenesis.RelayerSets = []types.RelayerSet{relayerSet}
	return defaultGenesis
}

//func setNewModuleParams(ctx sdk.Context, keepers *appkeepers.AppKeepers) {
//	bscState := bsctypes.DefaultGenesisState()
//	if err := keepers.BscKeeper.SetParams(ctx, &bscState.Params); err != nil {
//		panic(fmt.Sprintf("failed to set bsc module params during upgrade: %s", err.Error()))
//	}
//
//	tronState := trontypes.DefaultGenesisState()
//	if err := keepers.TronKeeper.SetParams(ctx, &tronState.Params); err != nil {
//		panic(fmt.Sprintf("failed to set tron module params during upgrade: %s", err.Error()))
//	}
//}
