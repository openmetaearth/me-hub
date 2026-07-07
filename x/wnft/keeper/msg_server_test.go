package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/x/wnft/keeper"
	"github.com/openmetaearth/me-hub/x/wnft/types"
	"github.com/stretchr/testify/suite"
)

type MsgServerTestSuite struct {
	apptesting.KeeperTestHelper
	msgServer types.MsgServer
	TestAccs  []sdk.AccAddress
}

func TestMsgServerTestSuite(t *testing.T) {
	suite.Run(t, new(MsgServerTestSuite))
}

func (s *MsgServerTestSuite) SetupTest() {
	app := apptesting.Setup(s.T())
	ctx := app.GetBaseApp().NewContext(false)

	s.App = app
	s.Ctx = ctx

	s.msgServer = keeper.NewMsgServerImpl(app.WNFTKeeper, app.WNFTKeeper.Keeper)

	s.TestAccs = s.NewAccounts(5)
}

// TestNewClass tests NFT class creation
func (s *MsgServerTestSuite) TestNewClass() {
	testCases := []struct {
		name      string
		msg       *types.MsgNewClass
		expectErr bool
		errMsg    string
	}{
		{
			name: "valid class creation",
			msg: &types.MsgNewClass{
				ClassId:     "test-class-1",
				Sender:      s.TestAccs[0].String(),
				Name:        "Test Class",
				Symbol:      "TEST",
				Description: "Test NFT Class",
				Uri:         "ipfs://test",
				UriHash:     "hash",
				TotalSupply: 1000,
			},
			expectErr: false,
		},
		{
			name: "duplicate class id",
			msg: &types.MsgNewClass{
				ClassId:     "test-class-1", // same as above
				Sender:      s.TestAccs[0].String(),
				Name:        "Test Class 2",
				Symbol:      "TEST2",
				Description: "Duplicate Test",
				Uri:         "ipfs://test2",
				UriHash:     "hash2",
				TotalSupply: 500,
			},
			expectErr: true,
			errMsg:    "already exists",
		},
		{
			name: "kyc class with zero supply",
			msg: &types.MsgNewClass{
				ClassId:     "test-kyc-class",
				Sender:      s.TestAccs[1].String(),
				Name:        "KYC",
				Symbol:      "KYC",
				Description: "KYC NFT",
				Uri:         "ipfs://kyc",
				UriHash:     "kyc-hash",
				TotalSupply: 0, // kyc can have zero supply
			},
			expectErr: false,
		},
		{
			name: "invalid region name as class id",
			msg: &types.MsgNewClass{
				ClassId:     "USA-NFT-CLASS-ID", // reserved region name pattern
				Sender:      s.TestAccs[2].String(),
				Name:        "Invalid",
				Symbol:      "INV",
				Description: "Invalid class",
				Uri:         "ipfs://invalid",
				UriHash:     "invalid-hash",
				TotalSupply: 100,
			},
			expectErr: true,
			errMsg:    "invalid class name",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.msgServer.NewClass(s.Ctx, tc.msg)

			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)

				// Verify class was created
				class, ok := s.App.WNFTKeeper.GetClass(s.Ctx, tc.msg.ClassId)
				s.Require().True(ok)
				s.Require().Equal(tc.msg.Name, class.Name)
				s.Require().Equal(tc.msg.Symbol, class.Symbol)

				// Verify total supply cap was set
				cap := s.App.WNFTKeeper.GetClassTotalSupplyCap(s.Ctx, tc.msg.ClassId)
				s.Require().Equal(tc.msg.TotalSupply, cap)
			}
		})
	}
}

// TestMintNFT tests NFT minting functionality
func (s *MsgServerTestSuite) TestMintNFT() {
	// First create a class
	creator := s.TestAccs[0]
	classMsg := &types.MsgNewClass{
		ClassId:     "mint-test-class",
		Sender:      creator.String(),
		Name:        "Mint Test",
		Symbol:      "MINT",
		Description: "Test minting",
		Uri:         "ipfs://mint",
		UriHash:     "mint-hash",
		TotalSupply: 100,
	}
	_, err := s.msgServer.NewClass(s.Ctx, classMsg)
	s.Require().NoError(err)

	testCases := []struct {
		name      string
		msg       *types.MsgMintNFT
		expectErr bool
		errMsg    string
	}{
		{
			name: "valid mint",
			msg: &types.MsgMintNFT{
				ClassId:  "mint-test-class",
				TokenId:  "1",
				Creator:  creator.String(),
				Receiver: s.TestAccs[1].String(),
				Uri:      "ipfs://token1",
				UriHash:  "token1-hash",
			},
			expectErr: false,
		},
		{
			name: "non-existent class",
			msg: &types.MsgMintNFT{
				ClassId:  "non-existent",
				TokenId:  "1",
				Creator:  creator.String(),
				Receiver: s.TestAccs[1].String(),
				Uri:      "ipfs://token",
				UriHash:  "hash",
			},
			expectErr: true,
			errMsg:    "class", // Check for "class" in error message
		},
		{
			name: "unauthorized creator",
			msg: &types.MsgMintNFT{
				ClassId:  "mint-test-class",
				TokenId:  "2",
				Creator:  s.TestAccs[2].String(), // not the class creator
				Receiver: s.TestAccs[1].String(),
				Uri:      "ipfs://token2",
				UriHash:  "token2-hash",
			},
			expectErr: true,
			errMsg:    "not the creator",
		},
		{
			name: "invalid token id - zero",
			msg: &types.MsgMintNFT{
				ClassId:  "mint-test-class",
				TokenId:  "0", // invalid
				Creator:  creator.String(),
				Receiver: s.TestAccs[1].String(),
				Uri:      "ipfs://token0",
				UriHash:  "token0-hash",
			},
			expectErr: true,
			errMsg:    "invalid token id",
		},
		{
			name: "invalid token id - exceeds supply",
			msg: &types.MsgMintNFT{
				ClassId:  "mint-test-class",
				TokenId:  "101", // exceeds total supply of 100
				Creator:  creator.String(),
				Receiver: s.TestAccs[1].String(),
				Uri:      "ipfs://token101",
				UriHash:  "token101-hash",
			},
			expectErr: true,
			errMsg:    "invalid token id",
		},
		{
			name: "valid mint at max supply",
			msg: &types.MsgMintNFT{
				ClassId:  "mint-test-class",
				TokenId:  "100", // exactly at max
				Creator:  creator.String(),
				Receiver: s.TestAccs[1].String(),
				Uri:      "ipfs://token100",
				UriHash:  "token100-hash",
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.msgServer.MintNFT(s.Ctx, tc.msg)

			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)

				// Verify NFT was minted
				nftItem, ok := s.App.WNFTKeeper.GetNFT(s.Ctx, tc.msg.ClassId, tc.msg.TokenId)
				s.Require().True(ok)
				s.Require().Equal(tc.msg.ClassId, nftItem.ClassId)
				s.Require().Equal(tc.msg.TokenId, nftItem.Id)

				// Verify owner
				owner := s.App.WNFTKeeper.GetOwner(s.Ctx, tc.msg.ClassId, tc.msg.TokenId)
				s.Require().Equal(tc.msg.Receiver, owner.String())
			}
		})
	}
}

