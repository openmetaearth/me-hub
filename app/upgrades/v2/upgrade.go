package v2

import (
	"encoding/json"
	"fmt"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
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
	"github.com/cosmos/cosmos-sdk/x/nft"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	epochtypes "github.com/osmosis-labs/osmosis/v15/x/epochs/types"
	gammtypes "github.com/osmosis-labs/osmosis/v15/x/gamm/types"
	poolmanagertypes "github.com/osmosis-labs/osmosis/v15/x/poolmanager/types"
	txfeestypes "github.com/osmosis-labs/osmosis/v15/x/txfees/types"
	"github.com/st-chain/me-hub/app/keepers"
	"github.com/st-chain/me-hub/app/upgrades"
	daokeeper "github.com/st-chain/me-hub/x/dao/keeper"
	daotypes "github.com/st-chain/me-hub/x/dao/types"
	delayedacktypes "github.com/st-chain/me-hub/x/delayedack/types"
	didkeeper "github.com/st-chain/me-hub/x/did/keeper"
	didtypes "github.com/st-chain/me-hub/x/did/types"
	eibctypes "github.com/st-chain/me-hub/x/eibc/types"
	incentivestypes "github.com/st-chain/me-hub/x/incentives/types"
	kyckeeper "github.com/st-chain/me-hub/x/kyc/keeper"
	kyctypes "github.com/st-chain/me-hub/x/kyc/types"
	rollapptypes "github.com/st-chain/me-hub/x/rollapp/types"
	sequencertypes "github.com/st-chain/me-hub/x/sequencer/types"
	wnftkeeper "github.com/st-chain/me-hub/x/wnft/keeper"
	wstakingkeeper "github.com/st-chain/me-hub/x/wstaking/keeper"
	"github.com/st-chain/me-hub/x/wstaking/types"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v4
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ upgrades.BaseAppParamManager,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		for n, m := range mm.Modules {
			if mod, ok := m.(module.HasConsensusVersion); ok {
				fromVM[n] = mod.ConsensusVersion()
			}
		}

		logger := ctx.Logger().With("upgrade", UpgradeName)

		ctx.Logger().Info("1.migrate module params")
		migrateModuleParams(ctx, keepers)

		ctx.Logger().Info("2.set new module params")
		setNewModuleParams(ctx, keepers)

		ctx.Logger().Info("3.migrate dao module")
		migrateDao(ctx, keepers.AccountKeeper, keepers.DaoKeeper)

		ctx.Logger().Info("4.migrate validators")
		migrateValidators(ctx, keepers.StakingKeeper)

		ctx.Logger().Info("5.init kyc and did module")
		homePath := GetPath(keepers.UpgradeKeeper)
		migrateKycModule(ctx, keepers.KycKeeper, homePath)

		ctx.Logger().Info("6.migrate kyc and did")
		migrateKycData(ctx, keepers.StakingKeeper, keepers.DidKeeper, keepers.KycKeeper, keepers.WNFTKeeper, homePath)

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
	//baseAppLegacySS := keepers.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable())
	//get, err := keepers.ConsensusParamsKeeper.Get(ctx)
	//if err != nil {
	//	panic(err)
	//}
	//baseapp.MigrateParams(ctx, baseAppLegacySS, &keepers.ConsensusParamsKeeper)
}

