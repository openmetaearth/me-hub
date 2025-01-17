package types_test

import (
	"testing"

	"github.com/st-chain/me-hub/x/megroup/types"

	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
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
			desc: "valid genesis state",
			genState: &types.GenesisState{

				GroupList: []types.Group{
					{
						Id: 0,
					},
					{
						Id: 1,
					},
				},
				GroupCount: 2,
				GroupMemberList: []types.GroupMember{
					{
						Id: 0,
					},
					{
						Id: 1,
					},
				},
				MemberJoinedList: []types.MemberJoined{
					{
						Address: "0",
					},
					{
						Address: "1",
					},
				},
				GroupMemberCountList: []types.GroupMemberCount{
					{
						GroupId: 0,
					},
					{
						GroupId: 1,
					},
				},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated group",
			genState: &types.GenesisState{
				GroupList: []types.Group{
					{
						Id: 0,
					},
					{
						Id: 0,
					},
				},
			},
			valid: false,
		},
		{
			desc: "invalid group count",
			genState: &types.GenesisState{
				GroupList: []types.Group{
					{
						Id: 1,
					},
				},
				GroupCount: 0,
			},
			valid: false,
		},
		{
			desc: "duplicated groupMember",
			genState: &types.GenesisState{
				GroupMemberList: []types.GroupMember{
					{
						Id: 0,
					},
					{
						Id: 0,
					},
				},
			},
			valid: false,
		},
		{
			desc: "invalid groupMember count",
			genState: &types.GenesisState{
				GroupMemberList: []types.GroupMember{
					{
						Id: 1,
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated memberJoined",
			genState: &types.GenesisState{
				MemberJoinedList: []types.MemberJoined{
					{
						Address: "0",
					},
					{
						Address: "0",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated groupMemberCount",
			genState: &types.GenesisState{
				GroupMemberCountList: []types.GroupMemberCount{
					{
						GroupId: 0,
					},
					{
						GroupId: 0,
					},
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	}
	for _, tc := range tests {
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
