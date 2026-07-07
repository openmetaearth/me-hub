package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	msgServer           types.MsgServer
	queryClient         types.QueryClient
	queryHelper         *baseapp.QueryServiceTestHelper
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

func (s *KeeperTestSuite) SetupTest() {
	app := apptesting.Setup(s.T())
	ctx := app.GetBaseApp().NewContext(false)

	stakingParams, err := app.StakingKeeper.GetParams(ctx)
	s.Require().NoError(err)
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
	s.queryHelper = queryHelper

	stakingKeeperMsgSrv := stakingkeeper.NewMsgServerImpl(app.StakingKeeper.Keeper)
	stakingMsgServer := wstakingkeeper.NewMsgServerImpl(app.StakingKeeper, app.TransferKeeper, stakingKeeperMsgSrv)

	s.InitializeDao()

	// Set up globalDao as an issuer for the KYC service
	globalDaoAddr := sdk.MustAccAddressFromBech32(s.Dao.GlobalDao)
	globalDaoDID := "0000000000001"
	s.App.KycKeeper.SetDID(s.Ctx, globalDaoAddr, globalDaoDID)
	s.App.KycKeeper.SetDidInfo(s.Ctx, globalDaoDID, didtypes.DidInfo{
		Did:    globalDaoDID,
		Status: didtypes.DID_STATUS_ACTIVE,
	})
	// Add globalDao DID to KYC service issuers
	svc, found := s.App.KycKeeper.GetService(s.Ctx)
	if !found {
		svc = didtypes.Service{
			Sid:         types.ModuleName,
			Name:        types.ModuleName,
			Description: "The KYC verifiable credential issuer based The DID(Decentralized Identity).",
			Status:      didtypes.SERVICE_STATUS_ACTIVE,
		}
	}
	svc.Issuers = append(svc.Issuers, globalDaoDID)
	s.App.KycKeeper.SetService(s.Ctx, svc)

	validators, err := s.App.StakingKeeper.GetValidators(s.Ctx, 10)
	s.Require().NoError(err)
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

	// simulate EndBlock cache refresh so GetRegionCache works in handlers (cache is populated in EndBlock in production)
	app.StakingKeeper.SetRegionsCache(s.Ctx, app.StakingKeeper.GetAllRegion(s.Ctx))

	// Update queryHelper context to include all setup state (service, DID, etc.)
	s.queryHelper.Ctx = s.Ctx
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

// NewAccountStr creates a funded test account and returns its address and pubkey as a JSON string.
func (s *KeeperTestSuite) NewAccountStr() (sdk.AccAddress, string) {
	privKey := ed25519.GenPrivKey()
	addr := sdk.AccAddress(privKey.PubKey().Address())
	apptesting.FundAccount(s.App, s.Ctx, addr, sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 1_000_000_000)))
	pubkeyBytes, err := s.App.AppCodec().MarshalInterfaceJSON(privKey.PubKey())
	if err != nil {
		panic(err)
	}
	return addr, string(pubkeyBytes)
}
