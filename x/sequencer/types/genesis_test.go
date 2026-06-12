package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	params := types.DefaultParams()

	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "bonded sequencer with min bond is valid",
			genState: &types.GenesisState{
				Params: params,
				SequencerList: []types.Sequencer{
					{
						SequencerAddress: "sequencer-1",
						Status:           types.Bonded,
						Tokens:           sdk.NewCoins(params.MinBond),
					},
				},
			},
			valid: true,
		},
		{
			desc: "zero min bond is invalid",
			genState: &types.GenesisState{
				Params: types.Params{
					MinBond:       sdk.NewCoin(params.MinBond.Denom, sdk.ZeroInt()),
					UnbondingTime: params.UnbondingTime,
				},
			},
			valid: false,
		},
		{
			desc: "bonded sequencer without tokens is invalid",
			genState: &types.GenesisState{
				Params: params,
				SequencerList: []types.Sequencer{
					{
						SequencerAddress: "sequencer-1",
						Status:           types.Bonded,
					},
				},
			},
			valid: false,
		},
		{
			desc: "unbonding sequencer below min bond is invalid",
			genState: &types.GenesisState{
				Params: params,
				SequencerList: []types.Sequencer{
					{
						SequencerAddress: "sequencer-1",
						Status:           types.Unbonding,
						Tokens: sdk.NewCoins(sdk.NewCoin(
							params.MinBond.Denom,
							params.MinBond.Amount.Sub(sdk.OneInt()),
						)),
					},
				},
			},
			valid: false,
		},
		{
			desc: "unbonded sequencer without tokens is valid",
			genState: &types.GenesisState{
				Params: params,
				SequencerList: []types.Sequencer{
					{
						SequencerAddress: "sequencer-1",
						Status:           types.Unbonded,
					},
				},
			},
			valid: true,
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
