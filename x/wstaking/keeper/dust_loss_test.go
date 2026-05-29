package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestReDelegationInterestDust() {
	s.SetupTest()

	// Create a new region
	newMeEarthRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newMeEarthRegion)
	s.Require().NoError(err)

	regionId := types.MeEarthRegionId
	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, regionId)
	s.Require().True(found)

	// Send initial coins to region treasure address
	err = s.App.BankKeeper.SendCoinsFromModuleToAccount(
		s.Ctx,
		mintypes.ModuleName,
		sdk.MustAccAddressFromBech32(region.GetRegionTreasureAddr()),
		sdk.Coins{sdk.NewInt64Coin(params.BaseDenom, 2000000000000)},
	)
	s.Require().NoError(err)

	// Set initial DelegateInterest in region to an integer value
	region.DelegateInterest = sdk.NewDec(100)
	s.App.StakingKeeper.SetRegion(s.Ctx, region)

	// Perform first delegation to establish staking
	msg := stakingtypes.MsgDelegate{
		DelegatorAddress: s.Dao.GlobalDao,
		ValidatorAddress: s.meEarthValidator.OperatorAddress,
		Amount:           sdk.NewCoin(params.BaseDenom, sdk.NewInt(1000000)),
	}
	_, err = s.msgServer.Delegate(s.Ctx, &msg)
	s.Require().NoError(err)

	// Advance block height to accumulate rewards with fractional part
	s.Ctx = s.Ctx.WithBlockHeight(10)

	// Read region before the second delegation
	regionBefore, found := s.App.StakingKeeper.GetRegion(s.Ctx, regionId)
	s.Require().True(found)

	// Perform another delegation to trigger rewards claiming and DelegateInterest deduction
	msg2 := stakingtypes.MsgDelegate{
		DelegatorAddress: s.Dao.GlobalDao,
		ValidatorAddress: s.meEarthValidator.OperatorAddress,
		Amount:           sdk.NewCoin(params.BaseDenom, sdk.NewInt(1000000)),
	}
	_, err = s.msgServer.Delegate(s.Ctx, &msg2)
	s.Require().NoError(err)

	// Read region after the second delegation
	regionAfter, found := s.App.StakingKeeper.GetRegion(s.Ctx, regionId)
	s.Require().True(found)

	// The difference in region.DelegateInterest must be exactly an integer.
	diff := regionBefore.DelegateInterest.Sub(regionAfter.DelegateInterest)
	s.Require().True(diff.IsInteger(), "The subtracted interest must be an integer, otherwise dust is lost. diff: %s", diff.String())
}
