package types

import (
	"testing"

	"github.com/evmos/ethermint/types"
	"github.com/stretchr/testify/require"
)

func TestChainIdWithEIP155From(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "legacy mainnet",
			input: "mechain",
			want:  "mechain_2404-1",
		},
		{
			name:  "mainnet eip155",
			input: "mechain_2404-1",
			want:  "mechain_2404-1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.want, ChainIdWithEIP155From(tc.input))
		})
	}
}

func TestChainIdWithEIP155(t *testing.T) {
	_, err := types.ParseChainID("me-chain_2404-1")
	require.ErrorIs(t, err, types.ErrInvalidChainID)
}
