package keeper_test

import (
	"strconv"
	"testing"

	keepertest "me-hub/testutil/keeper"
	"me-hub/testutil/nullify"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/megroup/keeper"
	"github.com/st-chain/me-hub/x/megroup/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNGroupMemberCount(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.GroupMemberCount {
	items := make([]types.GroupMemberCount, n)
	for i := range items {
		items[i].GroupId = uint64(i)

		keeper.SetGroupMemberCount(ctx, items[i])
	}
	return items
}

func TestGroupMemberCountGet(t *testing.T) {
	keeper, ctx := keepertest.MegroupKeeper(t)
	items := createNGroupMemberCount(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetGroupMemberCount(ctx,
			item.GroupId,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestGroupMemberCountRemove(t *testing.T) {
	keeper, ctx := keepertest.MegroupKeeper(t)
	items := createNGroupMemberCount(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveGroupMemberCount(ctx,
			item.GroupId,
		)
		_, found := keeper.GetGroupMemberCount(ctx,
			item.GroupId,
		)
		require.False(t, found)
	}
}

func TestGroupMemberCountGetAll(t *testing.T) {
	keeper, ctx := keepertest.MegroupKeeper(t)
	items := createNGroupMemberCount(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllGroupMemberCount(ctx)),
	)
}
