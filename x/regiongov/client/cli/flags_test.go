package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

func TestParseProposalKind(t *testing.T) {
	for _, tc := range []struct {
		value string
		want  types.ProposalKind
	}{
		{"", types.ProposalKindUnspecified},
		{"treasury", types.ProposalKindTreasuryWithdraw},
		{"replace-operator", types.ProposalKindReplaceRegionOperator},
	} {
		t.Run(tc.value, func(t *testing.T) {
			cmd := CmdQueryProposals()
			require.NoError(t, cmd.Flags().Set(flagKind, tc.value))
			kind, err := parseProposalKind(cmd)
			require.NoError(t, err)
			require.Equal(t, tc.want, kind)
		})
	}

	for _, value := range []string{"invalid", "update-min-voting-period", "PROPOSAL_KIND_UPDATE_MIN_VOTING_PERIOD", "3"} {
		t.Run(value, func(t *testing.T) {
			cmd := CmdQueryProposals()
			require.NoError(t, cmd.Flags().Set(flagKind, value))
			kind, err := parseProposalKind(cmd)
			require.ErrorContains(t, err, "unknown kind")
			require.Equal(t, types.ProposalKindUnspecified, kind)
		})
	}
}
