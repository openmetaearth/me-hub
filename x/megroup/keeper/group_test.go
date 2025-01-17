package keeper_test

import (
	"testing"

	keepertest "me-hub/testutil/keeper"
	"me-hub/testutil/nullify"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/megroup/keeper"
	"github.com/st-chain/me-hub/x/megroup/types"
	"github.com/stretchr/testify/require"
)

func createNGroup(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Group {
	items := make([]types.Group, n)
	for i := range items {
		items[i].Id = keeper.AppendGroup(ctx, items[i])
	}
	return items
}

func TestGroupGet(t *testing.T) {
	keeper, ctx := keepertest.MegroupKeeper(t)
	items := createNGroup(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetGroup(ctx, item.Id)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestGroupRemove(t *testing.T) {
	keeper, ctx := keepertest.MegroupKeeper(t)
	items := createNGroup(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveGroup(ctx, item.Id)
		_, found := keeper.GetGroup(ctx, item.Id)
		require.False(t, found)
	}
}

func TestGroupGetAll(t *testing.T) {
	keeper, ctx := keepertest.MegroupKeeper(t)
	items := createNGroup(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllGroup(ctx)),
	)
}

func TestGroupCount(t *testing.T) {
	keeper, ctx := keepertest.MegroupKeeper(t)
	items := createNGroup(keeper, ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, keeper.GetGroupCount(ctx))
}
