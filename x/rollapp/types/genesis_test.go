package types_test

import (
	"testing"

	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/stretchr/testify/require"
)

var (
	deployerAddr1 = sample.AccAddress()
	deployerAddr2 = sample.AccAddress()
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
			desc: "valid genesis state with empty DeployerWhitelist",
			genState: &types.GenesisState{
				Params: types.Params{
					DisputePeriodInBlocks: types.DefaultGenesis().Params.DisputePeriodInBlocks,
					DeployerWhitelist:     []types.DeployerParams{},
				},
				RollappList: []types.Rollapp{
					{
						RollappId: "rollapp-1",
						Creator:   deployerAddr1,
					},
					{
						RollappId: "otherrollapp-1",
						Creator:   deployerAddr2,
					},
				},
				StateInfoList: []types.StateInfo{
					{
						StateInfoIndex: types.StateInfoIndex{
							RollappId: "rollapp-1",
							Index:     0,
						},
					},
					{
						StateInfoIndex: types.StateInfoIndex{
							RollappId: "otherrollapp-1",
							Index:     1,
						},
					},
				},
				LatestStateInfoIndexList: []types.StateInfoIndex{
					{
						RollappId: "rollapp-1",
					},
					{
						RollappId: "otherrollapp-1",
					},
				},
				BlockHeightToFinalizationQueueList: []types.BlockHeightToFinalizationQueue{
					{
						CreationHeight: 0,
					},
					{
						CreationHeight: 1,
					},
				},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "valid genesis state with DeployerWhitelist",
			genState: &types.GenesisState{
				Params: types.Params{
					DisputePeriodInBlocks: types.DefaultGenesis().Params.DisputePeriodInBlocks,
					DeployerWhitelist:     []types.DeployerParams{{deployerAddr1}, {deployerAddr2}},
				},
				RollappList: []types.Rollapp{
					{
						RollappId: "rollapp-1",
						Creator:   deployerAddr1,
					},
					{
						RollappId: "otherrollapp-1",
						Creator:   deployerAddr2,
					},
				},
				StateInfoList: []types.StateInfo{
					{
						StateInfoIndex: types.StateInfoIndex{
							RollappId: "rollapp-1",
							Index:     0,
						},
					},
					{
						StateInfoIndex: types.StateInfoIndex{
							RollappId: "otherrollapp-1",
							Index:     1,
						},
					},
				},
				LatestStateInfoIndexList: []types.StateInfoIndex{
					{
						RollappId: "rollapp-1",
					},
					{
						RollappId: "otherrollapp-1",
					},
				},
				BlockHeightToFinalizationQueueList: []types.BlockHeightToFinalizationQueue{
					{
						CreationHeight: 0,
					},
					{
						CreationHeight: 1,
					},
				},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated deployer in whitelist",
			genState: &types.GenesisState{
				Params: types.Params{
					DisputePeriodInBlocks: types.DefaultGenesis().Params.DisputePeriodInBlocks,
					DeployerWhitelist:     []types.DeployerParams{{deployerAddr1}, {deployerAddr1}},
				},
				RollappList:                        []types.Rollapp{},
				StateInfoList:                      []types.StateInfo{},
				LatestStateInfoIndexList:           []types.StateInfoIndex{},
				BlockHeightToFinalizationQueueList: []types.BlockHeightToFinalizationQueue{},
			},
			valid: false,
		},
		{
			desc: "duplicated rollapp",
			genState: &types.GenesisState{
				Params: types.Params{
					DisputePeriodInBlocks: types.DefaultGenesis().Params.DisputePeriodInBlocks,
					DeployerWhitelist:     []types.DeployerParams{},
				},
				RollappList:                        []types.Rollapp{{RollappId: "rollapp-1", Creator: deployerAddr1}, {RollappId: "rollapp-1", Creator: deployerAddr1}},
				StateInfoList:                      []types.StateInfo{},
				LatestStateInfoIndexList:           []types.StateInfoIndex{},
				BlockHeightToFinalizationQueueList: []types.BlockHeightToFinalizationQueue{},
			},
			valid: false,
		},
		{
			desc: "invalid rollapp basic validation",
			genState: &types.GenesisState{
				Params: types.Params{
					DisputePeriodInBlocks: types.DefaultGenesis().Params.DisputePeriodInBlocks,
					DeployerWhitelist:     []types.DeployerParams{},
				},
				RollappList: []types.Rollapp{
					{
						RollappId:     "rollapp-1",
						Creator:       deployerAddr1,
						MaxSequencers: types.MaxAllowedSequencers + 1,
					},
				},
				StateInfoList:                      []types.StateInfo{},
				LatestStateInfoIndexList:           []types.StateInfoIndex{},
				BlockHeightToFinalizationQueueList: []types.BlockHeightToFinalizationQueue{},
			},
			valid: false,
		},
		{
			desc: "duplicated active eip155 rollapp",
			genState: &types.GenesisState{
				Params: types.Params{
					DisputePeriodInBlocks: types.DefaultGenesis().Params.DisputePeriodInBlocks,
					DeployerWhitelist:     []types.DeployerParams{},
				},
				RollappList: []types.Rollapp{
					{
						RollappId: "alpha_42-1",
						Creator:   deployerAddr1,
					},
					{
						RollappId: "beta_42-1",
						Creator:   deployerAddr2,
					},
				},
				StateInfoList:                      []types.StateInfo{},
				LatestStateInfoIndexList:           []types.StateInfoIndex{},
				BlockHeightToFinalizationQueueList: []types.BlockHeightToFinalizationQueue{},
			},
			valid: false,
		},
		{
			desc: "valid frozen eip155 revision replacement",
			genState: &types.GenesisState{
				Params: types.Params{
					DisputePeriodInBlocks: types.DefaultGenesis().Params.DisputePeriodInBlocks,
					DeployerWhitelist:     []types.DeployerParams{},
				},
				RollappList: []types.Rollapp{
					{
						RollappId: "rollapp_42-1",
						Creator:   deployerAddr1,
						Frozen:    true,
					},
					{
						RollappId: "rollapp_42-2",
						Creator:   deployerAddr2,
					},
				},
				StateInfoList:                      []types.StateInfo{},
				LatestStateInfoIndexList:           []types.StateInfoIndex{},
				BlockHeightToFinalizationQueueList: []types.BlockHeightToFinalizationQueue{},
			},
			valid: true,
		},
		{
			desc: "invalid frozen eip155 revision replacement name",
			genState: &types.GenesisState{
				Params: types.Params{
					DisputePeriodInBlocks: types.DefaultGenesis().Params.DisputePeriodInBlocks,
					DeployerWhitelist:     []types.DeployerParams{},
				},
				RollappList: []types.Rollapp{
					{
						RollappId: "rollapp_42-1",
						Creator:   deployerAddr1,
						Frozen:    true,
					},
					{
						RollappId: "otherrollapp_42-2",
						Creator:   deployerAddr2,
					},
				},
				StateInfoList:                      []types.StateInfo{},
				LatestStateInfoIndexList:           []types.StateInfoIndex{},
				BlockHeightToFinalizationQueueList: []types.BlockHeightToFinalizationQueue{},
			},
			valid: false,
		},
		{
			desc: "invalid frozen eip155 revision replacement gap",
			genState: &types.GenesisState{
				Params: types.Params{
					DisputePeriodInBlocks: types.DefaultGenesis().Params.DisputePeriodInBlocks,
					DeployerWhitelist:     []types.DeployerParams{},
				},
				RollappList: []types.Rollapp{
					{
						RollappId: "rollapp_42-1",
						Creator:   deployerAddr1,
						Frozen:    true,
					},
					{
						RollappId: "rollapp_42-3",
						Creator:   deployerAddr2,
					},
				},
				StateInfoList:                      []types.StateInfo{},
				LatestStateInfoIndexList:           []types.StateInfoIndex{},
				BlockHeightToFinalizationQueueList: []types.BlockHeightToFinalizationQueue{},
			},
			valid: false,
		},
		{
			desc: "invalid DisputePeriodInBlocks",
			genState: &types.GenesisState{
				Params: types.Params{
					DisputePeriodInBlocks: types.MinDisputePeriodInBlocks - 1,
					DeployerWhitelist:     []types.DeployerParams{},
				},
				RollappList:                        []types.Rollapp{{RollappId: "0"}},
				StateInfoList:                      []types.StateInfo{},
				LatestStateInfoIndexList:           []types.StateInfoIndex{},
				BlockHeightToFinalizationQueueList: []types.BlockHeightToFinalizationQueue{},
			},
			valid: false,
		},
		{
			desc: "invalid DeployerWhitelist",
			genState: &types.GenesisState{
				Params: types.Params{
					DisputePeriodInBlocks: types.MinDisputePeriodInBlocks,
					DeployerWhitelist:     []types.DeployerParams{{"asdad"}},
				},
				RollappList:                        []types.Rollapp{{RollappId: "0"}},
				StateInfoList:                      []types.StateInfo{},
				LatestStateInfoIndexList:           []types.StateInfoIndex{},
				BlockHeightToFinalizationQueueList: []types.BlockHeightToFinalizationQueue{},
			},
			valid: false,
		},
		{
			desc: "duplicated stateInfo",
			genState: &types.GenesisState{
				Params:                             types.Params{},
				RollappList:                        []types.Rollapp{},
				StateInfoList:                      []types.StateInfo{{StateInfoIndex: types.StateInfoIndex{RollappId: "0", Index: 0}}, {StateInfoIndex: types.StateInfoIndex{RollappId: "0", Index: 0}}},
				LatestStateInfoIndexList:           []types.StateInfoIndex{},
				BlockHeightToFinalizationQueueList: []types.BlockHeightToFinalizationQueue{},
			},
			valid: false,
		},
		{
			desc: "duplicated latestStateInfoIndex",
			genState: &types.GenesisState{
				Params:                             types.Params{},
				RollappList:                        []types.Rollapp{},
				StateInfoList:                      []types.StateInfo{},
				LatestStateInfoIndexList:           []types.StateInfoIndex{{RollappId: "0"}, {RollappId: "0"}},
				BlockHeightToFinalizationQueueList: []types.BlockHeightToFinalizationQueue{},
			},
			valid: false,
		},
		{
			desc: "duplicated blockHeightToFinalizationQueue",
			genState: &types.GenesisState{
				Params:                             types.Params{},
				RollappList:                        []types.Rollapp{},
				StateInfoList:                      []types.StateInfo{},
				LatestStateInfoIndexList:           []types.StateInfoIndex{},
				BlockHeightToFinalizationQueueList: []types.BlockHeightToFinalizationQueue{{CreationHeight: 0}, {CreationHeight: 0}},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
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
