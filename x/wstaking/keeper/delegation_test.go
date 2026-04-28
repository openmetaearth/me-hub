package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestIterateRegionKycDelegatins() {
	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	// Create delegations
	delegator1 := s.TestAccs[0]
	delegator2 := s.TestAccs[1]
	delegator3 := s.TestAccs[2]

	// Set KYC data for delegators
	s.App.KycKeeper.SetDID(s.Ctx, delegator1, "did1")
	s.App.KycKeeper.SetDidInfo(s.Ctx, "did1", didtypes.DidInfo{
		Did:      "did1",
		Address:  delegator1.String(),
		RegionId: types.MeEarthRegionId,
		KycLevel: 2,
		Status:   didtypes.DID_STATUS_ACTIVE,
	})
	vc1 := didtypes.Credential{Did: "did1", Data: []byte(types.MeEarthRegionId)}
	s.App.KycKeeper.SetKYC(s.Ctx, "did1", vc1)
	s.App.KycKeeper.AddFilters(s.Ctx, "did1", [][]byte{[]byte(types.MeEarthRegionId)}, vc1)

	s.App.KycKeeper.SetDID(s.Ctx, delegator2, "did2")
	s.App.KycKeeper.SetDidInfo(s.Ctx, "did2", didtypes.DidInfo{
		Did:      "did2",
		Address:  delegator2.String(),
		RegionId: types.MeEarthRegionId,
		KycLevel: 2,
		Status:   didtypes.DID_STATUS_ACTIVE,
	})
	vc1 = didtypes.Credential{Did: "did2", Data: []byte(types.MeEarthRegionId)}
	s.App.KycKeeper.SetKYC(s.Ctx, "did2", vc1)
	s.App.KycKeeper.AddFilters(s.Ctx, "did2", [][]byte{[]byte(types.MeEarthRegionId)}, vc1)

	s.App.KycKeeper.SetDID(s.Ctx, delegator3, "did3")
	vc1 = didtypes.Credential{Did: "did3", Data: []byte(types.MeEarthRegionId)}
	s.App.KycKeeper.SetDidInfo(s.Ctx, "did3", didtypes.DidInfo{
		Did:      "did3",
		Address:  delegator3.String(),
		RegionId: types.MeEarthRegionId,
		KycLevel: 2,
		Status:   didtypes.DID_STATUS_ACTIVE,
	})
	s.App.KycKeeper.SetKYC(s.Ctx, "did3", vc1)
	s.App.KycKeeper.AddFilters(s.Ctx, "did3", [][]byte{[]byte(types.MeEarthRegionId)}, vc1)

	s.App.StakingKeeper.SetDelegation(s.Ctx, stakingtypes.Delegation{
		DelegatorAddress: delegator1.String(),
		ValidatorAddress: s.meEarthValidator.OperatorAddress,
	})
	s.App.StakingKeeper.SetDelegation(s.Ctx, stakingtypes.Delegation{
		DelegatorAddress: delegator2.String(),
		ValidatorAddress: s.meEarthValidator.OperatorAddress,
	})
	s.App.StakingKeeper.SetDelegation(s.Ctx, stakingtypes.Delegation{
		DelegatorAddress: delegator3.String(),
		ValidatorAddress: s.meEarthValidator.OperatorAddress,
	})

	region, f := s.App.StakingKeeper.GetRegion(s.Ctx, types.MeEarthRegionId)
	s.Require().True(f)
	region.OperatorAddress = s.usaValidator.OperatorAddress
	s.App.StakingKeeper.SetRegion(s.Ctx, region)
	s.App.StakingKeeper.SetChangeDelegationValidator(s.Ctx, types.MeEarthRegionId)
	// Call ChangeDelegationValidator
	s.App.StakingKeeper.ChangeDelegationValidator(s.Ctx)

	// Verify delegations' validator addresses have been updated
	delegation1, found := s.App.StakingKeeper.GetDelegation(s.Ctx, delegator1, sdk.ValAddress{})
	s.Require().True(found)
	s.Require().Equal(s.usaValidator.OperatorAddress, delegation1.ValidatorAddress)

	delegation2, found := s.App.StakingKeeper.GetDelegation(s.Ctx, delegator2, sdk.ValAddress{})
	s.Require().True(found)
	s.Require().Equal(s.usaValidator.OperatorAddress, delegation2.ValidatorAddress)

	delegation3, found := s.App.StakingKeeper.GetDelegation(s.Ctx, delegator3, sdk.ValAddress{})
	s.Require().True(found)
	s.Require().Equal(s.usaValidator.OperatorAddress, delegation3.ValidatorAddress)
}
