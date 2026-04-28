package keepers

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	packetforwardmiddleware "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/types"
	ibctransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/evmos/ethermint/x/feemarket"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"github.com/openmetaearth/me-hub/x/bsc"
	bsctypes "github.com/openmetaearth/me-hub/x/bsc/types"
	"github.com/openmetaearth/me-hub/x/dao"
	daotypes "github.com/openmetaearth/me-hub/x/dao/types"
	"github.com/openmetaearth/me-hub/x/did"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	gravitytypes "github.com/openmetaearth/me-hub/x/gravity/types"
	"github.com/openmetaearth/me-hub/x/kyc"
	kyctypes "github.com/openmetaearth/me-hub/x/kyc/types"
	"github.com/openmetaearth/me-hub/x/tron"
	trontypes "github.com/openmetaearth/me-hub/x/tron/types"
	"github.com/openmetaearth/me-hub/x/wbank"
	wbanktypes "github.com/openmetaearth/me-hub/x/wbank/types"
	wdistr "github.com/openmetaearth/me-hub/x/wdistri"
	"github.com/openmetaearth/me-hub/x/wgov"
	"github.com/openmetaearth/me-hub/x/wmint"
	"github.com/openmetaearth/me-hub/x/wnft"
	"github.com/openmetaearth/me-hub/x/wstaking"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"

	appparams "github.com/openmetaearth/me-hub/app/params"
	delayedackmodule "github.com/openmetaearth/me-hub/x/delayedack"
	denommetadatamodule "github.com/openmetaearth/me-hub/x/denommetadata"
	eibcmodule "github.com/openmetaearth/me-hub/x/eibc"
	groupmodule "github.com/openmetaearth/me-hub/x/megroup"
	groupTypes "github.com/openmetaearth/me-hub/x/megroup/types"
	rollappmodule "github.com/openmetaearth/me-hub/x/rollapp"
	sequencermodule "github.com/openmetaearth/me-hub/x/sequencer"

	delayedacktypes "github.com/openmetaearth/me-hub/x/delayedack/types"
	denommetadatamoduletypes "github.com/openmetaearth/me-hub/x/denommetadata/types"
	eibcmoduletypes "github.com/openmetaearth/me-hub/x/eibc/types"
	meevm "github.com/openmetaearth/me-hub/x/evm"
	rollappmoduletypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	sequencermoduletypes "github.com/openmetaearth/me-hub/x/sequencer/types"
)

