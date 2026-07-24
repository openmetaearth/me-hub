package app

import (
	"context"
	"fmt"
	"path/filepath"

	"cosmossdk.io/log"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	ethflags "github.com/evmos/ethermint/server/flags"
	appparams "github.com/openmetaearth/me-hub/app/params"
	"github.com/spf13/cast"

	bsctypes "github.com/openmetaearth/me-hub/x/bsc/types"
	trontypes "github.com/openmetaearth/me-hub/x/tron/types"

	gravitykeeper "github.com/openmetaearth/me-hub/x/gravity/keeper"
	groupTypes "github.com/openmetaearth/me-hub/x/megroup/types"

	storetypes "cosmossdk.io/store/types"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	nftkeeper "cosmossdk.io/x/nft/keeper"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	wasmapp "github.com/CosmWasm/wasmd/app"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	packetforwardkeeper "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/keeper"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v8/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	ibcporttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	ibctestingtypes "github.com/cosmos/ibc-go/v8/testing/types"
	"github.com/evmos/ethermint/x/evm"
	ethermintevmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/evmos/ethermint/x/evm/vm/geth"
	feemarketkeeper "github.com/evmos/ethermint/x/feemarket/keeper"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	metypes "github.com/openmetaearth/me-hub/types"
	daokeeper "github.com/openmetaearth/me-hub/x/dao/keeper"
	daotypes "github.com/openmetaearth/me-hub/x/dao/types"
	delayedackmodule "github.com/openmetaearth/me-hub/x/delayedack"
	delayedackkeeper "github.com/openmetaearth/me-hub/x/delayedack/keeper"
	delayedacktypes "github.com/openmetaearth/me-hub/x/delayedack/types"
	denommetadatamodule "github.com/openmetaearth/me-hub/x/denommetadata"
	denommetadatamodulekeeper "github.com/openmetaearth/me-hub/x/denommetadata/keeper"
	denommetadatamoduletypes "github.com/openmetaearth/me-hub/x/denommetadata/types"
	didkeeper "github.com/openmetaearth/me-hub/x/did/keeper"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	eibckeeper "github.com/openmetaearth/me-hub/x/eibc/keeper"
	eibcmoduletypes "github.com/openmetaearth/me-hub/x/eibc/types"
	evmkeeper "github.com/openmetaearth/me-hub/x/evm/keeper"
	kyckeeper "github.com/openmetaearth/me-hub/x/kyc/keeper"
	kyctypes "github.com/openmetaearth/me-hub/x/kyc/types"
	lightclientmodulekeeper "github.com/openmetaearth/me-hub/x/lightclient/keeper"
	lightclientmoduletypes "github.com/openmetaearth/me-hub/x/lightclient/types"
	groupkeeper "github.com/openmetaearth/me-hub/x/megroup/keeper"
	"github.com/openmetaearth/me-hub/x/rollapp/genesisbridge"
	rollappmodulekeeper "github.com/openmetaearth/me-hub/x/rollapp/keeper"
	rollappmoduletypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	sequencermodulekeeper "github.com/openmetaearth/me-hub/x/sequencer/keeper"
	sequencermoduletypes "github.com/openmetaearth/me-hub/x/sequencer/types"
	vfchooks "github.com/openmetaearth/me-hub/x/vfc/hooks"
	wbankkeeper "github.com/openmetaearth/me-hub/x/wbank/keeper"
	wbanktypes "github.com/openmetaearth/me-hub/x/wbank/types"
	wdistrkeeper "github.com/openmetaearth/me-hub/x/wdistri/keeper"
	wgovkeeper "github.com/openmetaearth/me-hub/x/wgov/keeper"
	wmintkeeper "github.com/openmetaearth/me-hub/x/wmint/keeper"
	wnftkeeper "github.com/openmetaearth/me-hub/x/wnft/keeper"
	wstakingkeeper "github.com/openmetaearth/me-hub/x/wstaking/keeper"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

type GravityKeepers struct {
	BscKeeper  gravitykeeper.Keeper
	TronKeeper gravitykeeper.Keeper
}

type AppKeepers struct {
	// keepers
	AccountKeeper                 authkeeper.AccountKeeper
	AuthzKeeper                   authzkeeper.Keeper
	BankKeeper                    wbankkeeper.BaseKeeperWrapper
	CapabilityKeeper              *capabilitykeeper.Keeper
	StakingKeeper                 *wstakingkeeper.Keeper
	SlashingKeeper                slashingkeeper.Keeper
	MintKeeper                    wmintkeeper.Keeper
	DistrKeeper                   *wdistrkeeper.Keeper
	GovKeeper                     *wgovkeeper.Keeper
	CrisisKeeper                  *crisiskeeper.Keeper
	UpgradeKeeper                 *upgradekeeper.Keeper
	ParamsKeeper                  paramskeeper.Keeper
	IBCKeeper                     *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	TransferStack                 ibcporttypes.IBCModule
	DelayedAckMiddleware          *delayedackmodule.IBCMiddleware
	EvidenceKeeper                evidencekeeper.Keeper
	TransferKeeper                ibctransferkeeper.Keeper
	FeeGrantKeeper                feegrantkeeper.Keeper
	PacketForwardMiddlewareKeeper *packetforwardkeeper.Keeper
	ConsensusParamsKeeper         consensusparamkeeper.Keeper

	// Ethermint keepers
	EvmKeeper       *evmkeeper.Keeper
	FeeMarketKeeper feemarketkeeper.Keeper

	// did keeper
	DidKeeper *didkeeper.Keeper
	KycKeeper *kyckeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper

	RollappKeeper     *rollappmodulekeeper.Keeper
	SequencerKeeper   *sequencermodulekeeper.Keeper
	LightClientKeeper lightclientmodulekeeper.Keeper
	EIBCKeeper        eibckeeper.Keeper

	DelayedAckKeeper    delayedackkeeper.Keeper
	DenomMetadataKeeper *denommetadatamodulekeeper.Keeper
	DaoKeeper           daokeeper.Keeper
	WNFTKeeper          *wnftkeeper.Keeper
	WasmKeeper          wasmkeeper.Keeper
	GroupKeeper         *groupkeeper.Keeper

	GravityRouterKeeper gravitykeeper.RouterKeeper
	GravityKeepers

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey
}

// InitKeepers initializes all keepers for the app
func (a *AppKeepers) InitKeepers(
	appCodec codec.Codec,
	legacyAmino *codec.LegacyAmino,
	bApp *baseapp.BaseApp,
	logger log.Logger,
	moduleAccountAddrs map[string]bool,
	appOpts servertypes.AppOptions,
	wasmOpts []wasmkeeper.Option,
) {
	govModuleAddress := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	// get skipUpgradeHeights from the app options
	skipUpgradeHeights := map[int64]bool{}
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}
	homePath := cast.ToString(appOpts.Get(flags.FlagHome))
	tracer := cast.ToString(appOpts.Get(ethflags.EVMTracer))
	invCheckPeriod := cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod))

	// init keepers
	a.ParamsKeeper = initParamsKeeper(appCodec, legacyAmino, a.keys[paramstypes.StoreKey], a.tkeys[paramstypes.TStoreKey])
	// set the BaseApp's parameter store
	a.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(appCodec, runtime.NewKVStoreService(a.keys[consensusparamtypes.StoreKey]), govModuleAddress, runtime.EventService{})
	bApp.SetParamStore(a.ConsensusParamsKeeper.ParamsStore)

	// add capability keeper and ScopeToModule for ibc module
	a.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, a.keys[capabilitytypes.StoreKey], a.memKeys[capabilitytypes.MemStoreKey])
	// grant capabilities for the ibc and ibc-transfer modules
	a.ScopedIBCKeeper = a.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	a.ScopedTransferKeeper = a.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedWasmKeeper := a.CapabilityKeeper.ScopeToModule(wasmtypes.ModuleName)
	a.CapabilityKeeper.Seal()

	// set the governance module account as the authority for conducting upgrades
	a.UpgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, runtime.NewKVStoreService(a.keys[upgradetypes.StoreKey]), appCodec, homePath, bApp, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	a.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(a.keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		MaccPerms,
		authcodec.NewBech32Codec(appparams.AccountAddressPrefix),
		appparams.AccountAddressPrefix,
		authtypes.NewModuleAddress(govtypes.ModuleName).String())

	a.AuthzKeeper = authzkeeper.NewKeeper(runtime.NewKVStoreService(a.keys[authzkeeper.StoreKey]), appCodec, bApp.MsgServiceRouter(), a.AccountKeeper)

	a.DaoKeeper = daokeeper.NewKeeper(
		appCodec,
		a.keys[daotypes.StoreKey],
		a.AccountKeeper,
	)

	a.BankKeeper = wbankkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(a.keys[banktypes.StoreKey]),
		a.AccountKeeper,
		a.DaoKeeper,
		moduleAccountAddrs,
		govModuleAddress,
		logger,
	)

	a.CrisisKeeper = crisiskeeper.NewKeeper(appCodec, runtime.NewKVStoreService(a.keys[crisistypes.StoreKey]), invCheckPeriod,
		a.BankKeeper, authtypes.FeeCollectorName, authtypes.NewModuleAddress(govtypes.ModuleName).String(), a.AccountKeeper.AddressCodec())

	a.WNFTKeeper = wnftkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(a.keys[nftkeeper.StoreKey]),
		a.AccountKeeper,
		a.BankKeeper,
	)

	a.StakingKeeper = wstakingkeeper.NewKeeper(
		appCodec,
		a.keys[stakingtypes.StoreKey],
		runtime.NewKVStoreService(a.keys[stakingtypes.StoreKey]),
		a.AccountKeeper,
		a.BankKeeper,
		a.DaoKeeper,
		a.WNFTKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		authcodec.NewBech32Codec(appparams.AccountAddressPrefix+"valoper"),
		authcodec.NewBech32Codec(appparams.AccountAddressPrefix+"valcons"),
	)

	a.MintKeeper = wmintkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(a.keys[minttypes.StoreKey]),
		a.StakingKeeper,
		a.AccountKeeper,
		a.BankKeeper,
		wbanktypes.TreasuryPoolName,
		govModuleAddress,
	)
	a.StakingKeeper.SetMintKeeper(a.MintKeeper)

	a.DistrKeeper = wdistrkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(a.keys[distrtypes.StoreKey]),
		a.AccountKeeper,
		a.BankKeeper,
		a.StakingKeeper,
		a.StakingKeeper,
		wbanktypes.TreasuryPoolName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	a.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		legacyAmino,
		runtime.NewKVStoreService(a.keys[slashingtypes.StoreKey]),
		a.StakingKeeper,
		govModuleAddress,
	)
	a.StakingKeeper.SetSlashingKeeper(a.SlashingKeeper)

	a.FeeGrantKeeper = feegrantkeeper.NewKeeper(appCodec, runtime.NewKVStoreService(a.keys[feegrant.StoreKey]), a.AccountKeeper)

	// Create Ethermint keepers
	a.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		appCodec,
		sdk.MustAccAddressFromBech32(govModuleAddress),
		a.ConsensusParamsKeeper,
		a.keys[feemarkettypes.StoreKey],
		a.tkeys[feemarkettypes.TransientKey],
		a.GetSubspace(feemarkettypes.ModuleName),
	)

	// Seed global chain id before keeper init (supports legacy "mechain").
	initGlobalChainID(homePath, appOpts)

	// Create evmos keeper
	a.EvmKeeper = evmkeeper.NewKeeper(
		ethermintevmkeeper.NewKeeper(
			appCodec,
			a.keys[evmtypes.StoreKey],
			a.tkeys[evmtypes.TransientKey],
			authtypes.NewModuleAddress(govtypes.ModuleName),
			a.AccountKeeper,
			a.BankKeeper,
			a.StakingKeeper,
			a.FeeMarketKeeper,
			nil,
			geth.NewEVM,
			tracer,
			a.GetSubspace(evmtypes.ModuleName),
		))

	// Create IBC Keeper
	a.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		a.keys[ibcexported.StoreKey],
		a.GetSubspace(ibcexported.ModuleName),
		a.StakingKeeper,
		a.UpgradeKeeper,
		a.ScopedIBCKeeper,
		govModuleAddress,
	)

	// Rollapp/Sequencer/LightClient have circular deps; construct then wire setters.
	a.RollappKeeper = rollappmodulekeeper.NewKeeper(
		appCodec,
		a.keys[rollappmoduletypes.StoreKey],
		a.IBCKeeper.ChannelKeeper,
		nil,
		a.BankKeeper,
		a.TransferKeeper,
		govModuleAddress,
		nil,
	)

	a.SequencerKeeper = sequencermodulekeeper.NewKeeper(
		appCodec,
		a.keys[sequencermoduletypes.StoreKey],
		a.BankKeeper,
		a.AccountKeeper,
		a.RollappKeeper,
		govModuleAddress,
	)

	a.LightClientKeeper = *lightclientmodulekeeper.NewKeeper(
		appCodec,
		a.keys[lightclientmoduletypes.StoreKey],
		a.IBCKeeper.ClientKeeper,
		a.IBCKeeper.ConnectionKeeper,
		a.IBCKeeper.ChannelKeeper,
		a.SequencerKeeper,
		a.RollappKeeper,
	)

	a.SequencerKeeper.SetUnbondBlockers(a.RollappKeeper, &a.LightClientKeeper)
	a.SequencerKeeper.SetHooks(sequencermoduletypes.MultiHooks{
		rollappmodulekeeper.SequencerHooks{Keeper: a.RollappKeeper},
	})

	a.RollappKeeper.SetSequencerKeeper(a.SequencerKeeper)
	a.RollappKeeper.SetCanonicalClientKeeper(a.LightClientKeeper)
	a.RollappKeeper.SetDaoKeeper(a.DaoKeeper)

	a.DenomMetadataKeeper = denommetadatamodulekeeper.NewKeeper(
		a.BankKeeper,
		a.RollappKeeper,
	)

	a.EIBCKeeper = *eibckeeper.NewKeeper(
		appCodec,
		a.keys[eibcmoduletypes.StoreKey],
		a.memKeys[eibcmoduletypes.MemStoreKey],
		a.AccountKeeper,
		a.BankKeeper,
		a.DelayedAckKeeper,
		a.RollappKeeper,
		govModuleAddress,
	)

	// ICS4 stack (no ratelimit in phase 1): ChannelKeeper -> denommetadata -> genesisbridge
	ics4Wrapper := genesisbridge.NewICS4Wrapper(
		denommetadatamodule.NewICS4Wrapper(a.IBCKeeper.ChannelKeeper, a.RollappKeeper, a.BankKeeper),
		a.RollappKeeper,
		a.IBCKeeper.ChannelKeeper,
	)

	a.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		a.keys[ibctransfertypes.StoreKey],
		a.GetSubspace(ibctransfertypes.ModuleName),
		ics4Wrapper,
		a.IBCKeeper.ChannelKeeper,
		a.IBCKeeper.PortKeeper,
		a.AccountKeeper,
		BankKeeperWithoutSetMetadata{a.BankKeeper},
		a.ScopedTransferKeeper,
		govModuleAddress,
	)
	a.RollappKeeper.SetTransferKeeper(a.TransferKeeper)

	a.DelayedAckKeeper = *delayedackkeeper.NewKeeper(
		appCodec,
		a.keys[delayedacktypes.StoreKey],
		a.keys[ibcexported.StoreKey],
		govModuleAddress,
		a.RollappKeeper,
		a.IBCKeeper.ChannelKeeper,
		a.IBCKeeper.ChannelKeeper,
		&a.EIBCKeeper,
	)

	// Conditionally initialize WasmKeeper to avoid file lock issues in CLI commands
	// Only 'med start' needs WasmKeeper, other CLI commands can skip it
	skipWasmInit := cast.ToBool(appOpts.Get("skip-wasm-init"))

	if !skipWasmInit {
		wasmDir := filepath.Join(homePath, "wasm")
		wasmConfig, err := wasm.ReadWasmConfig(appOpts)
		if err != nil {
			panic(fmt.Sprintf("error while reading wasm config: %s", err))
		}

		// Add custom query plugins for KYC module
		allWasmOpts := append(wasmOpts, a.SetupCustomMsgs())

		// The last arguments can contain custom message handlers, and custom query handlers,
		// if we want to allow any custom callbacks
		a.WasmKeeper = wasmkeeper.NewKeeper(
			appCodec,
			runtime.NewKVStoreService(a.keys[wasmtypes.StoreKey]),
			a.AccountKeeper,
			a.BankKeeper,
			a.StakingKeeper,
			distrkeeper.NewQuerier(a.DistrKeeper.Keeper),
			ics4Wrapper, // ISC4 Wrapper: genesisbridge + denommetadata
			a.IBCKeeper.ChannelKeeper,
			a.IBCKeeper.PortKeeper,
			scopedWasmKeeper,
			a.TransferKeeper,
			bApp.MsgServiceRouter(),
			bApp.GRPCQueryRouter(),
			wasmDir,
			wasmConfig,
			wasmapp.AllCapabilities(),
			authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			allWasmOpts...,
		)
		logger.Info("WasmKeeper initialized")
	} else {
		logger.Info("WasmKeeper initialization skipped (CLI mode)")
		// WasmKeeper remains nil, wasm module will not be functional but won't block CLI
	}

	// Create did Keepers
	a.DidKeeper = didkeeper.NewKeeper(
		appCodec,
		a.keys[didtypes.StoreKey],
		a.DaoKeeper,
	)

	a.KycKeeper = kyckeeper.NewKeeper(
		appCodec,
		a.keys[kyctypes.StoreKey],
		a.StakingKeeper,
		a.AccountKeeper,
		a.DidKeeper,
		a.WNFTKeeper,
	)
	a.StakingKeeper.SetKycKeeper(a.KycKeeper)
	a.StakingKeeper.SetDidKeeper(a.DidKeeper)
	a.DaoKeeper.SetHook(a.KycKeeper)

	a.EIBCKeeper.SetDelayedAckKeeper(a.DelayedAckKeeper)
	a.GroupKeeper = groupkeeper.NewKeeper(
		appCodec,
		a.keys[groupTypes.StoreKey],
		a.GetSubspace(groupTypes.ModuleName),
		a.AccountKeeper,
		a.BankKeeper,
		a.StakingKeeper,
		a.DaoKeeper,
		a.KycKeeper,
	)
	a.StakingKeeper.SetGroupKeeper(a.GroupKeeper)

	a.BscKeeper = gravitykeeper.NewKeeper(
		bsctypes.ModuleName,
		appCodec,
		a.keys[bsctypes.StoreKey],
		a.BankKeeper,
		a.AccountKeeper,
		a.DaoKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	a.TronKeeper = gravitykeeper.NewKeeper(
		trontypes.ModuleName,
		appCodec,
		a.keys[trontypes.StoreKey],
		a.BankKeeper,
		a.AccountKeeper,
		a.DaoKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// add gravity router
	gravityRouter := gravitykeeper.NewRouter()
	gravityRouter.
		AddRoute(bsctypes.ModuleName, gravitykeeper.NewModuleHandler(a.BscKeeper)).
		AddRoute(trontypes.ModuleName, gravitykeeper.NewModuleHandler(a.TronKeeper))
	a.GravityRouterKeeper = gravitykeeper.NewRouterKeeper(gravityRouter)

	// Register the proposal types
	// Deprecated: Avoid adding new handlers, instead use the new proposal flow
	// by granting the governance module the right to execute the message.
	// See: https://github.com/cosmos/cosmos-sdk/blob/release/v0.46.x/x/gov/spec/01_concepts.md#proposal-messages
	govRouter := govv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(a.ParamsKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(a.IBCKeeper.ClientKeeper)).
		AddRoute(denommetadatamoduletypes.RouterKey, denommetadatamodule.NewDenomMetadataProposalHandler(a.DenomMetadataKeeper)).
		AddRoute(evmtypes.RouterKey, evm.NewEvmProposalHandler(a.EvmKeeper.Keeper))

	govConfig := govtypes.DefaultConfig()
	a.GovKeeper = wgovkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(a.keys[govtypes.StoreKey]),
		a.AccountKeeper,
		a.BankKeeper,
		a.StakingKeeper,
		a.DistrKeeper,
		bApp.MsgServiceRouter(),
		govConfig,
		govModuleAddress,
	)
	a.GovKeeper.SetLegacyRouter(govRouter)

	// create evidence keeper with router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(a.keys[evidencetypes.StoreKey]), a.StakingKeeper, a.SlashingKeeper, a.AccountKeeper.AddressCodec(), runtime.ProvideCometInfoService(),
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	a.EvidenceKeeper = *evidenceKeeper

	a.PacketForwardMiddlewareKeeper = packetforwardkeeper.NewKeeper(
		appCodec, a.keys[packetforwardtypes.StoreKey],
		a.TransferKeeper,
		a.IBCKeeper.ChannelKeeper,
		a.BankKeeper,
		a.IBCKeeper.ChannelKeeper,
		govModuleAddress,
	)
}

