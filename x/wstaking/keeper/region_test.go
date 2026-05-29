package keeper_test

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	megrouptypes "github.com/openmetaearth/me-hub/x/megroup/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestUnBondRegion() {
	s.SetupTest()

	// 1. Create a region and mock operator address
	regionId := strings.ToLower(types.ExperienceRegionName)
	region, found := s.Keeper().GetRegion(s.Ctx, regionId)
	s.Require().True(found)

	region.RegionShare = sdk.NewInt(12345)
	region.OperatorAddress = s.experienceValidator.OperatorAddress
	s.Keeper().SetRegion(s.Ctx, region)

	// 2. Set up group for region using GroupKeeper
	groupId := uint64(1)
	s.App.GroupKeeper.SetGroupToRegion(s.Ctx, regionId, groupId)

	groupInfo := megrouptypes.GroupInfo{
		Id:    groupId,
		Admin: s.Dao.GlobalDao,
	}
	s.App.GroupKeeper.SetGroupInfo(s.Ctx, groupInfo)

	// Verify initial state
	initialGroup, found := s.App.GroupKeeper.GetGroupInfo(s.Ctx, groupId)
	s.Require().True(found)
	s.Require().Equal(s.Dao.GlobalDao, initialGroup.Admin)

	// 3. Call UnBondRegion
	s.Keeper().UnBondRegion(s.Ctx, regionId)

	// 4. Verify post-conditions:
	// - region.RegionShare should be 0
	// - region.OperatorAddress should still be s.experienceValidator.OperatorAddress
	// - group admin should still be s.Dao.GlobalDao (groupKeeper.UpdateGroupAdmin was not called)
	updatedRegion, found := s.Keeper().GetRegion(s.Ctx, regionId)
	s.Require().True(found)
	s.Require().True(updatedRegion.RegionShare.IsZero())
	s.Require().Equal(s.experienceValidator.OperatorAddress, updatedRegion.OperatorAddress)

	updatedGroup, found := s.App.GroupKeeper.GetGroupInfo(s.Ctx, groupId)
	s.Require().True(found)
	s.Require().Equal(s.Dao.GlobalDao, updatedGroup.Admin)
}
