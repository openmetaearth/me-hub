package keeper_test

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/wdistri"
	"github.com/openmetaearth/me-hub/x/wmint"
	wmintTypes "github.com/openmetaearth/me-hub/x/wmint/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestKycReward_WithDelegation() {
	s.SetupTest()
	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            "USA",
		OperatorAddress: s.usaValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	s.Ctx = s.App.BaseApp.NewContext(false).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, *s.App.DistrKeeper)

	userAccount, _ := s.NewAccount()
	apptesting.FundAccount(s.App, s.Ctx, userAccount, sdk.Coins{sdk.NewInt64Coin(params.BaseDenom, 1000000000000)})

	delegateAmount := sdkmath.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(params.BaseDenomUnit), nil))
	_, err = s.msgServer.Delegate(s.Ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: userAccount.String(),
		ValidatorAddress: s.experienceValidator.OperatorAddress,
		Amount:           sdk.NewCoin(params.BaseDenom, delegateAmount),
	})
	s.Require().NoError(err)

	// check experience region DelegateAmount
	expRegion, found := s.Keeper().GetRegion(s.Ctx, s.experienceValidator.Description.RegionID)
	s.Require().True(found)
	s.Require().Equal(expRegion.DelegateAmount.String(), delegateAmount.String())

	// check experience validator DelegateAmount
	valAddress, err := sdk.ValAddressFromBech32(s.experienceValidator.OperatorAddress)
	s.Require().NoError(err)
	expVal, _ := s.Keeper().GetValidator(s.Ctx, valAddress)
	s.Require().NoError(err)
	s.Require().Equal(expVal.DelegationAmount.String(), delegateAmount.String())

	expValOpAddr, err := sdk.ValAddressFromBech32(expVal.GetOperator())
	s.Require().NoError(err)
	delegation, f := s.Keeper().GetDelegation(s.Ctx, userAccount, expValOpAddr)
	s.Require().NoError(f)
	s.Require().Equal(delegation.UnMeidAmount.String(), delegateAmount.String())
	s.Require().Equal(delegation.Unmovable.String(), sdkmath.NewInt(0).String())
	s.Require().Equal(delegation.Amount.String(), sdkmath.NewInt(0).String())

	// do kyc reward
	inviter, _ := s.NewAccount()
	inviterBalanceBefore := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(inviter.String()), params.BaseDenom)
	err = s.Keeper().KycReward(s.Ctx, userAccount, s.usaValidator.Description.RegionID, s.Dao.GlobalDao)
	s.Require().NoError(err)
	err = s.Keeper().SendInviteReward(s.Ctx, inviter.String(), userAccount.String(), s.usaValidator.Description.RegionID)
	s.Require().NoError(err)

	// check invite address (balance should increase by InviteReward)
	balance := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(inviter.String()), params.BaseDenom)
	s.Require().Equal(types.InviteReward.String(), balance.Amount.Sub(inviterBalanceBefore.Amount).String())

	// after kyc reward
	// check experience region DelegateAmount
	expRegion, found = s.Keeper().GetRegion(s.Ctx, s.experienceValidator.Description.RegionID)
	s.Require().True(found)
	s.Require().Equal(sdkmath.NewInt(0).String(), expRegion.DelegateAmount.String())

	// check experience validator DelegateAmount
	expVal, _ = s.Keeper().GetValidator(s.Ctx, valAddress)
	s.Require().NoError(err)
	s.Require().Equal(sdkmath.NewInt(0).String(), expVal.DelegationAmount.String())

	// check usa region DelegateAmount
	usaRegion, found := s.Keeper().GetRegion(s.Ctx, s.usaValidator.Description.RegionID)
	s.Require().True(found)
	s.Require().Equal(delegateAmount.Add(types.Bonus).String(), usaRegion.DelegateAmount.String())

	// check usa validator DelegateAmount
	usaValAddress, err := sdk.ValAddressFromBech32(s.usaValidator.OperatorAddress)
	usaVal, _ := s.Keeper().GetValidator(s.Ctx, usaValAddress)
	s.Require().NoError(err)
	s.Require().Equal(delegateAmount.String(), usaVal.DelegationAmount.String())

	delegation, f = s.Keeper().GetDelegation(s.Ctx, userAccount, usaValAddress)
	s.Require().NoError(f)
	s.Require().Equal(sdkmath.NewInt(0).String(), delegation.UnMeidAmount.String())
	s.Require().Equal(types.Bonus.String(), delegation.Unmovable.String())
	s.Require().Equal(delegateAmount.String(), delegation.Amount.String())
}

