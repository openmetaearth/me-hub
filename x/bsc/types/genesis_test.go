package types

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	gravitytypes "github.com/openmetaearth/me-hub/x/gravity/types"
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
					GravityId:                          "me-bsc-bridge",
					AverageBlockTime:                   7_000,
					ExternalBatchTimeout:               12 * 3600 * 1000,
					AverageExternalBlockTime:           3_000,
					SignedWindow:                       30_000,
					SlashFraction:                      sdk.NewDec(8).Quo(sdk.NewDec(10)),
					RelayerSetUpdatePowerChangePercent: sdk.MustNewDecFromStr("0.1"),
					MaxRelayers:                        10,
					MinDelegate:                        sdk.NewInt(100_000_000),
					MaxDelegate:                        sdk.NewInt(10_000_000_000),
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
