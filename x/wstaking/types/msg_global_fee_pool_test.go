package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgWithdrawFromGlobalDaoFeePool_ValidateBasic(t *testing.T) {
	validAddress := sample.AccAddress()

	tests := []struct {
		name string
		msg  *MsgWithdrawFromGlobalDaoFeePool
		err  error
	}{
		{
			name: "valid withdrawal",
			msg: NewMsgWithdrawFromGlobalDaoFeePool(
				validAddress,
				sdk.NewCoins(sdk.NewInt64Coin("adym", 1)),
			),
		},
		{
			name: "invalid withdrawer",
			msg: NewMsgWithdrawFromGlobalDaoFeePool(
				"invalid_address",
				sdk.NewCoins(sdk.NewInt64Coin("adym", 1)),
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "empty amount",
			msg: NewMsgWithdrawFromGlobalDaoFeePool(
				validAddress,
				sdk.Coins{},
			),
			err: sdkerrors.ErrInvalidCoins,
		},
		{
			name: "zero amount",
			msg: NewMsgWithdrawFromGlobalDaoFeePool(
				validAddress,
				sdk.Coins{sdk.NewInt64Coin("adym", 0)},
			),
			err: sdkerrors.ErrInvalidCoins,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
