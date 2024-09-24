package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/wstaking/types"
	"strings"
)

func (s *KeeperTestSuite) TestNewFixedDepositCfg() {
	s.SetupTest()
	newFixdDepositCfg := &types.MsgNewFixedDepositCfg{
		Admin:    s.Dao.GlobalDao,
		RegionId: strings.ToLower(types.MeEarthRegionName),
		Term:     30,
		Rate:     sdk.NewDec(1),
	}
	_, err := s.msgServer.NewFixedDepositCfg(s.Ctx, newFixdDepositCfg)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestRemoveFixedDepositCfg() {
	s.SetupTest()
	newFixdDepositCfg := &types.MsgNewFixedDepositCfg{
		Admin:    s.Dao.GlobalDao,
		RegionId: strings.ToLower(types.MeEarthRegionName),
		Term:     30,
		Rate:     sdk.NewDec(1),
	}
	_, err := s.msgServer.NewFixedDepositCfg(s.Ctx, newFixdDepositCfg)
	s.Require().NoError(err)

	_, err = s.msgServer.RemoveFixedDepositCfg(s.Ctx,
		&types.MsgRemoveFixedDepositCfg{
			Admin:    s.Dao.GlobalDao,
			RegionId: strings.ToLower(types.MeEarthRegionName),
			Term:     30,
		})
	s.Require().NoError(err)

	_, err = s.queryClient.FixedDepositCfg(s.Ctx, &types.QueryFixedDepositCfgRequest{RegionId: strings.ToLower(types.MeEarthRegionName)})
	s.Require().ErrorIs(err, nil)
}

func (s *KeeperTestSuite) TestWithdrawFixedDepositCfg() {
	s.SetupTest()
	newFixdDepositCfg := &types.MsgNewFixedDepositCfg{
		Admin:    s.Dao.GlobalDao,
		RegionId: strings.ToLower(types.MeEarthRegionName),
		Term:     30,
		Rate:     sdk.NewDec(1),
	}
	_, err := s.msgServer.NewFixedDepositCfg(s.Ctx, newFixdDepositCfg)
	s.Require().NoError(err)

	s.msgServer.DoFixedDeposit(s.Ctx, &types.MsgDoFixedDeposit{
		Account:   s,
		Principal: sdk.Coin{},
		Term:      0,
	})

	_, err = s.queryClient.FixedDepositCfg(s.Ctx, &types.QueryFixedDepositCfgRequest{RegionId: strings.ToLower(types.MeEarthRegionName)})
	s.Require().ErrorIs(err, nil)

}