func (a *AppKeepers) SetupHooks() {
	a.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			a.DistrKeeper.Hooks(),
			a.SlashingKeeper.Hooks()),
	)
	a.StakingKeeper.SetWstakingHooks(
		wstakingtypes.NewMultiWstakingHooks(wstakingtypes.NewMultiWstakingHooks(a.DistrKeeper.Hooks())),
	)
	// register the staking hooks
	a.DenomMetadataKeeper.SetHooks(
		denommetadatamoduletypes.NewMultiDenomMetadataHooks(
			vfchooks.NewVirtualFrontierBankContractRegistrationHook(*a.EvmKeeper.Keeper),
		),
	)

	a.DelayedAckKeeper.SetHooks(delayedacktypes.NewMultiDelayedAckHooks(
		// insert delayedAck hooks receivers here
		a.EIBCKeeper.GetDelayedAckHooks(),
	))

	a.EIBCKeeper.SetHooks(eibcmoduletypes.NewMultiEIBCHooks(
		// insert eibc hooks receivers here
		a.DelayedAckKeeper.GetEIBCHooks(),
	))

	// dependencies injected in InitTransferStack()
	a.DelayedAckMiddleware = delayedackmodule.NewIBCMiddleware()
	// register the rollapp hooks
	a.RollappKeeper.SetHooks(rollappmoduletypes.NewMultiRollappHooks(
		// insert rollapp hooks receivers here
		a.SequencerKeeper.RollappHooks(),
		a.DelayedAckKeeper, // OnHardFork; AfterStateFinalized is stub (packets finalized via Msg)
		a.LightClientKeeper.RollappHooks(),
	))
}

