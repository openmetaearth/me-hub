package megroup_test

import (
	"testing"

	keepertest "me-hub/testutil/keeper"
	"me-hub/testutil/nullify"

	"github.com/st-chain/me-hub/x/megroup"
	"github.com/st-chain/me-hub/x/megroup/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

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
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.MegroupKeeper(t)
	megroup.InitGenesis(ctx, *k, genesisState)
	got := megroup.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.GroupList, got.GroupList)
	require.Equal(t, genesisState.GroupCount, got.GroupCount)
	require.ElementsMatch(t, genesisState.GroupMemberList, got.GroupMemberList)
	require.ElementsMatch(t, genesisState.MemberJoinedList, got.MemberJoinedList)
	require.ElementsMatch(t, genesisState.GroupMemberCountList, got.GroupMemberCountList)
	// this line is used by starport scaffolding # genesis/test/assert
}
