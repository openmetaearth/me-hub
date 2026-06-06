package keeper_test

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *KeeperTestSuite) TestRouterQueriesRejectNilRequests() {
	s.SetupTest()

	_, err := s.App.GravityRouterKeeper.Relayer(s.Ctx, nil)
	s.Require().ErrorIs(err, status.Error(codes.InvalidArgument, "invalid request"))
}
