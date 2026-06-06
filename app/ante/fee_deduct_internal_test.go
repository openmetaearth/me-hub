package ante

import (
	"math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/stretchr/testify/require"
)

type feeCheckerTestTx struct {
	fee sdk.Coins
	gas uint64
}

func (tx feeCheckerTestTx) GetMsgs() []sdk.Msg       { return nil }
func (tx feeCheckerTestTx) ValidateBasic() error     { return nil }
func (tx feeCheckerTestTx) GetGas() uint64           { return tx.gas }
func (tx feeCheckerTestTx) GetFee() sdk.Coins        { return tx.fee }
func (tx feeCheckerTestTx) FeePayer() sdk.AccAddress { return nil }
func (tx feeCheckerTestTx) FeeGranter() sdk.AccAddress {
	return nil
}

func TestCheckTxFeeRejectsGasAboveMaxInt64(t *testing.T) {
	ctx := sdk.Context{}.
		WithIsCheckTx(true).
		WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(params.BaseDenom, sdk.OneInt())))
	tx := feeCheckerTestTx{
		fee: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(10_000))),
		gas: uint64(math.MaxInt64) + 1,
	}

	var err error
	require.NotPanics(t, func() {
		_, _, err = checkTxFeeWithValidatorMinGasPrices(ctx, tx)
	})
	require.ErrorIs(t, err, sdkerrors.ErrInvalidGasLimit)
}
