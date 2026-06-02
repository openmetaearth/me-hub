package types

import (
	"bytes"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestPowerReduction(t *testing.T) {
	t.Log(sdk.DefaultPowerReduction.String())

}

func TestMsgWithdrawFromGlobalDaoFeePoolValidateBasic(t *testing.T) {
	withdrawer := sdk.AccAddress(bytes.Repeat([]byte{0x01}, 20)).String()

	tests := []struct {
		name string
		msg  *MsgWithdrawFromGlobalDaoFeePool
		err  error
	}{
		{
			name: "valid",
			msg: NewMsgWithdrawFromGlobalDaoFeePool(
				withdrawer,
				sdk.NewCoins(sdk.NewInt64Coin("umec", 1)),
			),
		},
		{
			name: "invalid withdrawer",
			msg: NewMsgWithdrawFromGlobalDaoFeePool(
				"not-a-bech32-address",
				sdk.NewCoins(sdk.NewInt64Coin("umec", 1)),
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "empty amount",
			msg: NewMsgWithdrawFromGlobalDaoFeePool(
				withdrawer,
				sdk.Coins{},
			),
			err: sdkerrors.ErrInvalidCoins,
		},
		{
			name: "zero amount",
			msg: NewMsgWithdrawFromGlobalDaoFeePool(
				withdrawer,
				sdk.Coins{sdk.NewInt64Coin("umec", 0)},
			),
			err: sdkerrors.ErrInvalidCoins,
		},
		{
			name: "invalid denom",
			msg: NewMsgWithdrawFromGlobalDaoFeePool(
				withdrawer,
				sdk.Coins{{Denom: "bad denom", Amount: sdkmath.OneInt()}},
			),
			err: sdkerrors.ErrInvalidCoins,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