func (a *AppKeepers) SetupModules(
	appCodec codec.Codec,
	bApp *baseapp.BaseApp,
	encodingConfig appparams.EncodingConfig,
	skipGenesisInvariants bool,
) []module.AppModule {
	return []module.AppModule{
		genutil.NewAppModule(
			a.AccountKeeper, a.StakingKeeper, bApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, a.AccountKeeper, nil, a.GetSubspace(authtypes.ModuleName)),
		authzmodule.NewAppModule(appCodec, a.AuthzKeeper, a.AccountKeeper, a.BankKeeper, encodingConfig.InterfaceRegistry),
		vesting.NewAppModule(a.AccountKeeper, a.BankKeeper),
		wbank.NewAppModule(appCodec, a.BankKeeper, a.AccountKeeper, a.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *a.CapabilityKeeper, false),
		feegrantmodule.NewAppModule(appCodec, a.AccountKeeper, a.BankKeeper, a.FeeGrantKeeper, encodingConfig.InterfaceRegistry),
		crisis.NewAppModule(a.CrisisKeeper, skipGenesisInvariants, a.GetSubspace(crisistypes.ModuleName)),
		consensus.NewAppModule(appCodec, a.ConsensusParamsKeeper),
		wgov.NewAppModule(appCodec, a.GovKeeper, a.AccountKeeper, a.BankKeeper, a.GetSubspace(govtypes.ModuleName)),
		wmint.NewAppModule(appCodec, a.MintKeeper, a.AccountKeeper, nil, a.GetSubspace(minttypes.ModuleName)),
		slashing.NewAppModule(appCodec, a.SlashingKeeper, a.AccountKeeper, a.BankKeeper, a.StakingKeeper, a.GetSubspace(slashingtypes.ModuleName)),
		wdistr.NewAppModule(appCodec, *a.DistrKeeper, a.AccountKeeper, a.BankKeeper, a.StakingKeeper, a.GetSubspace(distrtypes.ModuleName)),
		wstaking.NewAppModule(appCodec, a.StakingKeeper, a.TransferKeeper, a.AccountKeeper, a.BankKeeper, a.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(a.UpgradeKeeper),
		evidence.NewAppModule(a.EvidenceKeeper),
		ibc.NewAppModule(a.IBCKeeper),
		params.NewAppModule(a.ParamsKeeper),
		packetforwardmiddleware.NewAppModule(a.PacketForwardMiddlewareKeeper, a.GetSubspace(packetforwardtypes.ModuleName)),
		ibctransfer.NewAppModule(a.TransferKeeper),
		rollappmodule.NewAppModule(appCodec, a.RollappKeeper, a.AccountKeeper, a.BankKeeper),
		sequencermodule.NewAppModule(appCodec, a.SequencerKeeper, a.AccountKeeper, a.BankKeeper),
		delayedackmodule.NewAppModule(appCodec, a.DelayedAckKeeper),
		denommetadatamodule.NewAppModule(a.DenomMetadataKeeper, *a.EvmKeeper, a.BankKeeper),
		eibcmodule.NewAppModule(appCodec, a.EIBCKeeper, a.AccountKeeper, a.BankKeeper),

		// Ethermint app modules
		meevm.NewAppModule(a.EvmKeeper, a.AccountKeeper, a.BankKeeper, a.GetSubspace(evmtypes.ModuleName)),
		feemarket.NewAppModule(a.FeeMarketKeeper, a.GetSubspace(feemarkettypes.ModuleName).WithKeyTable(feemarkettypes.ParamKeyTable())),

		// did app modules
		did.NewAppModule(appCodec, a.DidKeeper),
		kyc.NewAppModule(appCodec, a.KycKeeper),

		// me-group
		groupmodule.NewAppModule(appCodec, *a.GroupKeeper),

		dao.NewAppModule(appCodec, a.DaoKeeper),
		wnft.NewAppModule(appCodec, *a.WNFTKeeper, a.AccountKeeper, a.BankKeeper, encodingConfig.InterfaceRegistry),
		wasm.NewAppModule(appCodec, &a.WasmKeeper, a.StakingKeeper, a.AccountKeeper, a.BankKeeper, bApp.MsgServiceRouter(), a.GetSubspace(wasmtypes.ModuleName)),

		bsc.NewAppModule(a.BscKeeper),
		tron.NewAppModule(a.TronKeeper),
	}
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (*AppKeepers) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range MaccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}
	return modAccAddrs
}

