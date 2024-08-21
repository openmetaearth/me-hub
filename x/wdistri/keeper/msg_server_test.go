package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/st-chain/me-hub/testutil/keeper"
	"github.com/st-chain/me-hub/x/wdistri/keeper"
	"github.com/st-chain/me-hub/x/wdistri/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.WdistriKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}

func TestMsgServer(t *testing.T) {
	ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
}
