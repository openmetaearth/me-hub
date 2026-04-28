package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
	"github.com/spf13/cobra"
)

func CmdShowSequencersByRollapp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-sequencers-by-rollapp [rollapp-id]",
		Short: "shows a sequencers_by_rollapp",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argRollappId := args[0]

			params := &types.QueryGetSequencersByRollappRequest{
				RollappId: argRollappId,
			}

			res, err := queryClient.SequencersByRollapp(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowReplaceProposer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-replace-proposer [rollapp-id]",
		Short: "shows -replace-proposer ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argRollappId := args[0]

			params := &types.QueryReplaceProposerInfoRequest{
				RollappId: argRollappId,
			}

			res, err := queryClient.ReplaceProposerInfo(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
