package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wstaking"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"math/big"
	"strings"
)

func (s *KeeperTestSuite) TestStake() {
	s.SetupTest()

	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	moduleAddress := authtypes.NewModuleAddress(types.StakePoolName)
	stakePoolBalanceBefore := s.App.BankKeeper.GetBalance(s.Ctx, moduleAddress, params.BaseDenom)
	s.Require().Equal(stakePoolBalanceBefore.String(), "1000000000000000000umec")

	stakeAmount := sdk.NewCoin(params.BaseDenom, sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(params.BaseDenomUnit), nil)))

	valAddress, err := sdk.ValAddressFromBech32(s.meEarthValidator.OperatorAddress)
	s.Require().NoError(err)

	tests := []struct {
		name            string
		staker          string
		operatorAddress string
		amount          sdk.Coin
		expErr          error
		malleate        func()
	}{
		{
			name:            "Dao Permission",
			staker:          s.Dao.MeidDao,
			operatorAddress: s.usaValidator.OperatorAddress,
			amount:          stakeAmount,
			expErr:          types.ErrCheckGlobalDao,
		}, {
			name:            "wrong validator address",
			staker:          s.Dao.GlobalDao,
			operatorAddress: "mevaloper139mq752delxv78jvtmwxhasyrycufsvr707ate",
			amount:          stakeAmount,
			expErr:          stakingtypes.ErrNoValidatorFound,
		}, {
			name:            "small amount",
			staker:          s.Dao.GlobalDao,
			operatorAddress: s.meEarthValidator.OperatorAddress,
			amount:          sdk.NewCoin(params.BaseDenom, sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(params.BaseDenomUnit-1), nil))),
			expErr:          sdkerrors.ErrInvalidRequest,
		}, {
			name:            "No error",
			staker:          s.Dao.GlobalDao,
			operatorAddress: s.meEarthValidator.OperatorAddress,
			amount:          stakeAmount,
			expErr:          nil,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			msg := types.MsgStake{
				StakerAddress:    test.staker,
				ValidatorAddress: test.operatorAddress,
				Amount:           test.amount,
			}

			regionBefore, found := s.Keeper().GetRegion(s.Ctx, strings.ToLower(types.MeEarthRegionName))
			s.Require().True(found)

			validatorBefore, found := s.Keeper().GetValidator(s.Ctx, valAddress)
			s.Require().True(found)

			stakeBefore, _ := s.Keeper().GetStake(s.Ctx, sdk.MustAccAddressFromBech32(s.Dao.GlobalDao), valAddress)
			if stakeBefore.Shares.IsNil() {
				stakeBefore.Shares = sdk.ZeroDec()
			}

			_, err := s.msgServer.Stake(s.Ctx, &msg)
			s.Require().ErrorIs(err, test.expErr)

			if test.expErr == nil {
				// check stake pool balance
				stakePoolBalanceAfter := s.App.BankKeeper.GetBalance(s.Ctx, moduleAddress, params.BaseDenom)
				s.Require().Equal(stakePoolBalanceAfter.Amount.String(), stakePoolBalanceBefore.Sub(stakeAmount).Amount.String())

				// check region
				regionAfter, found := s.Keeper().GetRegion(s.Ctx, strings.ToLower(types.MeEarthRegionName))
				s.Require().True(found)
				s.Require().Equal(regionAfter.RegionShare.String(), regionBefore.RegionShare.Add(stakeAmount.Amount).String())

				// check validator
				validatorAfter, found := s.Keeper().GetValidator(s.Ctx, valAddress)
				s.Require().True(found)
				s.Require().Equal(validatorAfter.Tokens.String(), validatorBefore.Tokens.Add(stakeAmount.Amount).String())
				shares, err := validatorBefore.SharesFromTokens(stakeAmount.Amount)
				s.Require().NoError(err)
				s.Require().Equal(validatorAfter.DelegatorShares.String(), validatorBefore.DelegatorShares.Add(shares).String())

				// check stake
				stakeAfter, found := s.Keeper().GetStake(s.Ctx, sdk.MustAccAddressFromBech32(s.Dao.GlobalDao), valAddress)
				s.Require().True(found)
				s.Require().Equal(stakeAfter.Shares.String(), stakeBefore.Shares.Add(shares).String())
			}
		})
	}
}

