package cli

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/openmetaearth/me-hub/utils"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func CmdNewRegion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new-region [name] [validator]",
		Short: "Broadcast message new-region",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argName := args[0]
			argVal := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			name := strings.Trim(argName, " ")
			if name != "" {
				name = strings.ToUpper(name)
				_, err = utils.CheckRegionName(name)
				if err != nil {
					return types.ErrRegionName.Wrap(err.Error())
				}
			}

			msg := types.NewMsgNewRegion(
				clientCtx.GetFromAddress().String(),
				name,
				argVal,
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

func CmdRemoveRegion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-region [region-id]",
		Short: "Broadcast message remove-region",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argRegionId := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRemoveRegion(
				clientCtx.GetFromAddress().String(),
				argRegionId,
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

func CmdTransferRegion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-region  [from-region] [to-region] [address]",
		Short: "transfer meid to new region.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			from := args[0]
			to := args[1]
			addr := args[2]
			address := make([]string, 0, 1)
			address = append(address, addr)
			creatorAddr := clientCtx.GetFromAddress().String()
			msg := &types.MsgTransferRegion{FromRegion: from, ToRegion: to, Address: address, Creator: creatorAddr}

			// if err = msg.ValidateBasic(); err != nil {
			//	return err
			//}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
