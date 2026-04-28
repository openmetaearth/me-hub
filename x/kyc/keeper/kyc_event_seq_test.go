package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/openmetaearth/me-hub/testutil/keeper"
	"github.com/openmetaearth/me-hub/testutil/nullify"
	"github.com/openmetaearth/me-hub/x/kyc/keeper"
	"github.com/openmetaearth/me-hub/x/kyc/types"
)

func createTestKycEventSeq(keeper *keeper.Keeper, ctx sdk.Context) types.KycEventSeq {
	item := types.KycEventSeq{Seq: 0}
	keeper.SetKycEventSeq(ctx, item)
	return item
}

func TestKycEventSeqGet(t *testing.T) {
	keeper, ctx := keepertest.KycKeeper(t)
	item := createTestKycEventSeq(keeper, ctx)
	rst, found := keeper.GetKycEventSeq(ctx)
	require.True(t, found)
	require.Equal(t,
		nullify.Fill(&item),
		nullify.Fill(&rst),
	)
}

func TestKycEventSeqRemove(t *testing.T) {
	keeper, ctx := keepertest.KycKeeper(t)
	createTestKycEventSeq(keeper, ctx)
	keeper.RemoveKycEventSeq(ctx)
	_, found := keeper.GetKycEventSeq(ctx)
	require.False(t, found)
}
