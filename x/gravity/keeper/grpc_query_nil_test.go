package keeper_test

import (
	gravitykeeper "github.com/openmetaearth/me-hub/x/gravity/keeper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *KeeperTestSuite) TestDirectQueriesRejectNilRequests() {
	s.SetupTest()

	queryServer := gravitykeeper.NewQueryServerImpl(s.Keeper())
	expectedErr := status.Error(codes.InvalidArgument, "invalid request")

	testCases := []struct {
		name string
		call func() error
	}{
		{
			name: "RelayerSetRequest",
			call: func() error {
				_, err := queryServer.RelayerSetRequest(s.Ctx, nil)
				return err
			},
		},
		{
			name: "LastRelayerSetRequests",
			call: func() error {
				_, err := queryServer.LastRelayerSetRequests(s.Ctx, nil)
				return err
			},
		},
		{
			name: "BatchRequestByNonce",
			call: func() error {
				_, err := queryServer.BatchRequestByNonce(s.Ctx, nil)
				return err
			},
		},
		{
			name: "PendingOutgoingTxByAddr",
			call: func() error {
				_, err := queryServer.PendingOutgoingTxByAddr(s.Ctx, nil)
				return err
			},
		},
		{
			name: "BridgeTokens",
			call: func() error {
				_, err := queryServer.BridgeTokens(s.Ctx, nil)
				return err
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.Require().ErrorIs(tc.call(), expectedErr)
		})
	}
}
