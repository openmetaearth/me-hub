package cli

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/openmetaearth/me-hub/x/gravity/types"
	"github.com/spf13/cobra"
)

func GetTxCmd(moduleName string, subNames ...string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        moduleName,
		Short:                      fmt.Sprintf("%s%s transaction subcommands", strings.ToUpper(moduleName[:1]), moduleName[1:]),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	for _, chainName := range subNames {
		cmd.AddCommand(GetTxCmd(chainName))
	}
	if len(subNames) == 0 {
		cmd.AddCommand(getTxSubCmds(moduleName)...)
	}
	return cmd
}

func getTxSubCmds(chainName string) []*cobra.Command {
	cmds := []*cobra.Command{
		CmdBondedRelayer(chainName),
		CmdUnbondedRelayer(chainName),
		CmdAddDelegate(chainName),
		CmdProposalRelayers(chainName),

		// send to external chain
		CmdSendToExternal(chainName),
		CmdCancelSendToExternal(chainName),
		CmdIncreaseBridgeFee(chainName),
		CmdRequestBatch(chainName),

		// relayer consensus confirm
		CmdRelayerSetConfirm(chainName),
		CmdRequestBatchConfirm(chainName),
	}
	for _, command := range cmds {
		flags.AddTxFlagsToCmd(command)
	}
	return cmds
}

func CmdBondedRelayer(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bonded-relayer [external-address] [delegate-amount]",
		Short: "Allows relayer to delegate their voting responsibilities to a given key.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}
			msg := types.MsgBondedRelayer{
				ChainName:       chainName,
				RelayerAddress:  cliCtx.GetFromAddress().String(),
				ExternalAddress: args[0],
				DelegateAmount:  amount,
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}
	return cmd
}

func CmdUnbondedRelayer(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbonded-relayer",
		Short: "Quit the relayer",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.MsgUnbondedRelayer{
				RelayerAddress: cliCtx.GetFromAddress().String(),
				ChainName:      chainName,
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}
	return cmd
}

func CmdAddDelegate(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-delegate [delegate-amount]",
		Short: "Allows relayer add delegate.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}
			msg := types.MsgAddDelegate{
				RelayerAddress: cliCtx.GetFromAddress().String(),
				Amount:         amount,
				ChainName:      chainName,
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}
	return cmd
}

func CmdProposalRelayers(chainName string) *cobra.Command {
	var relayers []string
	cmd := &cobra.Command{
		Use:   "proposal-relayers --relayers addr1,addr2[,addrN]",
		Short: "Propose a new relayer set",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			if len(relayers) == 0 {
				return fmt.Errorf("at least one relayer is required")
			}
			// Clean the relayer list: trim spaces, remove duplicates and empty entries
			clean := make([]string, 0, len(relayers))
			seen := make(map[string]struct{})
			for _, r := range relayers {
				r = strings.TrimSpace(r)
				if r == "" {
					continue
				}
				if _, ok := seen[r]; ok {
					continue
				}
				seen[r] = struct{}{}
				clean = append(clean, r)
			}
			if len(clean) == 0 {
				return fmt.Errorf("at least one relayer is required")
			}

			msg := &types.MsgProposalRelayers{
				ChainName: chainName,
				Authority: cliCtx.GetFromAddress().String(),
				Relayers:  clean,
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().StringSliceVar(&relayers, "relayers", nil, "relayer addresses")
	return cmd
}

func CmdSendToExternal(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-to-external [external-dest] [amount] [bridge-fee]",
		Short: "Adds a new entry to the transaction pool to withdraw an amount from the bridge contract",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return errorsmod.Wrap(err, "amount")
			}

			bridgeFee, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return errorsmod.Wrap(err, "bridge fee")
			}

			msg := types.MsgSendToExternal{
				Sender:    cliCtx.GetFromAddress().String(),
				Dest:      args[0],
				Amount:    amount,
				BridgeFee: bridgeFee,
				ChainName: chainName,
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}
	return cmd
}

func CmdCancelSendToExternal(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-send-to-external [tx-ID]",
		Short: "Cancel transaction send to external",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			msg := &types.MsgCancelSendToExternal{
				TransactionId: txId,
				Sender:        cliCtx.GetFromAddress().String(),
				ChainName:     chainName,
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}
	return cmd
}

func CmdIncreaseBridgeFee(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "increase-bridge-fee [tx-ID] [add-bridge-fee]",
		Short: "Increase bridge fee for send to external transaction",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			addBridgeFee, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return errorsmod.Wrap(err, "add bridge fee")
			}

			msg := &types.MsgIncreaseBridgeFee{
				ChainName:     chainName,
				TransactionId: txId,
				Sender:        cliCtx.GetFromAddress().String(),
				AddBridgeFee:  addBridgeFee,
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}
	return cmd
}

