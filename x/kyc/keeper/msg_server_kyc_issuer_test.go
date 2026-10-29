package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	daokeeper "github.com/openmetaearth/me-hub/x/dao/keeper"
	daotypes "github.com/openmetaearth/me-hub/x/dao/types"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/kyc/types"
)

func (s *KeeperTestSuite) TestUpdateKycIssuerServiceState() {
	issuer, pubKey, _ := s.newSecp256k1UserAccount()
	s.setupActiveKyc(issuer, pubKey, "1000000000002")
	service, _ := s.Keeper().GetService(s.Ctx)
	service.Status = didtypes.SERVICE_STATUS_INACTIVE
	s.Keeper().SetService(s.Ctx, service)
	_, err := s.msgServer.UpdateKycIssuer(s.Ctx, types.NewMsgUpdateKycIssuer(s.Dao.GlobalDao, []string{issuer.String()}, types.KycIssuerAdd))
	s.Require().NoError(err)
	service, _ = s.Keeper().GetService(s.Ctx)
	s.Require().Equal(didtypes.SERVICE_STATUS_INACTIVE, service.Status)
	s.App.DidKeeper.DeleteService(s.Ctx, types.ModuleName)
	_, err = s.msgServer.UpdateKycIssuer(s.Ctx, types.NewMsgUpdateKycIssuer(s.Dao.GlobalDao, []string{issuer.String()}, types.KycIssuerRemove))
	s.Require().ErrorIs(err, didtypes.ErrServiceNotFound)
}

func (s *KeeperTestSuite) TestUpdateKycIssuerOnlyGlobalDao() {
	issuer, pubKey, _ := s.newSecp256k1UserAccount()
	s.setupActiveKyc(issuer, pubKey, "1000000000002")
	_, err := s.msgServer.UpdateKycIssuer(s.Ctx, types.NewMsgUpdateKycIssuer(s.Dao.GlobalDao, []string{issuer.String()}, types.KycIssuerAdd))
	s.Require().NoError(err)
	for _, creator := range []string{s.Dao.MeidDao, s.Dao.DevOperator, issuer.String()} {
		_, err := s.msgServer.UpdateKycIssuer(s.Ctx, types.NewMsgUpdateKycIssuer(creator, []string{issuer.String()}, types.KycIssuerRemove))
		s.Require().ErrorIs(err, types.ErrUnauthorized)
	}
	service, _ := s.Keeper().GetService(s.Ctx)
	s.Require().Contains(service.Issuers, "1000000000002")
}

func (s *KeeperTestSuite) TestUpdateKycIssuerRejectsInvalidIdentityWithoutPartialChanges() {
	valid, pubKeyValid, _ := s.newSecp256k1UserAccount()
	s.setupActiveKyc(valid, pubKeyValid, "1000000000002")
	invalid, _, _ := s.newSecp256k1UserAccount()
	original, _ := s.Keeper().GetService(s.Ctx)
	msg := types.NewMsgUpdateKycIssuer(s.Dao.GlobalDao, []string{valid.String(), invalid.String()}, types.KycIssuerAdd)
	_, err := s.msgServer.UpdateKycIssuer(s.Ctx, msg)
	s.Require().ErrorIs(err, didtypes.ErrDidNotFound)

	s.Keeper().SetDID(s.Ctx, invalid, "1000000000003")
	_, err = s.msgServer.UpdateKycIssuer(s.Ctx, msg)
	s.Require().ErrorIs(err, didtypes.ErrHolderNotFound)
	info := didtypes.DidInfo{Did: "1000000000003", Address: invalid.String(), Status: didtypes.DID_STATUS_INACTIVE}
	s.Keeper().SetDidInfo(s.Ctx, info.Did, info)
	_, err = s.msgServer.UpdateKycIssuer(s.Ctx, msg)
	s.Require().ErrorIs(err, didtypes.ErrDidNotActive)
	info.Status = didtypes.DID_STATUS_ACTIVE
	info.Address = valid.String()
	s.Keeper().SetDidInfo(s.Ctx, info.Did, info)
	_, err = s.msgServer.UpdateKycIssuer(s.Ctx, msg)
	s.Require().ErrorIs(err, didtypes.ErrInvalidIssuer)
	service, _ := s.Keeper().GetService(s.Ctx)
	s.Require().Equal(original, service)
}

