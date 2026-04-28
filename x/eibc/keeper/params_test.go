package keeper_test

import (
	"testing"

	testkeeper "github.com/openmetaearth/me-hub/testutil/keeper"
	"github.com/openmetaearth/me-hub/x/eibc/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.EibcKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
