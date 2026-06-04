package types_test

import (
	"testing"

	_ "github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
	"github.com/stretchr/testify/require"
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
		},
		{
			desc: "same rollapp cannot have multiple bonded proposers",
			genState: &types.GenesisState{
				SequencerList: []types.Sequencer{
					{
						SequencerAddress: "seq-1",
						RollappId:        "rollapp-1",
						Status:           types.Bonded,
						Proposer:         true,
					},
					{
						SequencerAddress: "seq-2",
						RollappId:        "rollapp-1",
						Status:           types.Bonded,
						Proposer:         true,
					},
				},
				Params: types.DefaultParams(),
			},
			valid: false,
		},
		{
			desc: "different rollapps can each have a bonded proposer",
			genState: &types.GenesisState{
				SequencerList: []types.Sequencer{
					{
						SequencerAddress: "seq-1",
						RollappId:        "rollapp-1",
						Status:           types.Bonded,
						Proposer:         true,
					},
					{
						SequencerAddress: "seq-2",
						RollappId:        "rollapp-2",
						Status:           types.Bonded,
						Proposer:         true,
					},
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
			}
		})
	}
}
