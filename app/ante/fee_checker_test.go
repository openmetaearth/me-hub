package ante

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/stretchr/testify/require"
)

type feeCheckerTx struct {
	fee sdk.Coins
	gas uint64
}

func (tx feeCheckerTx) GetMsgs() []sdk.Msg         { return nil }
func (tx feeCheckerTx) ValidateBasic() error       { return nil }
func (tx feeCheckerTx) GetGas() uint64             { return tx.gas }
func (tx feeCheckerTx) GetFee() sdk.Coins          { return tx.fee }
func (tx feeCheckerTx) FeePayer() sdk.AccAddress   { return nil }
func (tx feeCheckerTx) FeeGranter() sdk.AccAddress { return nil }

func TestFeeCheckerRequiresMinimumBaseFeeInDeliverTx(t *testing.T) {
	foreignFeeTx := feeCheckerTx{
		fee: sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(100))),
		gas: 200000,
	}

	for _, isCheckTx := range []bool{true, false} {
		t.Run(map[bool]string{true: "check_tx", false: "deliver_tx"}[isCheckTx], func(t *testing.T) {
			_, _, err := checkTxFeeWithValidatorMinGasPrices(sdk.Context{}.WithIsCheckTx(isCheckTx), foreignFeeTx)
			require.Error(t, err)
			require.Contains(t, err.Error(), "fee must greater than or equal 10000"+params.BaseDenom)
		})
	}

	baseFeeTx := feeCheckerTx{
		fee: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(10000))),
		gas: 200000,
	}

	fees, _, err := checkTxFeeWithValidatorMinGasPrices(sdk.Context{}.WithIsCheckTx(false), baseFeeTx)
	require.NoError(t, err)
	require.Equal(t, baseFeeTx.fee, fees)

	zeroFeeTx := feeCheckerTx{
		fee: sdk.Coins{sdk.NewCoin("stake", sdk.ZeroInt())},
		gas: 200000,
	}

	fees, _, err = checkTxFeeWithValidatorMinGasPrices(sdk.Context{}.WithIsCheckTx(false), zeroFeeTx)
	require.NoError(t, err)
	require.Equal(t, zeroFeeTx.fee, fees)
}
