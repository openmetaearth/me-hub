package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types" 
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/x/wnft/keeper"
	"github.com/openmetaearth/me-hub/x/wnft/types"
	"github.com/stretchr/testify/suite"
)

type GRPCQueryTestSuite struct {
	apptesting.KeeperTestHelper
	keeper    *keeper.Keeper
	msgServer types.MsgServer
	TestAccs  []sdk.AccAddress
}

func TestGRPCQueryTestSuite(t *testing.T) {
	suite.Run(t, new(GRPCQueryTestSuite))
}

func (s *GRPCQueryTestSuite) SetupTest() {
	app := apptesting.Setup(s.T())
	ctx := app.GetBaseApp().NewContext(false)

	s.App = app
	s.Ctx = ctx
	s.keeper = app.WNFTKeeper

	s.msgServer = keeper.NewMsgServerImpl(app.WNFTKeeper, app.WNFTKeeper.Keeper)

	s.TestAccs = s.NewAccounts(5)

	// Setup test data
	s.setupTestData()
}

func (s *GRPCQueryTestSuite) setupTestData() {
	creator := s.TestAccs[0]
	owner1 := s.TestAccs[1]
	owner2 := s.TestAccs[2]

	// Create class 1
	classMsg1 := &types.MsgNewClass{
		ClassId:     "test-class-1",
		Sender:      creator.String(),
		Name:        "Test Class 1",
		Symbol:      "TC1",
		Description: "Test class 1",
		Uri:         "ipfs://class1",
		UriHash:     "class1-hash",
		TotalSupply: 100,
	}
	_, err := s.msgServer.NewClass(s.Ctx, classMsg1)
	s.Require().NoError(err)

	// Create class 2
	classMsg2 := &types.MsgNewClass{
		ClassId:     "test-class-2",
		Sender:      creator.String(),
		Name:        "Test Class 2",
		Symbol:      "TC2",
		Description: "Test class 2",
		Uri:         "ipfs://class2",
		UriHash:     "class2-hash",
		TotalSupply: 50,
	}
	_, err = s.msgServer.NewClass(s.Ctx, classMsg2)
	s.Require().NoError(err)

	// Mint NFTs for owner1 in class 1
	for i := 1; i <= 3; i++ {
		mintMsg := &types.MsgMintNFT{
			ClassId:  "test-class-1",
			TokenId:  string(rune('0' + i)),
			Creator:  creator.String(),
			Receiver: owner1.String(),
			Uri:      "ipfs://nft" + string(rune('0'+i)),
			UriHash:  "nft-hash-" + string(rune('0'+i)),
		}
		_, err = s.msgServer.MintNFT(s.Ctx, mintMsg)
		s.Require().NoError(err)
	}

	// Mint NFTs for owner2 in class 1
	for i := 4; i <= 5; i++ {
		mintMsg := &types.MsgMintNFT{
			ClassId:  "test-class-1",
			TokenId:  string(rune('0' + i)),
			Creator:  creator.String(),
			Receiver: owner2.String(),
			Uri:      "ipfs://nft" + string(rune('0'+i)),
			UriHash:  "nft-hash-" + string(rune('0'+i)),
		}
		_, err = s.msgServer.MintNFT(s.Ctx, mintMsg)
		s.Require().NoError(err)
	}

	// Mint NFTs for owner1 in class 2
	mintMsg := &types.MsgMintNFT{
		ClassId:  "test-class-2",
		TokenId:  "1",
		Creator:  creator.String(),
		Receiver: owner1.String(),
		Uri:      "ipfs://class2-nft1",
		UriHash:  "class2-nft1-hash",
	}
	_, err = s.msgServer.MintNFT(s.Ctx, mintMsg)
	s.Require().NoError(err)
}

// TestClassAddress tests the ClassAddress query
func (s *GRPCQueryTestSuite) TestClassAddress() {
	owner1 := s.TestAccs[1]
	owner2 := s.TestAccs[2]

	testCases := []struct {
		name           string
		req            *types.QueryClassAddressRequest
		expectErr      bool
		expectedExists bool
		expectedCount  int
	}{
		{
			name: "valid query - owner1 class1",
			req: &types.QueryClassAddressRequest{
				ClassId: "test-class-1",
				Address: owner1.String(),
			},
			expectErr:      false,
			expectedExists: true,
			expectedCount:  3, // owner1 has 3 NFTs in class1
		},
		{
			name: "valid query - owner2 class1",
			req: &types.QueryClassAddressRequest{
				ClassId: "test-class-1",
				Address: owner2.String(),
			},
			expectErr:      false,
			expectedExists: true,
			expectedCount:  2, // owner2 has 2 NFTs in class1
		},
		{
			name: "valid query - owner1 class2",
			req: &types.QueryClassAddressRequest{
				ClassId: "test-class-2",
				Address: owner1.String(),
			},
			expectErr:      false,
			expectedExists: true,
			expectedCount:  1, // owner1 has 1 NFT in class2
		},
		{
			name: "non-existent class",
			req: &types.QueryClassAddressRequest{
				ClassId: "non-existent",
				Address: owner1.String(),
			},
			expectErr:      false,
			expectedExists: false,
			expectedCount:  0,
		},
		{
			name: "owner with no NFTs",
			req: &types.QueryClassAddressRequest{
				ClassId: "test-class-1",
				Address: s.TestAccs[3].String(),
			},
			expectErr:      false,
			expectedExists: true,
			expectedCount:  0,
		},
		{
			name:      "nil request",
			req:       nil,
			expectErr: true,
		},
		{
			name: "invalid address",
			req: &types.QueryClassAddressRequest{
				ClassId: "test-class-1",
				Address: "invalid-address",
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.keeper.ClassAddress(s.Ctx, tc.req)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(tc.expectedExists, resp.Exists)
				if tc.expectedExists {
					s.Require().Len(resp.Nfts, tc.expectedCount)
					if tc.expectedCount > 0 {
						s.Require().NotZero(resp.TotalSupply)
					}
				}
			}
		})
	}
}

