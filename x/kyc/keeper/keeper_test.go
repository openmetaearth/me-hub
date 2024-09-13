package keeper_test

import (
	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/st-chain/me-hub/app/apptesting"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/x/kyc/keeper"
	"github.com/st-chain/me-hub/x/kyc/types"
	wstakingkeeper "github.com/st-chain/me-hub/x/wstaking/keeper"
	wstakingtypes "github.com/st-chain/me-hub/x/wstaking/types"
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

func (suite *KeeperTestSuite) Keeper() *keeper.Keeper {
	return suite.App.KycKeeper
}

func (suite *KeeperTestSuite) nextBlock() {
	h := suite.Ctx.BlockHeight()
	suite.Ctx = suite.Ctx.WithBlockHeight(h + 1)
}

func (suite *KeeperTestSuite) SetupTest() {
	app := apptesting.Setup(suite.T(), false)
	ctx := app.GetBaseApp().NewContext(false, cometbftproto.Header{})

	err := app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	suite.Require().NoError(err)

	err = app.BankKeeper.SetParams(ctx, banktypes.DefaultParams())
	suite.Require().NoError(err)

	stakingParams := stakingtypes.DefaultParams()
	stakingParams.BondDenom = params.BaseDenom
	app.StakingKeeper.SetParams(ctx, stakingParams)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	nativeQuerier := keeper.Querier{Keeper: app.KycKeeper}
	types.RegisterQueryServer(queryHelper, nativeQuerier)
	queryClient := types.NewQueryClient(queryHelper)

	suite.App = app

	suite.msgServer = keeper.NewMsgServerImpl(*app.KycKeeper)
	suite.Ctx = ctx
	suite.queryClient = queryClient

	stakingKeeperMsgSrv := stakingkeeper.NewMsgServerImpl(app.StakingKeeper.Keeper)
	stakingMsgServer := wstakingkeeper.NewMsgServerImpl(app.StakingKeeper, stakingKeeperMsgSrv)

	suite.InitializeDao()

	validators := suite.App.StakingKeeper.GetValidators(suite.Ctx, 10)
	suite.Require().True(len(validators) >= 3)
	suite.meEarthValidator = validators[0]
	suite.experienceValidator = validators[1]
	suite.usaValidator = validators[2]

	newRegion := wstakingtypes.MsgNewRegion{
		Creator:         suite.Dao.GlobalDao,
		Name:            wstakingtypes.ExperienceRegionName,
		OperatorAddress: suite.experienceValidator.OperatorAddress,
	}
	_, err = stakingMsgServer.NewRegion(suite.Ctx, &newRegion)
	suite.Require().NoError(err)

	newRegion = wstakingtypes.MsgNewRegion{
		Creator:         suite.Dao.GlobalDao,
		Name:            "USA",
		OperatorAddress: suite.usaValidator.OperatorAddress,
	}
	_, err = stakingMsgServer.NewRegion(suite.Ctx, &newRegion)
	suite.Require().NoError(err)

	newRegion = wstakingtypes.MsgNewRegion{
		Creator:         suite.Dao.GlobalDao,
		Name:            wstakingtypes.MeEarthRegionName,
		OperatorAddress: suite.meEarthValidator.OperatorAddress,
	}
	_, err = stakingMsgServer.NewRegion(suite.Ctx, &newRegion)
	suite.Require().NoError(err)
}

func (s *KeeperTestSuite) TestPubKeyFromString() {
	s.SetupTest()
	pubkey := `{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A83z2Fnur8jc+tGvkCJjkZTBeJDLSObk8nVKOpY9P679"}`
	accAddr := s.Keeper().MustAccAddressFromPubkeyString(pubkey)
	s.Require().Equal("me13w3mxrd9tvq3r6gzheqjuzf8pnaruvug5787yu", accAddr.String())
}
