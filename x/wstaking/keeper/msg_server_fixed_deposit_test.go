package keeper_test

import (
	"fmt"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/app/apptesting"
	"github.com/st-chain/me-hub/app/params"
	wmintTypes "github.com/st-chain/me-hub/x/wmint/types"
	"github.com/st-chain/me-hub/x/wstaking/types"
	"strings"
	"time"
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
		RegionId: types.MeEarthRegionId,
		Term:     30,
		Rate:     sdk.NewDec(1),
	}
	_, err := s.msgServer.NewFixedDepositCfg(s.Ctx, newFixdDepositCfg)
	s.Require().NoError(err)

	_, err = s.msgServer.RemoveFixedDepositCfg(s.Ctx,
		&types.MsgRemoveFixedDepositCfg{
			Admin:    s.Dao.GlobalDao,
			RegionId: types.MeEarthRegionId,
			Term:     30,
		})
	s.Require().NoError(err)

	_, err = s.queryClient.FixedDepositCfg(s.Ctx, &types.QueryFixedDepositCfgRequest{RegionId: types.MeEarthRegionId})
	s.Require().ErrorIs(err, nil)
}

func (s *KeeperTestSuite) TestWithdrawFixedDepositCfg() {
	s.SetupTest()
	newFixdDepositCfg := &types.MsgNewFixedDepositCfg{
		Admin:    s.Dao.GlobalDao,
		RegionId: types.MeEarthRegionId,
		Term:     30,
		Rate:     sdk.NewDec(10),
	}
	_, err := s.msgServer.NewFixedDepositCfg(s.Ctx, newFixdDepositCfg)
	s.Require().NoError(err)

	fixDeposit, err := s.msgServer.DoFixedDeposit(s.Ctx, &types.MsgDoFixedDeposit{
		Account: s.Dao.GlobalDao,
		Principal: sdk.Coin{
			Denom:  params.BaseDenom,
			Amount: sdk.NewInt(100000000),
		},
		Term: 30,
	})
	s.Require().NoError(err)

	principalAddr := s.App.AccountKeeper.GetModuleAddress(types.FixedDepositPrincipalPool)

	poolBalance := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(principalAddr.String()), params.BaseDenom)
	s.T().Log(fmt.Sprintf("poolBalance balance: %s", poolBalance.String()))

	daoBalance := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(s.Dao.GlobalDao), params.BaseDenom)
	s.T().Log(fmt.Sprintf("daoBalance balance: %s", daoBalance.String()))

	region, _ := s.App.StakingKeeper.GetRegion(s.Ctx, types.MeEarthRegionId)
	regionInterestAddr, _ := sdk.AccAddressFromBech32(region.DepositInterestAddr)
	interestBalance := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(regionInterestAddr.String()), params.BaseDenom)
	s.T().Log(fmt.Sprintf("interestBalance balance: %s", interestBalance.String()))

	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wmintTypes.OneYearTotalBlocks).WithChainID(apptesting.TestChainID).WithBlockTime(s.Ctx.BlockTime().Add(7760 * time.Hour))
	_, err = s.msgServer.DoFixedWithdraw(s.Ctx, &types.MsgDoFixedWithdraw{
		Account: s.Dao.GlobalDao,
		Id:      fixDeposit.Id,
	})

	poolBalanceAfer := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(principalAddr.String()), params.BaseDenom)
	s.T().Log(fmt.Sprintf("poolBalanceAfer balance: %s", poolBalanceAfer.String()))

	daobalanceAfer := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(s.Dao.GlobalDao), params.BaseDenom)
	s.T().Log(fmt.Sprintf("daobalanceAfer balance: %s", daobalanceAfer.String()))

	interestBalanceAfter := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(regionInterestAddr.String()), params.BaseDenom)
	s.T().Log(fmt.Sprintf("interestBalanceAfter balance: %s", interestBalanceAfter.String()))

	s.Require().Equal(daobalanceAfer.String(), daoBalance.Add(poolBalance).Add(interestBalance).String())
	s.Require().NoError(err)

}
