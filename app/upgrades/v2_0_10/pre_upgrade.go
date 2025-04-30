package v2_0_10

import (
	"bufio"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

// preUpgradeCmd called by cosmovisor
func PreUpgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pre-upgrade-v2",
		Short: "pre-upgrade, called by cosmovisor, before migrations upgrade",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			serverCtx.Logger.Info("pre-upgrade-v2 update config starting")
			rootDir := serverCtx.Config.RootDir

			fileName := filepath.Join(rootDir, "config", "app.toml")
			evmConfig := "\n###############################################################################\n###                             EVM Configuration                           ###\n###############################################################################\n\n[evm]\n\n# Tracer defines the 'vm.Tracer' type that the EVM will use when the node is run in\n# debug mode. To enable tracing use the '--evm.tracer' flag when starting your node.\n# Valid types are: json|struct|access_list|markdown\ntracer = \"\"\n\n# MaxTxGasWanted defines the gas wanted for each eth tx returned in ante handler in check tx mode.\nmax-tx-gas-wanted = 0\n\n###############################################################################\n###                           JSON RPC Configuration                        ###\n###############################################################################\n\n[json-rpc]\n\n# Enable defines if the gRPC server should be enabled.\nenable = true\n\n# Address defines the EVM RPC HTTP server address to bind to.\naddress = \"0.0.0.0:8545\"\n\n# Address defines the EVM WebSocket server address to bind to.\nws-address = \"0.0.0.0:8546\"\n\n# API defines a list of JSON-RPC namespaces that should be enabled\n# Example: \"eth,txpool,personal,net,debug,web3\"\napi = \"eth,net,web3\"\n\n# GasCap sets a cap on gas that can be used in eth_call/estimateGas (0=infinite). Default: 25,000,000.\ngas-cap = 25000000\n\n# EVMTimeout is the global timeout for eth_call. Default: 5s.\nevm-timeout = \"5s\"\n\n# TxFeeCap is the global tx-fee cap for send transaction. Default: 1eth.\ntxfee-cap = 1\n\n# FilterCap sets the global cap for total number of filters that can be created\nfilter-cap = 200\n\n# FeeHistoryCap sets the global cap for total number of blocks that can be fetched\nfeehistory-cap = 100\n\n# LogsCap defines the max number of results can be returned from single 'eth_getLogs' query.\nlogs-cap = 10000\n\n# BlockRangeCap defines the max block range allowed for 'eth_getLogs' query.\nblock-range-cap = 10000\n\n# HTTPTimeout is the read/write timeout of http json-rpc server.\nhttp-timeout = \"30s\"\n\n# HTTPIdleTimeout is the idle timeout of http json-rpc server.\nhttp-idle-timeout = \"2m0s\"\n\n# AllowUnprotectedTxs restricts unprotected (non EIP155 signed) transactions to be submitted via\n# the node's RPC when the global parameter is disabled.\nallow-unprotected-txs = false\n\n# MaxOpenConnections sets the maximum number of simultaneous connections\n# for the server listener.\nmax-open-connections = 0\n\n# EnableIndexer enables the custom transaction indexer for the EVM (ethereum transactions).\nenable-indexer = false\n\n# MetricsAddress defines the EVM Metrics server address to bind to. Pass --metrics in CLI to enable\n# Prometheus metrics path: /debug/metrics/prometheus\nmetrics-address = \"127.0.0.1:6065\"\n\n# Upgrade height for fix of revert gas refund logic when transaction reverted.\nfix-revert-gas-refund-height = 0\n\n###############################################################################\n###                             TLS Configuration                           ###\n###############################################################################\n\n[tls]\n\n# Certificate path defines the cert.pem file path for the TLS configuration.\ncertificate-path = \"\"\n\n# Key path defines the key.pem file path for the TLS configuration.\nkey-path = \"\""
			err := appendEvmConfigIfMissing(fileName, evmConfig)
			if err != nil {
				return err
			}
			serverCtx.Logger.Info("pre-upgrade config app.toml success")
			return nil
		},
	}
	return cmd
}

func appendEvmConfigIfMissing(fileName, evmConfig string) error {
	// Open the file for reading
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Check if the file contains the [evm] section
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == "[evm]" {
			// [evm] section already exists, no need to append
			return nil
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// Append the evmConfig if [evm] section is missing
	return appendEvmConfig(fileName, evmConfig)
}

func appendEvmConfig(fileName, evmConfig string) error {
	// Open the file in append mode, create it if it doesn't exist
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the evmConfig string to the file
	_, err = file.WriteString(evmConfig)
	return err
}
