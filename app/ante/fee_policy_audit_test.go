package ante

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/stretchr/testify/require"
)

type auditFeeTx struct {
	fee sdk.Coins
	gas uint64
}

func (tx auditFeeTx) GetMsgs() []sdk.Msg         { return nil }
func (tx auditFeeTx) ValidateBasic() error       { return nil }
func (tx auditFeeTx) GetGas() uint64             { return tx.gas }
func (tx auditFeeTx) GetFee() sdk.Coins          { return tx.fee }
func (tx auditFeeTx) FeePayer() sdk.AccAddress   { return nil }
func (tx auditFeeTx) FeeGranter() sdk.AccAddress { return nil }

func TestAuditDeliverTxEnforcesMinimumFeePolicy(t *testing.T) {
	// Test case 1: Fee 100foo (foreign denom, less than minimum)
	txFoo := auditFeeTx{
		fee: sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(100))),
		gas: 200000,
	}

	// Should reject in CheckTx
	_, _, err := checkTxFeeWithValidatorMinGasPrices(sdk.Context{}.WithIsCheckTx(true), txFoo)
	require.Error(t, err)
	require.Contains(t, err.Error(), "fee must greater than or equal 10000"+params.BaseDenom)

	// Should now also reject in DeliverTx (IsCheckTx = false)
	_, _, err = checkTxFeeWithValidatorMinGasPrices(sdk.Context{}.WithIsCheckTx(false), txFoo)
	require.Error(t, err)
	require.Contains(t, err.Error(), "fee must greater than or equal 10000"+params.BaseDenom)

	// Test case 2: Fee 5000umec (insufficient base fee)
	txInsuff := auditFeeTx{
		fee: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(5000))),
		gas: 200000,
	}

	// Should reject in CheckTx
	_, _, err = checkTxFeeWithValidatorMinGasPrices(sdk.Context{}.WithIsCheckTx(true), txInsuff)
	require.Error(t, err)
	require.Contains(t, err.Error(), "fee must greater than or equal 10000"+params.BaseDenom)

	// Should reject in DeliverTx
	_, _, err = checkTxFeeWithValidatorMinGasPrices(sdk.Context{}.WithIsCheckTx(false), txInsuff)
	require.Error(t, err)
	require.Contains(t, err.Error(), "fee must greater than or equal 10000"+params.BaseDenom)

	// Test case 3: Fee 10000umec (valid minimum fee)
	txValid := auditFeeTx{
		fee: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(10000))),
		gas: 200000,
	}

	// Should pass in CheckTx
	_, _, err = checkTxFeeWithValidatorMinGasPrices(sdk.Context{}.WithIsCheckTx(true), txValid)
	require.NoError(t, err)

	// Should pass in DeliverTx
	_, _, err = checkTxFeeWithValidatorMinGasPrices(sdk.Context{}.WithIsCheckTx(false), txValid)
	require.NoError(t, err)
}
