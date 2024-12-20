package keeper_test

import (
	"testing"

	testkeeper "me-hub/testutil/keeper"

	"github.com/st-chain/me-hub/x/megroup/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.MegroupKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
