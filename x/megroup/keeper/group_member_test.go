package keeper_test

import (
	"testing"

	keepertest "me-hub/testutil/keeper"
	"me-hub/testutil/nullify"

	"github.com/st-chain/me-hub/x/megroup/keeper"
	"github.com/st-chain/me-hub/x/megroup/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func createNGroupMember(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.GroupMember {
	items := make([]types.GroupMember, n)
	for i := range items {
		items[i].Id = keeper.LoadMemberStoreByGroupID(ctx, items[i].GroupID).AppendGroupMember(ctx, items[i])
	}
	return items
}

func TestGroupMemberGet(t *testing.T) {
	keeper, ctx := keepertest.MegroupKeeper(t)
	items := createNGroupMember(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.LoadMemberStoreByGroupID(ctx, item.GroupID).GetGroupMember(item.Id)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestGroupMemberRemove(t *testing.T) {
	keeper, ctx := keepertest.MegroupKeeper(t)
	items := createNGroupMember(keeper, ctx, 10)
	for _, item := range items {
		keeper.LoadMemberStoreByGroupID(ctx, item.GroupID).RemoveGroupMember(item.Id)
		_, found := keeper.LoadMemberStoreByGroupID(ctx, item.GroupID).GetGroupMember(item.Id)
		require.False(t, found)
	}
}

func TestGroupMemberGetAll(t *testing.T) {
	keeper, ctx := keepertest.MegroupKeeper(t)
	items := createNGroupMember(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllGroupMember(ctx)),
	)
}
