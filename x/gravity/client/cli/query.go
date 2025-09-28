package cli

import (
	"fmt"
	abcitype "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/st-chain/me-hub/x/gravity/types"
	"strconv"
)

func GetQueryCmd(moduleName string, subNames ...string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        moduleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", moduleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	for _, chainName := range subNames {
		cmd.AddCommand(GetQueryCmd(chainName))
	}
	if len(subNames) == 0 {
		cmd.AddCommand(getQuerySubCmds(moduleName)...)
	}
	return cmd
}

func getQuerySubCmds(chainName string) []*cobra.Command {
	cmds := []*cobra.Command{
		// query module params
		CmdGetParams(chainName),

		// query Relayer
		CmdRelayer(chainName),
		CmdGetRelayers(chainName),
		CmdGetProposalRelayers(chainName),

		// query relayer set
		CmdGetCurrentRelayerSet(chainName),
		CmdGetRelayerSetRequest(chainName),

		// need relayer consensus sign
		// relayer set change confirm
		CmdGetLastRelayerSetRequests(chainName),
		CmdGetPendingRelayerSetRequest(chainName),
		CmdGetRelayerSetConfirm(chainName),
		CmdGetRelayerSetConfirms(chainName),
		// request batch confirm
		CmdGetPendingOutgoingTXBatchRequest(chainName),
		CmdBatchConfirm(chainName),
		CmdBatchConfirms(chainName),

		// send to external
		CmdBatchRequestByNonce(chainName),
		CmdGetPendingSendToExternal(chainName),
		CmdOutgoingTxBatches(chainName),

		CmdGetLastObservedBlockHeight(chainName),
		CmdProjectedBatchTimeoutHeight(chainName),

		// denom <-> external token
		CmdGetBridgeTokens(chainName),
		CmdGetBridgeCoinByDenom(chainName),

		// event nonce
		CmdGetRelayerEventNonce(chainName),
		CmdGetRelayerEventBlockHeight(chainName),
		CmdGetLastObservedEventNonce(chainName),
	}

	for _, command := range cmds {
		flags.AddQueryFlagsToCmd(command)
	}
	return cmds
}

func CmdGetParams(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current parameters information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(&res.Params)
		},
	}
	return cmd
}

