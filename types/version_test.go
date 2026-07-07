package types

import (
	"testing"

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
		{
			name:  "legacy testnet",
			input: "mechain_testnet",
			want:  "mechain_testnet_202404-1",
		},
		{
			name:  "testnet eip155",
			input: "mechain_testnet_202405-1",
			want:  "mechain_testnet_202405-1",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.want, ChainIdWithEIP155From(tc.input))
		})
	}
}
