package v2

import (
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	evmconfig "github.com/evmos/ethermint/server/config"
	"github.com/spf13/cobra"
	"path/filepath"
)

// preUpgradeCmd called by cosmovisor
func PreUpgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pre-upgrade-v2_0_1",
		Short: "pre-upgrade, called by cosmovisor, before migrations upgrade",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			serverCtx.Logger.Info("pre-upgrade-v2 starting")
			rootDir := serverCtx.Config.RootDir

			config.SetConfigTemplate(config.DefaultConfigTemplate)
			oldAppConfig := config.DefaultConfig()
			if err := serverCtx.Viper.Unmarshal(oldAppConfig); err != nil {
				return err
			}

			customAppConfig := evmconfig.Config{
				Config:  *oldAppConfig,
				EVM:     *evmconfig.DefaultEVMConfig(),
				JSONRPC: *evmconfig.DefaultJSONRPCConfig(),
				TLS:     *evmconfig.DefaultTLSConfig(),
			}
			config.SetConfigTemplate(config.DefaultConfigTemplate + evmconfig.DefaultConfigTemplate)

			fileName := filepath.Join(rootDir, "config", "app.toml")
			config.WriteConfigFile(fileName, customAppConfig)
			serverCtx.Logger.Info("pre-upgrade-v3 success")
			return nil
		},
	}
	return cmd
}
