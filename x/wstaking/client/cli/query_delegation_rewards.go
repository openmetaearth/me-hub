package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/st-chain/me-hub/x/wstaking/types"
	"strings"
)

// GetCmdQueryDelegatorRewards implements the query delegator rewards command.
func GetCmdQueryDelegatorRewards() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()
	//bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()

	cmd := &cobra.Command{
		Use:   "rewards [delegator-addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Query all distribution delegator rewards or rewards from a particular validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all rewards earned by a delegator, optionally restrict to rewards from a single validator.

Example:
$ %s query distribution rewards %s1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p`,
				version.AppName, bech32PrefixAccAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			delegatorAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			// query for rewards from a particular delegation
			ctx := cmd.Context()
			//if len(args) == 2 {
			//validatorAddr, err := sdk.ValAddressFromBech32(args[1])
			//if err != nil {
			//	return err
			//}

			res, err := queryClient.DelegationRewards(
				ctx,
				&types.QueryDelegationRewardsRequest{DelegatorAddress: delegatorAddr.String()},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
			//}

			//res, err := queryClient.DelegationTotalRewards(
			//	ctx,
			//	&types.QueryDelegationTotalRewardsRequest{DelegatorAddress: delegatorAddr.String()},
			//)
			//if err != nil {
			//	return err
			//}
			//
			//return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
