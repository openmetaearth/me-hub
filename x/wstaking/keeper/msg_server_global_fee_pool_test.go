package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

// TestWithdrawFromGlobalDaoFeePool tests withdrawing from global DAO fee pool
func (s *KeeperTestSuite) TestWithdrawFromGlobalDaoFeePool() {
	testCases := []struct {
		name      string
		setup     func() *types.MsgWithdrawFromGlobalDaoFeePool
		expectErr bool
		errMsg    string
	}{
		{
			name: "successful withdrawal",
			setup: func() *types.MsgWithdrawFromGlobalDaoFeePool {
				// Fund the global DAO fee pool
				feePoolAddr := s.App.DaoKeeper.GetGlobalDaoFeePoolAddr(s.Ctx)
				s.FundAcc(feePoolAddr, sdk.NewCoins(sdk.NewCoin(params.BaseDenom, math.NewInt(1000000000))))

				return &types.MsgWithdrawFromGlobalDaoFeePool{
					Withdrawer: s.Dao.GlobalDao,
					Amount:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, math.NewInt(500000000))),
				}
			},
			expectErr: false,
		},
		{
			name: "non-dao withdrawer",
			setup: func() *types.MsgWithdrawFromGlobalDaoFeePool {
				// Fund the global DAO fee pool
				feePoolAddr := s.App.DaoKeeper.GetGlobalDaoFeePoolAddr(s.Ctx)
				s.FundAcc(feePoolAddr, sdk.NewCoins(sdk.NewCoin(params.BaseDenom, math.NewInt(1000000000))))

				account := s.TestAccs[0]
				return &types.MsgWithdrawFromGlobalDaoFeePool{
					Withdrawer: account.String(),
					Amount:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, math.NewInt(500000000))),
				}
			},
			expectErr: true,
			errMsg:    "global dao",
		},
		{
			name: "insufficient balance in fee pool",
			setup: func() *types.MsgWithdrawFromGlobalDaoFeePool {
				// Don't fund the fee pool or fund with less than withdrawal amount
				feePoolAddr := s.App.DaoKeeper.GetGlobalDaoFeePoolAddr(s.Ctx)
				s.FundAcc(feePoolAddr, sdk.NewCoins(sdk.NewCoin(params.BaseDenom, math.NewInt(100000000))))

				return &types.MsgWithdrawFromGlobalDaoFeePool{
					Withdrawer: s.Dao.GlobalDao,
					Amount:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, math.NewInt(500000000))),
				}
			},
			expectErr: true,
			errMsg:    "insufficient",
		},
		{
			name: "invalid withdrawer address",
			setup: func() *types.MsgWithdrawFromGlobalDaoFeePool {
				return &types.MsgWithdrawFromGlobalDaoFeePool{
					Withdrawer: "invalid_address",
					Amount:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, math.NewInt(500000000))),
				}
			},
			expectErr: true,
			errMsg:    "global dao", // IsGlobalDao check happens before address format validation
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Use fresh context for each test
			s.SetupTest()
			msg := tc.setup()

			resp, err := s.msgServer.WithdrawFromGlobalDaoFeePool(s.Ctx, msg)

			if tc.expectErr {
				s.Require().Error(err)
				if tc.errMsg != "" {
					s.Require().Contains(err.Error(), tc.errMsg)
				}
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)

				// Verify balance was transferred
				withdrawerAddr, _ := sdk.AccAddressFromBech32(msg.Withdrawer)
				balance := s.App.BankKeeper.GetBalance(s.Ctx, withdrawerAddr, params.BaseDenom)
				s.Require().True(balance.Amount.GTE(msg.Amount[0].Amount))

				// Verify event was emitted
				events := s.Ctx.EventManager().Events()
				foundEvent := false
				for _, event := range events {
					if event.Type == types.EventTypeWithdrawFromGlobalDaoFeePool {
						foundEvent = true
						break
					}
				}
				s.Require().True(foundEvent, "WithdrawFromGlobalDaoFeePool event should be emitted")
			}
		})
	}
}

// TestWithdrawFromGlobalDaoFeePoolMultipleTimes tests multiple withdrawals
func (s *KeeperTestSuite) TestWithdrawFromGlobalDaoFeePoolMultipleTimes() {
	// Setup: Fund the global DAO fee pool
	feePoolAddr := s.App.DaoKeeper.GetGlobalDaoFeePoolAddr(s.Ctx)
	initialAmount := math.NewInt(1000000000)
	s.FundAcc(feePoolAddr, sdk.NewCoins(sdk.NewCoin(params.BaseDenom, initialAmount)))

	withdrawAmount := math.NewInt(100000000)

	// First withdrawal
	msg1 := &types.MsgWithdrawFromGlobalDaoFeePool{
		Withdrawer: s.Dao.GlobalDao,
		Amount:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, withdrawAmount)),
	}
	resp1, err := s.msgServer.WithdrawFromGlobalDaoFeePool(s.Ctx, msg1)
	s.Require().NoError(err)
	s.Require().NotNil(resp1)

	// Second withdrawal
	msg2 := &types.MsgWithdrawFromGlobalDaoFeePool{
		Withdrawer: s.Dao.GlobalDao,
		Amount:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, withdrawAmount)),
	}
	resp2, err := s.msgServer.WithdrawFromGlobalDaoFeePool(s.Ctx, msg2)
	s.Require().NoError(err)
	s.Require().NotNil(resp2)

	// Verify fee pool balance decreased correctly
	feePoolBalance := s.App.BankKeeper.GetBalance(s.Ctx, feePoolAddr, params.BaseDenom)
	expectedBalance := initialAmount.Sub(withdrawAmount).Sub(withdrawAmount)
	s.Require().Equal(expectedBalance, feePoolBalance.Amount)

	// Verify withdrawer balance increased correctly
	withdrawerAddr, _ := sdk.AccAddressFromBech32(s.Dao.GlobalDao)
	withdrawerBalance := s.App.BankKeeper.GetBalance(s.Ctx, withdrawerAddr, params.BaseDenom)
	s.Require().True(withdrawerBalance.Amount.GTE(withdrawAmount.Add(withdrawAmount)))
}
