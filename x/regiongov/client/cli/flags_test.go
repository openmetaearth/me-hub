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
		{"update-min-voting-period", types.ProposalKindUpdateMinVotingPeriod},
	} {
		t.Run(tc.value, func(t *testing.T) {
			cmd := CmdQueryProposals()
			require.NoError(t, cmd.Flags().Set(flagKind, tc.value))
			kind, err := parseProposalKind(cmd)
			require.NoError(t, err)
			require.Equal(t, tc.want, kind)
		})
	}

	cmd := CmdQueryProposals()
	require.NoError(t, cmd.Flags().Set(flagKind, "invalid"))
	_, err := parseProposalKind(cmd)
	require.Error(t, err)
}
