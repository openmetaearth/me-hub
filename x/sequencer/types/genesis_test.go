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
			desc: "duplicated replaceProposer",
			genState: &types.GenesisState{
				Params:        types.DefaultParams(),
				SequencerList: []types.Sequencer{},
				ReplaceProposerList: []types.MsgStoreReplaceProposer{
					{
						ReplaceProposer: types.MsgRepalceProposer{
							RollappId: "rollapp_1234-1",
						},
					},
					{
						ReplaceProposer: types.MsgRepalceProposer{
							RollappId: "rollapp_1234-1",
						},
					},
				},
			},
			valid: false,
		},
		{
			desc: "empty replaceProposer rollapp id",
			genState: &types.GenesisState{
				Params:        types.DefaultParams(),
				SequencerList: []types.Sequencer{},
				ReplaceProposerList: []types.MsgStoreReplaceProposer{
					{
						ReplaceProposer: types.MsgRepalceProposer{},
					},
				},
			},
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
