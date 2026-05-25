package keeper_test

import (
	"strings"

	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wdistri"
	"github.com/openmetaearth/me-hub/x/wmint"
	wmintTypes "github.com/openmetaearth/me-hub/x/wmint/types"
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

	// must have error
	_, err = s.msgServer.RemoveRegion(s.Ctx, &types.MsgRemoveRegion{
		Creator:  s.Dao.MeidDao,
		RegionId: "usa",
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

	// must have error
	_, err = s.msgServer.RemoveRegion(s.Ctx, &types.MsgRemoveRegion{
		Creator:  s.Dao.MeidDao,
		RegionId: "usa",
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

	// new region again
	newRegion = types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            "USA",
		OperatorAddress: s.usaValidator.OperatorAddress,
	}
	_, err = s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestWithdrawFromRegion() {
	s.SetupTest()

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

func (s *KeeperTestSuite) TestGrantRegionWithdraw() {
	s.SetupTest()

	regionId := types.ExperienceRegionId
	addr0 := s.TestAccs[0].String()
	addr1 := s.TestAccs[1].String()
	nonDao := s.TestAccs[2].String()

	tests := []struct {
		name     string
		creator  string
		regionId string
		address  string
		expErr   error
		malleate func()
		// afterCheck is called only when expErr == nil
		afterCheck func(res *types.QueryRegionWithdrawerResponse)
	}{
		{
			name:     "non-DAO cannot grant",
			creator:  nonDao,
			regionId: regionId,
			address:  addr0,
			expErr:   types.ErrCheckGlobalDao,
		},
		{
			name:     "region does not exist",
			creator:  s.Dao.GlobalDao,
			regionId: "nonexistent_region",
			address:  addr0,
			expErr:   types.ErrRegionNotExist,
		},
		{
			name:     "invalid grantee address",
			creator:  s.Dao.GlobalDao,
			regionId: regionId,
			address:  "not-a-valid-address",
			expErr:   sdkerrors.ErrInvalidAddress,
		},
		{
			name:     "grant success",
			creator:  s.Dao.GlobalDao,
			regionId: regionId,
			address:  addr0,
			expErr:   nil,
			afterCheck: func(res *types.QueryRegionWithdrawerResponse) {
				s.Require().Equal(addr0, res.Address)
				s.Require().True(s.Keeper().CanRegionWithdraw(s.Ctx, addr0, regionId))
			},
		},
		{
			name:     "overwrite with a different address",
			creator:  s.Dao.GlobalDao,
			regionId: regionId,
			address:  addr1,
			// addr0 already granted by the previous "grant success" case
			expErr: nil,
			afterCheck: func(res *types.QueryRegionWithdrawerResponse) {
				s.Require().Equal(addr1, res.Address)
				s.Require().False(s.Keeper().CanRegionWithdraw(s.Ctx, addr0, regionId))
				s.Require().True(s.Keeper().CanRegionWithdraw(s.Ctx, addr1, regionId))
			},
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			if test.malleate != nil {
				test.malleate()
			}
			_, err := s.msgServer.GrantRegionWithdraw(s.Ctx, &types.MsgGrantRegionWithdraw{
				Creator:  test.creator,
				RegionId: test.regionId,
				Address:  test.address,
			})
			s.Require().ErrorIs(err, test.expErr)
			if test.expErr == nil && test.afterCheck != nil {
				res, err := s.queryClient.RegionWithdrawer(s.Ctx, &types.QueryRegionWithdrawerRequest{
					RegionId: test.regionId,
				})
				s.Require().NoError(err)
				test.afterCheck(res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestRevokeRegionWithdraw() {
	s.SetupTest()

	regionId := types.ExperienceRegionId
	addr0 := s.TestAccs[0].String()
	nonDao := s.TestAccs[2].String()

	grantAddr0 := func() {
		_, err := s.msgServer.GrantRegionWithdraw(s.Ctx, &types.MsgGrantRegionWithdraw{
			Creator:  s.Dao.GlobalDao,
			RegionId: regionId,
			Address:  addr0,
		})
		s.Require().NoError(err)
	}

	tests := []struct {
		name     string
		creator  string
		regionId string
		expErr   error
		malleate func()
	}{
		{
			name:     "non-DAO cannot revoke",
			creator:  nonDao,
			regionId: regionId,
			expErr:   types.ErrCheckGlobalDao,
		},
		{
			name:     "region does not exist",
			creator:  s.Dao.GlobalDao,
			regionId: "nonexistent_region",
			expErr:   types.ErrRegionNotExist,
		},
		{
			name:     "no permission record to revoke",
			creator:  s.Dao.GlobalDao,
			regionId: regionId,
			expErr:   sdkerrors.ErrKeyNotFound,
		},
		{
			name:     "revoke success",
			creator:  s.Dao.GlobalDao,
			regionId: regionId,
			malleate: grantAddr0,
			expErr:   nil,
		},
		{
			name:     "revoke twice returns error",
			creator:  s.Dao.GlobalDao,
			regionId: regionId,
			malleate: func() {
				grantAddr0()
				_, err := s.msgServer.RevokeRegionWithdraw(s.Ctx, &types.MsgRevokeRegionWithdraw{
					Creator:  s.Dao.GlobalDao,
					RegionId: regionId,
				})
				s.Require().NoError(err)
			},
			expErr: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			if test.malleate != nil {
				test.malleate()
			}
			_, err := s.msgServer.RevokeRegionWithdraw(s.Ctx, &types.MsgRevokeRegionWithdraw{
				Creator:  test.creator,
				RegionId: test.regionId,
			})
			s.Require().ErrorIs(err, test.expErr)
			if test.expErr == nil {
				res, err := s.queryClient.RegionWithdrawer(s.Ctx, &types.QueryRegionWithdrawerRequest{
					RegionId: test.regionId,
				})
				s.Require().NoError(err)
				s.Require().Empty(res.Address)
			}
		})
	}
}
