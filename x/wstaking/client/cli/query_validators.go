package cli

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
)

const (
	FlagStatus = "status"
)

// GetCmdQueryValidators implements the query all validators command.
func GetCmdQueryValidators() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validators",
		Short: "Query for all validators",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(
			`Allows querying for all validators, optionally filtered by status.

Example:
$ med query staking validators
$ med query staking validators --status BOND_STATUS_BONDED
`,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := stakingtypes.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			statusFilter, _ := cmd.Flags().GetString(FlagStatus)
			res, err := queryClient.Validators(cmd.Context(), &stakingtypes.QueryValidatorsRequest{
				Status:     statusFilter,
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().String(FlagStatus, "", "Filter validators by status (BOND_STATUS_BONDED, BOND_STATUS_UNBONDING, BOND_STATUS_UNBONDED)")
	flags.AddPaginationFlagsToCmd(cmd, "validators")
	return cmd
}

// GetCmdQueryValidator implements the query single validator command.
func GetCmdQueryValidator() *cobra.Command {
	bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()
	cmd := &cobra.Command{
		Use:   "validator [validator-addr]",
		Short: "Query a specific validator",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			`Query details about an individual validator.

Example:
$ med query staking validator mevaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p
`,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := stakingtypes.NewQueryClient(clientCtx)
			if _, err := sdk.ValAddressFromBech32(args[0]); err != nil {
				return err
			}
			res, err := queryClient.Validator(cmd.Context(), &stakingtypes.QueryValidatorRequest{
				ValidatorAddr: args[0],
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	_ = bech32PrefixValAddr
	return cmd
}

// GetCmdQueryValidatorDelegations implements the command to query all the
// delegations to a specific validator.
func GetCmdQueryValidatorDelegations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegations-to [validator-addr]",
		Short: "Query all delegations made to one validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := stakingtypes.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			res, err := queryClient.ValidatorDelegations(cmd.Context(), &stakingtypes.QueryValidatorDelegationsRequest{
				ValidatorAddr: args[0],
				Pagination:    pageReq,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "delegations-to")
	return cmd
}

// GetCmdQueryPool implements the query staking pool command.
func GetCmdQueryPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool",
		Short: "Query the current staking pool values",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := stakingtypes.NewQueryClient(clientCtx)
			res, err := queryClient.Pool(cmd.Context(), &stakingtypes.QueryPoolRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryStakingParams implements the query staking params command.
func GetCmdQueryStakingParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current staking parameters information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := stakingtypes.NewQueryClient(clientCtx)
			res, err := queryClient.Params(cmd.Context(), &stakingtypes.QueryParamsRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
