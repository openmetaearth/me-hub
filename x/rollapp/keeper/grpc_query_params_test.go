package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	testkeeper "github.com/openmetaearth/me-hub/testutil/keeper"
	"github.com/openmetaearth/me-hub/x/rollapp/types"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := testkeeper.RollappKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	response, err := keeper.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.EqualValues(t, params.DisputePeriodInBlocks, response.Params.DisputePeriodInBlocks)
}
