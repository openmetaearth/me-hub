package cli

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/openmetaearth/me-hub/x/kyc/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
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
		CmdApprove(),
		CmdUpdate(),
		CmdRemove(),
		CmdCreateSBT(),
		CmdUpdateSBT(),
		CmdDeleteSBT(),
	)
	return cmd
}

func CmdApprove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve [DID] [region ID] [address] [pubkey] [level] [uri] [hash] [inviter address]",
		Short: "approve KYC information",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			did := args[0]
			regionId := args[1]
			address := args[2]
			pubkey := args[3]
			level, err := strconv.Atoi(args[4])
			if err != nil {
				return err
			}
			uri := args[5]
			hash := args[6]
			inviter := args[7]

			msg := types.NewMsgApprove(
				clientCtx.GetFromAddress().String(),
				did,
				regionId,
				address,
				pubkey,
				didtypes.KycLevel(level),
				uri,
				hash,
				inviter,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [DID] [region ID] [level] [uri] [hash] [inviter]",
		Short: "update KYC information",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			did := args[0]
			regionId := args[1]
			level, err := strconv.Atoi(args[2])
			if err != nil {
				return err
			}
			uri := args[3]
			hash := args[4]
			inviter := args[5]

			msg := types.NewMsgUpdate(
				clientCtx.GetFromAddress().String(),
				did,
				regionId,
				didtypes.KycLevel(level),
				uri,
				hash,
				inviter,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdRemove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [DID]",
		Short: "remove KYC information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			did := args[0]
			msg := types.NewMsgRemove(
				clientCtx.GetFromAddress().String(),
				did,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdCreateSBT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-sbt [DID] [uri] [uri hash] [data]",
		Short: "create SBT(Soul binding token)",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			did := args[0]
			uri := args[1]
			uriHash := args[2]
			data, err := hex.DecodeString(args[3])
			if err != nil {
				return fmt.Errorf("data is not a valid hex string")
			}

			msg := types.NewMsgCreateSBT(
				clientCtx.GetFromAddress().String(),
				did,
				uri,
				uriHash,
				data,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdUpdateSBT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-sbt [DID] [uri] [uri hash] [data]",
		Short: "update SBT(Soul binding token)",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			did := args[0]
			uri := args[1]
			uriHash := args[2]
			data, err := hex.DecodeString(args[3])
			if err != nil {
				return fmt.Errorf("data is not a valid hex string")
			}

			msg := types.NewMsgUpdateSBT(
				clientCtx.GetFromAddress().String(),
				did,
				uri,
				uriHash,
				data,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdDeleteSBT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-sbt [DID]",
		Short: "delete SBT(Soul binding token)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			did := args[0]
			msg := types.NewMsgDeleteSBT(
				clientCtx.GetFromAddress().String(),
				did,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
