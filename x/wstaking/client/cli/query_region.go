package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"github.com/spf13/cobra"
	"strings"
)

func GetCmdQueryRegion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "region [region-id]",
		Short: "query a region",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argRegionId := args[0]

			params := &types.QueryRegionRequest{
				RegionId: strings.ToLower(argRegionId),
			}
			res, err := queryClient.Region(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdQueryAllRegion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "regions",
		Short: "query all region",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllRegionRequest{}
			res, err := queryClient.AllRegion(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdQueryRegionWithdrawer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "region-withdrawer [region-id]",
		Short: "Query which address is granted withdraw for a region",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryRegionWithdrawerRequest{
				RegionId: args[0],
			}
			res, err := queryClient.RegionWithdrawer(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
