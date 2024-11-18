package cli

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/st-chain/me-hub/x/wstaking/types"

	"github.com/spf13/cobra"
	"strconv"
)

func CmdListFixedDepositCfg() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-fixed-deposit-cfg",
		Short: "show fixed deposit config",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			argRegionId := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryFixedDepositCfgRequest{
				RegionId: argRegionId,
			}

			res, err := queryClient.FixedDepositCfg(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryFixedDepositCfg() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-fixed-deposit-cfg-by-term",
		Short: "show fixed deposit config by term",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			argRegionId := args[0]
			argTerm := args[1]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			term, err := strconv.ParseInt(argTerm, 10, 64)
			if err != nil {
				return types.ErrParameter.Wrap("term error")
			}

			params := &types.QueryFixedDepositCfgByTermRequest{
				RegionId: argRegionId,
				Term:     term,
			}

			res, err := queryClient.FixedDepositCfgByTerm(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
