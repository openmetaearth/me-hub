package keeper

import (
	"testing"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"
)

func TestConsumeStaticCallResultChargesGas(t *testing.T) {
	tests := []struct {
		name      string
		vmError   string
		wantError bool
	}{
		{name: "success"},
		{name: "revert", vmError: "execution reverted", wantError: true},
		{name: "out of gas", vmError: "out of gas", wantError: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := sdk.Context{}.WithGasMeter(storetypes.NewGasMeter(100_000))
			ctx.GasMeter().ConsumeGas(5_000, "gas consumed before EVM call")
			result := &evmtypes.MsgEthereumTxResponse{
				GasUsed: 37_000,
				VmError: tc.vmError,
				Ret:     []byte{1, 2, 3},
			}

			ret, err := consumeStaticCallResult(ctx, result)

			require.Equal(t, uint64(42_000), ctx.GasMeter().GasConsumed())
			if tc.wantError {
				require.Error(t, err)
				require.Nil(t, ret)
				return
			}
			require.NoError(t, err)
			require.Equal(t, result.Ret, ret)
		})
	}
}

func TestConsumeStaticCallResultRejectsNilResult(t *testing.T) {
	ctx := sdk.Context{}.WithGasMeter(storetypes.NewGasMeter(100_000))

	ret, err := consumeStaticCallResult(ctx, nil)

	require.Error(t, err)
	require.Nil(t, ret)
	require.Zero(t, ctx.GasMeter().GasConsumed())
}
