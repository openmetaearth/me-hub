package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
)

func (s *KeeperTestSuite) TestSetKycIssersPreservesOtherIssuers() {
	s.SetupTest()

	oldGlobalAddr := sdk.MustAccAddressFromBech32(s.Dao.GlobalDao)
	oldMeidAddr := sdk.MustAccAddressFromBech32(s.Dao.MeidDao)
	oldGlobalDid, found := s.Keeper().GetDID(s.Ctx, oldGlobalAddr)
	s.Require().True(found)

	oldMeidDid := "0000000000002"
	s.Keeper().SetDID(s.Ctx, oldMeidAddr, oldMeidDid)
	s.Keeper().SetDidInfo(s.Ctx, oldMeidDid, didtypes.DidInfo{
		Did:     oldMeidDid,
		Address: oldMeidAddr.String(),
		Status:  didtypes.DID_STATUS_ACTIVE,
	})

	thirdPartyAddr, _ := s.NewAccount()
	thirdPartyDid := "third-party-issuer"
	s.Keeper().SetDID(s.Ctx, thirdPartyAddr, thirdPartyDid)
	s.Keeper().SetDidInfo(s.Ctx, thirdPartyDid, didtypes.DidInfo{
		Did:     thirdPartyDid,
		Address: thirdPartyAddr.String(),
		Status:  didtypes.DID_STATUS_ACTIVE,
	})

	service, found := s.Keeper().GetService(s.Ctx)
	s.Require().True(found)
	service.Issuers = []string{oldGlobalDid, thirdPartyDid, oldMeidDid}
	s.Keeper().SetService(s.Ctx, service)

	newGlobalAddr, _ := s.NewAccount()
	newMeidAddr, _ := s.NewAccount()
	err := s.Keeper().SetKycIssers(
		s.Ctx,
		[]string{s.Dao.GlobalDao, s.Dao.MeidDao},
		[]string{newGlobalAddr.String(), newMeidAddr.String()},
	)
	s.Require().NoError(err)

	updatedService, found := s.Keeper().GetService(s.Ctx)
	s.Require().True(found)
	s.Require().Equal([]string{oldGlobalDid, thirdPartyDid, oldMeidDid}, updatedService.Issuers)

	newGlobalDid, found := s.Keeper().GetDID(s.Ctx, newGlobalAddr)
	s.Require().True(found)
	s.Require().Equal(oldGlobalDid, newGlobalDid)

	newMeidDid, found := s.Keeper().GetDID(s.Ctx, newMeidAddr)
	s.Require().True(found)
	s.Require().Equal(oldMeidDid, newMeidDid)

	preservedThirdPartyDid, found := s.Keeper().GetDID(s.Ctx, thirdPartyAddr)
	s.Require().True(found)
	s.Require().Equal(thirdPartyDid, preservedThirdPartyDid)
}

func (s *KeeperTestSuite) TestSetKycIssersPreservesDistinctIssuersWhenDaoAddressesSwap() {
	s.SetupTest()

	oldGlobalAddr := sdk.MustAccAddressFromBech32(s.Dao.GlobalDao)
	oldMeidAddr := sdk.MustAccAddressFromBech32(s.Dao.MeidDao)
	oldGlobalDid, found := s.Keeper().GetDID(s.Ctx, oldGlobalAddr)
	s.Require().True(found)

	oldMeidDid := "0000000000002"
	s.Keeper().SetDID(s.Ctx, oldMeidAddr, oldMeidDid)
	s.Keeper().SetDidInfo(s.Ctx, oldMeidDid, didtypes.DidInfo{
		Did:     oldMeidDid,
		Address: oldMeidAddr.String(),
		Status:  didtypes.DID_STATUS_ACTIVE,
	})

	service, found := s.Keeper().GetService(s.Ctx)
	s.Require().True(found)
	service.Issuers = []string{oldGlobalDid, oldMeidDid}
	s.Keeper().SetService(s.Ctx, service)

	err := s.Keeper().SetKycIssers(
		s.Ctx,
		[]string{s.Dao.GlobalDao, s.Dao.MeidDao},
		[]string{s.Dao.MeidDao, s.Dao.GlobalDao},
	)
	s.Require().NoError(err)

	newGlobalDid, found := s.Keeper().GetDID(s.Ctx, oldMeidAddr)
	s.Require().True(found)
	s.Require().Equal(oldGlobalDid, newGlobalDid)

	newMeidDid, found := s.Keeper().GetDID(s.Ctx, oldGlobalAddr)
	s.Require().True(found)
	s.Require().Equal(oldMeidDid, newMeidDid)

	updatedGlobalInfo, found := s.Keeper().GetDidInfo(s.Ctx, oldGlobalDid)
	s.Require().True(found)
	s.Require().Equal(oldMeidAddr.String(), updatedGlobalInfo.Address)

	updatedMeidInfo, found := s.Keeper().GetDidInfo(s.Ctx, oldMeidDid)
	s.Require().True(found)
	s.Require().Equal(oldGlobalAddr.String(), updatedMeidInfo.Address)

	updatedService, found := s.Keeper().GetService(s.Ctx)
	s.Require().True(found)
	s.Require().Equal([]string{oldGlobalDid, oldMeidDid}, updatedService.Issuers)
}
