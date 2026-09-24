package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/openmetaearth/me-hub/x/sequencer/types"
)

func CmdReplaceProposerInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replace-proposer-info [rollapp-id]",
		Short: "Query pending ReplaceProposer request for a rollapp",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ReplaceProposerInfo(cmd.Context(), &types.QueryReplaceProposerInfoRequest{
				RollappId: args[0],
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
