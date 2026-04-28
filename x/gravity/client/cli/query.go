package cli

import (
	"encoding/json"
	"fmt"
	abcitype "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/gravity/types"
	"github.com/spf13/cobra"
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
		CmdGetPendingBatchRequest(chainName),
		CmdBatchConfirm(chainName),
		CmdBatchConfirms(chainName),

		// send to external
		CmdBatchRequestByNonce(chainName),
		CmdPendingOutgoingTxByAddr(chainName),
		CmdUnbatchedTxs(chainName),
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
		CmdClaims(chainName),
		CmdBridgeChainList(),
		CmdMeNonce(chainName),
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
		Short: "Query relayer for a given relayer or external address",
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
				if len(queryAbciResp.Value) == 0 {
					return fmt.Errorf("latest relayer-set nonce not found; please provide the nonce explicitly")
				}
				if len(queryAbciResp.Value) != 8 {
					return fmt.Errorf("unexpected relayer-set nonce encoding (got %d bytes)", len(queryAbciResp.Value))
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
		Use:   "relayer-set-requests",
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
		Use:   "pending-relayer-set-request [relayer]",
		Short: "Query the latest relayer-set request which has not been signed by a particular relayer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			relayerAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			res, err := queryClient.LastPendingRelayerSetRequestByAddr(cmd.Context(), &types.QueryLastPendingRelayerSetRequestByAddrRequest{
				RelayerAddress: relayerAddr.String(),
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
		Use:   "relayer-set-confirm [nonce] [relayer-address]",
		Short: "Query relayer-set confirmation with a particular nonce from a particular relayer",
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
			relayerAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			res, err := queryClient.RelayerSetConfirm(cmd.Context(), &types.QueryRelayerSetConfirmRequest{
				Nonce:          nonce,
				RelayerAddress: relayerAddr.String(),
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

func CmdGetPendingBatchRequest(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-batch-request [relayer-address]",
		Short: "Query the latest outgoing TX batch request which has not been signed by a particular relayer address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			relayerAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			res, err := queryClient.LastPendingBatchRequestByAddr(cmd.Context(), &types.QueryLastPendingBatchRequestByAddrRequest{
				RelayerAddress: relayerAddr.String(),
				ChainName:      chainName,
			})
			if err != nil {
				return err
			}
			if res.Batch == nil {
				return fmt.Errorf("no pending batch request found for relayer %s", relayerAddr.String())
			}
			return clientCtx.PrintProto(res.Batch)
		},
	}
	return cmd
}

