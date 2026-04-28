package keepers

import (
	"fmt"
	bsctypes "github.com/openmetaearth/me-hub/x/bsc/types"
	trontypes "github.com/openmetaearth/me-hub/x/tron/types"
	"path/filepath"
	"strings"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/cosmos/cosmos-sdk/x/nft"
	gravitykeeper "github.com/openmetaearth/me-hub/x/gravity/keeper"
	groupTypes "github.com/openmetaearth/me-hub/x/megroup/types"

	wasmapp "github.com/CosmWasm/wasmd/app"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	nftkeeper "github.com/cosmos/cosmos-sdk/x/nft/keeper"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	packetforwardmiddleware "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward"
	packetforwardkeeper "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/keeper"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/types"
	ibctransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v7/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	ibcporttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	ibctestingtypes "github.com/cosmos/ibc-go/v7/testing/types"
	"github.com/evmos/ethermint/x/evm"
	ethermintevmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/evmos/ethermint/x/evm/vm/geth"
	feemarketkeeper "github.com/evmos/ethermint/x/feemarket/keeper"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"github.com/openmetaearth/me-hub/x/bridgingfee"
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
	groupkeeper "github.com/openmetaearth/me-hub/x/megroup/keeper"
	rollappmodule "github.com/openmetaearth/me-hub/x/rollapp"
	rollappmodulekeeper "github.com/openmetaearth/me-hub/x/rollapp/keeper"
	"github.com/openmetaearth/me-hub/x/rollapp/transfergenesis"
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
	ICS4Wrapper                   ibcporttypes.ICS4Wrapper
	delayedAckMiddleware          *delayedackmodule.IBCMiddleware
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

	RollappKeeper   *rollappmodulekeeper.Keeper
	SequencerKeeper sequencermodulekeeper.Keeper
	EIBCKeeper      eibckeeper.Keeper

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
	cdc *codec.LegacyAmino,
	bApp *baseapp.BaseApp,
	moduleAccountAddrs map[string]bool,
	skipUpgradeHeights map[int64]bool,
	invCheckPeriod uint,
	tracer, homePath string,
	appOpts servertypes.AppOptions,
	wasmOpts []wasmkeeper.Option,
) {
	govModuleAddress := authtypes.NewModuleAddress(govtypes.ModuleName).String()
	// init keepers

	a.ParamsKeeper = initParamsKeeper(appCodec, cdc, a.keys[paramstypes.StoreKey], a.tkeys[paramstypes.TStoreKey])
	// set the BaseApp's parameter store
	a.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(appCodec, a.keys[consensusparamtypes.StoreKey], govModuleAddress)
	bApp.SetParamStore(&a.ConsensusParamsKeeper)

	// add capability keeper and ScopeToModule for ibc module
	a.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, a.keys[capabilitytypes.StoreKey], a.memKeys[capabilitytypes.MemStoreKey])

	// grant capabilities for the ibc and ibc-transfer modules
	a.ScopedIBCKeeper = a.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	a.ScopedTransferKeeper = a.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedWasmKeeper := a.CapabilityKeeper.ScopeToModule(wasmtypes.ModuleName)

	a.CapabilityKeeper.Seal()

	a.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		a.keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		bApp,
		govModuleAddress,
	)

	a.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		a.keys[authtypes.StoreKey],
		authtypes.ProtoBaseAccount,
		MaccPerms,
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		govModuleAddress,
	)

	a.AuthzKeeper = authzkeeper.NewKeeper(
		a.keys[authz.ModuleName],
		appCodec,
		bApp.MsgServiceRouter(),
		a.AccountKeeper,
	)

	a.DaoKeeper = daokeeper.NewKeeper(
		appCodec,
		a.keys[daotypes.StoreKey],
		a.AccountKeeper,
	)

	a.BankKeeper = wbankkeeper.NewKeeper(
		appCodec,
		a.keys[banktypes.StoreKey],
		a.AccountKeeper,
		a.DaoKeeper,
		moduleAccountAddrs,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	a.CrisisKeeper = crisiskeeper.NewKeeper(
		appCodec, a.keys[crisistypes.StoreKey], invCheckPeriod, a.BankKeeper, authtypes.FeeCollectorName, govModuleAddress,
	)

	a.WNFTKeeper = wnftkeeper.NewKeeper(
		appCodec,
		a.keys[nftkeeper.StoreKey],
		a.AccountKeeper,
		a.BankKeeper,
	)

	a.StakingKeeper = wstakingkeeper.NewKeeper(
		appCodec,
		a.keys[stakingtypes.StoreKey],
		a.AccountKeeper,
		a.BankKeeper,
		a.DaoKeeper,
		a.WNFTKeeper,
		govModuleAddress,
	)

	a.MintKeeper = wmintkeeper.NewKeeper(
		appCodec,
		a.keys[minttypes.StoreKey],
		a.StakingKeeper,
		a.AccountKeeper,
		a.BankKeeper,
		wbanktypes.TreasuryPoolName,
		govModuleAddress,
	)
	a.StakingKeeper.SetMintKeeper(a.MintKeeper)

	a.DistrKeeper = wdistrkeeper.NewKeeper(
		appCodec,
		a.keys[distrtypes.StoreKey],
		a.GetSubspace(distrtypes.ModuleName),
		a.AccountKeeper,
		a.BankKeeper,
		a.StakingKeeper,
		wbanktypes.TreasuryPoolName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	a.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		cdc,
		a.keys[slashingtypes.StoreKey],
		a.StakingKeeper,
		govModuleAddress,
	)
	a.StakingKeeper.SetSlashingKeeper(a.SlashingKeeper)

	a.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		a.keys[feegrant.StoreKey],
		a.AccountKeeper,
	)

	// Create Ethermint keepers
	a.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		appCodec,
		sdk.MustAccAddressFromBech32(govModuleAddress),
		a.ConsensusParamsKeeper,
		a.keys[feemarkettypes.StoreKey],
		a.tkeys[feemarkettypes.TransientKey],
		a.GetSubspace(feemarkettypes.ModuleName),
	)

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
	)

	a.DenomMetadataKeeper = denommetadatamodulekeeper.NewKeeper(
		a.BankKeeper,
	)

	a.RollappKeeper = rollappmodulekeeper.NewKeeper(
		appCodec,
		a.keys[rollappmoduletypes.StoreKey],
		a.GetSubspace(rollappmoduletypes.ModuleName),
		a.IBCKeeper.ChannelKeeper,
		a.IBCKeeper.ClientKeeper,
		a.DaoKeeper,
	)

	a.SequencerKeeper = *sequencermodulekeeper.NewKeeper(
		appCodec,
		a.keys[sequencermoduletypes.StoreKey],
		a.keys[sequencermoduletypes.MemStoreKey],
		a.GetSubspace(sequencermoduletypes.ModuleName),
		a.BankKeeper,
		a.RollappKeeper,
	)

	a.EIBCKeeper = *eibckeeper.NewKeeper(
		appCodec,
		a.keys[eibcmoduletypes.StoreKey],
		a.keys[eibcmoduletypes.MemStoreKey],
		a.GetSubspace(eibcmoduletypes.ModuleName),
		a.AccountKeeper,
		a.BankKeeper,
		nil,
	)

	// Create Transfer Keepers
	a.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		a.keys[ibctransfertypes.StoreKey],
		a.GetSubspace(ibctransfertypes.ModuleName),
		denommetadatamodule.NewICS4Wrapper(a.IBCKeeper.ChannelKeeper, a.RollappKeeper, a.BankKeeper),
		a.IBCKeeper.ChannelKeeper,
		&a.IBCKeeper.PortKeeper,
		a.AccountKeeper,
		a.BankKeeper,
		a.ScopedTransferKeeper,
	)

	a.DelayedAckKeeper = *delayedackkeeper.NewKeeper(
		appCodec,
		a.keys[delayedacktypes.StoreKey],
		a.GetSubspace(delayedacktypes.ModuleName),
		a.RollappKeeper,
		a.IBCKeeper.ChannelKeeper,
		a.IBCKeeper.ChannelKeeper,
		&a.EIBCKeeper,
	)

	wasmDir := filepath.Join(homePath, "wasm")
	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic(fmt.Sprintf("error while reading wasm config: %s", err))
	}

	// Add custom query plugins for KYC module

	allWasmOpts := append(wasmOpts, a.SetupCustomMsgs())

	// The last arguments can contain custom message handlers, and custom query handlers,
	// if we want to allow any custom callbacks
	availableCapabilities := strings.Join(wasmapp.AllCapabilities(), ",")
	a.WasmKeeper = wasmkeeper.NewKeeper(
		appCodec,
		a.keys[wasmtypes.StoreKey],
		a.AccountKeeper,
		a.BankKeeper,
		a.StakingKeeper,
		distrkeeper.NewQuerier(a.DistrKeeper.Keeper),
		a.ICS4Wrapper, // ISC4 Wrapper: fee IBC middleware
		a.IBCKeeper.ChannelKeeper,
		&a.IBCKeeper.PortKeeper,
		scopedWasmKeeper,
		a.TransferKeeper,
		bApp.MsgServiceRouter(),
		bApp.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		availableCapabilities,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		allWasmOpts...,
	)

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
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(a.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(a.IBCKeeper.ClientKeeper)).
		AddRoute(rollappmoduletypes.RouterKey, rollappmodule.NewRollappProposalHandler(a.RollappKeeper)).
		AddRoute(denommetadatamoduletypes.RouterKey, denommetadatamodule.NewDenomMetadataProposalHandler(a.DenomMetadataKeeper)).
		AddRoute(evmtypes.RouterKey, evm.NewEvmProposalHandler(a.EvmKeeper.Keeper))

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	// If evidence needs to be handled for the app, set routes in router here and seal
	a.EvidenceKeeper = *evidencekeeper.NewKeeper(
		appCodec, a.keys[evidencetypes.StoreKey], a.StakingKeeper, a.SlashingKeeper,
	)

	govConfig := govtypes.DefaultConfig()
	a.GovKeeper = wgovkeeper.NewKeeper(
		appCodec, a.keys[govtypes.StoreKey], a.AccountKeeper, a.BankKeeper,
		a.StakingKeeper, bApp.MsgServiceRouter(), govConfig, govModuleAddress,
	)
	a.GovKeeper.SetLegacyRouter(govRouter)

	a.PacketForwardMiddlewareKeeper = packetforwardkeeper.NewKeeper(
		appCodec, a.keys[packetforwardtypes.StoreKey],
		a.TransferKeeper,
		a.IBCKeeper.ChannelKeeper,
		a.DistrKeeper,
		a.BankKeeper,
		a.IBCKeeper.ChannelKeeper,
		govModuleAddress,
	)
}

