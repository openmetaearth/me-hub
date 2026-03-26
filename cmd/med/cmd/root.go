package cmd

import (
	"errors"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	txmodule "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/types/mempool"

	"github.com/evmos/ethermint/crypto/hd"
	ethermintserver "github.com/evmos/ethermint/server"
	mecli "github.com/st-chain/me-hub/client/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/rpc"

	"cosmossdk.io/log"
	cometbftcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"
	cometbftcli "github.com/cometbft/cometbft/libs/cli"
	dbm "github.com/cosmos/cosmos-db"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	// this line is used by starport scaffolding # root/moduleImport

	"github.com/st-chain/me-hub/app"
	appparams "github.com/st-chain/me-hub/app/params"

	ethermintclient "github.com/evmos/ethermint/client"
	ethservercfg "github.com/evmos/ethermint/server/config"
)

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) any {
	return nil
}

// NewRootCmd creates a new root command for me hub
func NewRootCmd() *cobra.Command {
	initSDKConfig()
	tempApp := app.New(log.NewNopLogger(), dbm.NewMemDB(), nil, true, EmptyAppOptions{})

	encodingConfig := appparams.EncodingConfig{
		InterfaceRegistry: tempApp.InterfaceRegistry(),
		Codec:             tempApp.AppCodec(),
		TxConfig:          tempApp.TxConfig(),
		Amino:             tempApp.LegacyAmino(),
	}

	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Codec).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithKeyringOptions(hd.EthSecp256k1Option()).
		WithHomeDir(app.DefaultNodeHome).
		WithViper("")

	rootCmd := &cobra.Command{
		Use: "med",
		Short: `
 __   __  _______  _______  _______    _______  _______  ______    _______  __   __ 
|  |_|  ||       ||       ||   _   |  |       ||   _   ||    _ |  |       ||  | |  |
|       ||    ___||_     _||  |_|  |  |    ___||  |_|  ||   | ||  |_     _||  |_|  |
|       ||   |___   |   |  |       |  |   |___ |       ||   |_||_   |   |  |       |
|       ||    ___|  |   |  |       |  |    ___||       ||    __  |  |   |  |       |
| ||_|| ||   |___   |   |  |   _   |  |   |___ |   _   ||   |  | |  |   |  |   _   |
|_|   |_||_______|  |___|  |__| |__|  |_______||__| |__||___|  |_|  |___|  |__| |__|
		`,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())
			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			// This needs to go after ReadFromClientConfig, as that function
			// sets the RPC client needed for SIGN_MODE_TEXTUAL.
			enabledSignModes := append(tx.DefaultSignModes, signing.SignMode_SIGN_MODE_TEXTUAL)
			txConfigOpts := tx.ConfigOptions{
				EnabledSignModes:           enabledSignModes,
				TextualCoinMetadataQueryFn: txmodule.NewGRPCCoinMetadataQueryFn(initClientCtx),
			}
			txConfigWithTextual, err := tx.NewTxConfigWithOptions(
				codec.NewProtoCodec(encodingConfig.InterfaceRegistry),
				txConfigOpts,
			)
			if err != nil {
				return err
			}
			initClientCtx = initClientCtx.WithTxConfig(txConfigWithTextual)

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := initAppConfig()
			customCMTConfig := initCometBFTConfig()

			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customCMTConfig)
		},
	}

	initRootCmd(rootCmd, encodingConfig, tempApp.BasicModuleManager)

	autoCliOpts := tempApp.AutoCliOpts()
	initClientCtx, _ = config.ReadFromClientConfig(initClientCtx)
	autoCliOpts.ClientCtx = initClientCtx

	// a workaround to wire the legacy proposals to the cli
	// autoCli uses AppModule, while the legacy proposals are registered on the AppModuleBasic
	govModule, ok := autoCliOpts.Modules["gov"].(gov.AppModule)
	if !ok {
		panic("gov module not found")
	}
	govBasicModule, ok := tempApp.BasicModuleManager["gov"].(gov.AppModuleBasic)
	if !ok {
		panic("gov module basic not found")
	}
	govModule.AppModuleBasic = govBasicModule
	autoCliOpts.Modules["gov"] = govModule

	if err := autoCliOpts.EnhanceRootCommand(rootCmd); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(cometbftcmd.RootCmd)
	return rootCmd
}

