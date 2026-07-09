package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
	"github.com/spf13/cobra"
)

func CmdListSequencer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-sequencer",
		Short: "list all sequencer",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QuerySequencersRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.Sequencers(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowSequencer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-sequencer [sequencer-address]",
		Short: "shows a sequencer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argSequencerAddress := args[0]

			params := &types.QueryGetSequencerRequest{
				SequencerAddress: argSequencerAddress,
			}

			res, err := queryClient.Sequencer(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
