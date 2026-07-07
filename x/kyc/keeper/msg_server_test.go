package keeper_test

import (
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/kyc/types"
	"github.com/openmetaearth/me-hub/x/wdistri"
	"github.com/openmetaearth/me-hub/x/wmint"
	wmintTypes "github.com/openmetaearth/me-hub/x/wmint/types"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestApprove() {
	s.Ctx = s.App.BaseApp.NewContext(false).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, *s.App.DistrKeeper)

	did := "1111111111111111"
	kycAccount, newUserPubkey := s.NewAccountStr()
	inviter, _ := s.NewAccount()
	msg := &types.MsgApprove{
		Issuer:   s.Dao.GlobalDao,
		Did:      did,
		RegionId: strings.ToLower(wstakingtypes.MeEarthRegionName),
		Address:  kycAccount.String(),
		Pubkey:   newUserPubkey,
		Uri:      "http://127.0.0.1/8001",
		Hash:     "aaaa",
		Inviter:  inviter.String(),
		Level:    2,
	}
	_, err := s.msgServer.Approve(s.Ctx, msg)
	s.Require().NoError(err)

	// check invite address - inviter starts with 1e9 pre-funded, verify delta = InviteReward
	balance := s.App.BankKeeper.GetBalance(s.Ctx, inviter, params.BaseDenom)
	s.Require().Equal(wstakingtypes.InviteReward.String(), balance.Amount.Sub(sdkmath.NewInt(1_000_000_000)).String())

	// check region DelegateAmount
	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, strings.ToLower(wstakingtypes.MeEarthRegionName))
	s.Require().True(found)
	s.Require().Equal(region.DelegateAmount.String(), wstakingtypes.Bonus.String())

	valAddress, err := sdk.ValAddressFromBech32(region.OperatorAddress)
	s.Require().NoError(err)

	// check user's delegation
	delegation, f := s.App.StakingKeeper.GetDelegation(s.Ctx, kycAccount, valAddress)
	s.Require().NoError(f)
	s.Require().Equal(delegation.Unmovable.String(), wstakingtypes.Bonus.String())

	// check kyc
	kyc, kycFound := s.Keeper().GetKYC(s.Ctx, did)
	s.Require().True(kycFound)
	s.Require().Equal(msg.Uri, kyc.Uri)
	s.Require().Equal(msg.Hash, kyc.Hash)
}

func (s *KeeperTestSuite) TestRemove() {
	s.SetupTest()

	s.Ctx = s.App.BaseApp.NewContext(false).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, *s.App.DistrKeeper)

	kycAccount, newUserPubkey := s.NewAccountStr()
	did := "1111111111111111"
	inviter, _ := s.NewAccount()
	msg := &types.MsgApprove{
		Issuer:   s.Dao.GlobalDao,
		Did:      did,
		RegionId: strings.ToLower(wstakingtypes.MeEarthRegionName),
		Address:  kycAccount.String(),
		Pubkey:   newUserPubkey,
		Uri:      "http://127.0.0.1/8001",
		Hash:     "aaaa",
		Inviter:  inviter.String(),
		Level:    2,
	}
	_, err := s.msgServer.Approve(s.Ctx, msg)
	s.Require().NoError(err)

	// check invite address - inviter starts with 1e9 pre-funded, verify delta = InviteReward
	balance := s.App.BankKeeper.GetBalance(s.Ctx, inviter, params.BaseDenom)
	s.Require().Equal(wstakingtypes.InviteReward.String(), balance.Amount.Sub(sdkmath.NewInt(1_000_000_000)).String())

	// check kyc
	kyc, f := s.Keeper().GetKYC(s.Ctx, did)
	s.Require().True(f)
	s.Require().Equal(msg.Uri, kyc.Uri)
	s.Require().Equal(msg.Hash, kyc.Hash)

	// remove kyc
	_, err = s.msgServer.Remove(s.Ctx, &types.MsgRemove{
		Issuer: s.Dao.GlobalDao,
		Did:    did,
	})
	s.Require().NoError(err)

	// check region DelegateAmount
	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, strings.ToLower(wstakingtypes.MeEarthRegionName))
	s.Require().True(found)
	s.Require().Equal(region.DelegateAmount.String(), sdkmath.NewInt(0).String())

	_, delErr := s.App.StakingKeeper.GetDelegation(s.Ctx, kycAccount, sdk.ValAddress{})
	s.Require().Error(delErr)

	// check kyc
	_, kycExists := s.Keeper().GetKYC(s.Ctx, did)
	s.Require().False(kycExists)
}

func (s *KeeperTestSuite) TestUpdate() {
	s.SetupTest()

	s.Ctx = s.App.BaseApp.NewContext(false).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, *s.App.DistrKeeper)

	kycAccount, newUserPubkey := s.NewAccountStr()
	did := "1111111111111111"
	inviter, _ := s.NewAccount()
	_, err := s.msgServer.Approve(s.Ctx, &types.MsgApprove{
		Issuer:   s.Dao.GlobalDao,
		Did:      did,
		RegionId: strings.ToLower(wstakingtypes.MeEarthRegionName),
		Address:  kycAccount.String(),
		Pubkey:   newUserPubkey,
		Uri:      "http://127.0.0.1/8001",
		Hash:     "aaaa",
		Inviter:  inviter.String(),
		Level:    2,
	})
	s.Require().NoError(err)

	// check invite address - inviter starts with 1e9 pre-funded, verify delta = InviteReward
	balance := s.App.BankKeeper.GetBalance(s.Ctx, inviter, params.BaseDenom)
	s.Require().Equal(wstakingtypes.InviteReward.String(), balance.Amount.Sub(sdkmath.NewInt(1_000_000_000)).String())

	// check region DelegateAmount
	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, strings.ToLower(wstakingtypes.MeEarthRegionName))
	s.Require().True(found)
	s.Require().Equal(region.DelegateAmount.String(), wstakingtypes.Bonus.String())

	meEarthValAddr, _ := sdk.ValAddressFromBech32(s.meEarthValidator.GetOperator())
	delegation, f := s.App.StakingKeeper.GetDelegation(s.Ctx, kycAccount, meEarthValAddr)
	s.Require().NoError(f)
	s.Require().Equal(delegation.Unmovable.String(), wstakingtypes.Bonus.String())
	s.Require().Equal(delegation.ValidatorAddress, s.meEarthValidator.OperatorAddress)

	s.Ctx = s.App.BaseApp.NewContext(false).WithBlockHeight(wmintTypes.OneDayTotalBlocks + 1).WithChainID(apptesting.TestChainID)
	// transfer kyc region
	_, err = s.msgServer.Update(s.Ctx, &types.MsgUpdate{
		Issuer:   s.Dao.GlobalDao,
		Did:      did,
		RegionId: "usa",
		Uri:      "http://127.0.0.1/8001",
		Hash:     "aaaa",
	})
	s.Require().NoError(err)

	usaValAddr, _ := sdk.ValAddressFromBech32(s.usaValidator.GetOperator())
	delegation, f = s.App.StakingKeeper.GetDelegation(s.Ctx, kycAccount, usaValAddr)
	s.Require().NoError(f)
	s.Require().Equal(delegation.Unmovable.String(), wstakingtypes.Bonus.String())
	s.Require().Equal(s.usaValidator.OperatorAddress, delegation.ValidatorAddress)
	s.Require().EqualValues(delegation.StartHeight, wmintTypes.OneDayTotalBlocks+1)
}