// module account permissions
var MaccPerms = map[string][]string{
	authtypes.FeeCollectorName:                         nil,
	distrtypes.ModuleName:                              nil,
	wbanktypes.TreasuryPoolName:                        nil,
	minttypes.ModuleName:                               {authtypes.Minter},
	stakingtypes.BondedPoolName:                        {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName:                     {authtypes.Burner, authtypes.Staking},
	wstakingtypes.BondedStakePoolName:                  {authtypes.Burner, authtypes.Staking},
	wstakingtypes.NotBondedStakePoolName:               {authtypes.Burner, authtypes.Staking},
	wstakingtypes.StakePoolName:                        {authtypes.Staking},
	govtypes.ModuleName:                                {authtypes.Burner},
	ibctransfertypes.ModuleName:                        {authtypes.Minter, authtypes.Burner},
	sequencermoduletypes.ModuleName:                    {authtypes.Minter, authtypes.Burner, authtypes.Staking},
	rollappmoduletypes.ModuleName:                      {},
	evmtypes.ModuleName:                                {authtypes.Minter, authtypes.Burner}, // used for secure addition and subtraction of balance using module account.
	evmtypes.ModuleVirtualFrontierContractDeployerName: nil,                                  // used for deploying virtual frontier bank contract.
	wstakingtypes.FixedDepositPrincipalPool:            nil,
	wstakingtypes.BridgeFeePool:                        nil,
	wasmtypes.ModuleName:                               {authtypes.Burner},
	groupTypes.ModuleName:                              {authtypes.Minter, authtypes.Burner},
	nft.ModuleName:                                     nil,
	bsctypes.ModuleName:                                {authtypes.Minter, authtypes.Burner},
	trontypes.ModuleName:                               {authtypes.Minter, authtypes.Burner},
	gravitytypes.SlashingModuleAccount:                 {authtypes.Minter, authtypes.Burner},
}

var BeginBlockers = []string{
	upgradetypes.ModuleName,
	capabilitytypes.ModuleName,
	minttypes.ModuleName,
	distrtypes.ModuleName,
	slashingtypes.ModuleName,
	evidencetypes.ModuleName,
	stakingtypes.ModuleName,
	vestingtypes.ModuleName,
	feemarkettypes.ModuleName,
	evmtypes.ModuleName,
	ibcexported.ModuleName,
	ibctransfertypes.ModuleName,
	packetforwardtypes.ModuleName,
	authtypes.ModuleName,
	authz.ModuleName,
	banktypes.ModuleName,
	govtypes.ModuleName,
	crisistypes.ModuleName,
	genutiltypes.ModuleName,
	feegrant.ModuleName,
	paramstypes.ModuleName,
	rollappmoduletypes.ModuleName,
	sequencermoduletypes.ModuleName,
	denommetadatamoduletypes.ModuleName,
	delayedacktypes.ModuleName,
	eibcmoduletypes.ModuleName,
	consensusparamtypes.ModuleName,
	daotypes.ModuleName,
	wasmtypes.ModuleName,
	didtypes.ModuleName,
	kyctypes.ModuleName,
	nft.ModuleName,
	groupTypes.ModuleName,
	bsctypes.ModuleName,
	trontypes.ModuleName,
}

var EndBlockers = []string{
	crisistypes.ModuleName,
	govtypes.ModuleName,
	stakingtypes.ModuleName,
	capabilitytypes.ModuleName,
	authtypes.ModuleName,
	authz.ModuleName,
	banktypes.ModuleName,
	distrtypes.ModuleName,
	feemarkettypes.ModuleName,
	evmtypes.ModuleName,
	slashingtypes.ModuleName,
	vestingtypes.ModuleName,
	minttypes.ModuleName,
	genutiltypes.ModuleName,
	evidencetypes.ModuleName,
	feegrant.ModuleName,
	paramstypes.ModuleName,
	upgradetypes.ModuleName,
	ibcexported.ModuleName,
	ibctransfertypes.ModuleName,
	packetforwardtypes.ModuleName,
	rollappmoduletypes.ModuleName,
	sequencermoduletypes.ModuleName,
	denommetadatamoduletypes.ModuleName,
	delayedacktypes.ModuleName,
	eibcmoduletypes.ModuleName,
	consensusparamtypes.ModuleName,
	daotypes.ModuleName,
	wasmtypes.ModuleName,
	didtypes.ModuleName,
	kyctypes.ModuleName,
	nft.ModuleName,
	groupTypes.ModuleName,
	bsctypes.ModuleName,
	trontypes.ModuleName,
}

var InitGenesis = []string{
	capabilitytypes.ModuleName,
	authtypes.ModuleName,
	authz.ModuleName,
	banktypes.ModuleName,
	distrtypes.ModuleName,
	daotypes.ModuleName,
	stakingtypes.ModuleName,
	vestingtypes.ModuleName,
	slashingtypes.ModuleName,
	feemarkettypes.ModuleName,
	evmtypes.ModuleName,
	govtypes.ModuleName,
	minttypes.ModuleName,
	crisistypes.ModuleName,
	ibcexported.ModuleName,
	genutiltypes.ModuleName,
	evidencetypes.ModuleName,
	paramstypes.ModuleName,
	upgradetypes.ModuleName,
	ibctransfertypes.ModuleName,
	packetforwardtypes.ModuleName,
	feegrant.ModuleName,
	rollappmoduletypes.ModuleName,
	sequencermoduletypes.ModuleName,
	denommetadatamoduletypes.ModuleName, // must after `x/bank` to trigger hooks
	delayedacktypes.ModuleName,
	eibcmoduletypes.ModuleName,
	consensusparamtypes.ModuleName,
	wasmtypes.ModuleName,
	didtypes.ModuleName,
	kyctypes.ModuleName,
	nft.ModuleName,
	groupTypes.ModuleName,
	bsctypes.ModuleName,
	trontypes.ModuleName,
}
