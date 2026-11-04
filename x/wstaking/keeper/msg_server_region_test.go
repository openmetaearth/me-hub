package keeper_test

import (
	"strings"

	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wdistri"
	"github.com/openmetaearth/me-hub/x/wmint"
	wminttypes "github.com/openmetaearth/me-hub/x/wmint/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestNewRegion() {
	s.SetupTest()

	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	tests := []struct {
		name            string
		creator         string
		regionName      string
		operatorAddress string
		expErr          error
		malleate        func()
	}{
		{
			name:            "Dao Permission",
			creator:         s.Dao.MeidDao,
			regionName:      "USA",
			operatorAddress: s.usaValidator.OperatorAddress,
			expErr:          types.ErrCheckGlobalDao,
		}, {
			name:            "have permission, but wrong region id",
			creator:         s.Dao.GlobalDao,
			regionName:      "USS",
			operatorAddress: s.usaValidator.OperatorAddress,
			expErr:          types.ErrRegionName,
		}, {
			name:            "wrong validator address",
			creator:         s.Dao.GlobalDao,
			regionName:      "USA",
			operatorAddress: "mevaloper139mq752delxv78jvtmwxhasyrycufsvr707ate",
			expErr:          types.ErrRegionValidatorNotExist,
		}, {
			name:            "wrong validator region id",
			creator:         s.Dao.GlobalDao,
			regionName:      "USA",
			operatorAddress: s.meEarthValidator.OperatorAddress,
			expErr:          types.ErrRegion,
		}, {
			name:            "No error",
			creator:         s.Dao.GlobalDao,
			regionName:      "USA",
			operatorAddress: s.usaValidator.OperatorAddress,
			expErr:          nil,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			newRegion := types.MsgNewRegion{
				Creator:         test.creator,
				Name:            test.regionName,
				OperatorAddress: test.operatorAddress,
			}
			_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
			s.Require().ErrorIs(err, test.expErr)

			// check nft class
			if test.expErr == nil {
				_, found := s.App.WNFTKeeper.GetClass(s.Ctx, types.GetClassId(test.regionName))
				s.Require().True(found)
			}
		})
	}
}

func (s *KeeperTestSuite) TestRemoveRegion() {
	s.SetupTest()
	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	newRegion = types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            "USA",
		OperatorAddress: s.usaValidator.OperatorAddress,
	}
	_, err = s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	// RemoveRegion is currently a no-op, should return no error
	_, err = s.msgServer.RemoveRegion(s.Ctx, &types.MsgRemoveRegion{
		Creator:  s.Dao.MeidDao,
		RegionId: "usa",
	})
	s.Require().NoError(err)

	_, err = s.msgServer.RemoveRegion(s.Ctx, &types.MsgRemoveRegion{
		Creator:  s.Dao.GlobalDao,
		RegionId: "usa",
	})
	s.Require().NoError(err)

	// Since RemoveRegion is a no-op, region still exists
	_, err = s.queryClient.Region(s.Ctx, &types.QueryRegionRequest{RegionId: "usa"})
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestRemoveRegionThenCreateRegion() {
	s.SetupTest()
	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	newRegion = types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            "USA",
		OperatorAddress: s.usaValidator.OperatorAddress,
	}
	_, err = s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	// RemoveRegion is currently a no-op, should return no error
	_, err = s.msgServer.RemoveRegion(s.Ctx, &types.MsgRemoveRegion{
		Creator:  s.Dao.MeidDao,
		RegionId: "usa",
	})
	s.Require().NoError(err)

	_, err = s.msgServer.RemoveRegion(s.Ctx, &types.MsgRemoveRegion{
		Creator:  s.Dao.GlobalDao,
		RegionId: "usa",
	})
	s.Require().NoError(err)

	// Since RemoveRegion is a no-op, region still exists, creating again should fail
	_, err = s.queryClient.Region(s.Ctx, &types.QueryRegionRequest{RegionId: "usa"})
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestWithdrawRegionTreasure() {
	s.SetupTest()

	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wminttypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}, *s.App.DistrKeeper)

	regionResp, err := s.queryClient.Region(s.Ctx, &types.QueryRegionRequest{RegionId: strings.ToLower(types.ExperienceRegionName)})
	s.Require().NoError(err)

	balance := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(regionResp.Region.RegionTreasureAddr), params.BaseDenom)
	s.T().Log(balance.String())

	amount := s.App.MintKeeper.GetMintedCoinAmount(s.Ctx)
	s.T().Log(amount.String())

	s.Require().Equal(balance.Amount.String(), amount.String())

	tests := []struct {
		name   string
		amount sdk.Coin
		expErr error
	}{
		{
			name: "over amount",
			amount: balance.Add(sdk.Coin{
				Denom:  params.BaseDenom,
				Amount: sdk.NewInt(1),
			}),
			expErr: sdkerrors.ErrInsufficientFunds,
		}, {
			name:   "pass",
			amount: balance,
			expErr: nil,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			daoBalanceBefore := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(s.Dao.GlobalDao), params.BaseDenom)

			err := s.Keeper().WithdrawRegionTreasure(
				s.Ctx,
				types.ExperienceRegionId,
				s.Dao.GlobalDao,
				sdk.NewCoins(test.amount),
			)
			s.Require().ErrorIs(err, test.expErr)

			daoBalanceAfter := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(s.Dao.GlobalDao), params.BaseDenom)
			if test.expErr == nil {
				s.Require().Equal(balance.String(), daoBalanceAfter.Sub(daoBalanceBefore).String())
			}
		})
	}
}

func (s *KeeperTestSuite) TestWithdrawRegionTreasureRejectsBlockedModuleReceiver() {
	s.SetupTest()

	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wminttypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}, *s.App.DistrKeeper)

	regionResp, err := s.queryClient.Region(s.Ctx, &types.QueryRegionRequest{RegionId: types.ExperienceRegionId})
	s.Require().NoError(err)

	regionTreasury := sdk.MustAccAddressFromBech32(regionResp.Region.RegionTreasureAddr)
	regionBalance := s.App.BankKeeper.GetBalance(s.Ctx, regionTreasury, params.BaseDenom)
	s.Require().True(regionBalance.IsPositive())

	blockedReceiver := authtypes.NewModuleAddress(authtypes.FeeCollectorName)
	s.Require().True(s.App.BankKeeper.BlockedAddr(blockedReceiver))

	receiverBefore := s.App.BankKeeper.GetBalance(s.Ctx, blockedReceiver, params.BaseDenom)
	treasuryBefore := s.App.BankKeeper.GetBalance(s.Ctx, regionTreasury, params.BaseDenom)

	err = s.Keeper().WithdrawRegionTreasure(
		s.Ctx,
		types.ExperienceRegionId,
		blockedReceiver.String(),
		sdk.NewCoins(regionBalance),
	)

	s.Require().ErrorIs(err, sdkerrors.ErrUnauthorized)
	s.Require().Contains(err.Error(), "not allowed to receive region treasury withdrawals")
	s.Require().Equal(receiverBefore.String(), s.App.BankKeeper.GetBalance(s.Ctx, blockedReceiver, params.BaseDenom).String())
	s.Require().Equal(treasuryBefore.String(), s.App.BankKeeper.GetBalance(s.Ctx, regionTreasury, params.BaseDenom).String())
}
