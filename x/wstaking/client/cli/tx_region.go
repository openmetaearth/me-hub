package cli

import (
	"github.com/st-chain/me-hub/utils"
	"strings"

	"github.com/st-chain/me-hub/x/wstaking/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdNewRegion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new-region [name] [validator]",
		Short: "Broadcast message new-region",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argRegionId := strings.ToLower(args[0])
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
				argRegionId,
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
