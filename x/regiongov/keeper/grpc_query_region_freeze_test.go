package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

func (s *KeeperTestSuite) TestRegionFreezeQueryCaseInsensitive() {
	k := s.App.RegionGovKeeper
	assertFrozen := func(want bool) {
		s.T().Helper()
		for _, id := range []string{"usa", "USA", "UsA"} {
			req := &types.QueryRegionFreezeRequest{RegionId: id}
			response, err := k.RegionFreeze(sdk.WrapSDKContext(s.Ctx), req)
			s.Require().NoError(err)
			s.Require().Equal(want, response.Frozen, "query region %s", id)
			s.Require().Equal(id, req.RegionId, "query must not mutate its request")
		}
	}
	assertFrozen(false)
	for _, id := range []string{"USA", "usa", "UsA"} {
		_, err := s.msgServer.FreezeRegion(s.Ctx, &types.MsgFreezeRegion{
			Authority: s.Dao.GlobalDao, RegionIds: []string{id},
		})
		s.Require().NoError(err)
		assertFrozen(true)
		state, found := k.GetRegionFreeze(s.Ctx, "usa")
		s.Require().True(found)
		s.Require().True(state.Frozen)
		for _, alias := range []string{"USA", "UsA"} {
			_, found := k.GetRegionFreeze(s.Ctx, alias)
			s.Require().False(found, "queries must not create alias records")
		}
		_, err = s.msgServer.UnfreezeRegion(s.Ctx, &types.MsgUnfreezeRegion{
			Authority: s.Dao.GlobalDao, RegionIds: []string{id},
		})
		s.Require().NoError(err)
		assertFrozen(false)
	}
	// Preserve the existing behavior for unknown regions and invalid requests.
	for _, id := range []string{"missing", "MISSING", "MiSsInG"} {
		response, err := k.RegionFreeze(sdk.WrapSDKContext(s.Ctx), &types.QueryRegionFreezeRequest{RegionId: id})
		s.Require().NoError(err)
		s.Require().False(response.Frozen)
	}
	for _, req := range []*types.QueryRegionFreezeRequest{nil, {}} {
		_, err := k.RegionFreeze(sdk.WrapSDKContext(s.Ctx), req)
		s.Require().Equal(codes.InvalidArgument, status.Code(err))
	}
}