func (s *KeeperTestSuite) TestKycReward_WithoutDelegation() {
	s.SetupTest()

	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            "USA",
		OperatorAddress: s.usaValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	s.Ctx = s.App.BaseApp.NewContext(false).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, *s.App.DistrKeeper)

	inviter, _ := s.NewAccount()
	err = s.Keeper().KycReward(s.Ctx, inviter, s.usaValidator.Description.RegionID, s.Dao.GlobalDao)
	s.Require().NoError(err)

	// check region DelegateAmount
	region, found := s.Keeper().GetRegion(s.Ctx, "usa")
	s.Require().True(found)
	s.Require().Equal(region.DelegateAmount.String(), types.Bonus.String())

	delegation, f := s.Keeper().GetDelegation(s.Ctx, inviter, sdk.ValAddress{})
	s.Require().NoError(f)
	s.Require().Equal(delegation.Unmovable.String(), types.Bonus.String())
}

func (s *KeeperTestSuite) TestRemoveKycReward() {
	s.SetupTest()

	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            "USA",
		OperatorAddress: s.usaValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	s.Ctx = s.App.BaseApp.NewContext(false).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, *s.App.DistrKeeper)

	inviter, _ := s.NewAccount()
	err = s.Keeper().KycReward(s.Ctx, inviter, s.usaValidator.Description.RegionID, s.Dao.GlobalDao)
	s.Require().NoError(err)

	// remove kyc
	err = s.Keeper().RemoveKycReward(s.Ctx, inviter, s.usaValidator.Description.RegionID)
	s.Require().NoError(err)

	// check region DelegateAmount
	region, found := s.Keeper().GetRegion(s.Ctx, "usa")
	s.Require().True(found)
	s.Require().Equal(region.DelegateAmount.String(), sdkmath.NewInt(0).String())

	_, f := s.Keeper().GetDelegation(s.Ctx, inviter, sdk.ValAddress{})
	s.Require().Error(f)
}

func (s *KeeperTestSuite) TestRemoveKycReward_WithDelegation() {
	s.SetupTest()

	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            "USA",
		OperatorAddress: s.usaValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	s.Ctx = s.App.BaseApp.NewContext(false).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, *s.App.DistrKeeper)

	// create user account
	userAccount, _ := s.NewAccount()
	apptesting.FundAccount(s.App, s.Ctx, userAccount, sdk.Coins{sdk.NewInt64Coin(params.BaseDenom, 1000000000000)})

	did := "1111111111111101"
	s.App.DidKeeper.SetDID(s.Ctx, userAccount, did)
	s.App.KycKeeper.SetKYC(s.Ctx, did, didtypes.Credential{
		Did:  did,
		Sid:  "",
		Hash: "",
		Uri:  "",
		Data: []byte(s.usaValidator.Description.RegionID),
	})

	inviter, _ := s.NewAccount()
	inviterBalanceBefore := s.App.BankKeeper.GetBalance(s.Ctx, inviter, params.BaseDenom)
	err = s.Keeper().KycReward(s.Ctx, userAccount, s.usaValidator.Description.RegionID, s.Dao.GlobalDao)
	s.Require().NoError(err)
	err = s.Keeper().SendInviteReward(s.Ctx, inviter.String(), userAccount.String(), s.usaValidator.Description.RegionID)
	s.Require().NoError(err)

	// check invite address (balance should increase by InviteReward)
	balance := s.App.BankKeeper.GetBalance(s.Ctx, inviter, params.BaseDenom)
	s.Require().Equal(types.InviteReward.String(), balance.Amount.Sub(inviterBalanceBefore.Amount).String())

	// check delegation after kyc
	del, f := s.Keeper().GetDelegation(s.Ctx, userAccount, sdk.ValAddress{})
	s.Require().NoError(f)
	s.Require().Equal(sdkmath.NewInt(0).String(), del.Amount.String())
	s.Require().Equal(types.Bonus.String(), del.Unmovable.String())
	s.Require().Equal(sdkmath.NewInt(0).String(), del.UnMeidAmount.String())

	// check region DelegateAmount
	expRegion, found := s.Keeper().GetRegion(s.Ctx, s.experienceValidator.Description.RegionID)
	s.Require().True(found)
	s.Require().Equal(sdkmath.NewInt(0).String(), expRegion.DelegateAmount.String())

	// delegate
	delegateAmount := sdkmath.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(params.BaseDenomUnit+1), nil))
	_, err = s.msgServer.Delegate(s.Ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: userAccount.String(),
		ValidatorAddress: s.usaValidator.OperatorAddress,
		Amount:           sdk.NewCoin(params.BaseDenom, delegateAmount),
	})
	s.Require().NoError(err)

	// check delegation after delegate
	del, f = s.Keeper().GetDelegation(s.Ctx, userAccount, sdk.ValAddress{})
	s.Require().NoError(f)
	s.Require().Equal(delegateAmount.String(), del.Amount.String())
	s.Require().Equal(types.Bonus.String(), del.Unmovable.String())
	s.Require().Equal(sdkmath.NewInt(0).String(), del.UnMeidAmount.String())

	// remove kyc
	err = s.Keeper().RemoveKycReward(s.Ctx, userAccount, s.usaValidator.Description.RegionID)
	s.Require().ErrorContains(err, types.ErrRemoveKyc.Error())
}

