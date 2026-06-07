package keeper_test

import (
	"strings"

	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestUpdateValidatorRejectsUncreatedRegionID() {
	s.SetupTest()

	_, err := s.msgServer.NewRegion(s.Ctx, &types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            "USA",
		OperatorAddress: s.usaValidator.OperatorAddress,
	})
	s.Require().NoError(err)

	oldRegionID := strings.ToLower(s.usaValidator.Description.RegionID)
	oldRegion, found := s.Keeper().GetRegion(s.Ctx, oldRegionID)
	s.Require().True(found)
	oldRegionOperator := oldRegion.OperatorAddress
	oldRegionShare := oldRegion.RegionShare

	description := s.usaValidator.Description
	description.RegionID = "CAN"
	_, err = s.msgServer.UpdateValidator(s.Ctx, &types.MsgUpdateValidator{
		StakerAddress:   s.Dao.GlobalDao,
		OperatorAddress: s.usaValidator.OperatorAddress,
		Description:     description,
	})
	s.Require().ErrorIs(err, types.ErrRegionNotExist)

	afterRegion, found := s.Keeper().GetRegion(s.Ctx, oldRegionID)
	s.Require().True(found)
	s.Require().Equal(oldRegionOperator, afterRegion.OperatorAddress)
	s.Require().Equal(oldRegionShare.String(), afterRegion.RegionShare.String())

	validator, found := s.Keeper().GetValidator(s.Ctx, s.usaValidator.GetOperator())
	s.Require().True(found)
	s.Require().Equal(oldRegionID, validator.Description.RegionID)
}