func (s *KeeperTestSuite) TestUpdateKycIssuerRemovalRejectsInvalidIdentityWithoutPartialChanges() {
	issuerA, pubKeyA, _ := s.newSecp256k1UserAccount()
	issuerB, pubKeyB, _ := s.newSecp256k1UserAccount()
	const didA = "1000000000002"
	const didB = "1000000000003"
	s.setupActiveKyc(issuerA, pubKeyA, didA)
	s.setupActiveKyc(issuerB, pubKeyB, didB)
	_, err := s.msgServer.UpdateKycIssuer(s.Ctx, types.NewMsgUpdateKycIssuer(
		s.Dao.GlobalDao, []string{issuerA.String(), issuerB.String()}, types.KycIssuerAdd,
	))
	s.Require().NoError(err)
	original, found := s.Keeper().GetService(s.Ctx)
	s.Require().True(found)
	infoA, found := s.Keeper().GetDidInfo(s.Ctx, didA)
	s.Require().True(found)
	infoB, found := s.Keeper().GetDidInfo(s.Ctx, didB)
	s.Require().True(found)
	unknown, _ := s.NewAccount()

	for _, tc := range []struct {
		name    string
		address string
		prepare func(sdk.Context)
		err     error
	}{
		{
			name: "missing address mapping with issuer record retained", address: issuerB.String(),
			prepare: func(ctx sdk.Context) { s.App.DidKeeper.DeleteDID(ctx, issuerB) },
			err:     didtypes.ErrDidNotFound,
		},
		{
			name: "unregistered address", address: unknown.String(),
			err: didtypes.ErrDidNotFound,
		},
		{
			name: "missing DID info", address: issuerB.String(),
			prepare: func(ctx sdk.Context) { s.App.DidKeeper.DeleteDidInfo(ctx, didB) },
			err:     didtypes.ErrHolderNotFound,
		},
		{
			name: "address mismatch", address: issuerB.String(),
			prepare: func(ctx sdk.Context) {
				info := infoB
				info.Address = issuerA.String()
				s.Keeper().SetDidInfo(ctx, didB, info)
			},
			err: didtypes.ErrInvalidIssuer,
		},
		{
			name: "DID mismatch", address: issuerB.String(),
			prepare: func(ctx sdk.Context) {
				info := infoB
				info.Did = "1000000000004"
				s.Keeper().SetDidInfo(ctx, didB, info)
			},
			err: didtypes.ErrInvalidIssuer,
		},
	} {
		s.Run(tc.name, func() {
			ctx, _ := s.Ctx.CacheContext()
			if tc.prepare != nil {
				tc.prepare(ctx)
			}
			eventCount := len(ctx.EventManager().Events())
			msg := types.NewMsgUpdateKycIssuer(s.Dao.GlobalDao, []string{issuerA.String(), tc.address}, types.KycIssuerRemove)
			_, err := s.msgServer.UpdateKycIssuer(ctx, msg)
			s.Require().ErrorIs(err, tc.err)
			s.Require().ErrorContains(err, tc.address)
			service, found := s.Keeper().GetService(ctx)
			s.Require().True(found)
			s.Require().Equal(original, service)
			afterInfo, found := s.Keeper().GetDidInfo(ctx, didA)
			s.Require().True(found)
			s.Require().Equal(infoA, afterInfo)
			s.Require().Len(ctx.EventManager().Events(), eventCount)
		})
	}
}

func (s *KeeperTestSuite) TestDaoRotationPreservesKycIssuerMemberships() {
	issuer, pubKey, _ := s.newSecp256k1UserAccount()
	s.setupActiveKyc(issuer, pubKey, "1000000000002")
	_, err := s.msgServer.UpdateKycIssuer(s.Ctx, types.NewMsgUpdateKycIssuer(s.Dao.GlobalDao, []string{issuer.String()}, types.KycIssuerAdd))
	s.Require().NoError(err)
	original, _ := s.Keeper().GetService(s.Ctx)
	meidAddr := sdk.MustAccAddressFromBech32(s.Dao.MeidDao)
	s.Keeper().SetDID(s.Ctx, meidAddr, "1000000000003")
	s.Keeper().SetDidInfo(s.Ctx, "1000000000003", didtypes.DidInfo{
		Did: "1000000000003", Address: s.Dao.MeidDao, Status: didtypes.DID_STATUS_ACTIVE,
	})
	newDao, _ := s.NewAccount()
	newAddresses := s.Dao
	newAddresses.GlobalDao = newDao.String()
	daoServer := daokeeper.NewMsgServerImpl(s.App.DaoKeeper)
	_, err = daoServer.UpdateDao(s.Ctx, daotypes.NewMsgUpdateDao(sdk.MustAccAddressFromBech32(s.Dao.GlobalDao), newAddresses))
	s.Require().NoError(err)
	service, _ := s.Keeper().GetService(s.Ctx)
	s.Require().Equal(original.Issuers, service.Issuers)
	_, found := s.Keeper().GetDID(s.Ctx, sdk.MustAccAddressFromBech32(s.Dao.GlobalDao))
	s.Require().False(found)
	_, err = s.msgServer.UpdateKycIssuer(s.Ctx, types.NewMsgUpdateKycIssuer(newDao.String(), []string{newDao.String()}, types.KycIssuerRemove))
	s.Require().NoError(err)
	_, err = daoServer.UpdateDao(s.Ctx, daotypes.NewMsgUpdateDao(newDao, s.Dao))
	s.Require().NoError(err)
	service, _ = s.Keeper().GetService(s.Ctx)
	s.Require().NotContains(service.Issuers, "1000000000001")
	s.Require().Contains(service.Issuers, "1000000000002")
}
