package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/spf13/cast"
	"github.com/st-chain/me-hub/x/megroup/types"
)

func CmdListGroupMemberCount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-group-member-count",
		Short: "list all group_member_count",
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

			params := &types.QueryAllGroupMemberCountRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.GroupMemberCountAll(cmd.Context(), params)
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

func CmdShowGroupMemberCount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-group-member-count [group-id]",
		Short: "shows a group_member_count",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argGroupId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			params := &types.QueryGetGroupMemberCountRequest{
				GroupId: argGroupId,
			}

			res, err := queryClient.GroupMemberCount(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
