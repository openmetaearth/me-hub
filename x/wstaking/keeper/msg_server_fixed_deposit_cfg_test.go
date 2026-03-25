package keeper_test

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/wstaking/types"
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
			name:     "invalid rate",
			creator:  s.Dao.GlobalDao,
			regionId: strings.ToLower(types.MeEarthRegionName),
			term:     1,
			rate:     sdk.MustNewDecFromStr("0"),
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
