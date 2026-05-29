package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

// TestGetValOwnerAddress tests GetValOwnerAddress with normal and empty OperatorAddress scenarios.
// This covers the bug fix for #102 where UnBondRegion empties OperatorAddress,
// causing the ante handler to block ALL transactions from users in that region.
func (s *KeeperTestSuite) TestGetValOwnerAddress() {
	s.SetupTest()

	// Create a region bound to a validator
	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	regionId := s.meEarthValidator.Description.RegionID

	s.Run("returns validator owner address for region with valid operator", func() {
		ownerAddr, err := s.Keeper().GetValOwnerAddress(s.Ctx, regionId)
		s.Require().NoError(err)
		s.Require().NotEmpty(ownerAddr)
		s.Require().Equal(s.meEarthValidator.OwnerAddress, ownerAddr)
	})

	s.Run("falls back to proposer when OperatorAddress is empty (after UnBondRegion)", func() {
		// Simulate UnBondRegion (as called when validator fully unstakes)
		s.Keeper().UnBondRegion(s.Ctx, regionId)

		// Verify the region now has empty OperatorAddress
		region, found := s.Keeper().GetRegion(s.Ctx, regionId)
		s.Require().True(found)
		s.Require().Empty(region.OperatorAddress)
		s.Require().True(region.RegionShare.IsZero())

		// GetValOwnerAddress should NOT error — it should fall back to proposer
		ownerAddr, err := s.Keeper().GetValOwnerAddress(s.Ctx, regionId)
		s.Require().NoError(err)
		s.Require().NotEmpty(ownerAddr)
	})

	s.Run("returns error for non-existent region", func() {
		_, err := s.Keeper().GetValOwnerAddress(s.Ctx, "nonexistent")
		s.Require().Error(err)
	})
}

// TestGetValOwnerAddressAfterFullUnstake tests the full flow: stake → unstake → verify
// that GetValOwnerAddress still works for users in the region.
func (s *KeeperTestSuite) TestGetValOwnerAddressAfterFullUnstake() {
	s.SetupTest()

	// Create region
	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	regionId := s.meEarthValidator.Description.RegionID

	// Verify it works before unbond
	ownerAddr, err := s.Keeper().GetValOwnerAddress(s.Ctx, regionId)
	s.Require().NoError(err)
	s.Require().Equal(s.meEarthValidator.OwnerAddress, ownerAddr)

	// Directly call UnBondRegion to simulate full unstake
	s.Keeper().UnBondRegion(s.Ctx, regionId)

	// The key assertion: GetValOwnerAddress must NOT return an error
	// This is what the ante handler calls for fee distribution
	ownerAddr, err = s.Keeper().GetValOwnerAddress(s.Ctx, regionId)
	s.Require().NoError(err, "GetValOwnerAddress must not fail after UnBondRegion - this would block all transactions")
	s.Require().NotEmpty(ownerAddr, "owner address must not be empty after fallback")
}

// TestUnBondRegionPreservesRegion tests that UnBondRegion zeroes shares and clears operator
// but does NOT remove the region from the store.
func (s *KeeperTestSuite) TestUnBondRegionPreservesRegion() {
	s.SetupTest()

	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	regionId := s.meEarthValidator.Description.RegionID

	// Verify region exists with operator
	region, found := s.Keeper().GetRegion(s.Ctx, regionId)
	s.Require().True(found)
	s.Require().NotEmpty(region.OperatorAddress)

	// UnBond
	s.Keeper().UnBondRegion(s.Ctx, regionId)

	// Region should still exist, but with empty operator and zero shares
	region, found = s.Keeper().GetRegion(s.Ctx, regionId)
	s.Require().True(found, "region should still exist after UnBondRegion")
	s.Require().Empty(region.OperatorAddress, "operator address should be empty after UnBondRegion")
	s.Require().True(region.RegionShare.IsZero(), "region share should be zero after UnBondRegion")
}
