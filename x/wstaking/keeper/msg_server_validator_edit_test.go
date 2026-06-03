package keeper_test

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestUpdateValidatorKeepsRegionBoundForCaseOnlyRegionID() {
	s.SetupTest()

	regionID := s.usaValidator.Description.RegionID
	_, err := s.msgServer.NewRegion(s.Ctx, &types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            strings.ToUpper(regionID),
		OperatorAddress: s.usaValidator.OperatorAddress,
	})
	s.Require().NoError(err)

	regionBefore, found := s.Keeper().GetRegion(s.Ctx, regionID)
	s.Require().True(found)
	s.Require().Equal(s.usaValidator.OperatorAddress, regionBefore.OperatorAddress)

	description := s.usaValidator.Description
	description.RegionID = strings.ToUpper(regionID)
	_, err = s.msgServer.UpdateValidator(s.Ctx, &types.MsgUpdateValidator{
		Description:     description,
		StakerAddress:   s.Dao.GlobalDao,
		OperatorAddress: s.usaValidator.OperatorAddress,
	})
	s.Require().NoError(err)

	regionAfter, found := s.Keeper().GetRegion(s.Ctx, regionID)
	s.Require().True(found)
	s.Require().Equal(regionBefore.OperatorAddress, regionAfter.OperatorAddress)
	s.Require().Equal(regionBefore.RegionShare.String(), regionAfter.RegionShare.String())

	valAddr, err := sdk.ValAddressFromBech32(s.usaValidator.OperatorAddress)
	s.Require().NoError(err)
	validatorAfter, found := s.Keeper().GetValidator(s.Ctx, valAddr)
	s.Require().True(found)
	s.Require().Equal(regionID, validatorAfter.Description.RegionID)
}
