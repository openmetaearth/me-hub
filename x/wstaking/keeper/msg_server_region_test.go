package keeper_test

import (
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/st-chain/me-hub/app/apptesting"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/x/wdistri"
	"github.com/st-chain/me-hub/x/wmint"
	wmintTypes "github.com/st-chain/me-hub/x/wmint/types"
	"github.com/st-chain/me-hub/x/wstaking/types"
	"strings"
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
				_, found := s.App.NFTKeeper.GetClass(s.Ctx, types.GetClassId(test.regionName))
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
		Name:            types.ExperienceRegionName,
		OperatorAddress: s.experienceValidator.OperatorAddress,
	}
	_, err = s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	newRegion = types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            "USA",
		OperatorAddress: s.usaValidator.OperatorAddress,
	}
	_, err = s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	// must have error
	_, err = s.msgServer.RemoveRegion(s.Ctx, &types.MsgRemoveRegion{
		Creator:  s.Dao.MeidDao,
		RegionId: "USA",
	})
	s.Require().ErrorIs(err, types.ErrCheckGlobalDao)

	// must no error
	_, err = s.msgServer.RemoveRegion(s.Ctx, &types.MsgRemoveRegion{
		Creator:  s.Dao.GlobalDao,
		RegionId: "usa",
	})
	s.Require().NoError(err)

	// must error
	_, err = s.queryClient.Region(s.Ctx, &types.QueryRegionRequest{RegionId: "usa"})
	s.Require().ErrorIs(err, types.ErrRegionNotExist)
}

func (s *KeeperTestSuite) TestWithdrawFromRegion() {
	s.SetupTest()

	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.ExperienceRegionName,
		OperatorAddress: s.experienceValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
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
		name       string
		withdrawer string
		amount     sdk.Coin
		expErr     error
	}{
		{
			name:       "Dao Permission",
			withdrawer: s.Dao.MeidDao,
			amount:     balance,
			expErr:     types.ErrCheckGlobalDao,
		}, {
			name:       "over amount",
			withdrawer: s.Dao.GlobalDao,
			amount: balance.Add(sdk.Coin{
				Denom:  params.BaseDenom,
				Amount: sdk.NewInt(1),
			}),
			expErr: sdkerrors.ErrInsufficientFunds,
		}, {
			name:       "pass",
			withdrawer: s.Dao.GlobalDao,
			amount:     balance,
			expErr:     nil,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			daoBalanceBefore := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(s.Dao.GlobalDao), params.BaseDenom)

			msg := types.MsgWithdrawFromRegion{
				Withdrawer: test.withdrawer,
				RegionId:   strings.ToLower(types.ExperienceRegionName),
				Receiver:   s.Dao.GlobalDao,
				Amount:     sdk.NewCoins(test.amount),
			}
			_, err = s.msgServer.WithdrawFromRegion(s.Ctx, &msg)
			s.Require().ErrorIs(err, test.expErr)

			daoBalanceAfter := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(s.Dao.GlobalDao), params.BaseDenom)
			if test.expErr == nil {
				s.Require().Equal(balance.String(), daoBalanceAfter.Sub(daoBalanceBefore).String())
			}
		})
	}
}
