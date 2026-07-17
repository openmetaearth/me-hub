package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/x/eibc/types"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		}, {
			desc: "valid genesis state",
			genState: &types.GenesisState{
				Params:       validParams,
				DemandOrders: []types.DemandOrder{validDemandOrder},
			},
			valid: true,
		}, {
			desc: "invalid params",
			genState: &types.GenesisState{
				Params: types.Params{
					TimeoutFee: sdkmath.LegacyNewDec(-1),
					ErrackFee:  sdkmath.LegacyNewDec(-1),
				},
			},
			valid: false,
		}, {
			desc:     "invalid demand order",
			genState: &types.GenesisState{DemandOrders: []types.DemandOrder{{}}, Params: types.DefaultParams()},
			valid:    false,
		}, {
			desc: "duplicate demand order",
			genState: &types.GenesisState{DemandOrders: []types.DemandOrder{
				validDemandOrder,
				validDemandOrder,
			}, Params: types.DefaultParams()},
			valid: false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

var validDemandOrder = types.DemandOrder{
	Id:        "1",
	Price:     sdk.Coins{sdk.NewInt64Coin("denom", 2)},
	Fee:       sdk.Coins{sdk.NewInt64Coin("denom", 1)},
	Recipient: "cosmos18wvvwfmq77a6d8tza4h5sfuy2yj3jj88yqg82a",
}

var validParams = types.Params{
	EpochIdentifier: "hour",
	TimeoutFee:      sdkmath.LegacyNewDecWithPrec(1, 1),
	ErrackFee:       sdkmath.LegacyNewDecWithPrec(1, 1),
}
