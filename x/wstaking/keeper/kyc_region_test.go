package keeper_test

import (
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
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

func (s *KeeperTestSuite) TestGetRegionIdByAccountRequiresActiveDID() {
	s.SetupTest()

	account := s.TestAccs[0]
	did := "1111111111111111"
	didInfo := didtypes.DidInfo{
		Did:      did,
		Address:  account.String(),
		RegionId: types.MeEarthRegionId,
		KycLevel: didtypes.KYC_LEVEL_TWO,
		Status:   didtypes.DID_STATUS_ACTIVE,
	}
	kyc := didtypes.Credential{Did: did, Data: []byte(types.MeEarthRegionId)}

	s.App.KycKeeper.SetDID(s.Ctx, account, did)
	s.App.KycKeeper.SetDidInfo(s.Ctx, did, didInfo)
	s.App.KycKeeper.SetKYC(s.Ctx, did, kyc)

	s.Require().Equal(types.MeEarthRegionId, s.Keeper().GetRegionIdByAccount(s.Ctx, account))
	regionId, err := s.Keeper().MustGetKycRegionIdByAccount(s.Ctx, account.String())
	s.Require().NoError(err)
	s.Require().Equal(types.MeEarthRegionId, regionId)

	didInfo.Status = didtypes.DID_STATUS_INACTIVE
	s.App.KycKeeper.SetDidInfo(s.Ctx, did, didInfo)

	s.Require().Equal(strings.ToLower(types.ExperienceRegionName), s.Keeper().GetRegionIdByAccount(s.Ctx, account))
	_, err = s.Keeper().MustGetKycRegionIdByAccount(s.Ctx, account.String())
	s.Require().Error(err)
}
