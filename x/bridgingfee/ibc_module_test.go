package bridgingfee

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestChargeableBridgingFee(t *testing.T) {
	amount := sdk.NewInt(10)

	testCases := []struct {
		name string
		fee  sdk.Int
		want sdk.Int
	}{
		{
			name: "zero fee",
			fee:  sdk.ZeroInt(),
			want: sdk.ZeroInt(),
		},
		{
			name: "negative fee",
			fee:  sdk.NewInt(-1),
			want: sdk.ZeroInt(),
		},
		{
			name: "fee below amount",
			fee:  sdk.NewInt(3),
			want: sdk.NewInt(3),
		},
		{
			name: "fee equal amount",
			fee:  sdk.NewInt(10),
			want: sdk.NewInt(10),
		},
		{
			name: "fee above amount",
			fee:  sdk.NewInt(11),
			want: sdk.ZeroInt(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.True(t, tc.want.Equal(chargeableBridgingFee(amount, tc.fee)))
		})
	}
}
