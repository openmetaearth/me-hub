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
		desc        string
		genState    *types.GenesisState
		valid       bool
		errContains string
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
		{
			desc: "duplicate sequencer address is invalid",
			genState: &types.GenesisState{
				SequencerList: []types.Sequencer{
					testSequencer("seq-1", "rollapp-1", types.Bonded, true),
					testSequencer("seq-1", "rollapp-2", types.Bonded, true),
				},
				Params: params,
			},
			valid:       false,
			errContains: "duplicated index for sequencer",
		},
		{
			desc: "same rollapp cannot have multiple bonded proposers",
			genState: &types.GenesisState{
				SequencerList: []types.Sequencer{
					testSequencer("seq-1", "rollapp-1", types.Bonded, true),
					testSequencer("seq-2", "rollapp-1", types.Bonded, true),
				},
				Params: params,
			},
			valid:       false,
			errContains: "multiple bonded proposers for rollapp rollapp-1",
		},
		{
			desc: "different rollapps can each have one bonded proposer",
			genState: &types.GenesisState{
				SequencerList: []types.Sequencer{
					testSequencer("seq-1", "rollapp-1", types.Bonded, true),
					testSequencer("seq-2", "rollapp-2", types.Bonded, true),
				},
				Params: params,
			},
			valid: true,
		},
		{
			desc: "unbonded proposer flag does not count as active rollapp proposer",
			genState: &types.GenesisState{
				SequencerList: []types.Sequencer{
					testSequencer("seq-1", "rollapp-1", types.Unbonded, true),
					testSequencer("seq-2", "rollapp-1", types.Bonded, true),
				},
				Params: params,
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
				require.Contains(t, err.Error(), tc.errContains)
			}
		})
	}
}

func testSequencer(address, rollappID string, status types.OperatingStatus, proposer bool) types.Sequencer {
	sequencer := types.Sequencer{
		SequencerAddress: address,
		RollappId:        rollappID,
		Status:           status,
		Proposer:         proposer,
	}
	if status == types.Bonded || status == types.Unbonding {
		sequencer.Tokens = sdk.NewCoins(types.DefaultParams().MinBond)
	}
	return sequencer
}
