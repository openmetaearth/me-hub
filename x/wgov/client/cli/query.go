package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	govcli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/spf13/cobra"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/openmetaearth/me-hub/x/wgov/types"
)

func GetQueryCmd() *cobra.Command {
	// Group gov queries under a subcommand
	govQueryCmd := &cobra.Command{
		Use:                        govtypes.ModuleName,
		Short:                      "Querying commands for the governance module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	govQueryCmd.AddCommand(
		govcli.GetCmdQueryProposal(),
		govcli.GetCmdQueryProposals(),
		govcli.GetCmdQueryVote(),
		govcli.GetCmdQueryVotes(),
		govcli.GetCmdQueryParams(),
		govcli.GetCmdQueryParam(),
		govcli.GetCmdQueryProposer(),
		govcli.GetCmdQueryDeposit(),
		govcli.GetCmdQueryDeposits(),
		GetCmdQueryMeTally(),
	)

	return govQueryCmd
}

// GetCmdQueryTally implements the command to query for proposal tally result.
func GetCmdQueryMeTally() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "me-tally [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Get the tally of a proposal vote",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query tally of votes on a proposal. You can find
the proposal-id by running "%s query gov proposals".

Example:
$ %s query gov tally 1
`,
				version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			ctx := cmd.Context()
			res, err := queryClient.MeTallyResult(
				ctx,
				&types.QueryMeTallyResultRequest{ProposalId: proposalID},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.Tally)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
