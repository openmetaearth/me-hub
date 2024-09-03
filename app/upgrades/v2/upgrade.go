package v2

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"github.com/st-chain/me-hub/app/keepers"
	"github.com/st-chain/me-hub/app/upgrades"
	daokeeper "github.com/st-chain/me-hub/x/dao/keeper"
	daotypes "github.com/st-chain/me-hub/x/dao/types"
	delayedacktypes "github.com/st-chain/me-hub/x/delayedack/types"
	rollapptypes "github.com/st-chain/me-hub/x/rollapp/types"
	wstakingkeeper "github.com/st-chain/me-hub/x/wstaking/keeper"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v4
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ upgrades.BaseAppParamManager,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)

		ctx.Logger().Info("1.migrate module params")
		migrateModuleParams(ctx, keepers)

		ctx.Logger().Info("2.set new module params")
		setNewModuleParams(ctx, keepers)

		ctx.Logger().Info("3.migrate dao module")
		migrateDao(ctx, keepers.AccountKeeper, keepers.DaoKeeper)

		ctx.Logger().Info("4.migrate validators")
		migrateValidators(ctx, keepers.StakingKeeper)

		// Start running the module migrations
		logger.Debug("running module migrations ...")
		//ctx = ctx.WithChainID(metypes.V2ChainId)
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

//nolint:staticcheck
func migrateModuleParams(ctx sdk.Context, keepers *keepers.AppKeepers) {
	// Set param key table for params module migration
	for _, subspace := range keepers.ParamsKeeper.GetSubspaces() {
		var keyTable paramstypes.KeyTable
		switch subspace.Name() {
		case authtypes.ModuleName:
			keyTable = authtypes.ParamKeyTable()
		case banktypes.ModuleName:
			keyTable = banktypes.ParamKeyTable()
		case stakingtypes.ModuleName:
			keyTable = stakingtypes.ParamKeyTable()
		case minttypes.ModuleName:
			keyTable = minttypes.ParamKeyTable()
		case distrtypes.ModuleName:
			keyTable = distrtypes.ParamKeyTable()
		case slashingtypes.ModuleName:
			keyTable = slashingtypes.ParamKeyTable()
		case govtypes.ModuleName:
			keyTable = govv1.ParamKeyTable()
		case crisistypes.ModuleName:
			keyTable = crisistypes.ParamKeyTable()

		// Ethermint  modules
		case evmtypes.ModuleName:
			keyTable = evmtypes.ParamKeyTable()
		case feemarkettypes.ModuleName:
			keyTable = feemarkettypes.ParamKeyTable()
		default:
			continue
		}

		if !subspace.HasKeyTable() {
			subspace.WithKeyTable(keyTable)
		}
	}
	// Migrate Tendermint consensus parameters from x/params module to a dedicated x/consensus module.
	baseAppLegacySS := keepers.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable())
	baseapp.MigrateParams(ctx, baseAppLegacySS, &keepers.ConsensusParamsKeeper)
}

func setNewModuleParams(ctx sdk.Context, keepers *keepers.AppKeepers) {
	// overwrite params for delayedack module due to added parameters
	delayedackParams := delayedacktypes.DefaultParams()
	keepers.DelayedAckKeeper.SetParams(ctx, delayedackParams)

	// overwrite params for rollapp module due to proto change
	rollappParams := rollapptypes.DefaultParams()
	keepers.RollappKeeper.SetParams(ctx, rollappParams)
}

func migrateDao(ctx sdk.Context, ak authkeeper.AccountKeeper, dk daokeeper.Keeper) {
	daoAddresses := daotypes.DaoAddresses{
		GlobalDao:      ak.GetAccountAddressByID(ctx, 0),
		MeidDao:        ak.GetAccountAddressByID(ctx, 1),
		DevOperator:    ak.GetAccountAddressByID(ctx, 2),
		AirdropAddress: ak.GetAccountAddressByID(ctx, 3),
	}
	dk.SetDaoAddresses(ctx, daoAddresses)
}

func migrateValidators(ctx sdk.Context, stakingKeeper *wstakingkeeper.Keeper) {
	validators := stakingKeeper.GetAllValidators(ctx)
	store := ctx.KVStore(stakingKeeper.GetStoreKey())

	iterator := sdk.KVStorePrefixIterator(store, types.ValidatorsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator := stakingtypes.MustUnmarshalValidator(stakingKeeper.GetCdc(), iterator.Value())
		validators = append(validators, validator)
	}

	for _, validator := range validators {
		stakingKeeper.SetValidator(ctx, validator)
	}
}
