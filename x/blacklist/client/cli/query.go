package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/st-chain/me-hub/x/blacklist/types"
)

var _ = strconv.Itoa(0)

func CmdBlacklist() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blacklist",
		Short: "Query blacklist",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			address, err := cmd.Flags().GetString("address")
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryBlacklistRequest{
				Address:    address,
				Pagination: pageReq,
			}

			res, err := queryClient.Blacklist(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String("address", "", "The address to query. If not specified, returns all blacklisted addresses.")
	flags.AddPaginationFlagsToCmd(cmd, "blacklist")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
