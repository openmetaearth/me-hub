package keeper_test

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/x/wdistri"
	"github.com/st-chain/me-hub/x/wmint"
	"github.com/st-chain/me-hub/x/wstaking/types"
	"strings"
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
		Admin:    s.Dao.GlobalDao,
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
