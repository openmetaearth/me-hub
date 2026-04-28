package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/version"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

// GetCmdQueryDelegatorRewards implements the query delegator rewards command.
func GetCmdQueryDelegatorRewards() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

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

			res, err := queryClient.DelegationRewards(
				ctx,
				&types.QueryDelegationRewardsRequest{DelegatorAddress: delegatorAddr.String()},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryDelegation the query delegation command.
func GetCmdQueryDelegation() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	cmd := &cobra.Command{
		Use:   "delegation [delegator-addr] ",
		Short: "Query a delegation based on address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query delegations for an individual delegator.

Example:
$ %s query staking delegation %s1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p
`,
				version.AppName, bech32PrefixAccAddr,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			delAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := &stakingtypes.QueryDelegationRequest{
				DelegatorAddr: delAddr.String(),
			}

			res, err := queryClient.Delegation(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.DelegationResponse)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryAllDelegations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-delegations",
		Short: "Query all delegations",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query delegations for an individual delegator on all validators.
Example:
$ %s query staking all-delegations 
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllDelegationsRequest{
				Pagination: &query.PageRequest{
					Key:        []byte(viper.GetString(flags.FlagPageKey)),
					Offset:     viper.GetUint64(flags.FlagOffset),
					Limit:      viper.GetUint64(flags.FlagLimit),
					CountTotal: true,
					Reverse:    false,
				},
			}

			res, err := queryClient.AllDelegations(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintObjectLegacy(res.Delegations)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all-delegations")
	return cmd
}
