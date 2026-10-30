package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

// GetTxCmd returns the regiongov transaction command.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		CmdSubmitProposal(), CmdVote(), CmdSubmitEquityUpdate(),
		CmdSubmitFreezeEquity(), CmdSubmitUnfreezeEquity(),
		CmdApprovePendingAction(), CmdRejectPendingAction(),
		CmdFreezeRegion(), CmdUnfreezeRegion(), CmdFreezeEquity(), CmdUnfreezeEquity(),
	)
	return cmd
}

type txMessageBuilder func(client.Context, []string) (sdk.Msg, error)

func newTxCommand(use, short string, args cobra.PositionalArgs, build txMessageBuilder) *cobra.Command {
	cmd := &cobra.Command{
		Use: use, Short: short, Args: args,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg, err := build(clientCtx, args)
			if err != nil {
				return err
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdSubmitProposal() *cobra.Command {
	cmd := newTxCommand("submit-proposal [proposal.json]", "Submit a treasury withdrawal or region-operator replacement proposal", cobra.ExactArgs(1),
		func(ctx client.Context, args []string) (sdk.Msg, error) {
			data, err := os.ReadFile(args[0])
			if err != nil {
				return nil, err
			}
			var msg types.MsgSubmitProposal
			if err := ctx.Codec.UnmarshalJSON(data, &msg); err != nil {
				return nil, fmt.Errorf("decode proposal: %w", err)
			}
			msg.Proposer = ctx.GetFromAddress().String()
			return &msg, nil
		})
	cmd.Long = `Submit a proposal from a MsgSubmitProposal JSON file. The proposer is taken from --from.
The file must include kind, title, voting_period (seconds), and either treasury or replace_operator.
Submitting a proposal does not cast a yes vote; use the vote command separately.`
	return cmd
}

func CmdVote() *cobra.Command {
	return newTxCommand("vote [proposal-id] [yes|no]", "Vote on a proposal", cobra.ExactArgs(2),
		func(ctx client.Context, args []string) (sdk.Msg, error) {
			id, err := parseTxID(args[0])
			if err != nil {
				return nil, err
			}
			var option types.VoteOption
			switch strings.ToLower(args[1]) {
			case "yes":
				option = types.VoteOptionYes
			case "no":
				option = types.VoteOptionNo
			default:
				return nil, fmt.Errorf("vote option must be yes or no")
			}
			return &types.MsgVote{Voter: ctx.GetFromAddress().String(), ProposalId: id, Option: option}, nil
		})
}

func parseTxID(value string) (uint64, error) {
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil || id == 0 {
		return 0, fmt.Errorf("invalid id %q: must be a positive uint64", value)
	}
	return id, nil
}

func CmdSubmitEquityUpdate() *cobra.Command {
	return newTxCommand("submit-equity-update [region-id] [address ...]", "Request replacement of a region's equity list (omit addresses to clear it)", cobra.MinimumNArgs(1),
		func(ctx client.Context, args []string) (sdk.Msg, error) {
			return &types.MsgSubmitEquityUpdate{Proposer: ctx.GetFromAddress().String(), RegionId: args[0], Addresses: args[1:]}, nil
		})
}

func CmdSubmitFreezeEquity() *cobra.Command {
	return newTxCommand("submit-freeze-equity [region-id] [address ...]", "Request DAO approval to freeze equity addresses", cobra.MinimumNArgs(2),
		func(ctx client.Context, args []string) (sdk.Msg, error) {
			return &types.MsgSubmitFreezeEquity{Proposer: ctx.GetFromAddress().String(), RegionId: args[0], Addresses: args[1:]}, nil
		})
}

func CmdSubmitUnfreezeEquity() *cobra.Command {
	return newTxCommand("submit-unfreeze-equity [region-id] [address ...]", "Request DAO approval to unfreeze equity addresses", cobra.MinimumNArgs(2),
		func(ctx client.Context, args []string) (sdk.Msg, error) {
			return &types.MsgSubmitUnfreezeEquity{Proposer: ctx.GetFromAddress().String(), RegionId: args[0], Addresses: args[1:]}, nil
		})
}

func CmdApprovePendingAction() *cobra.Command {
	return newTxCommand("approve-pending-action [action-id]", "Approve a pending action as the DAO authority", cobra.ExactArgs(1),
		func(ctx client.Context, args []string) (sdk.Msg, error) {
			id, err := parseTxID(args[0])
			if err != nil {
				return nil, err
			}
			return &types.MsgApprovePendingAction{Authority: ctx.GetFromAddress().String(), ActionId: id}, nil
		})
}

func CmdRejectPendingAction() *cobra.Command {
	return newTxCommand("reject-pending-action [action-id]", "Reject a pending action as the DAO authority", cobra.ExactArgs(1),
		func(ctx client.Context, args []string) (sdk.Msg, error) {
			id, err := parseTxID(args[0])
			if err != nil {
				return nil, err
			}
			return &types.MsgRejectPendingAction{Authority: ctx.GetFromAddress().String(), ActionId: id}, nil
		})
}

func CmdFreezeRegion() *cobra.Command {
	return newTxCommand("freeze-region [region-id ...]", "Freeze regions as the DAO authority", cobra.MinimumNArgs(1),
		func(ctx client.Context, args []string) (sdk.Msg, error) {
			return &types.MsgFreezeRegion{Authority: ctx.GetFromAddress().String(), RegionIds: args}, nil
		})
}

func CmdUnfreezeRegion() *cobra.Command {
	return newTxCommand("unfreeze-region [region-id ...]", "Unfreeze regions as the DAO authority", cobra.MinimumNArgs(1),
		func(ctx client.Context, args []string) (sdk.Msg, error) {
			return &types.MsgUnfreezeRegion{Authority: ctx.GetFromAddress().String(), RegionIds: args}, nil
		})
}

func CmdFreezeEquity() *cobra.Command {
	return newTxCommand("freeze-equity [region-id] [address]", "Immediately freeze an equity address as the DAO authority", cobra.ExactArgs(2),
		func(ctx client.Context, args []string) (sdk.Msg, error) {
			return &types.MsgFreezeEquity{Authority: ctx.GetFromAddress().String(), RegionId: args[0], Address: args[1]}, nil
		})
}

func CmdUnfreezeEquity() *cobra.Command {
	return newTxCommand("unfreeze-equity [region-id] [address]", "Immediately unfreeze an equity address as the DAO authority", cobra.ExactArgs(2),
		func(ctx client.Context, args []string) (sdk.Msg, error) {
			return &types.MsgUnfreezeEquity{Authority: ctx.GetFromAddress().String(), RegionId: args[0], Address: args[1]}, nil
		})
}
