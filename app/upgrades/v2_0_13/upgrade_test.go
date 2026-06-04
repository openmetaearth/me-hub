package v2_0_13

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGovPeriodsUseDeliberativeWindows(t *testing.T) {
	require.Equal(t, 14*24*time.Hour, govMaxDepositPeriod)
	require.Equal(t, 7*24*time.Hour, govVotingPeriod)
	require.GreaterOrEqual(t, govMaxDepositPeriod, 7*24*time.Hour)
	require.GreaterOrEqual(t, govVotingPeriod, 7*24*time.Hour)
}
