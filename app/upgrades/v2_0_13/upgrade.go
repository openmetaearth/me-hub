package v2_0_13

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	appkeepers "github.com/st-chain/me-hub/app/keepers"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/app/upgrades"
	"github.com/st-chain/me-hub/utils"
	bsctypes "github.com/st-chain/me-hub/x/bsc/types"
	gravitykeeper "github.com/st-chain/me-hub/x/gravity/keeper"
	"github.com/st-chain/me-hub/x/gravity/types"
	gravitytypes "github.com/st-chain/me-hub/x/gravity/types"
	trontypes "github.com/st-chain/me-hub/x/tron/types"
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

		logger.Info("1. upgrade for x/gravity module, set params")
		proposalRelayers := []string{
			"me19tuh989pxxq7wawmd2fer4ckd4vvl3a3sepez8",
			"me19jcd2970et5tg82qzd5nfltvcuxqtds6csztm7",
			"me1k99ppc456q29mmpf56hmqgnhma743h5w2dj0n2",
			"me1clsfspg3djv8em7u4zzj3z4jnpwl42ava2exrx",
			"me1al863lkzttl9kvtphlmn4z5ypjl83k7tk9hv27",
		}
		bscGenesis := bsctypes.DefaultGenesisState()
		InitGravityGenesis(ctx, proposalRelayers, bscGenesis, keepers, bsctypes.ModuleName)

		tronGenesis := trontypes.DefaultGenesisState()
		InitGravityGenesis(ctx, proposalRelayers, tronGenesis, keepers, trontypes.ModuleName)

		logger.Info("upgrade finished successfully.")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func InitGravityGenesis(ctx sdk.Context, proposalRelayers []string, defaultGenesis *gravitytypes.GenesisState, keepers *appkeepers.AppKeepers, moduleName string) {
	// 1. set proposal relayers
	defaultGenesis.ProposalRelayer = types.ProposalRelayer{
		Relayers: proposalRelayers,
	}

	// 2. set relayers
	delegateAmount := sdk.NewInt(1 * 1e8)
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
			StartHeight:     ctx.BlockHeight(),
			Online:          true,
			SlashTimes:      0,
		}
		defaultGenesis.Relayers = append(defaultGenesis.Relayers, relayer)

		// delegate total amount to module account
		if err := keepers.BankKeeper.SendCoinsFromAccountToModule(ctx, sdk.MustAccAddressFromBech32(relayerAddr), moduleName,
			sdk.NewCoins(sdk.NewCoin(params.BaseDenom, delegateAmount))); err != nil {
			panic(fmt.Sprintf("failed to delegate coins to relayer %s: %s", relayerAddr, err.Error()))
		}
	}

	// 3. relayer set
	var totalPower uint64
	var members []types.BridgeValidator
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
		defaultGenesis.RelayerSets[0].Members = append(defaultGenesis.RelayerSets[0].Members, bridgeVal)
	}
	for i := range members {
		members[i].Power = sdkmath.NewUint(members[i].Power).MulUint64(utils.PowerBase).QuoUint64(totalPower).Uint64()
	}
	defaultGenesis.RelayerSets = []types.RelayerSet{
		{
			Nonce:   1,
			Height:  uint64(ctx.BlockHeight()),
			Members: members,
		},
	}
	if moduleName == bsctypes.ModuleName {
		gravitykeeper.InitGenesis(ctx, keepers.BscKeeper, defaultGenesis)
	} else if moduleName == trontypes.ModuleName {
		gravitykeeper.InitGenesis(ctx, keepers.TronKeeper.Keeper, defaultGenesis)
	}
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
