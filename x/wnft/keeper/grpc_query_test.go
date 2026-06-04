package keeper_test

import (
	"fmt"
	"testing"

	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	// Setup a class
	classId := "testclass"
	err := k.SaveClass(ctx, nft.Class{
		Id:          classId,
		Name:        "Test Class",
		Symbol:      "TC",
		TotalSupply: 10,
	})
	s.Require().NoError(err)

	// Create a test owner address
	owner := s.NewAccounts(1)[0]

	// Mint 5 NFTs
	for i := 1; i <= 5; i++ {
		err = k.Mint(ctx, nft.NFT{
			ClassId: classId,
			Id:      fmt.Sprintf("nft%d", i),
		}, owner)
		s.Require().NoError(err)
	}

	// 1. Query without offset/limit (defaults to limit 1000)
	resp, err := k.ClassAddress(goCtx, &types.QueryClassAddressRequest{
		ClassId: classId,
		Address: owner.String(),
	})
	s.Require().NoError(err)
	s.Require().True(resp.Exists)
	s.Require().Len(resp.Nfts, 5)
	s.Require().Equal([]string{"nft1", "nft2", "nft3", "nft4", "nft5"}, resp.Nfts)

	// 2. Query with offset=1, limit=2
	resp, err = k.ClassAddress(goCtx, &types.QueryClassAddressRequest{
		ClassId: classId,
		Address: owner.String(),
		Offset:  1,
		Limit:   2,
	})
	s.Require().NoError(err)
	s.Require().Len(resp.Nfts, 2)
	s.Require().Equal([]string{"nft2", "nft3"}, resp.Nfts)

	// 3. Query with offset=4, limit=2 (should get last 1 NFT)
	resp, err = k.ClassAddress(goCtx, &types.QueryClassAddressRequest{
		ClassId: classId,
		Address: owner.String(),
		Offset:  4,
		Limit:   2,
	})
	s.Require().NoError(err)
	s.Require().Len(resp.Nfts, 1)
	s.Require().Equal([]string{"nft5"}, resp.Nfts)

	// 4. Query with invalid owner address
	resp, err = k.ClassAddress(goCtx, &types.QueryClassAddressRequest{
		ClassId: classId,
		Address: "invalidaddr",
	})
	s.Require().Error(err)
}

func (s *WnftKeeperTestSuite) TestNftFilterPagination() {
	k := s.App.WNFTKeeper
	ctx := s.Ctx
	goCtx := sdk.WrapSDKContext(ctx)

	classId1 := "testclass1"
	classId2 := "testclass2"

	err := k.SaveClass(ctx, nft.Class{
		Id:          classId1,
		Name:        "Test Class 1",
		TotalSupply: 10,
	})
	s.Require().NoError(err)

	err = k.SaveClass(ctx, nft.Class{
		Id:          classId2,
		Name:        "Test Class 2",
		TotalSupply: 10,
	})
	s.Require().NoError(err)

	owner := s.NewAccounts(1)[0]

	// Mint 3 NFTs in class1 and 3 in class2 for owner
	for i := 1; i <= 3; i++ {
		err = k.Mint(ctx, nft.NFT{ClassId: classId1, Id: fmt.Sprintf("nft-c1-%d", i)}, owner)
		s.Require().NoError(err)

		err = k.Mint(ctx, nft.NFT{ClassId: classId2, Id: fmt.Sprintf("nft-c2-%d", i)}, owner)
		s.Require().NoError(err)
	}

	// 1. Test NftFilter for specific class & owner
	resp, err := k.NftFilter(goCtx, &types.QueryNftFilterRequest{
		ClassId: classId1,
		Owner:   owner.String(),
		Offset:  1,
		Limit:   2,
	})
	s.Require().NoError(err)
	s.Require().Len(resp.Nfts, 2)
	s.Require().Equal("nft-c1-2", resp.Nfts[0].TokenId)
	s.Require().Equal("nft-c1-3", resp.Nfts[1].TokenId)

	// 2. Test NftFilter for owner across all classes (global pagination)
	resp, err = k.NftFilter(goCtx, &types.QueryNftFilterRequest{
		Owner:  owner.String(),
		Offset: 2,
		Limit:  3,
	})
	s.Require().NoError(err)
	s.Require().Len(resp.Nfts, 3)

	// 3. Test invalid owner address
	_, err = k.NftFilter(goCtx, &types.QueryNftFilterRequest{
		ClassId: classId1,
		Owner:   "invalidaddr",
	})
	s.Require().Error(err)
}
