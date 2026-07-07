package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdSkipDelayRollapp() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "skip-delay-rollapp [rollapp-id] [is-skip]",
		Short:   "skip delay rollapp",
		Example: "med tx rollapp skip-delay-rollapp ROLLAPP_CHAIN_ID true",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argRollappId := args[0]

			argIsSkip := cast.ToBool(args[1])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSkipDelayRollapp(clientCtx.GetFromAddress().String(), argRollappId, argIsSkip)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetCreateRollapp())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
