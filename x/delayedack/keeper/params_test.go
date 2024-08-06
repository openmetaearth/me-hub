package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	testkeeper "github.com/st-chain/me-hub/testutil/keeper"
	"github.com/st-chain/me-hub/x/delayedack/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.DelayedackKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
