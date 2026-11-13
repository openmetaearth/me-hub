package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

const (
	flagStatus = "status"
	flagKind   = "kind"
)

func parseProposalStatus(cmd *cobra.Command) (types.ProposalStatus, error) {
	value, err := cmd.Flags().GetString(flagStatus)
	if err != nil {
		return types.ProposalStatusUnspecified, err
	}
	switch value {
	case "":
		return types.ProposalStatusUnspecified, nil
	case "voting":
		return types.ProposalStatusVoting, nil
	case "passed":
		return types.ProposalStatusPassed, nil
	case "rejected":
		return types.ProposalStatusRejected, nil
	case "failed":
		return types.ProposalStatusFailed, nil
	default:
		return types.ProposalStatusUnspecified, fmt.Errorf("unknown status %q", value)
	}
}

func parseProposalKind(cmd *cobra.Command) (types.ProposalKind, error) {
	value, err := cmd.Flags().GetString(flagKind)
	if err != nil {
		return types.ProposalKindUnspecified, err
	}
	switch value {
	case "":
		return types.ProposalKindUnspecified, nil
	case "treasury":
		return types.ProposalKindTreasuryWithdraw, nil
	case "replace-operator":
		return types.ProposalKindReplaceRegionOperator, nil
	default:
		return types.ProposalKindUnspecified, fmt.Errorf("unknown kind %q", value)
	}
}

func parsePendingStatus(cmd *cobra.Command) (types.PendingActionStatus, error) {
	value, err := cmd.Flags().GetString(flagStatus)
	if err != nil {
		return types.PendingActionStatusUnspecified, err
	}
	switch value {
	case "":
		return types.PendingActionStatusUnspecified, nil
	case "pending":
		return types.PendingActionStatusPending, nil
	case "approved":
		return types.PendingActionStatusApproved, nil
	case "rejected":
		return types.PendingActionStatusRejected, nil
	default:
		return types.PendingActionStatusUnspecified, fmt.Errorf("unknown status %q", value)
	}
}