// GetIBCKeeper implements ibctesting.TestingApp
func (a *AppKeepers) GetIBCKeeper() *ibckeeper.Keeper {
	return a.IBCKeeper
}

// GetScopedIBCKeeper implements ibctesting.TestingApp
func (a *AppKeepers) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper {
	return a.ScopedIBCKeeper
}

// GetStakingKeeper implements ibctesting.TestingApp
func (a *AppKeepers) GetStakingKeeper() ibctestingtypes.StakingKeeper {
	return a.StakingKeeper
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	// ibc-go subspaces
	// register the key tables for legacy param subspaces
	keyTable := ibcclienttypes.ParamKeyTable()
	keyTable.RegisterParamSet(&ibcconnectiontypes.Params{})
	paramsKeeper.Subspace(ibcexported.ModuleName).WithKeyTable(keyTable)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName).WithKeyTable(ibctransfertypes.ParamKeyTable())

	// ethermint subspaces (keeper doesn't load key table so we do it manually)
	paramsKeeper.Subspace(evmtypes.ModuleName).WithKeyTable(evmtypes.ParamKeyTable())
	paramsKeeper.Subspace(feemarkettypes.ModuleName).WithKeyTable(feemarkettypes.ParamKeyTable())

	// Legacy gov subspace (needed for x/gov store migrations that still receive ParamSubspace).
	//nolint:staticcheck // ParamKeyTable is deprecated but required for upgrade migrations.
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govv1.ParamKeyTable())

	// Register subspaces for custom modules that still use x/params
	paramsKeeper.Subspace(groupTypes.ModuleName)
	paramsKeeper.Subspace(packetforwardtypes.ModuleName)

	return paramsKeeper
}

func initGlobalChainID(homePath string, appOpts servertypes.AppOptions) {
	if chainID := cast.ToString(appOpts.Get(flags.FlagChainID)); chainID != "" {
		metypes.SetChainId(chainID)
		return
	}

	genesisFile := filepath.Join(homePath, "config", "genesis.json")
	genDoc, err := tmtypes.GenesisDocFromFile(genesisFile)
	if err == nil && genDoc.ChainID != "" {
		metypes.SetChainId(genDoc.ChainID)
	}
}

// this is a workaround to get rid of the denommetadata set automatically by ibc-go v8.x
// it has 2 issues:
// - it's not valid metadata struct
// - it has no exponent
// we disable this feature by providing bank keeper that does nothing on SetDenomMetaData
type BankKeeperWithoutSetMetadata struct {
	ibctransfertypes.BankKeeper
}

func (bk BankKeeperWithoutSetMetadata) SetDenomMetaData(ctx context.Context, denomMetaData banktypes.Metadata) {
}
