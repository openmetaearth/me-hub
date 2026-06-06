package keeper

import (
	"testing"

	"github.com/openmetaearth/me-hub/x/gravity/types"
	"github.com/stretchr/testify/require"
)

func TestProjectBatchTimeoutHeightRejectsUnsafeProjection(t *testing.T) {
	params := types.DefaultParams()
	params.AverageBlockTime = uint64(1) << 63
	params.AverageExternalBlockTime = 100
	params.ExternalBatchTimeout = 60_000

	projected, timeout := projectBatchTimeoutHeight(12, types.LastObservedBlockHeight{
		BlockHeight:         10,
		ExternalBlockHeight: 1_000,
	}, params)
	require.Zero(t, projected)
	require.Zero(t, timeout)
}

func TestProjectBatchTimeoutHeightRejectsZeroExternalBlockWindow(t *testing.T) {
	params := types.DefaultParams()
	params.ExternalBatchTimeout = 60_000
	params.AverageExternalBlockTime = 120_000

	projected, timeout := projectBatchTimeoutHeight(20, types.LastObservedBlockHeight{
		BlockHeight:         10,
		ExternalBlockHeight: 1_000,
	}, params)
	require.Zero(t, projected)
	require.Zero(t, timeout)
}

func TestProjectBatchTimeoutHeightHandlesImportedLocalHeight(t *testing.T) {
	params := types.DefaultParams()

	projected, timeout := projectBatchTimeoutHeight(1, types.LastObservedBlockHeight{
		BlockHeight:         10,
		ExternalBlockHeight: 1_000,
	}, params)
	require.Equal(t, uint64(1_000), projected)
	require.Greater(t, timeout, projected)
}
