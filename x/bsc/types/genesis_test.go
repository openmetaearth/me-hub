package types

import (
	"reflect"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/st-chain/me-hub/types"
	gravitytypes "github.com/st-chain/me-hub/x/gravity/types"
)

func TestDefaultGenesisState(t *testing.T) {
	tests := []struct {
		name string
		want *gravitytypes.GenesisState
	}{
		{
			name: "bsc default genesis",
			want: &gravitytypes.GenesisState{
				Params: gravitytypes.Params{
					GravityId:                         "fx-bsc-bridge",
					AverageBlockTime:                  7_000,
					AverageExternalBlockTime:          3_000,
					ExternalBatchTimeout:              12 * 3600 * 1000,
					SignedWindow:                      30_000,
					SlashFraction:                     sdk.NewDec(8).Quo(sdk.NewDec(10)),
					OracleSetUpdatePowerChangePercent: sdk.NewDec(1).Quo(sdk.NewDec(10)),
					IbcTransferTimeoutHeight:          20_000,
					DelegateThreshold:                 sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10_000).MulRaw(1e18)),
					DelegateMultiple:                  10,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultGenesisState(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultGenesisState() = %v, want %v", got, tt.want)
			}
		})
	}
}
