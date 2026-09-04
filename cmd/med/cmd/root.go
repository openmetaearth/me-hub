package cmd

import (
	"errors"
	"io"
	"os"

	"cosmossdk.io/log"
	cometbftcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"
	cometbftcfg "github.com/cometbft/cometbft/config"
	cometbftcli "github.com/cometbft/cometbft/libs/cli"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/mempool"
	"github.com/cosmos/cosmos-sdk/types/module"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	ethermintclient "github.com/evmos/ethermint/client"
	"github.com/evmos/ethermint/crypto/hd"
	ethermintserver "github.com/evmos/ethermint/server"
	servercfg "github.com/evmos/ethermint/server/config"
	ipfslog "github.com/ipfs/go-log/v2"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/openmetaearth/me-hub/app"
	appparams "github.com/openmetaearth/me-hub/app/params"
	mecli "github.com/openmetaearth/me-hub/client/cli"
	"github.com/openmetaearth/me-hub/logger"
)

// NewRootCmd creates a new root command for me hub
func NewRootCmd() (*cobra.Command, appparams.EncodingConfig) {
	// Seal bech32 prefixes / coin type before constructing codecs or the temp app.
	initSDKConfig()

	encodingConfig := app.MakeEncodingConfig()
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

	// SDK 0.50 AppModuleBasic for distribution/bank/authz holds an unexported codec.
	// Package-level ModuleBasics uses zero-value structs, so GetTxCmd panics.
	// Construct a throwaway app to get a BasicManager (and AutoCLI) from real modules.
	tempDir, err := os.MkdirTemp("", "med-cli-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDir)

	tempApp := app.New(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		false,
		map[int64]bool{},
		tempDir,
		0,
		encodingConfig,
		sims.AppOptionsMap{
			flags.FlagHome:   tempDir,
			"skip-wasm-init": true,
		},
	)

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

			customAppTemplate, customAppConfig := initAppConfig()
			customTMConfig := initTendermintConfig()
			err = sdkserver.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customTMConfig)
			if err != nil {
				return err
			}
			enableMeLogger, _ := cmd.Flags().GetBool("enable_me_hub_logger")
			if os.Getenv("ENABLE_MEHUB_LOGGER") != "" || enableMeLogger {
				ctx := sdkserver.GetServerContextFromCmd(cmd)
				ctx.Logger = logger.NewLogger("me-hub").WithEnvLevelOr("info").WithStacktrace(ipfslog.LevelError)
			}
			return nil
		},
	}
	rootCmd.PersistentFlags().Bool("enable_me_hub_logger", false, "use me-hub logger instead of cosmos lib logger")
	initRootCmd(rootCmd, encodingConfig, tempApp.BasicModuleManager)

	autoCliOpts := tempApp.AutoCliOpts()
	autoCliOpts.ClientCtx = initClientCtx
	if err := autoCliOpts.EnhanceRootCommand(rootCmd); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(cometbftcmd.RootCmd)
	return rootCmd, encodingConfig
}

// initTendermintConfig helps to override default Tendermint Config values.
// return tmcfg.DefaultConfig if no custom configuration is required for the application.
func initTendermintConfig() *cometbftcfg.Config {
	cfg := cometbftcfg.DefaultConfig()

	// these values put a higher strain on node memory
	// cfg.P2P.MaxNumInboundPeers = 100
	// cfg.P2P.MaxNumOutboundPeers = 40

	return cfg
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	baseDenom, err := sdk.GetBaseDenom()
	if err != nil {
		panic(err)
	}

	customAppTemplate, customAppConfig := servercfg.AppConfig(baseDenom)
	return customAppTemplate, customAppConfig
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig appparams.EncodingConfig, basicManager module.BasicManager) {
	a := appCreator{encodingConfig}
	rootCmd.AddCommand(
		ethermintclient.ValidateChainID(
			genutilcli.InitCmd(basicManager, app.DefaultNodeHome),
		),
		genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome, app.GenTxMessageValidator, nil),
		GenTxCmd(
			basicManager,
			encodingConfig.TxConfig,
			banktypes.GenesisBalancesIterator{},
			app.DefaultNodeHome,
		),
		genutilcli.ValidateGenesisCmd(basicManager),
		AddGenesisAccountCmd(app.DefaultNodeHome),
		GenRelayersCmd(app.DefaultNodeHome),
		cometbftcli.NewCompletionCmd(rootCmd, true),
		debug.Cmd(),
		AddGenesisStakePoolAccountCmd(app.DefaultNodeHome),
		AddGenesisModuleAccountsCmd(app.DefaultNodeHome),
		SetDAOCmd(),
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
		sdkserver.StatusCommand(),
		queryCommand(basicManager),
		txCommand(basicManager),
		ethermintclient.KeyCommands(app.DefaultNodeHome),
	)
	rootCmd.AddCommand(mecli.Debug())
}

// queryCommand returns the sub-command to send queries to the app
func queryCommand(basicManager module.BasicManager) *cobra.Command {
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
		sdkserver.QueryBlockCmd(),
		rpc.ValidatorCommand(),
		sdkserver.QueryBlocksCmd(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
		sdkserver.QueryBlockResultsCmd(),
	)

	basicManager.AddQueryCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

// txCommand returns the sub-command to send transactions to the app
func txCommand(basicManager module.BasicManager) *cobra.Command {
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

	basicManager.AddTxCommands(cmd)
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
	// baseAppOptions = append(baseAppOptions, func(bapp *baseapp.BaseApp) {
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
	// })

	return app.New(
		logger,
		db,
		traceStore,
		true,
		skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		cast.ToUint(appOpts.Get(sdkserver.FlagInvCheckPeriod)),
		a.encodingConfig,
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
		skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		cast.ToUint(appOpts.Get(sdkserver.FlagInvCheckPeriod)),
		a.encodingConfig,
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
