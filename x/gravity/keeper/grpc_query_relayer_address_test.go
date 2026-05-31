package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestParseQueryRelayerAddress(t *testing.T) {
	t.Run("rejects malformed input", func(t *testing.T) {
		require.NotPanics(t, func() {
			address, err := parseQueryRelayerAddress("not-a-bech32-address")

			require.Nil(t, address)
			require.Equal(t, codes.InvalidArgument, status.Code(err))
		})
	})

	t.Run("accepts valid input", func(t *testing.T) {
		want := sdk.AccAddress([]byte("01234567890123456789"))

		address, err := parseQueryRelayerAddress(want.String())

		require.NoError(t, err)
		require.Equal(t, want, address)
	})
}