func setNewModuleParams(ctx sdk.Context, keepers *keepers.AppKeepers) {
	// overwrite params for delayedack module due to added parameters
	delayedackParams := delayedacktypes.DefaultParams()
	keepers.DelayedAckKeeper.SetParams(ctx, delayedackParams)

	eibcParams := eibctypes.DefaultParams()
	keepers.EIBCKeeper.SetParams(ctx, eibcParams)

	evmParams := evmtypes.DefaultParams()
	keepers.EvmKeeper.SetParams(ctx, evmParams)

	// overwrite params for rollapp module due to proto change
	rollappParams := rollapptypes.DefaultParams()
	keepers.RollappKeeper.SetParams(ctx, rollappParams)

	sequencerParams := sequencertypes.DefaultParams()
	keepers.SequencerKeeper.SetParams(ctx, sequencerParams)

	feemarketParams := feemarkettypes.DefaultParams()
	keepers.FeeMarketKeeper.SetParams(ctx, feemarketParams)

	incentivesParams := incentivestypes.DefaultGenesis()
	keepers.IncentivesKeeper.InitGenesis(ctx, *incentivesParams)

	gammParams := gammtypes.DefaultGenesis()
	keepers.GAMMKeeper.InitGenesis(ctx, *gammParams, nil)

	poolParams := poolmanagertypes.DefaultGenesis()
	keepers.PoolManagerKeeper.InitGenesis(ctx, poolParams)

	epochParams := epochtypes.DefaultGenesis()
	keepers.EpochsKeeper.InitGenesis(ctx, *epochParams)

	txfeeParams := txfeestypes.DefaultGenesis()
	keepers.TxFeesKeeper.InitGenesis(ctx, *txfeeParams)

	sequences := make([]wasmtypes.Sequence, 0)
	for _, k := range [][]byte{wasmtypes.KeyLastCodeID, wasmtypes.KeyLastInstanceID} {
		sequences = append(sequences, wasmtypes.Sequence{
			IDKey: k,
			Value: keepers.WasmKeeper.PeekAutoIncrementID(ctx, k),
		})
	}
	wasmDefault := wasmtypes.GenesisState{
		Params:    keepers.WasmKeeper.GetParams(ctx),
		Codes:     make([]wasmtypes.Code, 0),
		Contracts: make([]wasmtypes.Contract, 0),
		Sequences: sequences,
	}
	_, err := wasmkeeper.InitGenesis(ctx, &keepers.WasmKeeper, wasmDefault)
	if err != nil {
		panic(fmt.Sprintf("wasm init genesis: %v", err))
	}
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
		validator.DelegatorShares = validator.StakerShares
		validator.StakerShares = sdk.ZeroDec()
		validators = append(validators, validator)
	}

	for _, validator := range validators {
		stakingKeeper.SetValidator(ctx, validator)
	}
}

func GetPath(upgradeKeeper *upgradekeeper.Keeper) string {
	path, _ := upgradeKeeper.GetUpgradeInfoPath()
	return strings.TrimRight(path, "/data/upgrade-info.json")
}

func migrateKycModule(ctx sdk.Context, kycKeeper *kyckeeper.Keeper, path string) {
	issuer, err := ReadIssuer(path)
	if err != nil {
		panic(err)
	}
	address := kycKeeper.MustAccAddressFromPubkeyString(issuer.Pubkey)
	if _, found := kycKeeper.GetDID(ctx, address); found {
		panic(fmt.Errorf("issuer %s already exists", address))
	}

	kycKeeper.SetDID(ctx, address, issuer.Did)
	kycKeeper.SetDidInfo(ctx, issuer.Did, issuer)

	service := didtypes.Service{
		Sid:         kyctypes.ModuleName,
		Name:        kyctypes.ModuleName,
		Description: "The KYC verifiable credential issuer based The DID(Decentralized Identity).",
		Issuer:      issuer.Did,
		Status:      didtypes.SERVICE_STATUS_ACTIVE,
	}
	kycKeeper.SetService(ctx, service)
}