// TestSend tests NFT transfer functionality
func (s *MsgServerTestSuite) TestSend() {
	// Setup: create class and mint an NFT
	creator := s.TestAccs[0]
	owner := s.TestAccs[1]
	receiver := s.TestAccs[2]

	classMsg := &types.MsgNewClass{
		ClassId:     "send-test-class",
		Sender:      creator.String(),
		Name:        "Send Test",
		Symbol:      "SEND",
		Description: "Test sending",
		Uri:         "ipfs://send",
		UriHash:     "send-hash",
		TotalSupply: 50,
	}
	_, err := s.msgServer.NewClass(s.Ctx, classMsg)
	s.Require().NoError(err)

	mintMsg := &types.MsgMintNFT{
		ClassId:  "send-test-class",
		TokenId:  "1",
		Creator:  creator.String(),
		Receiver: owner.String(),
		Uri:      "ipfs://nft1",
		UriHash:  "nft1-hash",
	}
	_, err = s.msgServer.MintNFT(s.Ctx, mintMsg)
	s.Require().NoError(err)

	testCases := []struct {
		name      string
		msg       *types.MsgSend
		expectErr bool
		errMsg    string
	}{
		{
			name: "valid send",
			msg: &types.MsgSend{
				Sender:   owner.String(),
				Receiver: receiver.String(),
				ClassId:  "send-test-class",
				Id:       "1",
			},
			expectErr: false,
		},
		{
			name: "send non-owned NFT",
			msg: &types.MsgSend{
				Sender:   owner.String(), // no longer owner after previous transfer
				Receiver: s.TestAccs[3].String(),
				ClassId:  "send-test-class",
				Id:       "1",
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.msgServer.Send(s.Ctx, tc.msg)

			if tc.expectErr {
				s.Require().Error(err)
				if tc.errMsg != "" {
					s.Require().Contains(err.Error(), tc.errMsg)
				}
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)

				// Verify new owner
				newOwner := s.App.WNFTKeeper.GetOwner(s.Ctx, tc.msg.ClassId, tc.msg.Id)
				s.Require().Equal(tc.msg.Receiver, newOwner.String())
			}
		})
	}
}

// TestMintNFTEvents tests that events are emitted correctly
func (s *MsgServerTestSuite) TestMintNFTEvents() {
	creator := s.TestAccs[0]

	// Create class
	classMsg := &types.MsgNewClass{
		ClassId:     "event-test-class",
		Sender:      creator.String(),
		Name:        "Event Test",
		Symbol:      "EVENT",
		Description: "Test events",
		Uri:         "ipfs://event",
		UriHash:     "event-hash",
		TotalSupply: 10,
	}
	_, err := s.msgServer.NewClass(s.Ctx, classMsg)
	s.Require().NoError(err)

	// Mint NFT
	mintMsg := &types.MsgMintNFT{
		ClassId:  "event-test-class",
		TokenId:  "1",
		Creator:  creator.String(),
		Receiver: s.TestAccs[1].String(),
		Uri:      "ipfs://event-nft",
		UriHash:  "event-nft-hash",
	}
	_, err = s.msgServer.MintNFT(s.Ctx, mintMsg)
	s.Require().NoError(err)

	// Check events
	events := s.Ctx.EventManager().Events()
	s.Require().NotEmpty(events)

	// Find MintNFT event
	var foundMintEvent bool
	for _, event := range events {
		if event.Type == types.EventTypeMintNFT {
			foundMintEvent = true
			// Verify event attributes
			attrs := event.Attributes
			s.Require().NotEmpty(attrs)
		}
	}
	s.Require().True(foundMintEvent, "MintNFT event should be emitted")
}