func (s *KeeperTestSuite) TestRemoveKycReward_WithFixedDeposit() {
	s.SetupTest()

	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            "USA",
		OperatorAddress: s.usaValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	s.Ctx = s.App.BaseApp.NewContext(false).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, *s.App.DistrKeeper)

	// create user account
	userAccount, _ := s.NewAccount()
	apptesting.FundAccount(s.App, s.Ctx, userAccount, sdk.Coins{sdk.NewInt64Coin(params.BaseDenom, 1000000000000)})

	did := "1111111111111101"
	s.App.DidKeeper.SetDID(s.Ctx, userAccount, did)
	s.App.KycKeeper.SetKYC(s.Ctx, did, didtypes.Credential{
		Did:  did,
		Sid:  "",
		Hash: "",
		Uri:  "",
		Data: []byte(s.usaValidator.Description.RegionID),
	})

	inviter, _ := s.NewAccount()
	_ = inviter // inviter account created for potential SendInviteReward use
	err = s.Keeper().KycReward(s.Ctx, userAccount, s.usaValidator.Description.RegionID, s.Dao.GlobalDao)
	s.Require().NoError(err)

	// check delegation after kyc (userAccount was KYC'd)
	del, f := s.Keeper().GetDelegation(s.Ctx, userAccount, sdk.ValAddress{})
	s.Require().NoError(f)
	s.Require().Equal(sdkmath.NewInt(0).String(), del.Amount.String())
	s.Require().Equal(types.Bonus.String(), del.Unmovable.String())
	s.Require().Equal(sdkmath.NewInt(0).String(), del.UnMeidAmount.String())

	// check region DelegateAmount
	expRegion, found := s.Keeper().GetRegion(s.Ctx, s.experienceValidator.Description.RegionID)
	s.Require().True(found)
	s.Require().Equal(sdkmath.NewInt(0).String(), expRegion.DelegateAmount.String())

	// delegate
	delegateAmount := sdkmath.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(params.BaseDenomUnit+1), nil))
	_, err = s.msgServer.Delegate(s.Ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: userAccount.String(),
		ValidatorAddress: s.usaValidator.OperatorAddress,
		Amount:           sdk.NewCoin(params.BaseDenom, delegateAmount),
	})
	s.Require().NoError(err)

	// check delegation after delegate
	del, f = s.Keeper().GetDelegation(s.Ctx, userAccount, sdk.ValAddress{})
	s.Require().NoError(f)
	s.Require().Equal(delegateAmount.String(), del.Amount.String())
	s.Require().Equal(types.Bonus.String(), del.Unmovable.String())
	s.Require().Equal(sdkmath.NewInt(0).String(), del.UnMeidAmount.String())

	// remove kyc
	err = s.Keeper().RemoveKycReward(s.Ctx, userAccount, s.usaValidator.Description.RegionID)
	s.Require().ErrorContains(err, types.ErrRemoveKyc.Error())
}
