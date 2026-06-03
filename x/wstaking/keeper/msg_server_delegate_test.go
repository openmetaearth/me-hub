package keeper_test

import (
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wmint"
	"github.com/openmetaearth/me-hub/x/wstaking"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestDelegate() {
	s.SetupTest()

	newMeEarthRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newMeEarthRegion)
	s.Require().NoError(err)

	region, _ := s.App.StakingKeeper.GetRegion(s.Ctx, types.MeEarthRegionId)

	err = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, mintypes.ModuleName, sdk.MustAccAddressFromBech32(region.GetRegionTreasureAddr()), sdk.Coins{sdk.NewInt64Coin(params.BaseDenom, 2000000000000)})
	s.Require().NoError(err)

	tests := []struct {
		name             string
		account          string
		amount           sdk.Coin
		reward           float64
		height           int64
		validatorAddress string
		expErr           error
	}{
		{
			name:             "did delegate",
			account:          s.Dao.GlobalDao,
			amount:           sdk.NewCoin(params.BaseDenom, sdk.NewInt(1000000)),
			height:           5,
			reward:           0.1981862,
			validatorAddress: s.meEarthValidator.OperatorAddress,
			expErr:           nil,
		},
		{
			name:             "un did delegate",
			account:          s.Dao.AirdropAddress,
			amount:           sdk.NewCoin(params.BaseDenom, sdk.NewInt(1000000)),
			height:           5,
			reward:           0.1981862,
			validatorAddress: s.experienceValidator.OperatorAddress,
			expErr:           nil,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			msg := stakingtypes.MsgDelegate{
				DelegatorAddress: test.account,
				ValidatorAddress: test.validatorAddress,
				Amount:           test.amount,
			}
			_, err := s.msgServer.Delegate(s.Ctx, &msg)
			s.Require().ErrorIs(err, test.expErr)

			// check nft class
			if test.expErr == nil {
				withdrawRewardMsg := types.MsgWithdrawDelegatorReward{
					DelegatorAddress: test.account,
					ValidatorAddress: test.validatorAddress,
				}

				s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(test.height).WithChainID(apptesting.TestChainID)
				for i := 0; i < int(test.height); i++ {
					wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
					wstaking.BeginBlock(s.Ctx, s.App.StakingKeeper)
				}

				rewards, err := s.msgServer.DelegationRewards(s.Ctx, &types.QueryDelegationRewardsRequest{
					DelegatorAddress: test.account,
					ValidatorAddress: test.validatorAddress,
				})
				s.Require().NoError(err)
				_, err = s.msgServer.WithdrawDelegatorReward(s.Ctx, &withdrawRewardMsg)
				s.Require().NoError(err)
				s.Require().Equal(rewards.Rewards[0].Amount.MustFloat64(), test.reward)
			}
		})
	}
}

func (s *KeeperTestSuite) TestUnDelegate() {
	s.SetupTest()

	newMeEarthRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newMeEarthRegion)
	s.Require().NoError(err)

	region, _ := s.App.StakingKeeper.GetRegion(s.Ctx, types.MeEarthRegionId)

	err = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, mintypes.ModuleName, sdk.MustAccAddressFromBech32(region.GetRegionTreasureAddr()), sdk.Coins{sdk.NewInt64Coin(params.BaseDenom, 1000000000000)})
	s.Require().NoError(err)

	tests := []struct {
		name    string
		account string
		amount  sdk.Coin
		reward  float64
		height  int64
		expErr  error
	}{
		{
			name:    "un did undelegate",
			account: s.Dao.AirdropAddress,
			amount:  sdk.NewCoin(params.BaseDenom, sdk.NewInt(1000000)),
			height:  5,
			reward:  0.1981862,
			expErr:  nil,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			msg := stakingtypes.MsgDelegate{
				DelegatorAddress: test.account,
				ValidatorAddress: "",
				Amount:           test.amount,
			}
			_, err := s.msgServer.Delegate(s.Ctx, &msg)
			s.Require().ErrorIs(err, test.expErr)

			// check nft class
			if test.expErr == nil {
				undelegateRewardMsg := stakingtypes.MsgUndelegate{
					DelegatorAddress: test.account,
					ValidatorAddress: "",
					Amount:           test.amount,
				}

				s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(test.height).WithChainID(apptesting.TestChainID)
				for i := 0; i < int(test.height); i++ {
					wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
					wstaking.BeginBlock(s.Ctx, s.App.StakingKeeper)
				}

				rewards, err := s.msgServer.DelegationRewards(s.Ctx, &types.QueryDelegationRewardsRequest{
					DelegatorAddress: test.account,
					ValidatorAddress: "",
				})
				s.Require().NoError(err)
				_, err = s.msgServer.Undelegate(s.Ctx, &undelegateRewardMsg)
				s.Require().NoError(err)
				s.Require().Equal(rewards.Rewards[0].Amount.MustFloat64(), test.reward)
			}
		})
	}
}

func (s *KeeperTestSuite) TestDelegateRejectsInterestPayoutWhenRegionTreasuryUnderfunded() {
	s.SetupTest()

	delegator := s.Dao.AirdropAddress
	amount := sdk.NewCoin(params.BaseDenom, sdk.NewInt(1000000))
	_, err := s.msgServer.Delegate(s.Ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: delegator,
		ValidatorAddress: "",
		Amount:           amount,
	})
	s.Require().NoError(err)

	region, delegation := s.prepareUnderfundedDelegateInterest(delegator, 100)

	_, err = s.msgServer.Delegate(s.Ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: delegator,
		ValidatorAddress: "",
		Amount:           amount,
	})
	s.Require().ErrorIs(err, types.ErrPayInterest)

	regionAfter, found := s.Keeper().GetRegion(s.Ctx, region.RegionId)
	s.Require().True(found)
	s.Require().True(regionAfter.DelegateInterest.Equal(region.DelegateInterest))
	s.Require().Equal(region.DelegateAmount.String(), regionAfter.DelegateAmount.String())

	delegationAfter, found := s.Keeper().GetDelegation(s.Ctx, sdk.MustAccAddressFromBech32(delegator), sdk.ValAddress{})
	s.Require().True(found)
	s.Require().Equal(delegation.Amount.String(), delegationAfter.Amount.String())
}

