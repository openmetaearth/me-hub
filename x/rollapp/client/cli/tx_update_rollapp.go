package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"strings"
)

func CmdUpdateRollapp() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update-rollapp [rollapp-id] [channel-id] [max-sequencers] [permissioned-addresses]",
		Short:   "update rollapp",
		Example: "med tx rollapp update-rollapp ROLLAPP_CHAIN_ID CHANNEL_ID MAX_SEQUENCERS ADDR1,ADDR2",
		Args:    cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argRollappId := args[0]
			argChannelId := args[1]
			argMaxSequencers, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}
			argPermissionedAddresses := strings.Split(args[3], ",")

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateRollapp(
				clientCtx.GetFromAddress().String(),
				argRollappId,
				argChannelId,
				argMaxSequencers,
				argPermissionedAddresses)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetCreateRollapp())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
