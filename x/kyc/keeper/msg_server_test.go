package keeper_test

import (
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/kyc/types"
	"github.com/openmetaearth/me-hub/x/wdistri"
	"github.com/openmetaearth/me-hub/x/wmint"
	wmintTypes "github.com/openmetaearth/me-hub/x/wmint/types"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
	"strings"
)

func (s *KeeperTestSuite) TestApprove() {
	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}, *s.App.DistrKeeper)

	did := "1111111111111111"
	kycAccount, newUserPubkey := s.NewAccount()
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

	// check invite address
	balance := s.App.BankKeeper.GetBalance(s.Ctx, inviter, params.BaseDenom)
	s.Require().Equal(balance.Amount.String(), wstakingtypes.InviteReward.String())

	// check region DelegateAmount
	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, strings.ToLower(wstakingtypes.MeEarthRegionName))
	s.Require().True(found)
	s.Require().Equal(region.DelegateAmount.String(), wstakingtypes.Bonus.String())

	valAddress, err := sdk.ValAddressFromBech32(region.OperatorAddress)
	s.Require().NoError(err)

	// check user's delegation
	delegation, f := s.App.StakingKeeper.GetDelegation(s.Ctx, kycAccount, valAddress)
	s.Require().True(f)
	s.Require().Equal(delegation.Unmovable.String(), wstakingtypes.Bonus.String())

	// check kyc
	kyc, f := s.Keeper().GetKYC(s.Ctx, did)
	s.Require().True(f)
	s.Require().Equal(msg.Uri, kyc.Uri)
	s.Require().Equal(msg.Hash, kyc.Hash)
}

func (s *KeeperTestSuite) TestRemove() {
	s.SetupTest()

	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}, *s.App.DistrKeeper)

	kycAccount, newUserPubkey := s.NewAccount()
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

	// check invite address
	balance := s.App.BankKeeper.GetBalance(s.Ctx, inviter, params.BaseDenom)
	s.Require().Equal(balance.Amount.String(), wstakingtypes.InviteReward.String())

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
	s.Require().Equal(region.DelegateAmount.String(), sdk.NewInt(0).String())

	_, f = s.App.StakingKeeper.GetDelegation(s.Ctx, kycAccount, sdk.ValAddress{})
	s.Require().False(f)

	// check kyc
	_, f = s.Keeper().GetKYC(s.Ctx, did)
	s.Require().False(f)
}

func (s *KeeperTestSuite) TestKycRewardsCannotBePaidTwiceByRemoveAndReapprove() {
	s.SetupTest()

	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}, *s.App.DistrKeeper)

	kycAccount, newUserPubkey := s.NewAccount()
	did := "1111111111111111"
	regionID := strings.ToLower(wstakingtypes.MeEarthRegionName)

	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, regionID)
	s.Require().True(found)
	valAddress, err := sdk.ValAddressFromBech32(region.OperatorAddress)
	s.Require().NoError(err)
	validator, found := s.App.StakingKeeper.GetValidator(s.Ctx, valAddress)
	s.Require().True(found)

	ownerAddress := sdk.MustAccAddressFromBech32(validator.OwnerAddress)
	committeeAddress := sdk.MustAccAddressFromBech32(s.Dao.DevOperator)
	ownerBefore := s.App.BankKeeper.GetBalance(s.Ctx, ownerAddress, params.BaseDenom).Amount
	committeeBefore := s.App.BankKeeper.GetBalance(s.Ctx, committeeAddress, params.BaseDenom).Amount

	approve := func() {
		_, err = s.msgServer.Approve(s.Ctx, &types.MsgApprove{
			Issuer:   s.Dao.GlobalDao,
			Did:      did,
			RegionId: regionID,
			Address:  kycAccount.String(),
			Pubkey:   newUserPubkey,
			Uri:      "http://127.0.0.1/8001",
			Hash:     "aaaa",
			Level:    2,
		})
		s.Require().NoError(err)
	}

	approve()
	_, err = s.msgServer.Remove(s.Ctx, &types.MsgRemove{
		Issuer: s.Dao.GlobalDao,
		Did:    did,
	})
	s.Require().NoError(err)
	approve()

	ownerAfter := s.App.BankKeeper.GetBalance(s.Ctx, ownerAddress, params.BaseDenom).Amount
	committeeAfter := s.App.BankKeeper.GetBalance(s.Ctx, committeeAddress, params.BaseDenom).Amount

	s.Require().Equal(wstakingtypes.ValidatorReward.String(), ownerAfter.Sub(ownerBefore).String())
	s.Require().Equal(wstakingtypes.CommitteeReward.String(), committeeAfter.Sub(committeeBefore).String())
}

func (s *KeeperTestSuite) TestUpdate() {
	s.SetupTest()

	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}, *s.App.DistrKeeper)

	kycAccount, newUserPubkey := s.NewAccount()
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

	// check invite address
	balance := s.App.BankKeeper.GetBalance(s.Ctx, inviter, params.BaseDenom)
	s.Require().Equal(balance.Amount.String(), wstakingtypes.InviteReward.String())

	// check region DelegateAmount
	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, strings.ToLower(wstakingtypes.MeEarthRegionName))
	s.Require().True(found)
	s.Require().Equal(region.DelegateAmount.String(), wstakingtypes.Bonus.String())

	delegation, f := s.App.StakingKeeper.GetDelegation(s.Ctx, kycAccount, s.meEarthValidator.GetOperator())
	s.Require().True(f)
	s.Require().Equal(delegation.Unmovable.String(), wstakingtypes.Bonus.String())
	s.Require().Equal(delegation.ValidatorAddress, s.meEarthValidator.OperatorAddress)

	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wmintTypes.OneDayTotalBlocks + 1).WithChainID(apptesting.TestChainID)
	// transfer kyc region
	_, err = s.msgServer.Update(s.Ctx, &types.MsgUpdate{
		Issuer:   s.Dao.GlobalDao,
		Did:      did,
		RegionId: "usa",
		Uri:      "http://127.0.0.1/8001",
		Hash:     "aaaa",
	})
	s.Require().NoError(err)

	delegation, f = s.App.StakingKeeper.GetDelegation(s.Ctx, kycAccount, s.usaValidator.GetOperator())
	s.Require().True(f)
	s.Require().Equal(delegation.Unmovable.String(), wstakingtypes.Bonus.String())
	s.Require().Equal(s.usaValidator.OperatorAddress, delegation.ValidatorAddress)
	s.Require().EqualValues(delegation.StartHeight, wmintTypes.OneDayTotalBlocks+1)
}
