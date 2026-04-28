package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"github.com/spf13/cobra"
)

func CmdShowRecordByAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-record [address]",
		Short: "show a record by address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			address := args[0]
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			param := &types.QueryRecordsByAddress{Account: address}
			res, err := queryClient.QueryRecordByAddress(cmd.Context(), param)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdShowAllRecord() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-records",
		Short: "show all records",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			params := &types.QueryAllRecords{
				Pagination: pageReq,
			}
			res, err := queryClient.QueryAllRecord(cmd.Context(), params)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdShowReviewRecord() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-review-record [actionNumber]",
		Short: "show a review record result",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			id := args[0]
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			param := types.QueryReviewRecordByNumber{ActionNumber: id}
			res, err := queryClient.QueryReviewRecordByID(cmd.Context(), &param)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
