package cli

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/openmetaearth/me-hub/x/did/types"
)

const (
	FlagFilters = "filters"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		CmdCreateDid(),
		CmdUpdateDidStatus(),
		// CmdRemoveDid(),
		CmdCreateService(),
		CmdUpdateServiceStatus(),
		// CmdRemoveService(),
		CmdCreateVc(),
		CmdUpdateVc(),
		CmdRemoveVc(),
		CmdBondCoreID(),
		CmdSetZkContractAddress(),
	)
	return cmd
}

func CmdCreateDid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-did [did] [public-key-for-address]",
		Short: "create did",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			did := args[0]
			pubkey := args[1]
			msg := types.NewMsgCreateDid(clientCtx.GetFromAddress().String(), did, pubkey)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdUpdateDidStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-did-status [did] [status]",
		Short: "update did status, status: 1-active, 0-inactive",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			did := args[0]
			status, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			msg := types.NewMsgUpdateDidStatus(clientCtx.GetFromAddress().String(), did, types.DidStatus(status))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// func CmdRemoveDid() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "remove-did [did]",
// 		Short: "remove did",
// 		Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			clientCtx, err := client.GetClientTxContext(cmd)
// 			if err != nil {
// 				return err
// 			}

// 			did := args[0]
// 			msg := types.NewMsgRemoveDid(clientCtx.GetFromAddress().String(), did)
// 			if err := msg.ValidateBasic(); err != nil {
// 				return err
// 			}

// 			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
// 		},
// 	}

// 	flags.AddTxFlagsToCmd(cmd)
// 	return cmd
// }

func CmdCreateService() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-service [sid] [name] [description] [issuer]",
		Short: "create credential service",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sid := args[0]
			name := args[1]
			description := args[2]
			issuers := strings.Split(args[3], ",")
			msg := types.NewMsgCreateService(clientCtx.GetFromAddress().String(), sid, name, description, issuers)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdUpdateServiceStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-service-status [sid] [status]",
		Short: "update credential service status, status: 1-active, 0-inactive",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			sid := args[0]
			status, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			msg := types.NewMsgUpdateServiceStatus(clientCtx.GetFromAddress().String(), sid, types.ServiceStatus(status))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// func CmdRemoveService() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "remove-service [sid]",
// 		Short: "remove credential service",
// 		Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			clientCtx, err := client.GetClientTxContext(cmd)
// 			if err != nil {
// 				return err
// 			}

// 			sid := args[0]
// 			msg := types.NewMsgRemoveService(clientCtx.GetFromAddress().String(), sid)
// 			if err := msg.ValidateBasic(); err != nil {
// 				return err
// 			}

// 			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
// 		},
// 	}

// 	flags.AddTxFlagsToCmd(cmd)
// 	return cmd
// }

func CmdCreateVc() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-vc [holder-did] [sid] [credential-file-hash] [off-chain-credential-uri] [hex-data]",
		Short: "create verifiable credential",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			holder := args[0]
			sid := args[1]
			hash := args[2]
			uri := args[3]
			data, err := hex.DecodeString(args[4])
			if err != nil {
				return err
			}

			fs, _ := cmd.Flags().GetStringSlice(FlagFilters)
			var filters [][]byte
			for _, f := range fs {
				filter, err := hex.DecodeString(f)
				if err != nil {
					return err
				}

				filters = append(filters, filter)
			}

			msg := types.NewMsgCreateVC(clientCtx.GetFromAddress().String(), holder, sid, hash, uri, data, filters)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().StringSlice(FlagFilters, []string{}, "hex-filters")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateVc() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-vc [holder-did] [sid] [credential-file-hash] [off-chain-credential-uri] [hex-data]",
		Short: "update verifiable credential",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			holder := args[0]
			sid := args[1]
			hash := args[2]
			uri := args[3]
			data, err := hex.DecodeString(args[4])
			if err != nil {
				return err
			}

			fs, _ := cmd.Flags().GetStringSlice(FlagFilters)
			var filters [][]byte
			for _, f := range fs {
				filter, err := hex.DecodeString(f)
				if err != nil {
					return err
				}

				filters = append(filters, filter)
			}

			msg := types.NewMsgUpdateVC(clientCtx.GetFromAddress().String(), holder, sid, hash, uri, data, filters)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().StringSlice(FlagFilters, []string{}, "hex-filters")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdRemoveVc() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-vc [holder] [sid]",
		Short: "remove verifiable credential",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			holder := args[0]
			sid := args[1]
			msg := types.NewMsgRemoveVC(clientCtx.GetFromAddress().String(), holder, sid)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdBondCoreID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bond-core-id [core-id] [auth-proof] [nounce] [expires-at]",
		Short: "bond core ID",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			coreID := args[0]
			authProof,err := hex.DecodeString(args[1])
			if err != nil{
				return fmt.Errorf("invalid auth proof: %w",err)
			}
			nounce, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid nounce: %w", err)
			}
			expiresAt, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid expires-at: %w", err)
			}
			msg := types.NewMsgBondZkCoreID(clientCtx.GetFromAddress().String(),  &types.ZkAuthClaimCoreIdInfo{
				CoreId:    []byte(coreID),
				AuthProof: authProof,
				Nonce:     nounce,
				ExpiresAt: expiresAt})
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}



func CmdSetZkContractAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-zk-contract-address [zk-contract-address] [type]",
		Short: "set zk contract address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			zkContractAddress := args[0]
			iType, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			msg := types.NewMsgSetZkContractAddress(clientCtx.GetFromAddress().String(), &types.ZkContractAddressInfo{
				ContractAddress: zkContractAddress,
				Type:            types.ZkContractAddressType(iType)})
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
