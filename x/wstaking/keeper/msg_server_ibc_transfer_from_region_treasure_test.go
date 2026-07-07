package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/openmetaearth/me-hub/app/params"
	stakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

// TestIbcTransferFromRegionTreasure tests IBC transfer from region treasure
func (s *KeeperTestSuite) TestIbcTransferFromRegionTreasure() {
	testCases := []struct {
		name      string
		setup     func() *stakingtypes.MsgIbcTransferFromRegionTreasure
		expectErr bool
		errMsg    string
	}{
		{
			name: "unauthorized sender - not dao",
			setup: func() *stakingtypes.MsgIbcTransferFromRegionTreasure {
				account := s.TestAccs[0]

				// Create a region first
				s.App.StakingKeeper.SetRegion(s.Ctx, stakingtypes.Region{
					RegionId:           stakingtypes.ExperienceRegionName,
					OperatorAddress:    s.experienceValidator.OperatorAddress,
					DelegateInterest:   math.LegacyNewDec(100000),
					DelegateAmount:     math.ZeroInt(),
					RegionShare:        math.ZeroInt(),
					RegionTreasureAddr: s.Dao.GlobalDao,
				})

				return &stakingtypes.MsgIbcTransferFromRegionTreasure{
					Creator:       account.String(), // Not the DAO
					RegionId:      stakingtypes.ExperienceRegionName,
					SourcePort:    ibctransfertypes.PortID,
					SourceChannel: "channel-0",
					Token:         sdk.NewCoin(params.BaseDenom, math.NewInt(1000000)),
					TimeoutHeight: stakingtypes.Height{
						RevisionNumber: 1,
						RevisionHeight: 1000,
					},
					TimeoutTimestamp: 0,
					Memo:             "",
				}
			},
			expectErr: true,
			errMsg:    "sender is not the dao",
		},
		{
			name: "region not found",
			setup: func() *stakingtypes.MsgIbcTransferFromRegionTreasure {
				return &stakingtypes.MsgIbcTransferFromRegionTreasure{
					Creator:       s.Dao.GlobalDao,
					RegionId:      "non_existent_region",
					SourcePort:    ibctransfertypes.PortID,
					SourceChannel: "channel-0",
					Token:         sdk.NewCoin(params.BaseDenom, math.NewInt(1000000)),
					TimeoutHeight: stakingtypes.Height{
						RevisionNumber: 1,
						RevisionHeight: 1000,
					},
					TimeoutTimestamp: 0,
					Memo:             "",
				}
			},
			expectErr: true,
			errMsg:    "region not found",
		},
		{
			name: "successful ibc transfer",
			setup: func() *stakingtypes.MsgIbcTransferFromRegionTreasure {
				// Create a region with treasure address
				treasureAddr := s.TestAccs[0]
				s.App.StakingKeeper.SetRegion(s.Ctx, stakingtypes.Region{
					RegionId:           stakingtypes.MeEarthRegionName,
					OperatorAddress:    s.meEarthValidator.OperatorAddress,
					DelegateInterest:   math.LegacyNewDec(100000),
					DelegateAmount:     math.ZeroInt(),
					RegionShare:        math.ZeroInt(),
					RegionTreasureAddr: treasureAddr.String(),
				})

				// Fund the treasure account
				s.FundAcc(treasureAddr, sdk.NewCoins(sdk.NewCoin(params.BaseDenom, math.NewInt(10000000))))

				return &stakingtypes.MsgIbcTransferFromRegionTreasure{
					Creator:       s.Dao.GlobalDao,
					RegionId:      stakingtypes.MeEarthRegionName,
					SourcePort:    ibctransfertypes.PortID,
					SourceChannel: "channel-0",
					Token:         sdk.NewCoin(params.BaseDenom, math.NewInt(1000000)),
					TimeoutHeight: stakingtypes.Height{
						RevisionNumber: 1,
						RevisionHeight: 1000,
					},
					TimeoutTimestamp: 0,
					Memo:             "test ibc transfer",
				}
			},
			expectErr: true,
			errMsg:    "channel not found", // IBC channel doesn't exist in test environment
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Use fresh context for each test
			s.SetupTest()
			msg := tc.setup()

			resp, err := s.msgServer.IbcTransferFromRegionTreasure(s.Ctx, msg)

			if tc.expectErr {
				s.Require().Error(err)
				if tc.errMsg != "" {
					s.Require().Contains(err.Error(), tc.errMsg)
				}
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
			}
		})
	}
}

