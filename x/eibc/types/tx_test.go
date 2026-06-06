package types_test

import (
	"strings"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/x/eibc/types"
)

func TestMsgFulfillOrderValidateBasicRequiresPositiveExpectedFee(t *testing.T) {
	validAddress := "cosmos18wvvwfmq77a6d8tza4h5sfuy2yj3jj88yqg82a"
	validOrderID := strings.Repeat("a", 64)

	tests := []struct {
		name      string
		fee       string
		expectErr string
	}{
		{
			name: "positive fee",
			fee:  "1",
		},
		{
			name:      "zero fee",
			fee:       "0",
			expectErr: types.ErrNonPositiveFee.Error(),
		},
		{
			name:      "negative fee",
			fee:       "-1",
			expectErr: types.ErrNegativeFee.Error(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msg := types.NewMsgFulfillOrder(validAddress, validOrderID, tc.fee)

			err := msg.ValidateBasic()
			if tc.expectErr != "" {
				require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
				require.ErrorContains(t, err, tc.expectErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgUpdateDemandOrderValidateBasicAllowsZeroFee(t *testing.T) {
	validAddress := "cosmos18wvvwfmq77a6d8tza4h5sfuy2yj3jj88yqg82a"
	validOrderID := strings.Repeat("a", 64)

	msg := types.NewMsgUpdateDemandOrder(validAddress, validOrderID, "0")

	require.NoError(t, msg.ValidateBasic())
}