func CmdBatchConfirm(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-confirm [token-contract] [nonce] [relayer-address]",
		Short: "Query outgoing tx batches confirm by relayer address",
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
			relayerAddr, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}
			res, err := queryClient.BatchConfirm(cmd.Context(), &types.QueryBatchConfirmRequest{
				ChainName:      chainName,
				TokenContract:  tokenContract,
				Nonce:          nonce,
				RelayerAddress: relayerAddr.String(),
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
			nonce, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			res, err := queryClient.BatchConfirms(cmd.Context(), &types.QueryBatchConfirmsRequest{
				TokenContract: tokenContract,
				Nonce:         nonce,
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

func CmdPendingOutgoingTxByAddr(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-outgoing-tx-by-addr [address]",
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
			res, err := queryClient.PendingOutgoingTxByAddr(cmd.Context(), &types.QueryPendingOutgoingTxByAddrRequest{
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

func CmdUnbatchedTxs(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbatched-txs [token-contract]",
		Short: "Query unbatched send to external txs",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			var tokenContract string
			if len(args) > 0 {
				tokenContract = args[0]
				if err := types.ValidateExternalAddr(chainName, tokenContract); err != nil {
					return err
				}
			}
			pageReq, _ := client.ReadPageRequest(cmd.Flags())
			res, err := queryClient.UnbatchedTxs(cmd.Context(), &types.QueryUnbatchedTxsRequest{
				ChainName:     chainName,
				TokenContract: tokenContract,
				Pagination:    pageReq,
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
		Use:   "event-nonce [relayer-address]",
		Short: "Query last event nonce by relayer address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			relayerAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			res, err := queryClient.LastEventNonceByAddr(cmd.Context(), &types.QueryLastEventNonceByAddrRequest{
				ChainName:      chainName,
				RelayerAddress: relayerAddr.String(),
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
			if len(queryAbciResp.Value) == 0 {
				return clientCtx.PrintString("0\n")
			}
			if len(queryAbciResp.Value) != 8 {
				return fmt.Errorf("unexpected event nonce encoding (got %d bytes)", len(queryAbciResp.Value))
			}
			return clientCtx.PrintString(fmt.Sprintf("%d\n", sdk.BigEndianToUint64(queryAbciResp.Value)))
		},
	}
	return cmd
}

func CmdGetRelayerEventBlockHeight(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "event-block-height [relayer-address]",
		Short: "Query last event block height by relayer address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			relayerAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			res, err := queryClient.LastEventBlockHeightByAddr(cmd.Context(), &types.QueryLastEventBlockHeightByAddrRequest{
				RelayerAddress: relayerAddr.String(),
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
		Use:   "bridge-token [denom] [contract-address]",
		Short: "Query bridge coin from contract address",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			denom := args[0]
			var contract string
			if len(args) > 1 {
				contract = args[1]
			}
			res, err := queryClient.BridgeToken(cmd.Context(), &types.QueryBridgeTokenRequest{
				ChainName:       chainName,
				Denom:           denom,
				ContractAddress: contract,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdClaims(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claims [event-nonce]",
		Short: "Query claims by event nonce",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			eventNonce, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			res, err := queryClient.ClaimsByEventNonce(cmd.Context(), &types.QueryClaimsByEventNonceRequest{
				ChainName:  chainName,
				EventNonce: eventNonce,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdBridgeChainList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chain-list",
		Short: "Query bridge chain list",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.BridgeChainList(cmd.Context(), &types.QueryBridgeChainListRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdMeNonce(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "me-nonce",
		Short: "Query ME nonce",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			relayerSetNonce, err := clientCtx.QueryABCI(abcitype.RequestQuery{
				Path: fmt.Sprintf("store/%s/key", chainName),
				Data: types.LatestRelayerSetNonce,
			})
			if err != nil {
				return err
			}

			batchId, err := clientCtx.QueryABCI(abcitype.RequestQuery{
				Path: fmt.Sprintf("store/%s/key", chainName),
				Data: types.KeyLastOutgoingBatchID,
			})
			if err != nil {
				return err
			}

			txId, err := clientCtx.QueryABCI(abcitype.RequestQuery{
				Path: fmt.Sprintf("store/%s/key", chainName),
				Data: types.KeyLastTxPoolID,
			})
			if err != nil {
				return err
			}
			lastSlashedBatchBlock, err := clientCtx.QueryABCI(abcitype.RequestQuery{
				Path: fmt.Sprintf("store/%s/key", chainName),
				Data: types.LastSlashedBatchBlock,
			})
			if err != nil {
				return err
			}
			lastSlashedRelayerSetNonce, err := clientCtx.QueryABCI(abcitype.RequestQuery{
				Path: fmt.Sprintf("store/%s/key", chainName),
				Data: types.LastSlashedRelayerSetNonce,
			})
			if err != nil {
				return err
			}

			res, err := json.Marshal(struct {
				RelayerSetNonce            uint64 `json:"relayer_set_nonce"`
				BatchId                    uint64 `json:"batch_id"`
				TxId                       uint64 `json:"tx_id"`
				LastSlashedBatchBlock      uint64 `json:"last_slashed_batch_block"`
				LastSlashedRelayerSetNonce uint64 `json:"last_slashed_relayer_set_nonce"`
			}{
				RelayerSetNonce:            sdk.BigEndianToUint64(relayerSetNonce.Value),
				BatchId:                    sdk.BigEndianToUint64(batchId.Value),
				TxId:                       sdk.BigEndianToUint64(txId.Value),
				LastSlashedBatchBlock:      sdk.BigEndianToUint64(lastSlashedBatchBlock.Value),
				LastSlashedRelayerSetNonce: sdk.BigEndianToUint64(lastSlashedRelayerSetNonce.Value),
			})
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}
	return cmd
}