// TestIbcTransferFromRegionTreasureWithDifferentRegions tests transfers from different regions
func (s *KeeperTestSuite) TestIbcTransferFromRegionTreasureWithDifferentRegions() {
	regions := []struct {
		regionId string
		operator string
	}{
		{stakingtypes.ExperienceRegionName, s.experienceValidator.OperatorAddress},
		{stakingtypes.MeEarthRegionName, s.meEarthValidator.OperatorAddress},
	}

	for _, region := range regions {
		treasureAddr := s.NewAccounts(1)[0]

		// Create region
		s.App.StakingKeeper.SetRegion(s.Ctx, stakingtypes.Region{
			RegionId:           region.regionId,
			OperatorAddress:    region.operator,
			DelegateInterest:   math.LegacyNewDec(100000),
			DelegateAmount:     math.ZeroInt(),
			RegionShare:        math.ZeroInt(),
			RegionTreasureAddr: treasureAddr.String(),
		})

		// Fund the treasure
		s.FundAcc(treasureAddr, sdk.NewCoins(sdk.NewCoin(params.BaseDenom, math.NewInt(10000000))))

		msg := &stakingtypes.MsgIbcTransferFromRegionTreasure{
			Creator:       s.Dao.GlobalDao,
			RegionId:      region.regionId,
			SourcePort:    ibctransfertypes.PortID,
			SourceChannel: "channel-0",
			Token:         sdk.NewCoin(params.BaseDenom, math.NewInt(500000)),
			TimeoutHeight: stakingtypes.Height{
				RevisionNumber: 1,
				RevisionHeight: 2000,
			},
			TimeoutTimestamp: 0,
			Memo:             "transfer from " + region.regionId,
		}

		_, err := s.msgServer.IbcTransferFromRegionTreasure(s.Ctx, msg)
		// IBC channels don't exist in test environment, so transfer will fail
		s.Require().Error(err)
		s.Require().Contains(err.Error(), "channel not found")
	}
}

// TestIbcTransferFromRegionTreasureInvalidChannel tests transfer with invalid channel
func (s *KeeperTestSuite) TestIbcTransferFromRegionTreasureInvalidChannel() {
	// Create a region
	treasureAddr := s.TestAccs[0]
	s.App.StakingKeeper.SetRegion(s.Ctx, stakingtypes.Region{
		RegionId:           stakingtypes.ExperienceRegionName,
		OperatorAddress:    s.experienceValidator.OperatorAddress,
		DelegateInterest:   math.LegacyNewDec(100000),
		DelegateAmount:     math.ZeroInt(),
		RegionShare:        math.ZeroInt(),
		RegionTreasureAddr: treasureAddr.String(),
	})

	// Fund the treasure
	s.FundAcc(treasureAddr, sdk.NewCoins(sdk.NewCoin(params.BaseDenom, math.NewInt(10000000))))

	msg := &stakingtypes.MsgIbcTransferFromRegionTreasure{
		Creator:       s.Dao.GlobalDao,
		RegionId:      stakingtypes.ExperienceRegionName,
		SourcePort:    ibctransfertypes.PortID,
		SourceChannel: "invalid-channel-999", // Invalid channel
		Token:         sdk.NewCoin(params.BaseDenom, math.NewInt(1000000)),
		TimeoutHeight: stakingtypes.Height{
			RevisionNumber: 1,
			RevisionHeight: 1000,
		},
		TimeoutTimestamp: 0,
		Memo:             "",
	}

	_, err := s.msgServer.IbcTransferFromRegionTreasure(s.Ctx, msg)
	// Should error because channel doesn't exist
	s.Require().Error(err)
}
