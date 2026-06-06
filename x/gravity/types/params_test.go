package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParamsValidateBasicRejectsUnsafeBatchTiming(t *testing.T) {
	params := DefaultParams()
	params.ExternalBatchTimeout = 60_000
	params.AverageExternalBlockTime = 120_000
	require.ErrorContains(t, params.ValidateBasic(), "at least one average external block")

	params = DefaultParams()
	params.AverageBlockTime = uint64(1) << 63
	params.AverageExternalBlockTime = 100
	params.ExternalBatchTimeout = 60_000
	require.ErrorContains(t, params.ValidateBasic(), "at least one average ME block")
}
