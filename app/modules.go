package app

import (
	"cosmossdk.io/x/evidence"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	feegrantmodule "cosmossdk.io/x/feegrant/module"
	"cosmossdk.io/x/nft"
	"cosmossdk.io/x/upgrade"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"fmt"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	packetforwardmiddleware "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
	"github.com/cosmos/ibc-go/modules/capability"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ibctransfer "github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/evmos/ethermint/x/feemarket"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"github.com/st-chain/me-hub/x/bsc"
	bsctypes "github.com/st-chain/me-hub/x/bsc/types"
	"github.com/st-chain/me-hub/x/dao"
	daotypes "github.com/st-chain/me-hub/x/dao/types"
	"github.com/st-chain/me-hub/x/did"
	didtypes "github.com/st-chain/me-hub/x/did/types"
	gravitytypes "github.com/st-chain/me-hub/x/gravity/types"
	"github.com/st-chain/me-hub/x/kyc"
	kyctypes "github.com/st-chain/me-hub/x/kyc/types"
	"github.com/st-chain/me-hub/x/tron"
	trontypes "github.com/st-chain/me-hub/x/tron/types"
	"github.com/st-chain/me-hub/x/wbank"
	wbanktypes "github.com/st-chain/me-hub/x/wbank/types"
	wdistr "github.com/st-chain/me-hub/x/wdistri"
	"github.com/st-chain/me-hub/x/wgov"
	"github.com/st-chain/me-hub/x/wmint"
	"github.com/st-chain/me-hub/x/wnft"
	"github.com/st-chain/me-hub/x/wstaking"
	wstakingtypes "github.com/st-chain/me-hub/x/wstaking/types"
	"slices"

	appparams "github.com/st-chain/me-hub/app/params"
	delayedackmodule "github.com/st-chain/me-hub/x/delayedack"
	denommetadatamodule "github.com/st-chain/me-hub/x/denommetadata"
	eibcmodule "github.com/st-chain/me-hub/x/eibc"
	groupmodule "github.com/st-chain/me-hub/x/megroup"
	groupTypes "github.com/st-chain/me-hub/x/megroup/types"
	rollappmodule "github.com/st-chain/me-hub/x/rollapp"
	sequencermodule "github.com/st-chain/me-hub/x/sequencer"

	delayedacktypes "github.com/st-chain/me-hub/x/delayedack/types"
	denommetadatamoduletypes "github.com/st-chain/me-hub/x/denommetadata/types"
	eibcmoduletypes "github.com/st-chain/me-hub/x/eibc/types"
	meevm "github.com/st-chain/me-hub/x/evm"
	rollappmoduletypes "github.com/st-chain/me-hub/x/rollapp/types"
	sequencermoduletypes "github.com/st-chain/me-hub/x/sequencer/types"
)

func (app *App) SetupModules(
	appCodec codec.Codec,
	bApp *baseapp.BaseApp,
	encodingConfig appparams.EncodingConfig,
	skipGenesisInvariants bool,
) []module.AppModule {
	return []module.AppModule{
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app, app.txConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, nil, app.GetSubspace(authtypes.ModuleName)),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, encodingConfig.InterfaceRegistry),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		wbank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, encodingConfig.InterfaceRegistry),
		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)),
		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),
		wgov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		wmint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, nil, app.GetSubspace(minttypes.ModuleName)),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName), app.interfaceRegistry),
		wdistr.NewAppModule(appCodec, *app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		wstaking.NewAppModule(appCodec, app.StakingKeeper, app.TransferKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(app.UpgradeKeeper, app.AccountKeeper.AddressCodec()),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		params.NewAppModule(app.ParamsKeeper),
		packetforwardmiddleware.NewAppModule(app.PacketForwardMiddlewareKeeper, app.GetSubspace(packetforwardtypes.ModuleName)),
		ibctransfer.NewAppModule(app.TransferKeeper),
		rollappmodule.NewAppModule(appCodec, app.RollappKeeper),
		sequencermodule.NewAppModule(appCodec, app.SequencerKeeper),
		delayedackmodule.NewAppModule(appCodec, app.DelayedAckKeeper),
		denommetadatamodule.NewAppModule(app.DenomMetadataKeeper, *app.EvmKeeper.Keeper, app.BankKeeper),
		eibcmodule.NewAppModule(appCodec, app.EIBCKeeper),

		// Ethermint app modules
		meevm.NewAppModule(app.EvmKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(evmtypes.ModuleName)),
		feemarket.NewAppModule(app.FeeMarketKeeper, app.GetSubspace(feemarkettypes.ModuleName).WithKeyTable(feemarkettypes.ParamKeyTable())),

		// did app modules
		did.NewAppModule(appCodec, app.DidKeeper),
		kyc.NewAppModule(appCodec, app.KycKeeper),

		// me-group
		groupmodule.NewAppModule(appCodec, *app.GroupKeeper),

		dao.NewAppModule(appCodec, app.DaoKeeper),
		wnft.NewAppModule(appCodec, *app.WNFTKeeper, app.AccountKeeper, app.BankKeeper, encodingConfig.InterfaceRegistry),
		wasm.NewAppModule(appCodec, &app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, bApp.MsgServiceRouter(), app.GetSubspace(wasmtypes.ModuleName)),

		bsc.NewAppModule(app.BscKeeper),
		tron.NewAppModule(app.TronKeeper),
	}
}

// ModuleAccountAddrs returns all the app's module account addresses.
func ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}
	return modAccAddrs
}

// module account permissions
var maccPerms = map[string][]string{
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

var PreBlockers = []string{
	upgradetypes.ModuleName,
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

func GenTxMessageValidator(msgs []sdk.Msg) error {
	if len(msgs) == 0 {
		return fmt.Errorf("unexpected number of GenTx messages; got: %d, expected great than 0", len(msgs))
	}
	if _, ok := msgs[0].(*stakingtypes.MsgCreateValidator); !ok {
		return fmt.Errorf("unexpected GenTx message type; expected: MsgCreateValidator, got: %T", msgs[0])
	}
	return nil
}

// We have custom migration order to make sure we run txfees first (we need it for iro migrations)
func CustomMigrationOrder(modules []string) []string {
	slices.Sort(modules)

	// run txfees first (we need it for iro migrations)
	//txfeesIndex := slices.Index(modules, txfeestypes.ModuleName)
	//out := append(modules[:txfeesIndex], modules[txfeesIndex+1:]...)
	//out = append([]string{txfeestypes.ModuleName}, out...)

	return modules
}
