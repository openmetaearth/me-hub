package types_test

import (
	"testing"

	_ "github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
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
			desc: "duplicate sequencer address is invalid",
			genState: &types.GenesisState{
				SequencerList: []types.Sequencer{
					testSequencer("seq-1", "rollapp-1", types.Bonded, true),
					testSequencer("seq-1", "rollapp-2", types.Bonded, true),
				},
				Params: types.DefaultParams(),
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
				Params: types.DefaultParams(),
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
				Params: types.DefaultParams(),
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
				Params: types.DefaultParams(),
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
	return types.Sequencer{
		SequencerAddress: address,
		RollappId:        rollappID,
		Status:           status,
		Proposer:         proposer,
	}
}