func CmdGetRelayers(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relayers",
		Short: "Query current relayers",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Relayers(cmd.Context(), &types.QueryRelayersRequest{
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetProposalRelayers(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-relayers",
		Short: "Query proposal relayers address",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ProposalRelayers(cmd.Context(), &types.QueryProposalRelayersRequest{
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdRelayer(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relayer [relayer-address|external-address]",
		Short: "Query relayer for a given relayer address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			input := args[0]
			req := types.QueryRelayerRequest{ChainName: chainName}
			var relayerValid, externalValid bool
			if _, err := sdk.AccAddressFromBech32(input); err == nil {
				req.RelayerAddress = input
				relayerValid = true
			}
			if err := types.ValidateExternalAddr(chainName, input); err == nil {
				req.ExternalAddress = input
				externalValid = true
			}
			if !relayerValid && !externalValid {
				return fmt.Errorf("invalid input: %s, must be a valid relayer address or external address", input)
			}
			res, err := queryClient.Relayer(cmd.Context(), &req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res.Relayer)
		},
	}
	return cmd
}

func CmdGetCurrentRelayerSet(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current-relayer-set",
		Short: "Query current relayer-set",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CurrentRelayerSet(cmd.Context(), &types.QueryCurrentRelayerSetRequest{
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res.RelayerSet)
		},
	}
	return cmd
}

func CmdGetRelayerSetRequest(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relayer-set-request [nonce]",
		Short: "Query requested relayer-set with a particular nonce",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			var nonce uint64
			if len(args) == 0 {
				queryAbciResp, err := clientCtx.QueryABCI(abcitype.RequestQuery{
					Path: fmt.Sprintf("store/%s/key", chainName),
					Data: types.LatestRelayerSetNonce,
				})
				if err != nil {
					return err
				}
				nonce = sdk.BigEndianToUint64(queryAbciResp.Value)
			} else {
				var err error
				nonce, err = strconv.ParseUint(args[0], 10, 64)
				if err != nil {
					return err
				}
			}
			res, err := queryClient.RelayerSetRequest(cmd.Context(), &types.QueryRelayerSetRequestRequest{
				ChainName: chainName,
				Nonce:     nonce,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res.RelayerSet)
		},
	}
	return cmd
}

func CmdGetLastRelayerSetRequests(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last-relayer-set-requests",
		Short: "Query last relayer set requests",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.LastRelayerSetRequests(cmd.Context(), &types.QueryLastRelayerSetRequestsRequest{
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetPendingRelayerSetRequest(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-relayer-set-request [bridger]",
		Short: "Query the latest relayer-set request which has not been signed by a particular relayer bridger",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			bridgerAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			res, err := queryClient.LastPendingRelayerSetRequestByAddr(cmd.Context(), &types.QueryLastPendingRelayerSetRequestByAddrRequest{
				RelayerAddress: bridgerAddr.String(),
				ChainName:      chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetRelayerSetConfirm(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relayer-set-confirm [nonce] [bridger-address]",
		Short: "Query relayer-set confirmation with a particular nonce from a particular relayer bridger",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			nonce, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			bridgerAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			res, err := queryClient.RelayerSetConfirm(cmd.Context(), &types.QueryRelayerSetConfirmRequest{
				Nonce:          nonce,
				RelayerAddress: bridgerAddr.String(),
				ChainName:      chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res.Confirm)
		},
	}
	return cmd
}

func CmdGetRelayerSetConfirms(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relayer-set-confirms [nonce]",
		Short: "Query relayer-set confirmations with a particular nonce",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			nonce, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			res, err := queryClient.RelayerSetConfirmsByNonce(cmd.Context(), &types.QueryRelayerSetConfirmsByNonceRequest{
				Nonce:     nonce,
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetPendingOutgoingTXBatchRequest(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-batch-request [bridger-address]",
		Short: "Query the latest outgoing TX batch request which has not been signed by a particular relayer bridger address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			bridgerAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			res, err := queryClient.LastPendingBatchRequestByAddr(cmd.Context(), &types.QueryLastPendingBatchRequestByAddrRequest{
				RelayerAddress: bridgerAddr.String(),
				ChainName:      chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res.Batch)
		},
	}
	return cmd
}

func CmdBatchConfirm(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-confirm [token-contract] [nonce] [bridger-address]",
		Short: "Query outgoing tx batches confirm by relayer bridger address",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			tokenContract := args[0]
			if err := types.ValidateExternalAddr(chainName, tokenContract); err != nil {
				return err
			}
			nonce, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			bridgerAddr, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}
			res, err := queryClient.BatchConfirm(cmd.Context(), &types.QueryBatchConfirmRequest{
				ChainName:      chainName,
				TokenContract:  tokenContract,
				Nonce:          nonce,
				RelayerAddress: bridgerAddr.String(),
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res.Confirm)
		},
	}
	return cmd
}

func CmdBatchConfirms(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-confirms [token-contract] [nonce]",
		Short: "Query outgoing tx batches confirms",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			tokenContract := args[0]
			if err := types.ValidateExternalAddr(chainName, tokenContract); err != nil {
				return err
			}
			nonce, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			res, err := queryClient.BatchConfirms(cmd.Context(), &types.QueryBatchConfirmsRequest{
				TokenContract: tokenContract,
				Nonce:         uint64(nonce),
				ChainName:     chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdBatchRequestByNonce(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-request [token-contract] [nonce]",
		Short: "Query outgoing tx batches",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			tokenContract := args[0]
			if err := types.ValidateExternalAddr(chainName, tokenContract); err != nil {
				return err
			}
			nonce, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			res, err := queryClient.BatchRequestByNonce(cmd.Context(), &types.QueryBatchRequestByNonceRequest{
				ChainName:     chainName,
				TokenContract: tokenContract,
				Nonce:         nonce,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res.Batch)
		},
	}
	return cmd
}

func CmdGetPendingSendToExternal(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-send-to-external [address]",
		Short: "Query pending send to external txs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			res, err := queryClient.GetPendingSendToExternal(cmd.Context(), &types.QueryPendingSendToExternalRequest{
				ChainName:     chainName,
				SenderAddress: addr.String(),
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdOutgoingTxBatches(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "outgoing-tx-batches",
		Short: "Query outgoing tx batches",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.OutgoingTxBatches(cmd.Context(), &types.QueryOutgoingTxBatchesRequest{
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetLastObservedBlockHeight(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last-observed-block-height",
		Short: "Query last observed block height",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.LastObservedBlockHeight(cmd.Context(), &types.QueryLastObservedBlockHeightRequest{
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdProjectedBatchTimeoutHeight(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "projected-batch-timeout-height",
		Short: "Query projected batch timeout height",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.ProjectedBatchTimeoutHeight(cmd.Context(), &types.QueryProjectedBatchTimeoutHeightRequest{
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetBridgeTokens(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bridge-tokens",
		Short: "Query bridge token list",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.BridgeTokens(cmd.Context(), &types.QueryBridgeTokensRequest{
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetRelayerEventNonce(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "event-nonce [bridger-address]",
		Short: "Query last event nonce by bridger address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			bridgerAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			res, err := queryClient.LastEventNonceByAddr(cmd.Context(), &types.QueryLastEventNonceByAddrRequest{
				ChainName:      chainName,
				RelayerAddress: bridgerAddr.String(),
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetLastObservedEventNonce(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last-observed-nonce",
		Short: "Query last observed event nonce",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryAbciResp, err := clientCtx.QueryABCI(abcitype.RequestQuery{
				Path: fmt.Sprintf("store/%s/key", chainName),
				Data: types.LastObservedEventNonceKey,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintString(fmt.Sprintf("%d\n", sdk.BigEndianToUint64(queryAbciResp.Value)))
		},
	}
	return cmd
}

func CmdGetRelayerEventBlockHeight(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "event-block-height [bridger-address]",
		Short: "Query last event block height by bridger address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			bridgerAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			res, err := queryClient.LastEventBlockHeightByAddr(cmd.Context(), &types.QueryLastEventBlockHeightByAddrRequest{
				RelayerAddress: bridgerAddr.String(),
				ChainName:      chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetBridgeCoinByDenom(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bridge-token [denom]",
		Short: "Query bridge coin from contract address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			denom := args[0]
			res, err := queryClient.BridgeCoinByDenom(cmd.Context(), &types.QueryBridgeCoinByDenomRequest{
				ChainName: chainName,
				Denom:     denom,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}
