package keeper_test

import (
	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/kyc/keeper"
	"github.com/openmetaearth/me-hub/x/kyc/types"
	wstakingkeeper "github.com/openmetaearth/me-hub/x/wstaking/keeper"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
	"github.com/stretchr/testify/suite"
	"testing"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	msgServer           types.MsgServer
	queryClient         types.QueryClient
	meEarthValidator    stakingtypes.Validator
	experienceValidator stakingtypes.Validator
	usaValidator        stakingtypes.Validator
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) Keeper() *keeper.Keeper {
	return s.App.KycKeeper
}

func (s *KeeperTestSuite) TestGetAllRegions() {
	s.SetupTest()
	regions := s.Keeper().GetAllRegions(s.Ctx)
	// We expect exactly 3 regions that were set up in SetupTest
	s.Require().Equal(3, len(regions))
	for _, reg := range regions {
		s.Require().NotEmpty(reg.Id)
		s.Require().NotEmpty(reg.Name)
	}
}

func (s *KeeperTestSuite) TestKeeper_SetKycIssers() {
	s.SetupTest()
	k := s.Keeper()

	// Initial addresses and DIDs
	oldGlobalDao := "me1frjhlw9slyy7mrhmk0r4vytkyldxqtkf326amv"
	oldMeidDao := "me1c5zp26c0gq2klk87nrpff3y52u34zn4ydug2yd"
	newGlobalDao := "me1hrxxjeqae2y5wx3kxcljzns9f2lguygu9qngxh"
	newMeidDao := "me14jazxhme3ptv00k52fza5rravx4xn27qs0slz2"

	globalDid := "did:me:global"
	meidDid := "did:me:meid"
	thirdPartyDid := "did:me:thirdparty"

	// Setup DIDs and DID Info
	k.SetDID(s.Ctx, sdk.MustAccAddressFromBech32(oldGlobalDao), globalDid)
	k.SetDID(s.Ctx, sdk.MustAccAddressFromBech32(oldMeidDao), meidDid)

	k.SetDidInfo(s.Ctx, globalDid, didtypes.DidInfo{Did: globalDid, Address: oldGlobalDao, Status: didtypes.DID_STATUS_ACTIVE})
	k.SetDidInfo(s.Ctx, meidDid, didtypes.DidInfo{Did: meidDid, Address: oldMeidDao, Status: didtypes.DID_STATUS_ACTIVE})

	// Setup KYC service with Global DAO, MEID DAO, and a third-party issuer
	k.SetService(s.Ctx, didtypes.Service{
		Sid:     "kyc",
		Status:  didtypes.SERVICE_STATUS_ACTIVE,
		Issuers: []string{globalDid, meidDid, thirdPartyDid},
	})

	// Execute SetKycIssers
	err := k.SetKycIssers(
		s.Ctx,
		[]string{oldGlobalDao, oldMeidDao},
		[]string{newGlobalDao, newMeidDao},
	)
	s.Require().NoError(err)

	// Verify that the KYC service issuers preserved the third-party issuer
	svc, found := k.GetService(s.Ctx)
	s.Require().True(found)
	s.Require().Equal(3, len(svc.Issuers))
	s.Require().Contains(svc.Issuers, thirdPartyDid)
	s.Require().Contains(svc.Issuers, globalDid)
	s.Require().Contains(svc.Issuers, meidDid)

	// Verify old DID maps are deleted and new maps are created
	_, found = k.GetDID(s.Ctx, sdk.MustAccAddressFromBech32(oldGlobalDao))
	s.Require().False(found)
	resolvedGlobal, found := k.GetDID(s.Ctx, sdk.MustAccAddressFromBech32(newGlobalDao))
	s.Require().True(found)
	s.Require().Equal(globalDid, resolvedGlobal)
}

func (s *KeeperTestSuite) SetupTest() {
	app := apptesting.Setup(s.T(), false)
	ctx := app.GetBaseApp().NewContext(false, cometbftproto.Header{})

	err := app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	s.Require().NoError(err)

	err = app.BankKeeper.SetParams(ctx, banktypes.DefaultParams())
	s.Require().NoError(err)

	stakingParams := stakingtypes.DefaultParams()
	stakingParams.BondDenom = params.BaseDenom
	err = app.StakingKeeper.SetParams(ctx, stakingParams)
	s.Require().NoError(err)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	nativeQuerier := keeper.Querier{Keeper: app.KycKeeper}
	types.RegisterQueryServer(queryHelper, nativeQuerier)
	queryClient := types.NewQueryClient(queryHelper)

	s.App = app

	s.msgServer = keeper.NewMsgServerImpl(*app.KycKeeper)
	s.Ctx = ctx
	s.queryClient = queryClient

	stakingKeeperMsgSrv := stakingkeeper.NewMsgServerImpl(app.StakingKeeper.Keeper)
	stakingMsgServer := wstakingkeeper.NewMsgServerImpl(app.StakingKeeper, app.TransferKeeper, stakingKeeperMsgSrv)

	s.InitializeDao()

	validators := s.App.StakingKeeper.GetValidators(s.Ctx, 10)
	s.Require().True(len(validators) >= 3)
	s.meEarthValidator = validators[0]
	s.experienceValidator = validators[1]
	s.usaValidator = validators[2]

	newRegion := wstakingtypes.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            wstakingtypes.ExperienceRegionName,
		OperatorAddress: s.experienceValidator.OperatorAddress,
	}
	_, err = stakingMsgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	newRegion = wstakingtypes.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            "USA",
		OperatorAddress: s.usaValidator.OperatorAddress,
	}
	_, err = stakingMsgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	newRegion = wstakingtypes.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            wstakingtypes.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err = stakingMsgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestPubKeyFromString() {
	s.SetupTest()
	pubkey := `{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A83z2Fnur8jc+tGvkCJjkZTBeJDLSObk8nVKOpY9P679"}`
	accAddr, _ := s.Keeper().MustAccAddressFromPubkeyString(pubkey)
	s.Require().Equal("me13w3mxrd9tvq3r6gzheqjuzf8pnaruvug5787yu", accAddr.String())

	secp256k1Pubkey := `{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"A9Iwsz0CXw/AEVGq7wyM4wuNbcoeB1dXTBje1lRXvKBD"}`
	secpAccAddr, _ := s.Keeper().MustAccAddressFromPubkeyString(secp256k1Pubkey)
	s.Require().Equal("me1kj3emedrrq66vdqf3pzpfjmytympl4j2a4xd0c", secpAccAddr.String())
}