func (a *AppKeepers) InitTransferStack() {
	a.TransferStack = ibctransfer.NewIBCModule(a.TransferKeeper)
	a.TransferStack = bridgingfee.NewIBCModule(
		a.TransferStack.(ibctransfer.IBCModule),
		a.DelayedAckKeeper,
		a.TransferKeeper,
		a.AccountKeeper.GetModuleAddress(wstakingtypes.BridgeFeePool),
		*a.RollappKeeper,
	)
	a.TransferStack = packetforwardmiddleware.NewIBCMiddleware(
		a.TransferStack,
		a.PacketForwardMiddlewareKeeper,
		0,
		packetforwardkeeper.DefaultForwardTransferPacketTimeoutTimestamp,
		packetforwardkeeper.DefaultRefundTransferPacketTimeoutTimestamp,
	)

	a.TransferStack = denommetadatamodule.NewIBCModule(a.TransferStack, a.DenomMetadataKeeper, a.RollappKeeper)
	// already instantiated in SetupHooks()
	a.delayedAckMiddleware.Setup(
		delayedackmodule.WithIBCModule(a.TransferStack),
		delayedackmodule.WithKeeper(a.DelayedAckKeeper),
		delayedackmodule.WithRollappKeeper(a.RollappKeeper),
	)
	a.TransferStack = a.delayedAckMiddleware
	a.TransferStack = transfergenesis.NewIBCModule(a.TransferStack, a.DelayedAckKeeper, *a.RollappKeeper, a.TransferKeeper, a.DenomMetadataKeeper)
	a.TransferStack = transfergenesis.NewIBCModuleCanonicalChannelHack(a.TransferStack, *a.RollappKeeper, a.IBCKeeper.ChannelKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := ibcporttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, a.TransferStack)
	a.IBCKeeper.SetRouter(ibcRouter)
}

