package keeper_test

import (
	"fmt"
	"testing"

	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/x/wnft/types"
	"github.com/stretchr/testify/suite"
)

type WnftKeeperTestSuite struct {
	apptesting.KeeperTestHelper
}

func TestWnftKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(WnftKeeperTestSuite))
}

func (s *WnftKeeperTestSuite) SetupTest() {
	app := apptesting.Setup(s.T(), false)
	ctx := app.GetBaseApp().NewContext(false, cometbftproto.Header{})

	s.App = app
	s.Ctx = ctx
}

func (s *WnftKeeperTestSuite) TestClassAddressPagination() {
	k := s.App.WNFTKeeper
	ctx := s.Ctx
	goCtx := sdk.WrapSDKContext(ctx)

	classID := "testclass"
	err := k.SaveClass(ctx, nft.Class{
		Id:          classID,
		Name:        "Test Class",
		Symbol:      "TC",
		TotalSupply: 10,
	})
	s.Require().NoError(err)

	owner := s.NewAccounts(1)[0]
	for i := 1; i <= 5; i++ {
		err = k.Mint(ctx, nft.NFT{
			ClassId: classID,
			Id:      fmt.Sprintf("nft%d", i),
		}, owner)
		s.Require().NoError(err)
	}

	resp, err := k.ClassAddress(goCtx, &types.QueryClassAddressRequest{
		ClassId: classID,
		Address: owner.String(),
	})
	s.Require().NoError(err)
	s.Require().True(resp.Exists)
	s.Require().Len(resp.Nfts, 5)
	s.Require().NotNil(resp.Pagination)

	resp, err = k.ClassAddress(goCtx, &types.QueryClassAddressRequest{
		ClassId: classID,
		Address: owner.String(),
		Pagination: &query.PageRequest{
			Offset: 1,
			Limit:  2,
		},
	})
	s.Require().NoError(err)
	s.Require().Len(resp.Nfts, 2)
	s.Require().Equal([]string{"nft2", "nft3"}, resp.Nfts)
	if s.NotNil(resp.Pagination) {
		s.Require().NotNil(resp.Pagination.NextKey)
	}

	resp, err = k.ClassAddress(goCtx, &types.QueryClassAddressRequest{
		ClassId: classID,
		Address: owner.String(),
		Pagination: &query.PageRequest{
			Offset: 4,
			Limit:  2,
		},
	})
	s.Require().NoError(err)
	s.Require().Len(resp.Nfts, 1)
	s.Require().Equal([]string{"nft5"}, resp.Nfts)

	_, err = k.ClassAddress(goCtx, &types.QueryClassAddressRequest{
		ClassId: classID,
		Address: "invalidaddr",
	})
	s.Require().Error(err)
}

func (s *WnftKeeperTestSuite) TestNftFilterPagination() {
	k := s.App.WNFTKeeper
	ctx := s.Ctx
	goCtx := sdk.WrapSDKContext(ctx)

	classID1 := "testclass1"
	classID2 := "testclass2"

	err := k.SaveClass(ctx, nft.Class{Id: classID1, Name: "Test Class 1", TotalSupply: 10})
	s.Require().NoError(err)
	err = k.SaveClass(ctx, nft.Class{Id: classID2, Name: "Test Class 2", TotalSupply: 10})
	s.Require().NoError(err)

	owner := s.NewAccounts(1)[0]
	for i := 1; i <= 3; i++ {
		err = k.Mint(ctx, nft.NFT{ClassId: classID1, Id: fmt.Sprintf("nft-c1-%d", i)}, owner)
		s.Require().NoError(err)
		err = k.Mint(ctx, nft.NFT{ClassId: classID2, Id: fmt.Sprintf("nft-c2-%d", i)}, owner)
		s.Require().NoError(err)
	}

	resp, err := k.NftFilter(goCtx, &types.QueryNftFilterRequest{
		ClassId: classID1,
		Owner:   owner.String(),
		Pagination: &query.PageRequest{
			Offset: 1,
			Limit:  2,
		},
	})
	s.Require().NoError(err)
	s.Require().Len(resp.Nfts, 2)
	s.Require().Equal("nft-c1-2", resp.Nfts[0].TokenId)
	s.Require().Equal("nft-c1-3", resp.Nfts[1].TokenId)
	s.Require().NotNil(resp.Pagination)

	resp, err = k.NftFilter(goCtx, &types.QueryNftFilterRequest{
		Owner: owner.String(),
		Pagination: &query.PageRequest{
			Offset: 2,
			Limit:  3,
		},
	})
	s.Require().NoError(err)
	s.Require().Len(resp.Nfts, 3)
	s.Require().NotNil(resp.Pagination)

	_, err = k.NftFilter(goCtx, &types.QueryNftFilterRequest{
		ClassId: classID1,
		Owner:   "invalidaddr",
	})
	s.Require().Error(err)
}
