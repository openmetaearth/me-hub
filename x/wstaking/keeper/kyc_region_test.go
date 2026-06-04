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
)

func (s *KeeperTestSuite) TestTransferKycRegion() {
	s.SetupTest()
	accounts := s.NewAccounts(1)

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

	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}, *s.App.DistrKeeper)

	kycAccount := sdk.MustAccAddressFromBech32(s.Dao.DevOperator)
	inviter := accounts[0]
	err = s.Keeper().KycReward(s.Ctx, inviter, s.meEarthValidator.Description.RegionID, s.Dao.GlobalDao)
	s.Require().NoError(err)

	// check invite address
	balance := s.App.BankKeeper.GetBalance(s.Ctx, inviter, params.BaseDenom)
	s.Require().Equal(balance.Amount.String(), types.InviteReward.String())

	// check region DelegateAmount
	region, found := s.Keeper().GetRegion(s.Ctx, strings.ToLower(types.MeEarthRegionName))
	s.Require().True(found)
	s.Require().Equal(region.DelegateAmount.String(), types.Bonus.String())

	delegation, f := s.Keeper().GetDelegation(s.Ctx, kycAccount, sdk.ValAddress{})
	s.Require().True(f)
	s.Require().Equal(delegation.Unmovable.String(), types.Bonus.String())
	s.Require().Equal(delegation.ValidatorAddress, s.meEarthValidator.OperatorAddress)

	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wmintTypes.OneDayTotalBlocks + 1).WithChainID(apptesting.TestChainID)

	// transfer kyc region
	err = s.Keeper().TransferKycRegion(s.Ctx, kycAccount, s.Dao.GlobalDao, s.meEarthValidator.Description.RegionID, s.usaValidator.Description.RegionID)
	s.Require().NoError(err)

	delegation, f = s.Keeper().GetDelegation(s.Ctx, kycAccount, sdk.ValAddress{})
	s.Require().True(f)
	s.Require().Equal(delegation.Unmovable.String(), types.Bonus.String())
	s.Require().Equal(delegation.ValidatorAddress, s.usaValidator.OperatorAddress)
	s.Require().EqualValues(delegation.StartHeight, wmintTypes.OneDayTotalBlocks+1)
}

func (s *KeeperTestSuite) TestTransferKycRegionIgnoresInactiveUnusedFixedDepositCfg() {
	s.SetupTest()

	sourceRegionID := strings.ToLower(types.MeEarthRegionName)
	targetRegionID := "usa"

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

	for _, msg := range []types.MsgNewFixedDepositCfg{
		{
			Dao:      s.Dao.GlobalDao,
			RegionId: sourceRegionID,
			Term:     30,
			Rate:     sdk.MustNewDecFromStr("0.1"),
		},
		{
			Dao:      s.Dao.GlobalDao,
			RegionId: targetRegionID,
			Term:     30,
			Rate:     sdk.MustNewDecFromStr("0.1"),
		},
		{
			Dao:      s.Dao.GlobalDao,
			RegionId: targetRegionID,
			Term:     365,
			Rate:     sdk.MustNewDecFromStr("0.12"),
		},
	} {
		_, err = s.msgServer.NewFixedDepositCfg(s.Ctx, &msg)
		s.Require().NoError(err)
	}

	_, err = s.msgServer.SetFixedDepositCfgStatus(s.Ctx, &types.MsgSetFixedDepositCfgStatus{
		Admin:    s.Dao.GlobalDao,
		RegionId: targetRegionID,
		Term:     365,
		Status:   types.RegionFixedDepositCfgStatusInactive,
	})
	s.Require().NoError(err)

	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}, *s.App.DistrKeeper)

	userAccount, _ := s.NewAccount()
	err = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx,
		mintypes.ModuleName,
		userAccount,
		sdk.Coins{sdk.NewInt64Coin(params.BaseDenom, 10000000000)},
	)
	s.Require().NoError(err)

	s.InitKyc(userAccount, "1111111111111192", sourceRegionID)
	err = s.Keeper().KycReward(s.Ctx, userAccount, sourceRegionID, s.Dao.GlobalDao)
	s.Require().NoError(err)

	_, err = s.msgServer.DoFixedDeposit(s.Ctx, &types.MsgDoFixedDeposit{
		Account:   userAccount.String(),
		Principal: sdk.NewCoin(params.BaseDenom, sdk.NewInt(100000000)),
		Term:      30,
	})
	s.Require().NoError(err)

	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wmintTypes.OneDayTotalBlocks + 1).WithChainID(apptesting.TestChainID)
	err = s.Keeper().TransferKycRegion(s.Ctx, userAccount, s.Dao.GlobalDao, sourceRegionID, targetRegionID)
	s.Require().NoError(err)

	delegation, found := s.Keeper().GetDelegation(s.Ctx, userAccount, sdk.ValAddress{})
	s.Require().True(found)
	s.Require().Equal(s.usaValidator.OperatorAddress, delegation.ValidatorAddress)
	s.Require().Equal(uint64(0), s.Keeper().GetFixedDepositCountOfCfg(s.Ctx, sourceRegionID, 30))
	s.Require().Equal(uint64(1), s.Keeper().GetFixedDepositCountOfCfg(s.Ctx, targetRegionID, 30))
	s.Require().Equal(uint64(0), s.Keeper().GetFixedDepositCountOfCfg(s.Ctx, targetRegionID, 365))
}
