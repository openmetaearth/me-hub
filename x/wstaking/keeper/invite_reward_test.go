package keeper_test

import (
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wdistri"
	"github.com/openmetaearth/me-hub/x/wmint"
	wmintTypes "github.com/openmetaearth/me-hub/x/wmint/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"math/big"
)

// TestSendInviteReward_NoReplay verifies that calling SendInviteReward twice
// for the same invitee does NOT pay the inviter a second time (issue #11).
func (s *KeeperTestSuite) TestSendInviteReward_NoReplay() {
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

	// Fund inviter account (the one who will receive the reward)
	inviter, _ := s.NewAccount()
	// Fund invitee account (the one who gets KYC'd)
	invitee, _ := s.NewAccount()
	err = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, mintypes.ModuleName, invitee, sdk.Coins{sdk.NewInt64Coin(params.BaseDenom, 1000000000000)})
	s.Require().NoError(err)

	// Delegate some tokens for invitee in experience region so KycReward works
	delegateAmount := sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(params.BaseDenomUnit), nil))
	_, err = s.msgServer.Delegate(s.Ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: invitee.String(),
		ValidatorAddress: s.experienceValidator.OperatorAddress,
		Amount:           sdk.NewCoin(params.BaseDenom, delegateAmount),
	})
	s.Require().NoError(err)

	// Approve KYC for invitee
	err = s.Keeper().KycReward(s.Ctx, invitee, s.usaValidator.Description.RegionID, s.Dao.GlobalDao)
	s.Require().NoError(err)

	// First call: should succeed and pay the inviter
	err = s.Keeper().SendInviteReward(s.Ctx, inviter.String(), invitee.String(), s.usaValidator.Description.RegionID)
	s.Require().NoError(err)

	balanceAfterFirst := s.App.BankKeeper.GetBalance(s.Ctx, inviter, params.BaseDenom)
	s.Require().Equal(types.InviteReward.String(), balanceAfterFirst.Amount.String(),
		"inviter should receive exactly one InviteReward after first call")

	// Verify HasInviterReward returns true now
	s.Require().True(s.Keeper().HasInviterReward(s.Ctx, invitee.String()),
		"HasInviterReward should return true for the invitee after reward is paid")

	// Second call: should be a no-op (return nil, no extra payment)
	err = s.Keeper().SendInviteReward(s.Ctx, inviter.String(), invitee.String(), s.usaValidator.Description.RegionID)
	s.Require().NoError(err)

	balanceAfterSecond := s.App.BankKeeper.GetBalance(s.Ctx, inviter, params.BaseDenom)
	s.Require().Equal(types.InviteReward.String(), balanceAfterSecond.Amount.String(),
		"inviter balance must NOT increase on second call — replay exploit is fixed")
}

// TestSendInviteReward_EmptyInviter verifies that an empty inviter string
// is treated as a no-op without error.
func (s *KeeperTestSuite) TestSendInviteReward_EmptyInviter() {
	s.SetupTest()
	err := s.Keeper().SendInviteReward(s.Ctx, "", "invitee123", "nonexistent")
	s.Require().NoError(err)
}

// TestSetAndHasInviterReward verifies the low-level Set/Has helpers.
func (s *KeeperTestSuite) TestSetAndHasInviterReward() {
	s.SetupTest()

	addr := "cosmos1testaddr"
	s.Require().False(s.Keeper().HasInviterReward(s.Ctx, addr),
		"HasInviterReward should be false before Set")

	s.Keeper().SetInviterReward(s.Ctx, addr)
	s.Require().True(s.Keeper().HasInviterReward(s.Ctx, addr),
		"HasInviterReward should be true after Set")
}
