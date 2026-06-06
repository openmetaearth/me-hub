package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMergeKycIssuerDIDsPreservesExistingNonDaoIssuers(t *testing.T) {
	merged := mergeKycIssuerDIDs(
		[]string{"did:dao:global", "did:third-party", "did:dao:global"},
		[]string{"did:dao:global", "did:dao:meid"},
	)

	require.Equal(t, []string{
		"did:dao:global",
		"did:third-party",
		"did:dao:meid",
	}, merged)
}
