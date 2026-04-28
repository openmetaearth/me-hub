package cli

import (
	"github.com/openmetaearth/me-hub/x/megroup/types"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdListGroupMember() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-group-member [groupID]",
		Short: "list  group_member",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			strGrpID := args[0]
			grpId, _ := strconv.ParseUint(strGrpID, 10, 64)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGroupAllMemberRequest{
				GroupID:    grpId,
				Pagination: pageReq,
			}

			res, err := queryClient.GroupMemberAll(cmd.Context(), params)
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

func CmdShowGroupMember() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-group-member [address]",
		Short: "shows a group_member",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			address := args[0]

			params := &types.QueryGetGroupMemberRequest{
				Address: address,
			}

			res, err := queryClient.GroupMember(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