func (a *AppKeepers) SetupHooks() {
	a.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(a.DistrKeeper.Hooks(), a.SlashingKeeper.Hooks()),
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
	a.delayedAckMiddleware = delayedackmodule.NewIBCMiddleware()
	// register the rollapp hooks
	a.RollappKeeper.SetHooks(rollappmoduletypes.NewMultiRollappHooks(
		// insert rollapp hooks receivers here
		a.SequencerKeeper.RollappHooks(),
		a.delayedAckMiddleware,
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

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govv1.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(packetforwardtypes.ModuleName).WithKeyTable(packetforwardtypes.ParamKeyTable())
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibcexported.ModuleName)
	paramsKeeper.Subspace(rollappmoduletypes.ModuleName)
	paramsKeeper.Subspace(sequencermoduletypes.ModuleName)
	paramsKeeper.Subspace(denommetadatamoduletypes.ModuleName)
	paramsKeeper.Subspace(delayedacktypes.ModuleName)
	paramsKeeper.Subspace(eibcmoduletypes.ModuleName)

	// ethermint subspaces
	paramsKeeper.Subspace(evmtypes.ModuleName)
	paramsKeeper.Subspace(feemarkettypes.ModuleName)

	// did subspace
	paramsKeeper.Subspace(didtypes.ModuleName)
	paramsKeeper.Subspace(kyctypes.ModuleName)

	paramsKeeper.Subspace(wasmtypes.ModuleName)
	paramsKeeper.Subspace(nft.ModuleName)
	paramsKeeper.Subspace(groupTypes.ModuleName)
	return paramsKeeper
}
