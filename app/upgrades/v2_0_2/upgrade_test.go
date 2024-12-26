package v2_0_2_test

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/baseapp"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/app/upgrades/v2_0_2"
	wstakingkeeper "github.com/st-chain/me-hub/x/wstaking/keeper"
	wstakingtypes "github.com/st-chain/me-hub/x/wstaking/types"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/stretchr/testify/suite"

	"github.com/st-chain/me-hub/app/apptesting"
)

type UpgradeTestSuite struct {
	apptesting.KeeperTestHelper

	msgServer           wstakingkeeper.MsgServer
	queryClient         wstakingtypes.QueryClient
	meEarthValidator    stakingtypes.Validator
	experienceValidator stakingtypes.Validator
	usaValidator        stakingtypes.Validator
}

func TestUpgradeTestSuite(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))
}

func (suite *UpgradeTestSuite) Keeper() *wstakingkeeper.Keeper {
	return suite.App.StakingKeeper
}

func (suite *UpgradeTestSuite) nextBlock() {
	h := suite.Ctx.BlockHeight()
	suite.Ctx = suite.Ctx.WithBlockHeight(h + 1)
}

func (suite *UpgradeTestSuite) SetupTest() {
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
	nativeQuerier := wstakingkeeper.Querier{Keeper: app.StakingKeeper}
	wstakingtypes.RegisterQueryServer(queryHelper, nativeQuerier)
	queryClient := wstakingtypes.NewQueryClient(queryHelper)
	suite.queryClient = queryClient

	suite.App = app
	suite.Ctx = ctx

	stakingKeeperMsgSrv := stakingkeeper.NewMsgServerImpl(app.StakingKeeper.Keeper)
	suite.msgServer = wstakingkeeper.NewMsgServerImpl(app.StakingKeeper, stakingKeeperMsgSrv)

	suite.InitializeDao()

	validators := suite.Keeper().GetValidators(suite.Ctx, 10)
	suite.Require().True(len(validators) >= 3)
	suite.meEarthValidator = validators[0]
	suite.experienceValidator = validators[1]
	suite.usaValidator = validators[2]

	newRegion := wstakingtypes.MsgNewRegion{
		Creator:         suite.Dao.GlobalDao,
		Name:            wstakingtypes.ExperienceRegionName,
		OperatorAddress: suite.experienceValidator.OperatorAddress,
	}
	_, err = suite.msgServer.NewRegion(suite.Ctx, &newRegion)

	newRegion = wstakingtypes.MsgNewRegion{
		Creator:         suite.Dao.GlobalDao,
		Name:            wstakingtypes.MeEarthRegionName,
		OperatorAddress: suite.meEarthValidator.OperatorAddress,
	}
	_, err = suite.msgServer.NewRegion(suite.Ctx, &newRegion)

	newRegion = wstakingtypes.MsgNewRegion{
		Creator:         suite.Dao.GlobalDao,
		Name:            "USA",
		OperatorAddress: suite.usaValidator.OperatorAddress,
	}
	_, err = suite.msgServer.NewRegion(suite.Ctx, &newRegion)
	suite.Require().NoError(err)
}

const (
	dummyUpgradeHeight = 5
)

// TestUpgrade is a method of UpgradeTestSuite to test the upgrade process.
func (s *UpgradeTestSuite) TestUpgrade() {
	upgrade := func() {
		// Run upgrade
		s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight - 1)
		plan := upgradetypes.Plan{Name: v2_0_2.UpgradeName, Height: dummyUpgradeHeight}
		err := s.App.UpgradeKeeper.ScheduleUpgrade(s.Ctx, plan)
		s.Require().NoError(err)
		_, exists := s.App.UpgradeKeeper.GetUpgradePlan(s.Ctx)
		s.Require().True(exists)

		s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight)
		// simulate the upgrade process not panic.
		s.Require().NotPanics(func() {
			// simulate the upgrade process.
			s.App.BeginBlocker(s.Ctx, abci.RequestBeginBlock{})
		})
	}

	postUpgrade := func() error {
		regionMap := make(map[string]wstakingtypes.Region)
		regions := s.App.StakingKeeper.GetAllRegion(s.Ctx)
		for _, region := range regions {
			regionMap[region.RegionId] = region
		}
		return nil
	}

	upgrade()
	err := postUpgrade()
	s.Require().NoError(err)
}

// TestUpgrade is a method of UpgradeTestSuite to test the upgrade process.
func (s *UpgradeTestSuite) TestMigrateValidator() {
	upgrade := func() {
		// Run upgrade
		s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight - 1)
		plan := upgradetypes.Plan{Name: v2_0_2.UpgradeName, Height: dummyUpgradeHeight}
		err := s.App.UpgradeKeeper.ScheduleUpgrade(s.Ctx, plan)
		s.Require().NoError(err)
		_, exists := s.App.UpgradeKeeper.GetUpgradePlan(s.Ctx)
		s.Require().True(exists)

		s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight)
		// simulate the upgrade process not panic.
		s.Require().NotPanics(func() {
			// simulate the upgrade process.
			s.App.BeginBlocker(s.Ctx, abci.RequestBeginBlock{})
		})
	}

	postUpgrade := func() error {
		regionMap := make(map[string]wstakingtypes.Region)
		regions := s.App.StakingKeeper.GetAllRegion(s.Ctx)
		for _, region := range regions {
			regionMap[region.RegionId] = region
		}

		validators := s.App.StakingKeeper.GetAllValidators(s.Ctx)
		for _, validator := range validators {
			region := regionMap[validator.Description.RegionID]
			if region.OperatorAddress != validator.OperatorAddress {
				return fmt.Errorf("validator %v does not have the correct operator address, expected: %s, got: %s", validator.Description.RegionID,
					validator.OperatorAddress, region.OperatorAddress)
			}
		}
		return nil
	}

	upgrade()
	err := postUpgrade()
	s.Require().NoError(err)
}