func (s *KeeperTestSuite) TestUnStake() {
	s.SetupTest()

	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	moduleAddress := authtypes.NewModuleAddress(types.StakePoolName)
	stakePoolBalanceBefore := s.App.BankKeeper.GetBalance(s.Ctx, moduleAddress, params.BaseDenom)
	s.Require().Equal("1000000000000000000umec", stakePoolBalanceBefore.String())

	stakeAmount := sdk.NewCoin(params.BaseDenom, sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(params.BaseDenomUnit), nil)))

	valAddress, err := sdk.ValAddressFromBech32(s.meEarthValidator.OperatorAddress)
	s.Require().NoError(err)

	_, err = s.msgServer.Stake(s.Ctx, &types.MsgStake{
		StakerAddress:    s.Dao.GlobalDao,
		ValidatorAddress: s.meEarthValidator.OperatorAddress,
		Amount:           stakeAmount,
	})
	s.Require().NoError(err)
	stakePoolBalanceBefore = s.App.BankKeeper.GetBalance(s.Ctx, moduleAddress, params.BaseDenom)
	s.Require().Equal("999999999900000000umec", stakePoolBalanceBefore.String())

	tests := []struct {
		name            string
		staker          string
		operatorAddress string
		amount          sdk.Coin
		expErr          error
		malleate        func()
	}{
		{
			name:            "Dao Permission",
			staker:          s.Dao.MeidDao,
			operatorAddress: s.usaValidator.OperatorAddress,
			amount:          stakeAmount,
			expErr:          types.ErrCheckGlobalDao,
		}, {
			name:            "wrong validator address",
			staker:          s.Dao.GlobalDao,
			operatorAddress: "mevaloper139mq752delxv78jvtmwxhasyrycufsvr707ate",
			amount:          stakeAmount,
			expErr:          stakingtypes.ErrNoValidatorFound,
		}, {
			name:            "No error",
			staker:          s.Dao.GlobalDao,
			operatorAddress: s.meEarthValidator.OperatorAddress,
			amount:          stakeAmount,
			expErr:          nil,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			msg := types.MsgUnstake{
				StakerAddress:    test.staker,
				ValidatorAddress: test.operatorAddress,
				Amount:           test.amount,
			}

			regionBefore, found := s.Keeper().GetRegion(s.Ctx, strings.ToLower(types.MeEarthRegionName))
			s.Require().True(found)

			validatorBefore, found := s.Keeper().GetValidator(s.Ctx, valAddress)
			s.Require().True(found)

			_, err := s.msgServer.Unstake(s.Ctx, &msg)
			s.Require().ErrorIs(err, test.expErr)

			// call endblock for complete unstake
			wstaking.EndBlock(s.Ctx, s.Keeper())
			if test.expErr == nil {
				// check stake pool balance
				stakePoolBalanceAfter := s.App.BankKeeper.GetBalance(s.Ctx, moduleAddress, params.BaseDenom)
				s.Require().Equal("1000000000000000000umec", stakePoolBalanceAfter.String())

				// check region
				regionAfter, found := s.Keeper().GetRegion(s.Ctx, strings.ToLower(types.MeEarthRegionName))
				s.Require().True(found)
				s.Require().Equal(regionAfter.RegionShare.String(), regionBefore.RegionShare.Sub(stakeAmount.Amount).String())

				// check validator
				validatorAfter, found := s.Keeper().GetValidator(s.Ctx, valAddress)
				s.Require().True(found)
				s.Require().Equal(validatorBefore.Tokens.Sub(stakeAmount.Amount).String(), validatorAfter.Tokens.String())

				shares, err := validatorBefore.SharesFromTokens(stakeAmount.Amount)
				s.Require().NoError(err)
				s.Require().Equal(validatorBefore.DelegatorShares.Sub(shares).String(), validatorAfter.DelegatorShares.String())

				// check stake
				_, found = s.Keeper().GetStake(s.Ctx, sdk.MustAccAddressFromBech32(s.Dao.GlobalDao), valAddress)
				s.Require().False(found)
			}
		})
	}
}
