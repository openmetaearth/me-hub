package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"strings"
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
		rate     sdk.Dec
		expErr   error
	}{
		{
			name:     "Dao Permission",
			creator:  s.Dao.MeidDao,
			regionId: strings.ToLower(types.MeEarthRegionName),
			term:     1,
			rate:     sdk.MustNewDecFromStr("0.1"),
			expErr:   types.ErrCheckGlobalDao,
		}, {
			name:     "have permission, but wrong region id",
			creator:  s.Dao.GlobalDao,
			regionId: types.MeEarthRegionName,
			term:     1,
			rate:     sdk.MustNewDecFromStr("0.1"),
			expErr:   types.ErrRegionName,
		}, {
			name:     "invalid term",
			creator:  s.Dao.GlobalDao,
			regionId: strings.ToLower(types.MeEarthRegionName),
			term:     0,
			rate:     sdk.MustNewDecFromStr("0.1"),
			expErr:   types.ErrAddFixedDepositConfig,
		}, {
			name:     "invalid rate (zero)",
			creator:  s.Dao.GlobalDao,
			regionId: strings.ToLower(types.MeEarthRegionName),
			term:     1,
			rate:     sdk.MustNewDecFromStr("0"),
			expErr:   types.ErrAddFixedDepositConfig,
		}, {
			name:     "invalid rate (negative)",
			creator:  s.Dao.GlobalDao,
			regionId: strings.ToLower(types.MeEarthRegionName),
			term:     1,
			rate:     sdk.MustNewDecFromStr("-0.1"),
			expErr:   types.ErrAddFixedDepositConfig,
		}, {
			name:     "invalid rate (too small)",
			creator:  s.Dao.GlobalDao,
			regionId: strings.ToLower(types.MeEarthRegionName),
			term:     1,
			rate:     sdk.MustNewDecFromStr("0.00009"),
			expErr:   types.ErrAddFixedDepositConfig,
		}, {
			name:     "invalid rate (too large)",
			creator:  s.Dao.GlobalDao,
			regionId: strings.ToLower(types.MeEarthRegionName),
			term:     1,
			rate:     sdk.MustNewDecFromStr("10000.0001"),
			expErr:   types.ErrAddFixedDepositConfig,
		}, {
			name:     "No error",
			creator:  s.Dao.GlobalDao,
			regionId: strings.ToLower(types.MeEarthRegionName),
			term:     1,
			rate:     sdk.MustNewDecFromStr("0.1"),
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
				s.Require().True(cfg.RegionFixedDepositCfgs[0].RegionFixedDepositCfg[0].Rate.Equal(sdk.MustNewDecFromStr("0.1")))
			}
		})
	}
}

func (s *KeeperTestSuite) TestSetFixedDepositCfgRateRejectsInvalidRates() {
	s.SetupTest()

	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	newCfg := types.MsgNewFixedDepositCfg{
		Dao:      s.Dao.GlobalDao,
		RegionId: strings.ToLower(types.MeEarthRegionName),
		Term:     1,
		Rate:     sdk.MustNewDecFromStr("0.1"),
	}
	_, err = s.msgServer.NewFixedDepositCfg(s.Ctx, &newCfg)
	s.Require().NoError(err)

	testCases := []struct {
		name string
		rate sdk.Dec
	}{
		{
			name: "negative rate",
			rate: sdk.MustNewDecFromStr("-0.1"),
		},
		{
			name: "zero rate",
			rate: sdk.ZeroDec(),
		},
		{
			name: "too small rate",
			rate: sdk.MustNewDecFromStr("0.00009"),
		},
		{
			name: "too large rate",
			rate: sdk.MustNewDecFromStr("10000.0001"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := s.msgServer.SetFixedDepositCfgRate(s.Ctx, &types.MsgSetFixedDepositCfgRate{
				Admin:    s.Dao.GlobalDao,
				RegionId: strings.ToLower(types.MeEarthRegionName),
				Term:     1,
				Rate:     tc.rate,
			})
			s.Require().Error(err)
			s.Require().ErrorIs(err, types.ErrSetFixedDepositConfigRate)
		})
	}
}
