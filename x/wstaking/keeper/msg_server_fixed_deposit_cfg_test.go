package keeper_test

import (
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestNewFixedDepositCfg() {
	s.SetupTest()

	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	tests := []struct {
		name     string
		creator  string
		regionId string
		term     int64
		rate     sdkmath.LegacyDec
		expErr   error
	}{
		{
			name:     "Dao Permission",
			creator:  s.Dao.MeidDao,
			regionId: strings.ToLower(types.MeEarthRegionName),
			term:     1,
			rate:     sdkmath.LegacyMustNewDecFromStr("0.1"),
			expErr:   types.ErrCheckGlobalDao,
		}, {
			name:     "have permission, but wrong region id",
			creator:  s.Dao.GlobalDao,
			regionId: types.MeEarthRegionName,
			term:     1,
			rate:     sdkmath.LegacyMustNewDecFromStr("0.1"),
			expErr:   types.ErrRegionName,
		}, {
			name:     "invalid term",
			creator:  s.Dao.GlobalDao,
			regionId: strings.ToLower(types.MeEarthRegionName),
			term:     0,
			rate:     sdkmath.LegacyMustNewDecFromStr("0.1"),
			expErr:   types.ErrAddFixedDepositConfig,
		}, {
			name:     "invalid rate",
			creator:  s.Dao.GlobalDao,
			regionId: strings.ToLower(types.MeEarthRegionName),
			term:     1,
			rate:     sdkmath.LegacyMustNewDecFromStr("0"),
			expErr:   types.ErrAddFixedDepositConfig,
		}, {
			name:     "No error",
			creator:  s.Dao.GlobalDao,
			regionId: strings.ToLower(types.MeEarthRegionName),
			term:     1,
			rate:     sdkmath.LegacyMustNewDecFromStr("0.1"),
			expErr:   nil,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			msg := types.MsgNewFixedDepositCfg{
				Dao:      test.creator,
				RegionId: test.regionId,
				Term:     test.term,
				Rate:     test.rate,
			}
			_, err := s.msgServer.NewFixedDepositCfg(s.Ctx, &msg)
			s.Require().ErrorIs(err, test.expErr)

			if test.expErr == nil {
				cfg, err := s.queryClient.FixedDepositCfg(s.Ctx, &types.QueryFixedDepositCfgRequest{RegionIds: []string{strings.ToLower(types.MeEarthRegionName)}})
				s.Require().NoError(err)
				s.Require().Equal(1, len(cfg.RegionFixedDepositCfgs))
				s.Require().Equal(strings.ToLower(types.MeEarthRegionName), cfg.RegionFixedDepositCfgs[0].RegionId)
				s.Require().Equal(int64(1), cfg.RegionFixedDepositCfgs[0].RegionFixedDepositCfg[0].Term)
				s.Require().True(cfg.RegionFixedDepositCfgs[0].RegionFixedDepositCfg[0].Rate.Equal(sdkmath.LegacyMustNewDecFromStr("0.1")))
			}
		})
	}
}

// TestSetFixedDepositCfgStatus tests setting fixed deposit config status  
func (s *KeeperTestSuite) TestSetFixedDepositCfgStatus() {
	s.SetupTest()

	// Create a region first
	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	// Create a fixed deposit config
	newCfg := types.MsgNewFixedDepositCfg{
		Dao:      s.Dao.GlobalDao,
		RegionId: strings.ToLower(types.MeEarthRegionName),
		Term:     1,
		Rate:     sdkmath.LegacyMustNewDecFromStr("0.1"),
	}
	_, err = s.msgServer.NewFixedDepositCfg(s.Ctx, &newCfg)
	s.Require().NoError(err)

	tests := []struct {
		name     string
		admin    string
		regionId string
		term     int64
		status   types.FIXED_DEPOSIT_CFG_STATUS
		expErr   error
	}{
		{
			name:     "successful status update",
			admin:    s.Dao.GlobalDao,
			regionId: strings.ToLower(types.MeEarthRegionName),
			term:     1,
			status:   types.RegionFixedDepositCfgStatusInactive,
			expErr:   nil,
		},
		{
			name:     "non-dao admin",
			admin:    s.TestAccs[0].String(),
			regionId: strings.ToLower(types.MeEarthRegionName),
			term:     1,
			status:   types.RegionFixedDepositCfgStatusInactive,
			expErr:   types.ErrCheckGlobalDao,
		},
		{
			name:     "config not found",
			admin:    s.Dao.GlobalDao,
			regionId: strings.ToLower(types.MeEarthRegionName),
			term:     999,
			status:   types.RegionFixedDepositCfgStatusInactive,
			expErr:   types.ErrSetFixedDepositConfigStatus,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			msg := types.MsgSetFixedDepositCfgStatus{
				Admin:    test.admin,
				RegionId: test.regionId,
				Term:     test.term,
				Status:   test.status,
			}
			resp, err := s.msgServer.SetFixedDepositCfgStatus(s.Ctx, &msg)

			if test.expErr != nil {
				s.Require().ErrorIs(err, test.expErr)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)

				cfg, found := s.App.StakingKeeper.GetFixedDepositCfg(s.Ctx, test.regionId, test.term)
				s.Require().True(found)
				s.Require().Equal(test.status, cfg.Status)
			}
		})
	}
}

// TestSetFixedDepositCfgRate tests setting fixed deposit config rate
func (s *KeeperTestSuite) TestSetFixedDepositCfgRate() {
	s.SetupTest()

	// Create a region first
	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	// Create a fixed deposit config
	initialRate := sdkmath.LegacyMustNewDecFromStr("0.1")
	newCfg := types.MsgNewFixedDepositCfg{
		Dao:      s.Dao.GlobalDao,
		RegionId: strings.ToLower(types.MeEarthRegionName),
		Term:     1,
		Rate:     initialRate,
	}
	_, err = s.msgServer.NewFixedDepositCfg(s.Ctx, &newCfg)
	s.Require().NoError(err)

	tests := []struct {
		name     string
		admin    string
		regionId string
		term     int64
		rate     sdkmath.LegacyDec
		expErr   error
	}{
		{
			name:     "successful rate update",
			admin:    s.Dao.GlobalDao,
			regionId: strings.ToLower(types.MeEarthRegionName),
			term:     1,
			rate:     sdkmath.LegacyMustNewDecFromStr("0.15"),
			expErr:   nil,
		},
		{
			name:     "non-dao admin",
			admin:    s.TestAccs[0].String(),
			regionId: strings.ToLower(types.MeEarthRegionName),
			term:     1,
			rate:     sdkmath.LegacyMustNewDecFromStr("0.2"),
			expErr:   types.ErrCheckGlobalDao,
		},
		{
			name:     "config not found",
			admin:    s.Dao.GlobalDao,
			regionId: strings.ToLower(types.MeEarthRegionName),
			term:     999,
			rate:     sdkmath.LegacyMustNewDecFromStr("0.2"),
			expErr:   types.ErrSetFixedDepositConfigRate,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			msg := types.MsgSetFixedDepositCfgRate{
				Admin:    test.admin,
				RegionId: test.regionId,
				Term:     test.term,
				Rate:     test.rate,
			}
			resp, err := s.msgServer.SetFixedDepositCfgRate(s.Ctx, &msg)

			if test.expErr != nil {
				s.Require().ErrorIs(err, test.expErr)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)

				cfg, found := s.App.StakingKeeper.GetFixedDepositCfg(s.Ctx, test.regionId, test.term)
				s.Require().True(found)
				s.Require().True(cfg.Rate.Equal(test.rate))
			}
		})
	}
}
