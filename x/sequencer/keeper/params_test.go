package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testkeeper "github.com/openmetaearth/me-hub/testutil/keeper"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.SequencerKeeper(t)
	params := types.DefaultParams()
	params.MinBond = sdk.NewCoin("testdenom", sdk.NewInt(100))

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
