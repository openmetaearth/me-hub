package keeper_test

import (
	"github.com/openmetaearth/me-hub/x/wstaking/keeper"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *KeeperTestSuite) TestNonRegionQueriesRejectNilRequests() {
	s.SetupTest()

	expectedErr := status.Error(codes.InvalidArgument, "invalid request")
	querier := keeper.Querier{Keeper: s.Keeper()}

	_, err := s.Keeper().QueryAllRecord(s.Ctx, nil)
	s.Require().ErrorIs(err, expectedErr)

	_, err = querier.QueryRecordByAddress(s.Ctx, (*types.QueryRecordsByAddress)(nil))
	s.Require().ErrorIs(err, expectedErr)

	_, err = querier.QueryReviewRecordByID(s.Ctx, (*types.QueryReviewRecordByNumber)(nil))
	s.Require().ErrorIs(err, expectedErr)

	_, err = querier.AllDelegations(s.Ctx, (*types.QueryAllDelegationsRequest)(nil))
	s.Require().ErrorIs(err, expectedErr)
}
