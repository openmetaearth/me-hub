package cli

import (
	"github.com/openmetaearth/me-hub/x/wstaking/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdNewRecord() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new-record [activity-number] [url]",
		Short: "create a new record",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			actionNum := args[0]
			url := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRecord(
				actionNum,
				url,
				clientCtx.GetFromAddress().String(),
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

func CmdNewReviewRecord() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "review-record [hash] [result] [recordNumber] [reviewed-address]",
		Short: "send a review record result",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			hash := args[0]
			result := args[1]
			id := args[2]
			reviewedAddress := args[3]
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgReviewRecord(
				hash,
				result,
				clientCtx.GetFromAddress().String(),
				id,
				reviewedAddress,
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