func migrateKycData(ctx sdk.Context,
	stakingKeeper *wstakingkeeper.Keeper,
	didKeeper *didkeeper.Keeper,
	kycKeeper *kyckeeper.Keeper,
	nftKeeper *wnftkeeper.Keeper,
	homePath string) {
	// get all data from old module
	meids := stakingKeeper.GetAllMeid(ctx)

	service, found := kycKeeper.GetService(ctx)
	if !found {
		panic("kyc: service not found")
	}

	dids, err := ReadDID(homePath)
	if err != nil {
		panic(err)
	}

	accountPubkey, err := ReadKycPubkey(homePath)
	if err != nil {
		panic(err)
	}

	// Iterate over old data and transform it into new data structure
	for _, oldRecord := range meids {
		did := dids[oldRecord.Account]
		if len(did) > 0 {
			didInfo := didtypes.NewDidInfo(did, oldRecord.Account, accountPubkey[oldRecord.Account], didtypes.DID_STATUS_ACTIVE)
			// write new data to the new module s storage
			didKeeper.SetDID(ctx, sdk.MustAccAddressFromBech32(oldRecord.Account), did)
			didKeeper.SetDidInfo(ctx, didInfo.Did, didInfo)
			didKeeper.SetCredential(
				ctx,
				didInfo.Did,
				service.Sid,
				didtypes.Credential{
					Did:  did,
					Sid:  service.Sid,
					Hash: "",
					Uri:  "",
					Data: []byte(oldRecord.RegionId),
				},
			)
			migrateNFTtoSBT(ctx, stakingKeeper, oldRecord, nftKeeper, kycKeeper, did)
		}
	}

	// If the old module is no longer used, delete the data of the old module
	//oldKeeper.ClearAllData(ctx)
}

func migrateNFTtoSBT(ctx sdk.Context, stakingKeeper *wstakingkeeper.Keeper, oldRecord types.Meid, nftKeeper *wnftkeeper.Keeper, kycKeeper *kyckeeper.Keeper, did string) {
	region, found := stakingKeeper.GetRegion(ctx, oldRecord.RegionId)
	if !found {
		panic(fmt.Sprintf("kyc: region %s not found", oldRecord.RegionId))
	}

	_, classExist := nftKeeper.GetClass(ctx, kyctypes.ModuleName)
	if !classExist {
		err := nftKeeper.SaveClass(ctx, nft.Class{
			Id:          kyctypes.ModuleName,
			Name:        "Soul Bound Token",
			Symbol:      "SBT",
			Description: "",
			Uri:         "",
			UriHash:     "",
			Data:        nil,
		})
		if err != nil {
			panic(err)
		}
	}

	meidNFT, nftFound := stakingKeeper.GetMeidNFTByAccount(ctx, oldRecord.Account)
	if nftFound {
		oldNft, f := nftKeeper.GetNFT(ctx, region.NftClassId, meidNFT.NftId)
		if f {
			if err := kycKeeper.SetSBT(
				ctx,
				nft.NFT{
					ClassId: kyctypes.ModuleName,
					Id:      did,
					Uri:     oldNft.Uri,
					UriHash: oldNft.UriHash,
					Data:    oldNft.Data,
				},
				sdk.MustAccAddressFromBech32(oldRecord.Account),
			); err != nil {
				panic(fmt.Sprintf("account: %s, did: %s, error: %v", oldRecord.Account, did, err))
			}
		}
	}
	//if err := nftKeeper.Burn(ctx, nftInfo.ClassId, nftInfo.Id); err != nil {
	//	panic(err)
	//}
}

func ReadKycPubkey(homePath string) (map[string]string, error) {
	data, err := ioutil.ReadFile(filepath.Join(homePath, "kyc_pubkey.json"))
	if err != nil {
		return nil, err
	}
	accounts := make(map[string]string)
	err = json.Unmarshal(data, &accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func ReadIssuer(path string) (issuer didtypes.DidInfo, err error) {
	data, err := ioutil.ReadFile(filepath.Join(path, "issuer.json"))
	if err != nil {
		return issuer, err
	}
	err = json.Unmarshal(data, &issuer)
	if err != nil {
		return issuer, err
	}
	return issuer, nil
}

func ReadDID(path string) (map[string]string, error) {
	data, err := ioutil.ReadFile(filepath.Join(path, "did.json"))
	if err != nil {
		return nil, err
	}
	list := make(map[string]string)
	err = json.Unmarshal(data, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
