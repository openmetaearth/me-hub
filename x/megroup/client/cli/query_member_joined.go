package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/openmetaearth/me-hub/x/megroup/types"
)

/*
func CmdListMemberJoined() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-member-joined",
		Short: "list all member_joined",
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

			params := &types.QueryAllMemberJoinedRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.MemberJoinedAll(cmd.Context(), params)
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

*/

func CmdShowMemberJoinedGroup() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-member-joined-group [address]",
		Short: "shows a member_joined",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argAddress := args[0]

			params := &types.QueryGroupByMemberRequest{
				Address: argAddress,
			}

			res, err := queryClient.GroupByMember(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
