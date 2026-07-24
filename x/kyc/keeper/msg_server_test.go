package keeper_test

import sdkmath "cosmossdk.io/math"

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/kyc/types"
	"github.com/openmetaearth/me-hub/x/wdistri"
	"github.com/openmetaearth/me-hub/x/wmint"
	wminttypes "github.com/openmetaearth/me-hub/x/wmint/types"
	"github.com/openmetaearth/me-hub/x/wstaking"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestApprove() {
	s.Ctx = s.App.BaseApp.NewContext(false).WithBlockHeight(wminttypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, *s.App.DistrKeeper)

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
	delegation, err := s.App.StakingKeeper.GetDelegation(s.Ctx, kycAccount, valAddress)
	s.Require().NoError(err)
	s.Require().Equal(delegation.Unmovable.String(), wstakingtypes.Bonus.String())

	// check kyc
	kyc, f := s.Keeper().GetKYC(s.Ctx, did)
	s.Require().True(f)
	s.Require().Equal(msg.Uri, kyc.Uri)
	s.Require().Equal(msg.Hash, kyc.Hash)
}

func (s *KeeperTestSuite) setupApproveCtx() {
	s.Ctx = s.App.BaseApp.NewContext(false).WithBlockHeight(wminttypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, *s.App.DistrKeeper)
}

func (s *KeeperTestSuite) TestApproveRejectsSubAccountAddress() {
	s.setupApproveCtx()

	s.Run("sub account already registered in map", func() {
		subAddr, subPubkey := s.newEthSubAccount()
		const parentDid = "parent-did-0001"
		s.App.DidKeeper.SetSubAccountDidMap(s.Ctx, subAddr.String(), parentDid)

		inviter, _ := s.NewAccount()
		_, err := s.msgServer.Approve(s.Ctx, &types.MsgApprove{
			Issuer:   s.Dao.GlobalDao,
			Did:      "2222222222222222",
			RegionId: strings.ToLower(wstakingtypes.MeEarthRegionName),
			Address:  subAddr.String(),
			Pubkey:   subPubkey,
			Uri:      "http://127.0.0.1/8001",
			Hash:     "aaaa",
			Inviter:  inviter.String(),
			Level:    2,
		})
		s.Require().ErrorIs(err, didtypes.ErrSubAccountAlreadyRegistered)

		_, found := s.Keeper().GetDID(s.Ctx, subAddr)
		s.Require().False(found)
	})

	s.Run("address bound via create sub account", func() {
		const parentDid = "3333333333333333"
		userAddr, userPubkey, userPrivKey := s.newSecp256k1UserAccount()
		s.setupActiveKyc(userAddr, userPubkey, parentDid)

		subAddr, subPubkey := s.newEthSubAccountFromPrivKey(userPrivKey)
		_, err := s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:          userAddr.String(),
			SubAccount:       subAddr.String(),
			SubAccountPubkey: subPubkey,
		})
		s.Require().NoError(err)

		inviter, _ := s.NewAccount()
		_, err = s.msgServer.Approve(s.Ctx, &types.MsgApprove{
			Issuer:   s.Dao.GlobalDao,
			Did:      "4444444444444444",
			RegionId: strings.ToLower(wstakingtypes.MeEarthRegionName),
			Address:  subAddr.String(),
			Pubkey:   subPubkey,
			Uri:      "http://127.0.0.1/8001",
			Hash:     "bbbb",
			Inviter:  inviter.String(),
			Level:    2,
		})
		s.Require().ErrorIs(err, didtypes.ErrSubAccountAlreadyRegistered)

		_, found := s.Keeper().GetDID(s.Ctx, subAddr)
		s.Require().False(found)
	})
}

func (s *KeeperTestSuite) TestCreateSBTRejectsDuplicateDid() {
	s.SetupTest()

	s.Ctx = s.App.BaseApp.NewContext(false).WithBlockHeight(wminttypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, *s.App.DistrKeeper)

	did := "duplicate-sbt-did"
	kycAccount, newUserPubkey := s.NewAccount()
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

	createMsg := &types.MsgCreateSBT{
		Issuer:  s.Dao.GlobalDao,
		Did:     did,
		Uri:     "sbt-uri-1",
		UriHash: "sbt-hash-1",
		Data:    []byte("sbt-data-1"),
	}
	_, err = s.msgServer.CreateSBT(s.Ctx, createMsg)
	s.Require().NoError(err)

	_, err = s.msgServer.CreateSBT(s.Ctx, &types.MsgCreateSBT{
		Issuer:  s.Dao.GlobalDao,
		Did:     did,
		Uri:     "sbt-uri-2",
		UriHash: "sbt-hash-2",
		Data:    []byte("sbt-data-2"),
	})
	s.Require().ErrorIs(err, types.ErrSbtExists)

	sbt, found := s.Keeper().GetSBT(s.Ctx, did)
	s.Require().True(found)
	s.Require().Equal(createMsg.Uri, sbt.Uri)
	s.Require().Equal(createMsg.UriHash, sbt.UriHash)
}

func (s *KeeperTestSuite) TestRemove() {
	s.SetupTest()

	s.Ctx = s.App.BaseApp.NewContext(false).WithBlockHeight(wminttypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, *s.App.DistrKeeper)

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
	s.Require().Equal(region.DelegateAmount.String(), sdkmath.NewInt(0).String())

	_, err = s.App.StakingKeeper.GetDelegation(s.Ctx, kycAccount, sdk.ValAddress{})
	s.Require().Error(err)

	// check kyc
	_, f = s.Keeper().GetKYC(s.Ctx, did)
	s.Require().Error(err)
}

func (s *KeeperTestSuite) TestUpdate() {
	s.SetupTest()

	s.Ctx = s.App.BaseApp.NewContext(false).WithBlockHeight(wminttypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, *s.App.DistrKeeper)

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

	delegation, err := s.App.StakingKeeper.GetDelegation(s.Ctx, kycAccount, sdk.ValAddress(sdk.MustAccAddressFromBech32(s.meEarthValidator.OperatorAddress)))
	s.Require().NoError(err)
	s.Require().Equal(delegation.Unmovable.String(), wstakingtypes.Bonus.String())
	s.Require().Equal(delegation.ValidatorAddress, s.meEarthValidator.OperatorAddress)

	s.Ctx = s.App.BaseApp.NewContext(false).WithBlockHeight(wminttypes.OneDayTotalBlocks + 1).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wstaking.BeginBlock(s.Ctx, s.App.StakingKeeper)
	// transfer kyc region
	_, err = s.msgServer.Update(s.Ctx, &types.MsgUpdate{
		Issuer:   s.Dao.GlobalDao,
		Did:      did,
		RegionId: "usa",
		Uri:      "http://127.0.0.1/8001",
		Hash:     "aaaa",
	})
	s.Require().NoError(err)

	delegation, err = s.App.StakingKeeper.GetDelegation(s.Ctx, kycAccount, sdk.ValAddress(sdk.MustAccAddressFromBech32(s.usaValidator.OperatorAddress)))
	s.Require().NoError(err)
	s.Require().Equal(delegation.Unmovable.String(), wstakingtypes.Bonus.String())
	s.Require().Equal(s.usaValidator.OperatorAddress, delegation.ValidatorAddress)
	s.Require().EqualValues(delegation.StartHeight, wminttypes.OneDayTotalBlocks+1)
}