func (s *KeeperTestSuite) TestUndelegateRejectsInterestPayoutWhenRegionTreasuryUnderfunded() {
	s.SetupTest()

	delegator := s.Dao.AirdropAddress
	amount := sdk.NewCoin(params.BaseDenom, sdk.NewInt(1000000))
	_, err := s.msgServer.Delegate(s.Ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: delegator,
		ValidatorAddress: "",
		Amount:           amount,
	})
	s.Require().NoError(err)

	region, delegation := s.prepareUnderfundedDelegateInterest(delegator, 100)

	_, err = s.msgServer.Undelegate(s.Ctx, &stakingtypes.MsgUndelegate{
		DelegatorAddress: delegator,
		ValidatorAddress: "",
		Amount:           amount,
	})
	s.Require().ErrorIs(err, types.ErrPayInterest)

	regionAfter, found := s.Keeper().GetRegion(s.Ctx, region.RegionId)
	s.Require().True(found)
	s.Require().True(regionAfter.DelegateInterest.Equal(region.DelegateInterest))
	s.Require().Equal(region.DelegateAmount.String(), regionAfter.DelegateAmount.String())

	delegationAfter, found := s.Keeper().GetDelegation(s.Ctx, sdk.MustAccAddressFromBech32(delegator), sdk.ValAddress{})
	s.Require().True(found)
	s.Require().Equal(delegation.Amount.String(), delegationAfter.Amount.String())
}

func (s *KeeperTestSuite) TestUndelegateRejectsInterestPayoutWhenRegionReserveUnderfunded() {
	s.SetupTest()

	delegator := s.Dao.AirdropAddress
	amount := sdk.NewCoin(params.BaseDenom, sdk.NewInt(1000000))
	_, err := s.msgServer.Delegate(s.Ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: delegator,
		ValidatorAddress: "",
		Amount:           amount,
	})
	s.Require().NoError(err)

	region, delegation := s.prepareUnderfundedDelegateInterest(delegator, 100)

	rewards, err := s.Keeper().CalculateInterest(
		s.Ctx,
		delegation.Amount.Add(delegation.UnMeidAmount).Add(delegation.Unmovable),
		delegation.StartHeight,
	)
	s.Require().NoError(err)
	rewardCoin := sdk.NewCoin(params.BaseDenom, rewards.TruncateInt())
	s.Require().True(rewardCoin.Amount.IsPositive())

	region.DelegateInterest = rewards.Add(sdk.NewDec(100))
	s.Keeper().SetRegion(s.Ctx, region)

	err = s.App.BankKeeper.SendCoinsFromModuleToAccount(
		s.Ctx,
		mintypes.ModuleName,
		sdk.MustAccAddressFromBech32(region.RegionTreasureAddr),
		sdk.NewCoins(rewardCoin),
	)
	s.Require().NoError(err)

	_, err = s.msgServer.Undelegate(s.Ctx, &stakingtypes.MsgUndelegate{
		DelegatorAddress: delegator,
		ValidatorAddress: "",
		Amount:           amount,
	})
	s.Require().ErrorIs(err, types.ErrPayInterest)

	regionAfter, found := s.Keeper().GetRegion(s.Ctx, region.RegionId)
	s.Require().True(found)
	s.Require().True(regionAfter.DelegateInterest.Equal(region.DelegateInterest))
	s.Require().Equal(region.DelegateAmount.String(), regionAfter.DelegateAmount.String())
}

func (s *KeeperTestSuite) prepareUnderfundedDelegateInterest(delegator string, blockHeight int64) (types.Region, stakingtypes.Delegation) {
	delegatorAddr := sdk.MustAccAddressFromBech32(delegator)
	delegation, found := s.Keeper().GetDelegation(s.Ctx, delegatorAddr, sdk.ValAddress{})
	s.Require().True(found)

	s.Ctx = s.Ctx.WithBlockHeight(blockHeight)
	rewards, err := s.Keeper().CalculateInterest(
		s.Ctx,
		delegation.Amount.Add(delegation.UnMeidAmount).Add(delegation.Unmovable),
		delegation.StartHeight,
	)
	s.Require().NoError(err)
	s.Require().True(rewards.TruncateInt().IsPositive())

	region, found := s.Keeper().GetRegion(s.Ctx, s.experienceValidator.Description.RegionID)
	s.Require().True(found)
	region.DelegateInterest = rewards
	s.Keeper().SetRegion(s.Ctx, region)

	treasuryAddr := sdk.MustAccAddressFromBech32(region.RegionTreasureAddr)
	treasuryBalance := s.App.BankKeeper.GetBalance(s.Ctx, treasuryAddr, params.BaseDenom)
	if treasuryBalance.Amount.IsPositive() {
		err = s.App.BankKeeper.SendCoins(
			s.Ctx,
			treasuryAddr,
			sdk.MustAccAddressFromBech32(s.Dao.GlobalDao),
			sdk.NewCoins(treasuryBalance),
		)
		s.Require().NoError(err)
	}

	return region, delegation
}
