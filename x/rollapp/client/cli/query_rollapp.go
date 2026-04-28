package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/spf13/cobra"
)

func CmdListRollapp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Query all rollapps currently registered in the hub",
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

			params := &types.QueryAllRollappRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.RollappAll(cmd.Context(), params)
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

func CmdShowRollapp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show [rollapp-id]",
		Short: "Query the rollapp associated with the specified rollapp-id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argRollappId := args[0]

			params := &types.QueryGetRollappRequest{
				RollappId: argRollappId,
			}

			res, err := queryClient.Rollapp(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowSkipDelayRollapp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skip-delay-rollapps",
		Short: "Query skip delayed rollapp",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			params := &types.QuerySkipDelayRollappRequest{}
			res, err := queryClient.SkipDelayRollapp(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
