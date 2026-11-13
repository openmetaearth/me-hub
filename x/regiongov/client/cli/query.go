package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

// GetQueryCmd returns the cli query commands for regiongov.
func GetQueryCmd(_ string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryParams(),
		CmdQueryProposal(),
		CmdQueryProposals(),
		CmdQueryVotes(),
		CmdQueryElectors(),
		CmdQueryEquity(),
		CmdQueryRegionFreeze(),
		CmdQueryEquityFreezes(),
		CmdQueryPendingAction(),
		CmdQueryPendingActions(),
	)
	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query regiongov params",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(clientCtx).Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal [id]",
		Short: "Query one proposal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid proposal id: %w", err)
			}
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(clientCtx).Proposal(cmd.Context(), &types.QueryProposalRequest{Id: id})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryProposals() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposals",
		Short: "Query proposals",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			status, err := parseProposalStatus(cmd)
			if err != nil {
				return err
			}
			kind, err := parseProposalKind(cmd)
			if err != nil {
				return err
			}
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(clientCtx).Proposals(cmd.Context(), &types.QueryProposalsRequest{
				Status:     status,
				Kind:       kind,
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	cmd.Flags().String(flagStatus, "", "filter by status: voting, passed, rejected, failed")
	cmd.Flags().String(flagKind, "", "filter by kind: treasury, replace-operator")
	flags.AddPaginationFlagsToCmd(cmd, "proposals")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryVotes() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "votes [proposal-id]",
		Short: "Query votes on a proposal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid proposal id: %w", err)
			}
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(clientCtx).Votes(cmd.Context(), &types.QueryVotesRequest{
				ProposalId: id,
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddPaginationFlagsToCmd(cmd, "votes")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryElectors() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "electors [proposal-id]",
		Short: "Query the snapshotted electorate of a proposal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid proposal id: %w", err)
			}
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(clientCtx).Electors(cmd.Context(), &types.QueryElectorsRequest{ProposalId: id})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryEquity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "equity [region-id]",
		Short: "Query the public equity addresses of a region",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(clientCtx).Equity(cmd.Context(), &types.QueryEquityRequest{RegionId: args[0]})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryRegionFreeze() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "region-freeze [region-id]",
		Short: "Query whether a region is frozen",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(clientCtx).RegionFreeze(cmd.Context(), &types.QueryRegionFreezeRequest{RegionId: args[0]})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryEquityFreezes() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "equity-freezes [region-id]",
		Short: "Query frozen equity addresses of a region",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(clientCtx).EquityFreezes(cmd.Context(), &types.QueryEquityFreezesRequest{RegionId: args[0]})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryPendingAction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-action [id]",
		Short: "Query one pending dao approval",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid action id: %w", err)
			}
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(clientCtx).PendingAction(cmd.Context(), &types.QueryPendingActionRequest{Id: id})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryPendingActions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-actions",
		Short: "Query pending dao approvals",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			status, err := parsePendingStatus(cmd)
			if err != nil {
				return err
			}
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(clientCtx).PendingActions(cmd.Context(), &types.QueryPendingActionsRequest{
				Status:     status,
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	cmd.Flags().String(flagStatus, "", "filter by status: pending, approved, rejected")
	flags.AddPaginationFlagsToCmd(cmd, "pending actions")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
