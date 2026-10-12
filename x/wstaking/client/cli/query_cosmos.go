package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
)

const flagBondStatus = "status"

// Cosmos SDK 0.50 moved staking query CLIs to AutoCLI. wstaking replaces the
// staking module and exposes a custom GetQueryCmd, so those AutoCLI commands
// never appear. Re-add the v2 commands that called cosmos.staking.v1beta1.Query.

func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current staking parameters information",
		Long:  "Query values set as staking parameters.",
		Example: fmt.Sprintf(
			"$ %s query staking params",
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			res, err := stakingtypes.NewQueryClient(clientCtx).Params(cmd.Context(), &stakingtypes.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdQueryValidators() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validators",
		Args:  cobra.NoArgs,
		Short: "Query for all validators",
		Long:  "Query details about all validators on a network.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			status, _ := cmd.Flags().GetString(flagBondStatus)
			res, err := stakingtypes.NewQueryClient(clientCtx).Validators(cmd.Context(), &stakingtypes.QueryValidatorsRequest{
				Status:     status,
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String(flagBondStatus, "", "The validator bond status to filter by (BOND_STATUS_BONDED, BOND_STATUS_UNBONDED, or BOND_STATUS_UNBONDING)")
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "validators")
	return cmd
}

func GetCmdQueryValidator() *cobra.Command {
	bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()

	cmd := &cobra.Command{
		Use:   "validator [validator-addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Query a validator",
		Long:  "Query details about an individual validator.",
		Example: fmt.Sprintf(
			"$ %s query staking validator %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj",
			version.AppName, bech32PrefixValAddr,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			res, err := stakingtypes.NewQueryClient(clientCtx).Validator(cmd.Context(), &stakingtypes.QueryValidatorRequest{
				ValidatorAddr: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdQueryUnbondingDelegations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbonding-delegations [delegator-addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Query all unbonding-delegations records for one delegator",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := stakingtypes.NewQueryClient(clientCtx).DelegatorUnbondingDelegations(cmd.Context(), &stakingtypes.QueryDelegatorUnbondingDelegationsRequest{
				DelegatorAddr: args[0],
				Pagination:    pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "unbonding delegations")
	return cmd
}

func GetCmdQueryValidatorUnbondingDelegations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbonding-delegations-from [validator-addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Query all unbonding delegations from a validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := stakingtypes.NewQueryClient(clientCtx).ValidatorUnbondingDelegations(cmd.Context(), &stakingtypes.QueryValidatorUnbondingDelegationsRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "unbonding delegations")
	return cmd
}

func GetCmdQueryHistoricalInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "historical-info [height]",
		Args:  cobra.ExactArgs(1),
		Short: "Query historical info at given height",
		Example: fmt.Sprintf(
			"$ %s query staking historical-info 5",
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			height, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid height: %w", err)
			}

			res, err := stakingtypes.NewQueryClient(clientCtx).HistoricalInfo(cmd.Context(), &stakingtypes.QueryHistoricalInfoRequest{
				Height: height,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
