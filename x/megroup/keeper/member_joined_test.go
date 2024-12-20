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

func createNMemberJoined(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.MemberJoined {
	items := make([]types.MemberJoined, n)
	for i := range items {
		items[i].Address = strconv.Itoa(i)

		keeper.SetMemberJoined(ctx, items[i])
	}
	return items
}

func TestMemberJoinedGet(t *testing.T) {
	keeper, ctx := keepertest.MegroupKeeper(t)
	items := createNMemberJoined(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetMemberJoined(ctx,
			item.Address,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestMemberJoinedRemove(t *testing.T) {
	keeper, ctx := keepertest.MegroupKeeper(t)
	items := createNMemberJoined(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveMemberJoined(ctx,
			item.Address,
		)
		_, found := keeper.GetMemberJoined(ctx,
			item.Address,
		)
		require.False(t, found)
	}
}

func TestMemberJoinedGetAll(t *testing.T) {
	keeper, ctx := keepertest.MegroupKeeper(t)
	items := createNMemberJoined(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllMemberJoined(ctx)),
	)
}
