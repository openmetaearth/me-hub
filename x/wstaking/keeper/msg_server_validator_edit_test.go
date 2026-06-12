package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestUpdateValidatorRejectsRegionMoveWithDelegations() {
	s.SetupTest()

	oldRegionID := s.usaValidator.Description.RegionID
	oldRegion, found := s.Keeper().GetRegion(s.Ctx, oldRegionID)
	if !found {
		oldRegion = wstakingtypes.Region{
			RegionId:        oldRegionID,
			Name:            "USA",
			Creator:         s.Dao.GlobalDao,
			OperatorAddress: s.usaValidator.OperatorAddress,
			RegionShare:     s.usaValidator.Tokens,
		}
	}
	oldRegion.OperatorAddress = s.usaValidator.OperatorAddress
	oldRegion.RegionShare = s.usaValidator.Tokens
	oldRegion.DelegateAmount = sdk.NewInt(100)
	oldRegion.DelegateInterest = sdk.ZeroDec()
	oldRegion.FixedDepositAmount = sdk.ZeroInt()
	s.Keeper().SetRegion(s.Ctx, oldRegion)

	const newRegionID = "can"
	s.Keeper().SetRegion(s.Ctx, wstakingtypes.Region{
		RegionId:            newRegionID,
		Name:                "CAN",
		Creator:             s.Dao.GlobalDao,
		OperatorAddress:     "",
		RegionShare:         sdk.ZeroInt(),
		DelegateInterest:    sdk.ZeroDec(),
		DelegateAmount:      sdk.ZeroInt(),
		FixedDepositAmount:  sdk.ZeroInt(),
		RegionTreasureAddr:  s.TestAccs[0].String(),
		DepositInterestAddr: s.TestAccs[1].String(),
	})

	description := s.usaValidator.Description
	description.RegionID = newRegionID
	_, err := s.msgServer.UpdateValidator(s.Ctx, &wstakingtypes.MsgUpdateValidator{
		Description:     description,
		StakerAddress:   s.Dao.GlobalDao,
		OperatorAddress: s.usaValidator.OperatorAddress,
	})
	s.Require().Error(err)
	s.Require().ErrorContains(err, "active delegations")

	oldRegionAfter, found := s.Keeper().GetRegion(s.Ctx, oldRegionID)
	s.Require().True(found)
	s.Require().Equal(s.usaValidator.OperatorAddress, oldRegionAfter.OperatorAddress)
	s.Require().Equal(sdk.NewInt(100).String(), oldRegionAfter.DelegateAmount.String())

	newRegionAfter, found := s.Keeper().GetRegion(s.Ctx, newRegionID)
	s.Require().True(found)
	s.Require().Empty(newRegionAfter.OperatorAddress)

	valAddr, err := sdk.ValAddressFromBech32(s.usaValidator.OperatorAddress)
	s.Require().NoError(err)
	validatorAfter, found := s.Keeper().GetValidator(s.Ctx, valAddr)
	s.Require().True(found)
	s.Require().Equal(oldRegionID, validatorAfter.Description.RegionID)
}

func (s *KeeperTestSuite) TestUpdateValidatorMovesRegionWithoutDelegations() {
	s.SetupTest()

	oldRegionID := s.usaValidator.Description.RegionID
	oldRegion, found := s.Keeper().GetRegion(s.Ctx, oldRegionID)
	if !found {
		oldRegion = wstakingtypes.Region{
			RegionId:        oldRegionID,
			Name:            "USA",
			Creator:         s.Dao.GlobalDao,
			OperatorAddress: s.usaValidator.OperatorAddress,
			RegionShare:     s.usaValidator.Tokens,
		}
	}
	oldRegion.OperatorAddress = s.usaValidator.OperatorAddress
	oldRegion.RegionShare = s.usaValidator.Tokens
	oldRegion.DelegateAmount = sdk.ZeroInt()
	oldRegion.DelegateInterest = sdk.ZeroDec()
	oldRegion.FixedDepositAmount = sdk.ZeroInt()
	s.Keeper().SetRegion(s.Ctx, oldRegion)

	const newRegionID = "can"
	s.Keeper().SetRegion(s.Ctx, wstakingtypes.Region{
		RegionId:            newRegionID,
		Name:                "CAN",
		Creator:             s.Dao.GlobalDao,
		OperatorAddress:     "",
		RegionShare:         sdk.ZeroInt(),
		DelegateInterest:    sdk.ZeroDec(),
		DelegateAmount:      sdk.ZeroInt(),
		FixedDepositAmount:  sdk.ZeroInt(),
		RegionTreasureAddr:  s.TestAccs[0].String(),
		DepositInterestAddr: s.TestAccs[1].String(),
	})

	description := s.usaValidator.Description
	description.RegionID = newRegionID
	_, err := s.msgServer.UpdateValidator(s.Ctx, &wstakingtypes.MsgUpdateValidator{
		Description:     description,
		StakerAddress:   s.Dao.GlobalDao,
		OperatorAddress: s.usaValidator.OperatorAddress,
	})
	s.Require().NoError(err)

	oldRegionAfter, found := s.Keeper().GetRegion(s.Ctx, oldRegionID)
	s.Require().True(found)
	s.Require().Empty(oldRegionAfter.OperatorAddress)

	newRegionAfter, found := s.Keeper().GetRegion(s.Ctx, newRegionID)
	s.Require().True(found)
	s.Require().Equal(s.usaValidator.OperatorAddress, newRegionAfter.OperatorAddress)

	valAddr, err := sdk.ValAddressFromBech32(s.usaValidator.OperatorAddress)
	s.Require().NoError(err)
	validatorAfter, found := s.Keeper().GetValidator(s.Ctx, valAddr)
	s.Require().True(found)
	s.Require().Equal(newRegionID, validatorAfter.Description.RegionID)
}
