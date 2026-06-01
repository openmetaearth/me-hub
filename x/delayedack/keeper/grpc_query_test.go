package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	commontypes "github.com/openmetaearth/me-hub/x/common/types"
	"github.com/openmetaearth/me-hub/x/delayedack/keeper"
	"github.com/openmetaearth/me-hub/x/delayedack/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *DelayedAckTestSuite) TestGetPacketsRejectsInvalidStatusWithoutPanic() {
	s.SetupTest()

	queryServer := keeper.NewQuerier(s.App.DelayedAckKeeper)
	var (
		res *types.QueryRollappPacketListResponse
		err error
	)

	s.Require().NotPanics(func() {
		res, err = queryServer.GetPackets(sdk.WrapSDKContext(s.Ctx), &types.QueryRollappPacketsRequest{
			Status: commontypes.Status(99),
		})
	})

	s.Require().Nil(res)
	s.Require().Error(err)
	s.Require().Equal(codes.InvalidArgument, status.Code(err))
	s.Require().ErrorContains(err, "invalid packet status")
}