func CmdRequestBatch(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build-batch [token-denom] [minimum-fee] [external-fee-receive] [base-fee]",
		Short: "Build a new batch on the fx side for pooled withdrawal transactions",
		Args:  cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			minimumFee, ok := sdkmath.NewIntFromString(args[1])
			if !ok || minimumFee.IsNegative() {
				return fmt.Errorf("minimum fee is invalid, %v", args[1])
			}
			baseFee := sdkmath.ZeroInt()
			if len(args) == 4 {
				baseFee, ok = sdkmath.NewIntFromString(args[3])
				if !ok {
					return fmt.Errorf("invalid base fee: %v", args[3])
				}
			}

			msg := &types.MsgRequestBatch{
				Sender:     clientCtx.GetFromAddress().String(),
				Denom:      args[0],
				MinimumFee: minimumFee,
				FeeReceive: args[2],
				ChainName:  chainName,
				BaseFee:    baseFee,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	return cmd
}

func CmdRequestBatchConfirm(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request-batch-confirm [contract-address] [nonce] [private-key]",
		Short: "Send batch confirm msg",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fromAddress := clientCtx.GetFromAddress()

			tokenContract := args[0]
			nonce, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			privateKey, err := recoveryPrivateKeyByKeystore(args[2])
			if err != nil {
				return err
			}
			externalAddress := ethcrypto.PubkeyToAddress(privateKey.PublicKey)

			queryClient := types.NewQueryClient(clientCtx)
			batchRequestByNonceResp, err := queryClient.BatchRequestByNonce(cmd.Context(), &types.QueryBatchRequestByNonceRequest{
				Nonce:         nonce,
				TokenContract: tokenContract,
				ChainName:     chainName,
			})
			if err != nil {
				return err
			}
			if batchRequestByNonceResp.Batch == nil {
				return fmt.Errorf("not found batch request by nonce, tokenContract: %v, nonce: %v", tokenContract, nonce)
			}
			// Determine whether it has been confirmed
			batchConfirmResp, err := queryClient.BatchConfirm(cmd.Context(), &types.QueryBatchConfirmRequest{
				Nonce:          nonce,
				TokenContract:  tokenContract,
				RelayerAddress: fromAddress.String(),
				ChainName:      chainName,
			})
			if err != nil {
				return err
			}
			if batchConfirmResp.GetConfirm() != nil {
				confirm := batchConfirmResp.GetConfirm()
				return clientCtx.PrintProto(confirm)
			}
			paramsResp, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			checkpoint, err := batchRequestByNonceResp.GetBatch().GetCheckpoint(paramsResp.Params.GetGravityId())
			if err != nil {
				return err
			}
			signature, err := types.NewEthereumSignature(checkpoint, privateKey)
			if err != nil {
				return err
			}
			msg := &types.MsgConfirmBatch{
				Nonce:           nonce,
				TokenContract:   tokenContract,
				ExternalAddress: externalAddress.String(),
				RelayerAddress:  fromAddress.String(),
				Signature:       hex.EncodeToString(signature),
				ChainName:       chainName,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	return cmd
}

func CmdRelayerSetConfirm(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relayer-set-confirm [nonce] [private-key]",
		Short: "Send relayer-set confirm msg",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fromAddress := clientCtx.GetFromAddress()

			nonce, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			privateKey, err := recoveryPrivateKeyByKeystore(args[1])
			if err != nil {
				return err
			}
			externalAddress := ethcrypto.PubkeyToAddress(privateKey.PublicKey)

			queryClient := types.NewQueryClient(clientCtx)
			relayerSetRequestResp, err := queryClient.RelayerSetRequest(cmd.Context(), &types.QueryRelayerSetRequestRequest{
				Nonce: nonce, ChainName: chainName,
			})
			if err != nil {
				return err
			}
			// Determine whether it has been confirmed
			relayerSetConfirmResp, err := queryClient.RelayerSetConfirm(cmd.Context(), &types.QueryRelayerSetConfirmRequest{
				Nonce:          nonce,
				RelayerAddress: fromAddress.String(),
				ChainName:      chainName,
			})
			if err != nil {
				return err
			}
			if relayerSetConfirmResp.GetConfirm() != nil {
				confirm := relayerSetConfirmResp.GetConfirm()
				return clientCtx.PrintProto(confirm)
			}
			paramsResp, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			checkpoint, err := relayerSetRequestResp.GetRelayerSet().GetCheckpoint(paramsResp.Params.GetGravityId())
			if err != nil {
				return err
			}
			signature, err := types.NewEthereumSignature(checkpoint, privateKey)
			if err != nil {
				return err
			}
			msg := &types.MsgRelayerSetConfirm{
				Nonce:           nonce,
				RelayerAddress:  fromAddress.String(),
				ExternalAddress: externalAddress.String(),
				Signature:       hex.EncodeToString(signature),
				ChainName:       chainName,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	return cmd
}

func recoveryPrivateKeyByKeystore(privateKey string) (*ecdsa.PrivateKey, error) {
	var ethPrivateKey *ecdsa.PrivateKey
	if _, err := os.Stat(privateKey); err == nil {
		file, err := os.ReadFile(privateKey)
		if err != nil {
			return nil, err
		}
		stdinReader, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return nil, err
		}
		password := strings.TrimSpace(stdinReader)
		key, err := keystore.DecryptKey(file, password)
		if err != nil {
			return nil, err
		}
		ethPrivateKey = key.PrivateKey
	} else {
		key, err := ethcrypto.HexToECDSA(privateKey)
		if err != nil {
			return nil, fmt.Errorf("invalid eth private key: %s", err.Error())
		}
		ethPrivateKey = key
	}
	return ethPrivateKey, nil
}
