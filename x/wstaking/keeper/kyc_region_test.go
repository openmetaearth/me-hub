package keeper_test

import (
	"strings"

	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wdistri"
	"github.com/openmetaearth/me-hub/x/wmint"
	wmintTypes "github.com/openmetaearth/me-hub/x/wmint/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
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
	err = s.Keeper().KycReward(s.Ctx, kycAccount, s.meEarthValidator.Description.RegionID, s.Dao.GlobalDao)
	s.Require().NoError(err)

	// check invite address - inviter reward logic was removed
	balance := s.App.BankKeeper.GetBalance(s.Ctx, inviter, params.BaseDenom)
	s.Require().Equal(balance.Amount.String(), "0")

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

func (s *KeeperTestSuite) TestTransferKycRegionDoesNotCommitPartialStateOnError() {
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

	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}, *s.App.DistrKeeper)

	kycAccount := sdk.MustAccAddressFromBech32(s.Dao.DevOperator)
	err = s.Keeper().KycReward(s.Ctx, kycAccount, s.meEarthValidator.Description.RegionID, s.Dao.GlobalDao)
	s.Require().NoError(err)

	fromRegion, found := s.Keeper().GetRegion(s.Ctx, strings.ToLower(types.MeEarthRegionName))
	s.Require().True(found)
	toRegion, found := s.Keeper().GetRegion(s.Ctx, s.usaValidator.Description.RegionID)
	s.Require().True(found)

	delegationBefore, found := s.Keeper().GetDelegation(s.Ctx, kycAccount, sdk.ValAddress{})
	s.Require().True(found)

	fromValAddr, err := sdk.ValAddressFromBech32(fromRegion.OperatorAddress)
	s.Require().NoError(err)
	toValAddr, err := sdk.ValAddressFromBech32(toRegion.OperatorAddress)
	s.Require().NoError(err)

	fromValidatorBefore, found := s.Keeper().GetValidator(s.Ctx, fromValAddr)
	s.Require().True(found)
	toValidatorBefore, found := s.Keeper().GetValidator(s.Ctx, toValAddr)
	s.Require().True(found)

	fromRegion.DelegateAmount = sdk.ZeroInt()
	s.Keeper().SetRegion(s.Ctx, fromRegion)
	corruptedFromRegion := fromRegion

	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wmintTypes.OneDayTotalBlocks + 1).WithChainID(apptesting.TestChainID)

	err = s.Keeper().TransferKycRegion(s.Ctx, kycAccount, s.Dao.GlobalDao, s.meEarthValidator.Description.RegionID, s.usaValidator.Description.RegionID)
	s.Require().Error(err)

	delegationAfter, found := s.Keeper().GetDelegation(s.Ctx, kycAccount, sdk.ValAddress{})
	s.Require().True(found)
	s.Require().Equal(delegationBefore, delegationAfter)

	fromRegionAfter, found := s.Keeper().GetRegion(s.Ctx, strings.ToLower(types.MeEarthRegionName))
	s.Require().True(found)
	s.Require().Equal(corruptedFromRegion, fromRegionAfter)
	toRegionAfter, found := s.Keeper().GetRegion(s.Ctx, s.usaValidator.Description.RegionID)
	s.Require().True(found)
	s.Require().Equal(toRegion, toRegionAfter)

	fromValidatorAfter, found := s.Keeper().GetValidator(s.Ctx, fromValAddr)
	s.Require().True(found)
	s.Require().Equal(fromValidatorBefore, fromValidatorAfter)
	toValidatorAfter, found := s.Keeper().GetValidator(s.Ctx, toValAddr)
	s.Require().True(found)
	s.Require().Equal(toValidatorBefore, toValidatorAfter)
}
