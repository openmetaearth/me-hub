package keeper_test

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/app/apptesting"
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

	s.Ctx = s.App.BaseApp.NewContext(false).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, *s.App.DistrKeeper)

	inviter := accounts[0]
	err = s.Keeper().KycReward(s.Ctx, inviter, s.meEarthValidator.Description.RegionID, s.Dao.GlobalDao)
	s.Require().NoError(err)

	// check region DelegateAmount
	region, found := s.Keeper().GetRegion(s.Ctx, strings.ToLower(types.MeEarthRegionName))
	s.Require().True(found)
	s.Require().Equal(region.DelegateAmount.String(), types.Bonus.String())

	// check delegation for inviter (who was KYC'd)
	delegation, f := s.Keeper().GetDelegation(s.Ctx, inviter, sdk.ValAddress{})
	s.Require().NoError(f)
	s.Require().Equal(delegation.Unmovable.String(), types.Bonus.String())
	s.Require().Equal(delegation.ValidatorAddress, s.meEarthValidator.OperatorAddress)

	s.Ctx = s.App.BaseApp.NewContext(false).WithBlockHeight(wmintTypes.OneDayTotalBlocks + 1).WithChainID(apptesting.TestChainID)

	// transfer kyc region
	err = s.Keeper().TransferKycRegion(s.Ctx, inviter, s.Dao.GlobalDao, s.meEarthValidator.Description.RegionID, s.usaValidator.Description.RegionID)
	s.Require().NoError(err)

	delegation, f = s.Keeper().GetDelegation(s.Ctx, inviter, sdk.ValAddress{})
	s.Require().NoError(f)
	s.Require().Equal(delegation.Unmovable.String(), types.Bonus.String())
	s.Require().Equal(delegation.ValidatorAddress, s.usaValidator.OperatorAddress)
	s.Require().EqualValues(delegation.StartHeight, wmintTypes.OneDayTotalBlocks+1)
}
