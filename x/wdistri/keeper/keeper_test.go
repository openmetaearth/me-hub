package keeper_test

import (
	"fmt"
	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/golang/mock/gomock"
	"github.com/openmetaearth/me-hub/app/apptesting"
	wbanktypes "github.com/openmetaearth/me-hub/x/wbank/types"
	"github.com/openmetaearth/me-hub/x/wdistri/types"
	"github.com/openmetaearth/me-hub/x/wdistri/types/mock"
	"testing"

	"github.com/openmetaearth/me-hub/app/params"

	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/openmetaearth/me-hub/testutil/mocks"
	"github.com/openmetaearth/me-hub/x/wdistri/keeper"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/cometbft/cometbft/libs/log"
	tmtime "github.com/cometbft/cometbft/types/time"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	queryClient         distrtypes.QueryClient
	meEarthValidator    stakingtypes.Validator
	experienceValidator stakingtypes.Validator
	usaValidator        stakingtypes.Validator
	TestAccs            []sdk.AccAddress

	authKeeper    *mock.MockAccountKeeper
	bankKeeper    *mock.MockBankKeeper
	stakingKeeper *mock.MockStakingKeeper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) Keeper() *keeper.Keeper {
	return s.App.DistrKeeper
}

func (s *KeeperTestSuite) SetupTest() {
	app := apptesting.Setup(s.T(), false)
	ctx := app.GetBaseApp().NewContext(false, cometbftproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	nativeQuerier := distrkeeper.Querier{Keeper: app.DistrKeeper.Keeper}
	distrtypes.RegisterQueryServer(queryHelper, nativeQuerier)
	queryClient := distrtypes.NewQueryClient(queryHelper)
	s.queryClient = queryClient

	s.App = app
	s.Ctx = ctx

	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
	s.authKeeper = mock.NewMockAccountKeeper(ctrl)
	s.authKeeper.EXPECT().GetModuleAddress(distrtypes.ModuleName).Return(authtypes.NewModuleAddress(distrtypes.ModuleName))
	s.bankKeeper = mock.NewMockBankKeeper(ctrl)
	s.stakingKeeper = mock.NewMockStakingKeeper(ctrl)

	s.App.DistrKeeper = keeper.NewKeeper(
		s.App.AppCodec(),
		s.App.GetKey(distrtypes.StoreKey),
		s.App.GetSubspace(distrtypes.ModuleName),
		s.authKeeper,
		s.bankKeeper,
		s.stakingKeeper,
		wbanktypes.TreasuryPoolName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	s.InitializeDao()

	//validators := s.Keeper().GetValidators(s.Ctx, 10)
	//s.Require().True(len(validators) >= 3)
	//s.meEarthValidator = validators[0]
	//s.experienceValidator = validators[1]
	//s.usaValidator = validators[2]

	s.TestAccs = s.NewAccounts(3)
}

func (s *KeeperTestSuite) TestGetAuthority() {
	testsCases := []struct {
		name     string
		expected string
		success  bool
	}{
		{
			name:     "invalid account",
			expected: "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
			success:  false,
		},
		{
			name:     "valid account",
			expected: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			success:  true,
		},
	}

	for _, testCase := range testsCases {
		s.T().Run(testCase.name, func(t *testing.T) {
			actual := s.App.DistrKeeper.GetAuthority()
			if testCase.success {
				s.Require().Equal(testCase.expected, actual)
			} else {
				s.Require().NotEqual(testCase.expected, actual)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGetTreasuryModuleAccount() {
	testsCases := []struct {
		name     string
		expected string
		success  bool
	}{
		{
			name:     "invalid account",
			expected: "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
			success:  false,
		},
		{
			name:     "valida account",
			expected: wbanktypes.TreasuryPoolName,
			success:  true,
		},
	}

	for _, testCase := range testsCases {
		s.T().Run(testCase.name, func(t *testing.T) {
			actual := s.App.DistrKeeper.GetTreasuryModuleAccount()
			if testCase.success {
				s.Require().Equal(testCase.expected, actual)
			} else {
				s.Require().NotEqual(testCase.expected, actual)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestEndBlocker() {
	/*
		first year per block reward is :792.7448 mec 79274480000 umec
		first year daily reward is :13698630.1440 mec 1369863014400000 umec
		second year per block reward is :396.3724 mec 39637240000 umec
		second year daily reward is :6849315.0720 mec 684931507200000 umec
	*/
	testsCases := []struct {
		name                string
		height              int
		regionShares        []int
		regionWantGetReward []int
		totalReward         int
	}{
		{
			name:                "two region with equal share",
			height:              types.OneDayTotalBlocks,
			regionShares:        []int{1, 1},
			regionWantGetReward: []int{684931507200000, 684931507200000}, // umec
		},
		{
			name:                "one region should get all reward",
			height:              types.OneDayTotalBlocks * 2,
			regionShares:        []int{1},
			regionWantGetReward: []int{1369863014400000}, // umec
		},
		{
			name:                "one region should get all reward",
			height:              types.OneDayTotalBlocks * 3,
			regionShares:        []int{1, 2, 2},
			regionWantGetReward: []int{273972602880000, 547945205760000, 547945205760000}, // umec
		},
		{
			name:                "not trigger distribution",
			height:              types.OneDayTotalBlocks / 2,
			regionShares:        []int{},
			regionWantGetReward: []int{}, // umec
		},
		{
			name:                "second year first day",
			height:              366 * types.OneDayTotalBlocks,
			regionShares:        []int{1, 1},
			regionWantGetReward: []int{342465753600000, 342465753600000}, // umec
		},
		{
			name:                "second year first half of day",
			height:              366*types.OneDayTotalBlocks + types.OneDayTotalBlocks/2,
			regionShares:        []int{},
			regionWantGetReward: []int{}, // umec
		},
	}
	runCase := func(index int) {
		testcase := testsCases[index]
		ctx := suite.HelperNewContextWith(int64(testcase.height))
		addrs := suite.mockGetRegionI(ctx, testcase.regionShares...)
		var wantReward []coinAndAddr
		totalWantReward := 0
		for i, addr := range addrs {
			wantReward = append(wantReward, coinAndAddr{
				num:  int64(testcase.regionWantGetReward[i]),
				addr: addr,
			})
			totalWantReward += testcase.regionWantGetReward[i]
		}
		if totalWantReward != 0 {
			suite.SetMockGetBalance(ctx, sdk.NewInt(int64(totalWantReward)))
		}
		suite.setMockSendCoinsFromModuleToAccountExpect(ctx, wantReward...)

		err := suite.App.DistrKeeper.AllocateBlockRewardEveryday(ctx, abci.RequestEndBlock{Height: ctx.BlockHeight()})
		events := ctx.EventManager().ABCIEvents()
		suite.Require().NoError(err, "case %d: %s", index, testcase.name)
		assert.Equal(suite.T(), len(addrs), len(events))
	}
	for i := range testsCases {
		suite.Run(testsCases[i].name, func() {
			runCase(i)
		})
	}
}

func (suite *KeeperTestSuite) mockGetRegionI(ctx sdk.Context, regionShare ...int) []string {
	var addrs []string
	if len(regionShare) == 0 {
		return addrs
	}
	var regions []wstakingtypes.RegionI
	for i, share := range regionShare {
		region := mocks.NewMockRegionI(suite.T())
		region.EXPECT().GetRegionShare().Return(sdk.NewInt(int64(share)))
		addr := authtypes.NewModuleAddress(fmt.Sprintf("region_%d", i)).String()
		addrs = append(addrs, addr)
		region.EXPECT().GetRegionTreasureAddr().Return(addr)
		region.EXPECT().GetRegionId().Return(fmt.Sprintf("region_ID_%d", i))
		regions = append(regions, region)
	}
	suite.stakingKeeper.EXPECT().GetAllRegionI(ctx).Return(regions)
	return addrs
}

func (suite *KeeperTestSuite) SetMockGetBalance(ctx sdk.Context, amount sdkmath.Int) {
	acc := authtypes.NewModuleAddress(suite.App.DistrKeeper.GetTreasuryModuleAccount())
	suite.authKeeper.EXPECT().GetModuleAddress(suite.App.DistrKeeper.GetTreasuryModuleAccount()).Return(acc)
	suite.bankKeeper.EXPECT().GetAllBalances(ctx, acc).Return(sdk.NewCoins(sdk.NewCoin(params.BaseDenom, amount)))
}

func (suite *KeeperTestSuite) HelperNewContextWith(height int64) sdk.Context {
	return sdk.NewContext(suite.Ctx.MultiStore(), tmproto.Header{Time: tmtime.Now(), Height: height}, false, log.NewNopLogger())
}

type coinAndAddr struct {
	num  int64
	addr string
}

func (suite *KeeperTestSuite) setMockSendCoinsFromModuleToAccountExpect(ctx sdk.Context, want ...coinAndAddr) {
	baseDenom := params.BaseDenom
	for _, w := range want {
		suite.bankKeeper.EXPECT().
			SendCoinsFromModuleToAccount(
				ctx,
				suite.App.DistrKeeper.GetTreasuryModuleAccount(),
				sdk.MustAccAddressFromBech32(w.addr),
				sdk.NewCoins(sdk.NewCoin(baseDenom, sdk.NewInt(w.num))),
			).Return(nil)
	}
}
