package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/openmetaearth/me-hub/x/rollapp/types"
)

func CmdShowSkipDelayRollapp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skip-delay-rollapps",
		Short: "Query skip delayed rollapps",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.SkipDelayRollapp(cmd.Context(), &types.QuerySkipDelayRollappRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
