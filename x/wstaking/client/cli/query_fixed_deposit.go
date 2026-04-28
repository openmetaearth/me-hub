package cli

import (
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"github.com/spf13/cobra"
)

func CmdListFixedDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-fixed-deposit",
		Short: "list all fixed_deposit",
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

			params := &types.QueryAllFixedDepositRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.FixedDepositAll(cmd.Context(), params)
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

func CmdShowFixedDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-fixed-deposit [id]",
		Short: "show a fixed_deposit",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			params := &types.QueryGetFixedDepositRequest{
				Id: id,
			}

			res, err := queryClient.FixedDeposit(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

var _ = strconv.Itoa(0)

func CmdFixedDepositByAcct() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-fixed-deposit-by-acct [account] [query_type]",
		Short: "show all fixed-deposits of an account",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqAccount := args[0]
			reqQueryType := args[1]

			queryType, ok := types.FixedDepositState_value[strings.ToUpper(strings.Trim(reqQueryType, " "))]
			if !ok {
				return types.ErrParameter.Wrap("period error")
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryFixedDepositByAcctRequest{
				Account:   reqAccount,
				QueryType: types.FixedDepositState(queryType),
			}

			res, err := queryClient.FixedDepositByAcct(cmd.Context(), params)
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

var _ = strconv.Itoa(0)

func CmdFixedDepositByRegion() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fixed-deposit-by-region [region-id] [query_type]",
		Short:   "show fixed_deposit-by-region",
		Example: "med q staking fixed-deposit-by-region me_earth ALL_STATE",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqRegionId := args[0]
			reqQueryType := args[1]

			queryType, ok := types.FixedDepositState_value[strings.ToUpper(strings.Trim(reqQueryType, " "))]
			if !ok {
				return types.ErrParameter.Wrap("query type invalid")
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryFixedDepositByRegionRequest{
				RegionId:   reqRegionId,
				QueryType:  types.FixedDepositState(queryType),
				Pagination: pageReq,
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.FixedDepositByRegion(cmd.Context(), params)
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

func CmdShowFixedDepositAmountByAcct() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-fixed-deposit-amount-by-acct [account]",
		Short: "show the fixed deposit amount of an account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argAccount := args[0]

			params := &types.QueryFixedDepositAmountByMeidRequest{
				Account: argAccount,
			}

			res, err := queryClient.FixedDepositAmountByMeid(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowFixedDepositTotalAmount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-fixed-deposit-total-amount",
		Short: "show the total amount of all fixed_Deposits",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryFixedDepositTotalAmountRequest{}

			res, err := queryClient.FixedDepositTotalAmount(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