// TestNftFilter tests the NftFilter query with different filter combinations
func (s *GRPCQueryTestSuite) TestNftFilter() {
	owner1 := s.TestAccs[1]

	testCases := []struct {
		name          string
		req           *types.QueryNftFilterRequest
		expectErr     bool
		expectedCount int
	}{
		{
			name: "query specific NFT - classId + tokenId",
			req: &types.QueryNftFilterRequest{
				ClassId: "test-class-1",
				TokenId: "1",
				Owner:   "",
			},
			expectErr:     false,
			expectedCount: 1,
		},
		{
			name: "query owner's NFTs in specific class - classId + owner",
			req: &types.QueryNftFilterRequest{
				ClassId: "test-class-1",
				TokenId: "",
				Owner:   owner1.String(),
			},
			expectErr:     false,
			expectedCount: 3, // owner1 has 3 NFTs in class1
		},
		{
			name: "query all NFTs owned by address - owner only",
			req: &types.QueryNftFilterRequest{
				ClassId: "",
				TokenId: "",
				Owner:   owner1.String(),
			},
			expectErr:     false,
			expectedCount: 4, // owner1 has 3 in class1 + 1 in class2
		},
		{
			name: "non-existent NFT",
			req: &types.QueryNftFilterRequest{
				ClassId: "test-class-1",
				TokenId: "999",
				Owner:   "",
			},
			expectErr:     false,
			expectedCount: 0,
		},
		{
			name: "non-existent class with owner",
			req: &types.QueryNftFilterRequest{
				ClassId: "non-existent",
				TokenId: "",
				Owner:   owner1.String(),
			},
			expectErr:     false,
			expectedCount: 0,
		},
		{
			name: "owner with no NFTs",
			req: &types.QueryNftFilterRequest{
				ClassId: "",
				TokenId: "",
				Owner:   s.TestAccs[4].String(),
			},
			expectErr:     false,
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.keeper.NftFilter(s.Ctx, tc.req)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				if tc.expectedCount > 0 {
					s.Require().NotNil(resp)
					s.Require().Len(resp.Nfts, tc.expectedCount)

					// Verify NFT list structure
					for _, nft := range resp.Nfts {
						s.Require().NotEmpty(nft.ClassId)
						s.Require().NotEmpty(nft.TokenId)
						s.Require().NotEmpty(nft.Owner)
						s.Require().NotEmpty(nft.Uri)
					}
				}
			}
		})
	}
}

// TestNftFilterSpecificNFT tests querying a specific NFT by classId and tokenId
func (s *GRPCQueryTestSuite) TestNftFilterSpecificNFT() {
	owner1 := s.TestAccs[1]

	req := &types.QueryNftFilterRequest{
		ClassId: "test-class-1",
		TokenId: "1",
		Owner:   "",
	}

	resp, err := s.keeper.NftFilter(s.Ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Len(resp.Nfts, 1)

	nft := resp.Nfts[0]
	s.Require().Equal("test-class-1", nft.ClassId)
	s.Require().Equal("1", nft.TokenId)
	s.Require().Equal(owner1.String(), nft.Owner)
	s.Require().NotEmpty(nft.Uri)
}

// TestNftFilterOwnerAllClasses tests querying all NFTs owned by an address across all classes
func (s *GRPCQueryTestSuite) TestNftFilterOwnerAllClasses() {
	owner1 := s.TestAccs[1]

	req := &types.QueryNftFilterRequest{
		ClassId: "",
		TokenId: "",
		Owner:   owner1.String(),
	}

	resp, err := s.keeper.NftFilter(s.Ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Len(resp.Nfts, 4) // 3 from class1 + 1 from class2

	// Verify we have NFTs from both classes
	classIds := make(map[string]int)
	for _, nft := range resp.Nfts {
		classIds[nft.ClassId]++
		s.Require().Equal(owner1.String(), nft.Owner)
	}

	s.Require().Equal(3, classIds["test-class-1"])
	s.Require().Equal(1, classIds["test-class-2"])
}

// TestQueryWithEmptyDatabase tests queries when no data exists
func (s *GRPCQueryTestSuite) TestQueryWithEmptyDatabase() {
	// Create a fresh context with no data
	freshApp := apptesting.Setup(s.T())
	freshCtx := freshApp.GetBaseApp().NewContext(false)

	keeper := freshApp.WNFTKeeper

	// Test ClassAddress with empty database
	req1 := &types.QueryClassAddressRequest{
		ClassId: "non-existent",
		Address: s.TestAccs[0].String(),
	}
	resp1, err := keeper.ClassAddress(freshCtx, req1)
	s.Require().NoError(err)
	s.Require().NotNil(resp1)
	s.Require().False(resp1.Exists)

	// Test NftFilter with empty database
	req2 := &types.QueryNftFilterRequest{
		ClassId: "",
		TokenId: "",
		Owner:   s.TestAccs[0].String(),
	}
	resp2, err := keeper.NftFilter(freshCtx, req2)
	s.Require().NoError(err)
	// Should return empty result, not error
	if resp2 != nil {
		s.Require().Empty(resp2.Nfts)
	}
}
