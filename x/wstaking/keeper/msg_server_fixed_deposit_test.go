package keeper_test

import (
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wdistri"
	"github.com/openmetaearth/me-hub/x/wmint"
	wmintTypes "github.com/openmetaearth/me-hub/x/wmint/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"strings"
	"time"
)

func (s *KeeperTestSuite) TestFixedDeposit() {
	s.SetupTest()

	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	msg := types.MsgNewFixedDepositCfg{
		Dao:      s.Dao.GlobalDao,
		RegionId: strings.ToLower(types.MeEarthRegionName),
		Term:     1,
		Rate:     sdk.MustNewDecFromStr("0.1"),
	}
	_, err = s.msgServer.NewFixedDepositCfg(s.Ctx, &msg)
	s.Require().NoError(err)

	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}, *s.App.DistrKeeper)

	amount := sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(10000000)))
	_, err = s.msgServer.WithdrawFromRegion(s.Ctx, &types.MsgWithdrawFromRegion{
		Withdrawer: s.Dao.GlobalDao,
		RegionId:   strings.ToLower(types.MeEarthRegionName),
		Receiver:   s.Dao.GlobalDao,
		Amount:     amount,
	})
	s.Require().NoError(err)

	tests := []struct {
		name      string
		account   string
		principal sdk.Coin
		term      int64
		expErr    error
	}{
		{
			name:      "invalid term",
			account:   s.Dao.GlobalDao,
			term:      0,
			principal: sdk.NewCoin(params.BaseDenom, sdk.NewInt(1)),
			expErr:    types.ErrDoFixedDeposit,
		}, {
			name:      "invalid principal",
			account:   s.Dao.GlobalDao,
			term:      1,
			principal: sdk.NewCoin(params.BaseDenom, sdk.NewInt(0)),
			expErr:    types.ErrDoFixedDeposit,
		}, {
			name:      "invalid kyc and regionId",
			account:   s.Dao.MeidDao,
			term:      1,
			principal: sdk.NewCoin(params.BaseDenom, sdk.NewInt(1)),
			expErr:    types.ErrDidNotExists,
		}, {
			name:      "insufficient principal",
			account:   s.Dao.GlobalDao,
			term:      1,
			principal: sdk.NewCoin(params.BaseDenom, sdk.NewInt(1)),
			expErr:    types.ErrDoFixedDeposit,
		}, {
			name:      "No error",
			account:   s.Dao.GlobalDao,
			term:      1,
			principal: amount[0],
			expErr:    nil,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			msg := types.MsgDoFixedDeposit{
				Account:   test.account,
				Principal: test.principal,
				Term:      test.term,
			}
			_, err := s.msgServer.DoFixedDeposit(s.Ctx, &msg)
			s.Require().ErrorIs(err, test.expErr)

			// check nft class
			if test.expErr == nil {
				deposit, err := s.queryClient.FixedDeposit(s.Ctx, &types.QueryGetFixedDepositRequest{
					Address: s.Dao.GlobalDao,
					Id:      0,
				})
				s.T().Log(deposit.FixedDeposit.String())
				s.Require().NoError(err)
				s.Require().Equal(int64(1), deposit.FixedDeposit.Term)
				s.Require().True(deposit.FixedDeposit.Principal.Equal(amount[0]))
			}
		})
	}
}

func (s *KeeperTestSuite) TestNewFixedDepositCfgs() {
	s.SetupTest()

	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	newFixdDepositCfg := &types.MsgNewFixedDepositCfg{
		Dao:      s.Dao.GlobalDao,
		RegionId: strings.ToLower(types.MeEarthRegionName),
		Term:     30,
		Rate:     sdk.NewDec(1),
	}
	_, err = s.msgServer.NewFixedDepositCfg(s.Ctx, newFixdDepositCfg)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestRemoveFixedDepositCfg() {
	s.SetupTest()
	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	newFixdDepositCfg := &types.MsgNewFixedDepositCfg{
		Dao:      s.Dao.GlobalDao,
		RegionId: types.MeEarthRegionId,
		Term:     30,
		Rate:     sdk.NewDec(1),
	}
	_, err = s.msgServer.NewFixedDepositCfg(s.Ctx, newFixdDepositCfg)
	s.Require().NoError(err)

	_, err = s.msgServer.RemoveFixedDepositCfg(s.Ctx,
		&types.MsgRemoveFixedDepositCfg{
			Admin:    s.Dao.GlobalDao,
			RegionId: types.MeEarthRegionId,
			Term:     30,
		})
	s.Require().NoError(err)

	_, err = s.queryClient.FixedDepositCfg(s.Ctx, &types.QueryFixedDepositCfgRequest{RegionIds: []string{types.MeEarthRegionName}})
	s.Require().ErrorIs(err, nil)
}

func (s *KeeperTestSuite) TestWithdrawFixedDeposit() {
	s.SetupTest()

	newMeEarthRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newMeEarthRegion)

	s.Require().NoError(err)

	err = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, mintypes.ModuleName, s.App.StakingKeeper.GetRegionAccount(s.Ctx, types.RegionAccountTypeBase, types.MeEarthRegionId).GetAddress(), sdk.Coins{sdk.NewInt64Coin(params.BaseDenom, 1000000000000)})
	s.Require().NoError(err)

	newFixdDepositCfg := &types.MsgNewFixedDepositCfg{
		Dao:      s.Dao.GlobalDao,
		RegionId: types.MeEarthRegionId,
		Term:     30,
		Rate:     sdk.NewDec(10),
	}
	_, err = s.msgServer.NewFixedDepositCfg(s.Ctx, newFixdDepositCfg)
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
	s.T().Logf("poolBalance balance: %s", poolBalance.String())

	daoBalance := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(s.Dao.GlobalDao), params.BaseDenom)
	s.T().Logf("daoBalance balance: %s", daoBalance.String())

	region, _ := s.App.StakingKeeper.GetRegion(s.Ctx, types.MeEarthRegionId)
	regionInterestAddr, _ := sdk.AccAddressFromBech32(region.DepositInterestAddr)
	interestBalance := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(regionInterestAddr.String()), params.BaseDenom)
	s.T().Logf("interestBalance balance: %s", interestBalance.String())

	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wmintTypes.OneYearTotalBlocks).WithChainID(apptesting.TestChainID).WithBlockTime(s.Ctx.BlockTime().Add(7760 * time.Hour))
	_, err = s.msgServer.WithdrawFixedDeposit(s.Ctx, &types.MsgWithdrawFixedDeposit{
		Account: s.Dao.GlobalDao,
		Id:      fixDeposit.Id,
	})

	poolBalanceAfer := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(principalAddr.String()), params.BaseDenom)
	s.T().Logf("poolBalanceAfer balance: %s", poolBalanceAfer.String())

	daobalanceAfer := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(s.Dao.GlobalDao), params.BaseDenom)
	s.T().Logf("daobalanceAfer balance: %s", daobalanceAfer.String())

	interestBalanceAfter := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(regionInterestAddr.String()), params.BaseDenom)
	s.T().Logf("interestBalanceAfter balance: %s", interestBalanceAfter.String())

	s.Require().Equal(daobalanceAfer.String(), daoBalance.Add(poolBalance).Add(interestBalance).String())
	s.Require().NoError(err)
}
