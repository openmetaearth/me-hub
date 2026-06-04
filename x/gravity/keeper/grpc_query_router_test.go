package keeper_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	bsctypes "github.com/openmetaearth/me-hub/x/bsc/types"
	"github.com/openmetaearth/me-hub/x/gravity/keeper"
	gravitytypes "github.com/openmetaearth/me-hub/x/gravity/types"
	trontypes "github.com/openmetaearth/me-hub/x/tron/types"
)

type routerQueryServer struct {
	gravitytypes.UnimplementedQueryServer

	relayerSet *gravitytypes.RelayerSet
	calls      int
}

func (s *routerQueryServer) LastObservedRelayer(
	_ context.Context,
	_ *gravitytypes.QueryLastObservedRelayer,
) (*gravitytypes.QueryLastObservedRelayerResponse, error) {
	s.calls++
	return &gravitytypes.QueryLastObservedRelayerResponse{RelayerSet: s.relayerSet}, nil
}

func TestRouterKeeperBridgeChainListDoesNotRequireBscRoute(t *testing.T) {
	rtr := keeper.NewRouter()
	rtr.AddRoute(trontypes.ModuleName, &keeper.ModuleHandler{
		QueryServer: &routerQueryServer{},
	})

	routerKeeper := keeper.NewRouterKeeper(rtr)
	res, err := routerKeeper.BridgeChainList(context.Background(), &gravitytypes.QueryBridgeChainListRequest{})

	require.NoError(t, err)
	require.ElementsMatch(t, gravitytypes.GetSupportChains(), res.ChainNames)
}

func TestRouterKeeperLastObservedRelayerRoutesByRequestChainName(t *testing.T) {
	bscServer := &routerQueryServer{
		relayerSet: &gravitytypes.RelayerSet{Nonce: 1},
	}
	tronServer := &routerQueryServer{
		relayerSet: &gravitytypes.RelayerSet{Nonce: 2},
	}
	rtr := keeper.NewRouter()
	rtr.AddRoute(bsctypes.ModuleName, &keeper.ModuleHandler{
		QueryServer: bscServer,
	})
	rtr.AddRoute(trontypes.ModuleName, &keeper.ModuleHandler{
		QueryServer: tronServer,
	})

	routerKeeper := keeper.NewRouterKeeper(rtr)
	res, err := routerKeeper.LastObservedRelayer(
		context.Background(),
		&gravitytypes.QueryLastObservedRelayer{ChainName: trontypes.ModuleName},
	)

	require.NoError(t, err)
	require.Equal(t, uint64(2), res.RelayerSet.Nonce)
	require.Zero(t, bscServer.calls)
	require.Equal(t, 1, tronServer.calls)
}
