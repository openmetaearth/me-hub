package keeper_test

import (
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/st-chain/me-hub/app/apptesting"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/x/wdistri"
	"github.com/st-chain/me-hub/x/wmint"
	wmintTypes "github.com/st-chain/me-hub/x/wmint/types"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestKycReward_WithoutDelegation() {
	s.SetupTest()

	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            "USA",
		OperatorAddress: s.usaValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}, *s.App.DistrKeeper)

	kycAccount := sdk.MustAccAddressFromBech32(s.Dao.DevOperator)
	invitee := s.Dao.GlobalDao
	err = s.Keeper().KycReward(s.Ctx, kycAccount, invitee, s.usaValidator.Description.RegionId, s.Dao.GlobalDao)
	s.Require().NoError(err)

	// check invite address
	balance := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(invitee), params.BaseDenom)
	s.Require().Equal(balance.Amount.String(), types.InviteReward.String())

	// check region DelegateAmount
	region, found := s.Keeper().GetRegion(s.Ctx, "usa")
	s.Require().True(found)
	s.Require().Equal(region.DelegateAmount.String(), types.Bonus.String())

	delegation, f := s.Keeper().GetDelegation(s.Ctx, kycAccount, sdk.ValAddress{})
	s.Require().True(f)
	s.Require().Equal(delegation.Unmovable.String(), types.Bonus.String())
}

func (s *KeeperTestSuite) TestRemoveKycReward_WithoutDelegation() {
	s.SetupTest()

	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            "USA",
		OperatorAddress: s.usaValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	// must have experience region
	newRegion = types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.ExperienceRegion,
		OperatorAddress: s.experienceValidator.OperatorAddress,
	}
	_, err = s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}, *s.App.DistrKeeper)

	kycAccount := sdk.MustAccAddressFromBech32(s.Dao.DevOperator)
	invitee := s.Dao.GlobalDao
	err = s.Keeper().KycReward(s.Ctx, kycAccount, invitee, s.usaValidator.Description.RegionId, s.Dao.GlobalDao)
	s.Require().NoError(err)

	// check invite address
	balance := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(invitee), params.BaseDenom)
	s.Require().Equal(balance.Amount.String(), types.InviteReward.String())

	// remove kyc
	err = s.Keeper().RemoveKycReward(s.Ctx, kycAccount, s.usaValidator.Description.RegionId)
	s.Require().NoError(err)

	// check region DelegateAmount
	region, found := s.Keeper().GetRegion(s.Ctx, "usa")
	s.Require().True(found)
	s.Require().Equal(region.DelegateAmount.String(), sdk.NewInt(0).String())

	_, f := s.Keeper().GetDelegation(s.Ctx, kycAccount, sdk.ValAddress{})
	s.Require().False(f)
}

func (s *KeeperTestSuite) TestKycReward_WithDelegation() {
	s.SetupTest()

	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.ExperienceRegion,
		OperatorAddress: s.experienceValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}, *s.App.DistrKeeper)

	// TODO: msg delegate
	_, err = s.msgServer.Delegate(s.Ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: "",
		ValidatorAddress: "",
		Amount:           sdk.Coin{},
	})
	s.Require().NoError(err)

	// do kyc reward
	kycAccount := sdk.MustAccAddressFromBech32(s.Dao.DevOperator)
	invitee := s.Dao.GlobalDao
	err = s.Keeper().KycReward(s.Ctx, kycAccount, invitee, s.usaValidator.Description.RegionId, s.Dao.GlobalDao)
	s.Require().NoError(err)

	// check invite address
	balance := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(invitee), params.BaseDenom)
	s.Require().Equal(balance.Amount.String(), types.InviteReward.String())
}