// initCometBFTConfig helps to override default CometBFT Config values.
// return cmtcfg.DefaultConfig if no custom configuration is required for the application.
func initCometBFTConfig() *cmtcfg.Config {
	cfg := cmtcfg.DefaultConfig()
	// Set consensus timeouts to support fast block time
	cfg.Consensus.TimeoutPropose = 1800 * time.Millisecond
	cfg.Consensus.TimeoutCommit = 500 * time.Millisecond
	return cfg
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	baseDenom, err := sdk.GetBaseDenom()
	if err != nil {
		panic(err)
	}

	customAppTemplate, customAppConfig := ethservercfg.AppConfig(baseDenom)
	return customAppTemplate, customAppConfig
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig appparams.EncodingConfig, basicManager module.BasicManager) {
	a := appCreator{encodingConfig}
	rootCmd.AddCommand(
		ethermintclient.ValidateChainID(genutilcli.InitCmd(basicManager, app.DefaultNodeHome)),
		genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome, app.GenTxMessageValidator, nil),
		GenTxCmd(basicManager, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome),
		genutilcli.ValidateGenesisCmd(basicManager),
		AddGenesisAccountCmd(app.DefaultNodeHome),
		GenRelayersCmd(app.DefaultNodeHome),
		cometbftcli.NewCompletionCmd(rootCmd, true),
		debug.Cmd(),
		AddGenesisStakePoolAccountCmd(app.DefaultNodeHome),
		AddGenesisModuleAccountsCmd(app.DefaultNodeHome),
		SetDAOCmd(),
		mecli.Debug(),
	)

	// add server commands
	ethermintserver.AddCommands(
		rootCmd,
		ethermintserver.NewDefaultStartOptions(a.newApp, app.DefaultNodeHome),
		a.appExport,
		addModuleInitFlags,
	)

	rootCmd.AddCommand(InspectCmd(a.appExport, a.newApp, app.DefaultNodeHome))

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		server.StatusCommand(),
		queryCommand(),
		txCommand(),
		ethermintclient.KeyCommands(app.DefaultNodeHome),
		cometbftcli.NewCompletionCmd(rootCmd, true),
	)
}

// queryCommand returns the sub-command to send queries to the app
func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		rpc.QueryEventForTxCmd(),
		server.QueryBlockCmd(),
		authcmd.QueryTxsByEventsCmd(),
		server.QueryBlocksCmd(),
		authcmd.QueryTxCmd(),
		server.QueryBlockResultsCmd(),
		rpc.ValidatorCommand(),
	)

	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

// txCommand returns the sub-command to send transactions to the app
func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		GetEncodeToRawTxCommand(),
		GetDecodeRawTxCommand(),
		authcmd.GetDecodeCommand(),
		authcmd.GetSimulateCmd(),
	)

	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")
	return cmd
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
	// this line is used by starport scaffolding # root/arguments
}

type appCreator struct {
	encodingConfig appparams.EncodingConfig
}

// newApp creates a new Cosmos SDK app
func (a appCreator) newApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) servertypes.Application {
	baseAppOptions := sdkserver.DefaultBaseappOptions(appOpts)

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(sdkserver.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	baseAppOptions = append(baseAppOptions, func(bapp *baseapp.BaseApp) {
		bapp.SetMempool(mempool.NoOpMempool{})
	})

	// NOTE we use custom transaction decoder that supports the sdk.Tx interface instead of sdk.StdTx
	// Setup Mempool and Proposal Handlers
	//baseAppOptions = append(baseAppOptions, func(bapp *baseapp.BaseApp) {
	//	maxTxs := cast.ToInt(appOpts.Get(sdkserver.FlagMempoolMaxTxs))
	//	if maxTxs <= 0 {
	//		maxTxs = 5000
	//	}
	//	priorityMempool := mempool.NewPriorityMempool(
	//		mempool.PriorityNonceWithMaxTx(maxTxs),
	//		mempool.PriorityNonceWithTxReplacement(func(op, np int64, oTx, nTx sdk.Tx) bool {
	//			threshold := int64(100 + 1)
	//			return np >= op*threshold/100
	//		}),
	//	)
	//	//handler := baseapp.NewDefaultProposalHandler(priorityMempool, bapp)
	//	bapp.SetMempool(priorityMempool)
	//	bapp.SetPrepareProposal(baseapp.NoOpPrepareProposal())
	//	bapp.SetProcessProposal(baseapp.NoOpProcessProposal())
	//})

	return app.New(
		logger,
		db,
		traceStore,
		true,
		appOpts,
		baseAppOptions...,
	)
}

// appExport creates a new simapp (optionally at a given height)
func (a appCreator) appExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home not set")
	}

	baseAppOptions := sdkserver.DefaultBaseappOptions(appOpts)

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(sdkserver.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	newApp := app.New(
		logger,
		db,
		traceStore,
		height == -1,
		appOpts,
		baseAppOptions...,
	)

	if height != -1 {
		if err := newApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	}

	return newApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}
