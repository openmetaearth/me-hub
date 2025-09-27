package v2_0_10

import (
	"cosmossdk.io/math"
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
	appkeepers "github.com/st-chain/me-hub/app/keepers"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/app/upgrades"
	"github.com/st-chain/me-hub/utils"
	daokeeper "github.com/st-chain/me-hub/x/dao/keeper"
	daotypes "github.com/st-chain/me-hub/x/dao/types"
	delayedacktypes "github.com/st-chain/me-hub/x/delayedack/types"
	didkeeper "github.com/st-chain/me-hub/x/did/keeper"
	didtypes "github.com/st-chain/me-hub/x/did/types"
	eibctypes "github.com/st-chain/me-hub/x/eibc/types"
	kyckeeper "github.com/st-chain/me-hub/x/kyc/keeper"
	kyctypes "github.com/st-chain/me-hub/x/kyc/types"
	groupkeeper "github.com/st-chain/me-hub/x/megroup/keeper"
	megrouptypes "github.com/st-chain/me-hub/x/megroup/types"
	rollappkeeper "github.com/st-chain/me-hub/x/rollapp/keeper"
	rollapptypes "github.com/st-chain/me-hub/x/rollapp/types"
	sequencertypes "github.com/st-chain/me-hub/x/sequencer/types"
	wbankkeeper "github.com/st-chain/me-hub/x/wbank/keeper"
	wnftkeeper "github.com/st-chain/me-hub/x/wnft/keeper"
	wstakingkeeper "github.com/st-chain/me-hub/x/wstaking/keeper"
	"github.com/st-chain/me-hub/x/wstaking/types"
	wstakingtypes "github.com/st-chain/me-hub/x/wstaking/types"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v4
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ upgrades.BaseAppParamManager,
	keepers *appkeepers.AppKeepers,
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
		MigrateDao(ctx, keepers.AccountKeeper, keepers.DaoKeeper, keepers.RollappKeeper)

		ctx.Logger().Info("4.migrate validators")
		migrateValidators(ctx, keepers.StakingKeeper)

		ctx.Logger().Info("5.fixed deposit")
		MigrateFixedDeposit(ctx, keepers.StakingKeeper, keepers.KycKeeper, keepers.BankKeeper)

		ctx.Logger().Info("6.init kyc and did module")
		homePath := GetPath(keepers.UpgradeKeeper)
		migrateKycModule(ctx, keepers.KycKeeper, homePath)

		ctx.Logger().Info("7.migrate kyc and did")
		MigrateKycData(ctx,
			keepers.StakingKeeper,
			keepers.DidKeeper,
			keepers.KycKeeper,
			keepers.WNFTKeeper,
			keepers.GroupKeeper,
			homePath,
			RealDIDReader{},
			RealKycPubkeyReader{})

		ctx.Logger().Info("8.migrate nft ipfs uri")
		MigrateNftUri(ctx, keepers.WNFTKeeper, homePath, RealNftReader{})

		// Start running the module migrations
		logger.Debug("running module migrations ...")
		//ctx = ctx.WithChainID(metypes.V2ChainId)

		ctx.Logger().Info("9.migrate region class id, fix name...")
		migrateRegionClassName(ctx, keepers.StakingKeeper, keepers.WNFTKeeper)

		ctx.Logger().Info("10.migrate group")
		migrateGroup(ctx, homePath, keepers.GroupKeeper, keepers.StakingKeeper, keepers.KycKeeper)

		ctx.Logger().Info("11.migrate delegation")
		MigrateDelegation(ctx, keepers.StakingKeeper, keepers.KycKeeper)

		// check
		CheckDao(ctx, keepers.AccountKeeper, keepers.DaoKeeper)

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

//nolint:staticcheck
func migrateModuleParams(ctx sdk.Context, keepers *appkeepers.AppKeepers) {
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

func setNewModuleParams(ctx sdk.Context, keepers *appkeepers.AppKeepers) {
	// overwrite params for delayedack module due to added parameters
	delayedackParams := delayedacktypes.DefaultParams()
	keepers.DelayedAckKeeper.SetParams(ctx, delayedackParams)

	eibcParams := eibctypes.DefaultParams()
	keepers.EIBCKeeper.SetParams(ctx, eibcParams)

	evmParams := evmtypes.DefaultParams()
	keepers.EvmKeeper.SetParams(ctx, evmParams)

	//rollappParams := rollapptypes.DefaultParams()
	//keepers.RollappKeeper.SetParams(ctx, rollappParams)

	sequencerParams := sequencertypes.DefaultParams()
	keepers.SequencerKeeper.SetParams(ctx, sequencerParams)

	feemarketParams := feemarkettypes.DefaultParams()
	keepers.FeeMarketKeeper.SetParams(ctx, feemarketParams)

	//gammParams := gammtypes.DefaultGenesis()
	//keepers.GAMMKeeper.InitGenesis(ctx, *gammParams, nil)

	//poolParams := poolmanagertypes.DefaultGenesis()
	//keepers.PoolManagerKeeper.InitGenesis(ctx, poolParams)

	//epochParams := epochtypes.DefaultGenesis()
	//keepers.EpochsKeeper.InitGenesis(ctx, *epochParams)

	//txfeeParams := txfeestypes.DefaultGenesis()
	//keepers.TxFeesKeeper.InitGenesis(ctx, *txfeeParams)

	sequences := make([]wasmtypes.Sequence, 0)
	for _, k := range [][]byte{wasmtypes.KeySequenceCodeID, wasmtypes.KeySequenceInstanceID} {
		sequences = append(sequences, wasmtypes.Sequence{
			IDKey: k,
			Value: keepers.WasmKeeper.PeekAutoIncrementID(ctx, k),
		})
	}
	params := keepers.WasmKeeper.GetParams(ctx)
	if params.InstantiateDefaultPermission == wasmtypes.AccessTypeUnspecified {
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

	govParams := keepers.GovKeeper.GetParams(ctx)
	govParams.BurnVoteQuorum = false
	govParams.BurnVoteVeto = false
	govParams.BurnProposalDepositPrevote = false

	keepers.GovKeeper.SetParams(ctx, govParams)
}

func MigrateDao(ctx sdk.Context, ak authkeeper.AccountKeeper, dk daokeeper.Keeper, rk *rollappkeeper.Keeper) {
	daoAddresses := daotypes.DaoAddresses{
		GlobalDao:      ak.GetAccountAddressByID(ctx, 0),
		MeidDao:        ak.GetAccountAddressByID(ctx, 1),
		DevOperator:    ak.GetAccountAddressByID(ctx, 2),
		AirdropAddress: ak.GetAccountAddressByID(ctx, 3),
	}
	dk.SetDaoAddresses(ctx, daoAddresses)

	rollappParams := rollapptypes.DefaultParams()
	rollappParams.DisputePeriodInBlocks = 50
	rollappParams.DeployerWhitelist = []rollapptypes.DeployerParams{
		{Address: daoAddresses.GlobalDao},
		{Address: daoAddresses.MeidDao},
	}
	rk.SetParams(ctx, rollappParams)
}

func CheckDao(ctx sdk.Context, ak authkeeper.AccountKeeper, dk daokeeper.Keeper) {
	daoAddresses := daotypes.DaoAddresses{
		GlobalDao:      ak.GetAccountAddressByID(ctx, 0),
		MeidDao:        ak.GetAccountAddressByID(ctx, 1),
		DevOperator:    ak.GetAccountAddressByID(ctx, 2),
		AirdropAddress: ak.GetAccountAddressByID(ctx, 3),
	}
	dao, found := dk.GetDaoAddresses(ctx)
	if !found {
		panic("dao set failed, not found")
	}
	if dao.GlobalDao != daoAddresses.GlobalDao {
		panic("dao set failed, global dao")
	}
	if dao.MeidDao != dk.GetMeidDao(ctx) {
		panic("dao set failed, meid dao")
	}
}

func migrateValidators(ctx sdk.Context, stakingKeeper *wstakingkeeper.Keeper) {
	validators := stakingKeeper.GetAllValidators(ctx)
	store := ctx.KVStore(stakingKeeper.GetStoreKey())

	iterator := sdk.KVStorePrefixIterator(store, stakingtypes.ValidatorsKey)
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
	pathList := strings.Split(path, "/data")
	return pathList[0]
}

func migrateKycModule(ctx sdk.Context, kycKeeper *kyckeeper.Keeper, path string) {
	issuers, err := ReadIssuer(path)
	if err != nil {
		panic(err)
	}
	issuerDids := []string{}
	for _, issuer := range issuers {
		//address := kycKeeper.MustAccAddressFromPubkeyString(issuer.Pubkey) // global admin has multi pubkey
		if _, found := kycKeeper.GetDID(ctx, sdk.MustAccAddressFromBech32(issuer.Address)); found {
			panic(fmt.Errorf("issuer %s already exists", sdk.MustAccAddressFromBech32(issuer.Address)))
		}
		kycKeeper.SetDID(ctx, sdk.MustAccAddressFromBech32(issuer.Address), issuer.Did)
		kycKeeper.SetDidInfo(ctx, issuer.Did, issuer)

		vc := didtypes.Credential{
			Did:  issuer.Did,
			Sid:  kyctypes.ModuleName,
			Hash: "",
			Uri:  "",
			Data: []byte(issuer.RegionId),
		}
		kycKeeper.SetKYC(ctx, issuer.Did, vc)
		kycKeeper.AddFilters(ctx, issuer.Did, [][]byte{[]byte(issuer.RegionId)}, vc)
		issuerDids = append(issuerDids, issuer.Did)
	}
	service := didtypes.Service{
		Sid:         kyctypes.ModuleName,
		Name:        kyctypes.ModuleName,
		Description: "The KYC verifiable credential issuer based The DID(Decentralized Identity).",
		Issuers:     issuerDids,
		Status:      didtypes.SERVICE_STATUS_ACTIVE,
	}
	kycKeeper.SetService(ctx, service)
}

func MigrateKycData(ctx sdk.Context,
	stakingKeeper *wstakingkeeper.Keeper,
	didKeeper *didkeeper.Keeper,
	kycKeeper *kyckeeper.Keeper,
	nftKeeper *wnftkeeper.Keeper,
	gk *groupkeeper.Keeper,
	homePath string,
	didReader DIDReader,
	kycPubkeyReader KycPubkeyReader) {
	// get all data from old module
	meids := stakingKeeper.GetAllMeid(ctx)

	service, found := kycKeeper.GetService(ctx)
	if !found {
		panic("kyc: service not found")
	}

	didData, err := didReader.ReadDID(filepath.Join(homePath, didFilePath))
	if err != nil {
		panic(fmt.Sprintf("read did: %v", err))
	}

	//accountPubkey, err := kycPubkeyReader.ReadKycPubkey(filepath.Join(homePath, kycPubkeyFilePath))
	//if err != nil {
	//	panic(err)
	//}

	_, classExist := nftKeeper.GetClass(ctx, kyctypes.ModuleName)
	if !classExist {
		err := nftKeeper.SaveClass(ctx, nft.Class{
			Id:          kyctypes.ModuleName,
			Name:        kyctypes.ModuleName,
			Symbol:      "SBT",
			Description: "Soul Bound Token",
			Uri:         "",
			UriHash:     "",
			Data:        nil,
			TotalSupply: 0,
		})
		if err != nil {
			panic(err)
		}
	}

	// Iterate over old data and transform it into new data structure
	for _, meid := range meids {
		did, ok := didData[meid.Account]
		if !(ok && len(did.Did) > 0) {
			panic(fmt.Sprintf("did not found, account: %s, please upgrade later.", meid.Account))
		}
		if ok && len(did.Did) > 0 {
			didLevel := didtypes.KycLevel(did.Level)
			if did.Level == 0 {
				didLevel = didtypes.KYC_LEVEL_ONE
			}
			didInfo := didtypes.DidInfo{
				Did:      did.Did,
				Address:  meid.Account,
				Pubkey:   did.PubKey,
				RegionId: meid.RegionId,
				KycLevel: didLevel,
				Status:   didtypes.DID_STATUS_ACTIVE,
			}
			vc := didtypes.Credential{
				Did:  did.Did,
				Sid:  service.Sid,
				Uri:  did.KycUri,
				Hash: did.KycUriHash,
				Data: []byte(meid.RegionId),
			}
			if did.Level == 1 {
				stakingKeeper.SetInviterReward(ctx, meid.Account)
			}
			// write new data to the new module s storage
			didKeeper.SetDID(ctx, sdk.MustAccAddressFromBech32(meid.Account), did.Did)
			didKeeper.SetDidInfo(ctx, didInfo.Did, didInfo)
			didKeeper.SetCredential(
				ctx,
				didInfo.Did,
				service.Sid,
				vc,
			)
			didKeeper.AddFilters(ctx, did.Did, service.Sid, [][]byte{[]byte(meid.RegionId)}, vc)
			migrateNFTtoSBT(ctx, stakingKeeper, meid, nftKeeper, kycKeeper, did)
		} else {
			didNumber := 9998887776660
			didStr := fmt.Sprintf("%d", didNumber)
			for kycKeeper.HasDidInfo(ctx, didStr) {
				didNumber++
				didStr = fmt.Sprintf("%d", didNumber)
			}
			didNumber++
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
			stakingKeeper.SetInviterReward(ctx, meid.Account)
			// write new data to the new module s storage
			didKeeper.SetDID(ctx, sdk.MustAccAddressFromBech32(meid.Account), didStr)
			didKeeper.SetDidInfo(ctx, didInfo.Did, didInfo)
			didKeeper.SetCredential(
				ctx,
				didInfo.Did,
				service.Sid,
				vc,
			)
			didKeeper.AddFilters(ctx, didStr, service.Sid, [][]byte{[]byte(meid.RegionId)}, vc)
			migrateNFTtoSBT(ctx, stakingKeeper, meid, nftKeeper, kycKeeper, DidData{
				Did:        didStr,
				Level:      1,
				Uri:        "",
				UriHash:    "",
				KycUri:     "",
				KycUriHash: "",
			})
		}

		if meid.RewardType == 1 {
			gk.SetMemberJoined(ctx, megrouptypes.MemberJoined{
				Address: meid.Account,
				GroupId: 0,
			})
		}
		stakingKeeper.RemoveMeid(ctx, meid.Account, meid.RegionId)
	}
}

func migrateNFTtoSBT(ctx sdk.Context,
	stakingKeeper *wstakingkeeper.Keeper,
	oldRecord types.Meid,
	nftKeeper *wnftkeeper.Keeper,
	kycKeeper *kyckeeper.Keeper,
	did DidData) {
	if err := kycKeeper.SetSBT(
		ctx,
		nft.NFT{
			ClassId: kyctypes.ModuleName,
			Id:      did.Did,
			Uri:     did.Uri,
			UriHash: did.UriHash,
			Data:    nil,
		},
		sdk.MustAccAddressFromBech32(oldRecord.Account),
	); err != nil {
		panic(fmt.Sprintf("account: %s, did: %s, error: %v", oldRecord.Account, did.Did, err))
	}
	meidNft, found := stakingKeeper.GetMeidNFT(ctx, oldRecord.Account)
	if found {
		nftKeeper.Burn(ctx, strings.ToUpper(oldRecord.RegionId)+"-NFT-CLASS-ID-", meidNft.NftId)
		stakingKeeper.RemoveMeidNFT(ctx, oldRecord.Account, oldRecord.RegionId)
	}
}

func MigrateNftUri(ctx sdk.Context, nftKeeper *wnftkeeper.Keeper, homePath string, nftReader NftReader) {
	nftData, err := nftReader.ReadNft(filepath.Join(homePath, nftFilePath))
	if err != nil {
		panic(fmt.Sprintf("read nft: %v", err))
	}
	classlist := nftKeeper.GetClasses(ctx)
	for _, class := range classlist {
		if class == nil || class.Id == kyctypes.ModuleName {
			continue
		}

		classData, ok := nftData[class.Id]
		if !ok {
			continue
		}

		if classData.ClassURI != "" {
			class.Uri = classData.ClassURI
		}
		if classData.ClassURIHash != "" {
			class.UriHash = classData.ClassURIHash
		}
		err = nftKeeper.UpdateClass(ctx, *class)
		if err != nil {
			panic(fmt.Errorf("update class in migrate nft: %v", err))
		}

		nftList := nftKeeper.GetNFTsOfClass(ctx, class.Id)
		for _, nft := range nftList {
			nftUriData, ok := classData.NftData[nft.Id]
			if !ok {
				continue
			}
			if nftUriData.URI != "" {
				nft.Uri = nftUriData.URI
			}
			if nftUriData.URIHash != "" {
				nft.UriHash = nftUriData.URIHash
			}
			err = nftKeeper.Update(ctx, nft)
			if err != nil {
				panic(err)
			}
		}
	}
}

func ReadIssuer(path string) (issuer []didtypes.DidInfo, err error) {
	data, err := ioutil.ReadFile(filepath.Join(path, issuerFilePath))
	if err != nil {
		return issuer, err
	}
	err = json.Unmarshal(data, &issuer)
	if err != nil {
		return issuer, err
	}
	return issuer, nil
}

func migrateRegionClassName(ctx sdk.Context, stakingKeeper *wstakingkeeper.Keeper, nftKeeper *wnftkeeper.Keeper) {
	regions := stakingKeeper.GetAllRegion(ctx)
	for _, regionObj := range regions {
		newClassId := regionObj.NftClassId[:len(regionObj.NftClassId)-1]
		class, found := nftKeeper.GetClass(ctx, regionObj.NftClassId)
		if found {
			nftKeeper.DeleteClass(ctx, class.Id)
			class.Id = newClassId
			class.Uri = utils.CalculateUriHash(class.Uri)
			err := nftKeeper.SaveClass(ctx, class)
			if err != nil {
				panic(err)
			}
		}
		regionObj.NftClassId = newClassId
		stakingKeeper.SetRegion(ctx, regionObj)
	}
}

func migrateGroup(ctx sdk.Context, path string, gk *groupkeeper.Keeper, sk *wstakingkeeper.Keeper, kk *kyckeeper.Keeper) {
	file, err := os.Open(filepath.Join(path, groupFilePath))
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var data struct {
		AppState struct {
			Group struct {
				GroupMembers []struct {
					GroupId string `json:"group_id,omitempty"`
					Member  struct {
						Address  string    `json:"address,omitempty"`
						Weight   string    `json:"weight,omitempty"`
						Metadata string    `json:"metadata,omitempty"`
						AddedAt  time.Time `json:"added_at"`
					} `json:"member,omitempty"`
				} `json:"group_members"`
				Groups []struct {
					Id          string    `json:"id,omitempty"`
					Admin       string    `json:"admin,omitempty"`
					Metadata    string    `json:"metadata,omitempty"`
					Version     string    `json:"version,omitempty"`
					TotalWeight string    `json:"total_weight,omitempty"`
					CreatedAt   time.Time `json:"created_at"`
					RegionID    string    `json:"regionID,omitempty"`
				} `json:"groups"`
			} `json:"group"`
		} `json:"app_state"`
	}
	if err := decoder.Decode(&data); err != nil {
		panic(fmt.Sprintf("Failed to decode JSON: %v", err))
	}

	lastGroupId := uint64(0)
	groupExist := make(map[string]string)
	groupAdmin := make(map[string]string)
	for _, groupInfoV1 := range data.AppState.Group.Groups {
		regionId := strings.ToLower(groupInfoV1.RegionID)
		if _, ok := groupExist[regionId]; ok {
			continue
		}

		groupExist[regionId] = groupInfoV1.Id

		if _, ok := groupAdmin[regionId]; !ok {
			region, found := sk.GetRegion(ctx, regionId)
			if !found {
				continue
			} else {
				addr, err := sdk.ValAddressFromBech32(region.OperatorAddress)
				if err != nil {
					panic(fmt.Sprintf("Failed to get operator address: %v", err))
				}
				groupAdmin[regionId] = sdk.AccAddress(addr).String()
			}
		}

		id, err := strconv.ParseUint(groupInfoV1.Id, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("Parse group id: %v", err))
		}
		if lastGroupId <= id {
			lastGroupId = id
		}
		version, err := strconv.ParseUint(groupInfoV1.Version, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("Parse group version: %v", err))
		}
		groupInfoV2 := megrouptypes.GroupInfo{
			Id:          id,
			Admin:       groupAdmin[regionId],
			Metadata:    groupInfoV1.Metadata,
			Version:     version,
			TotalWeight: groupInfoV1.TotalWeight,
			CreatedAt:   groupInfoV1.CreatedAt,
			RegionID:    regionId,
		}
		err = gk.AppendGroup(ctx, &groupInfoV2)
		if err != nil {
			panic(fmt.Sprintf("Failed to append group: %v", err))
		}
		gk.SetGroupToRegion(ctx, groupInfoV2.RegionID, groupInfoV2.Id)

		gk.SetMemberJoined(ctx, megrouptypes.MemberJoined{
			Address: groupInfoV2.Admin,
			GroupId: groupInfoV2.Id})

		gk.AddGroupMember(ctx, &megrouptypes.GroupMember{
			GroupId: groupInfoV2.Id,
			Member: &megrouptypes.Member{
				Address:  groupInfoV2.Admin,
				Weight:   math.NewInt(0).String(),
				Metadata: "",
				AddedAt:  groupInfoV1.CreatedAt,
			},
		})
		gk.SetGroupMemberCount(ctx, groupInfoV2.Id, 1)
	}

	gk.SetLastGroupID(ctx, lastGroupId)

	for _, memberV1 := range data.AppState.Group.GroupMembers {
		//groupId, err := strconv.ParseUint(memberV1.GroupId, 10, 64)
		//if err != nil {
		//	panic(fmt.Sprintf("Parse group id: %v", err))
		//}

		did, ok := kk.GetDID(ctx, sdk.MustAccAddressFromBech32(memberV1.Member.Address))
		if !ok {
			ctx.Logger().Error("adding group member has no did", "member address", memberV1.Member.Address)
			continue
		}
		kycData, ok := kk.GetKYC(ctx, did)
		if !ok {
			ctx.Logger().Error("adding group member has no kyc", "member address", memberV1.Member.Address)
			continue
		}
		regionId := string(kycData.Data)
		groupId, f := gk.GetGroupIdByRegion(ctx, regionId)
		if !f {
			continue
		}

		memberV2 := megrouptypes.GroupMember{
			GroupId: groupId,
			Member: &megrouptypes.Member{
				Address:  memberV1.Member.Address,
				Weight:   memberV1.Member.Weight,
				Metadata: memberV1.Member.Metadata,
				AddedAt:  memberV1.Member.AddedAt,
			},
		}

		gk.SetMemberJoined(ctx, megrouptypes.MemberJoined{
			Address: memberV2.Member.Address,
			GroupId: groupId,
		})

		err = gk.AddGroupMember(ctx, &memberV2)
		if err != nil {
			panic(fmt.Sprintf("Failed to add group member: %v", err))
		}

		grpNumber, found := gk.GetGroupMemberCount(ctx, groupId)
		if !found {
			grpNumber = 0
		}

		gk.SetGroupMemberCount(ctx, groupId, grpNumber+1)
	}
}

func MigrateDelegation(ctx sdk.Context, stakingKeeper *wstakingkeeper.Keeper, kk *kyckeeper.Keeper) {
	// experience region has no validator address
	expRegion, isFound := stakingKeeper.GetRegion(ctx, wstakingtypes.ExperienceRegionId)
	if !isFound {
		panic(fmt.Errorf("should have experience region"))
	}
	stakingKeeper.IterateAllDelegations(ctx, func(del stakingtypes.Delegation) (stop bool) {
		did, didFound := kk.GetDID(ctx, sdk.MustAccAddressFromBech32(del.DelegatorAddress))
		if didFound {
			kyc, kycFound := kk.GetKYC(ctx, did)
			if kycFound {
				region, regionFound := stakingKeeper.GetRegion(ctx, string(kyc.Data))
				if !regionFound {
					panic(fmt.Errorf("region not found: %s", string(kyc.Data)))
				}
				del.ValidatorAddress = region.OperatorAddress
			}
		} else {
			del.ValidatorAddress = expRegion.OperatorAddress
		}
		stakingKeeper.SetDelegation(ctx, del)
		return false
	})
}

func MigrateFixedDeposit(ctx sdk.Context, stakingKeeper *wstakingkeeper.Keeper, kk *kyckeeper.Keeper, bk wbankkeeper.BaseKeeperWrapper) {
	balance := bk.GetBalance(ctx, authtypes.NewModuleAddress(wstakingtypes.FixedDepositPrincipalPool), params.BaseDenom)

	fixedDeposits := stakingKeeper.GetAllFixedDeposit(ctx)
	totalDeposit := sdk.ZeroInt()
	for _, fixedDeposit := range fixedDeposits {
		totalDeposit = totalDeposit.Add(fixedDeposit.Principal.Amount)
		meid, ok := stakingKeeper.GetMeid(ctx, fixedDeposit.Account)
		if !ok {
			panic(fmt.Errorf("meid not found: %s", fixedDeposit.Account))
		}
		region, found := stakingKeeper.GetRegion(ctx, meid.RegionId)
		if !found {
			panic(fmt.Errorf("region not found: %s", meid.RegionId))
		}
		if region.FixedDepositAmount.IsNil() {
			region.FixedDepositAmount = sdk.ZeroInt()
		}
		region.FixedDepositAmount = region.FixedDepositAmount.Add(fixedDeposit.Principal.Amount)
		stakingKeeper.SetRegion(ctx, region)
	}

	if !balance.Amount.Equal(totalDeposit) {
		panic(fmt.Sprintf("total deposit amount is not equal to the balance: %s, %s", totalDeposit, balance.Amount))
	}

	totalDepositInRegion := sdk.ZeroInt()
	regions := stakingKeeper.GetAllRegion(ctx)
	for _, region := range regions {
		if region.FixedDepositAmount.IsNil() {
			region.FixedDepositAmount = sdk.ZeroInt()
		}
		totalDepositInRegion = totalDepositInRegion.Add(region.FixedDepositAmount)
	}
	if !balance.Amount.Equal(totalDepositInRegion) {
		panic(fmt.Sprintf("total deposit amount in region is not equal to the balance: %s, %s", totalDeposit, balance.Amount))
	}
}
