package keeper

import (
	"fmt"
	wbanktypes "github.com/st-chain/me-hub/x/wbank/types"
	"testing"

	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/st-chain/me-hub/testutil/mocks"
	"github.com/st-chain/me-hub/x/wdistri/types"
	"github.com/st-chain/me-hub/x/wdistri/types/mock_types"
	wstakingtypes "github.com/st-chain/me-hub/x/wstaking/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmtime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
)

type KeeperTestSuite struct {
	suite.Suite
	ctx            sdk.Context
	wdistriKeeper  *Keeper
	authKeeper     *mock_types.MockAccountKeeper
	bankKeeper     *mock_types.MockBankKeeper
	stakingKeeper  *mock_types.MockStakingKeeper
	queryClient    types.QueryClient
	msgServer      types.MsgServer
	encCfg         moduletestutil.TestEncodingConfig
	paramsSubspace typesparams.Subspace
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	t := suite.T()
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey("transient")
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"WdistriParams",
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{Time: tmtime.Now()}, false, log.NewNopLogger())
	encCfg := moduletestutil.MakeTestEncodingConfig()

	// gomock initializations

	authKeeper := mock_types.NewMockAccountKeeper(t)
	authKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(authtypes.NewModuleAddress(types.ModuleName))
	bankKeeper := mock_types.NewMockBankKeeper(t)
	stakingKeeper := mock_types.NewMockStakingKeeper(t)

	suite.ctx = ctx
	suite.encCfg = encCfg
	suite.paramsSubspace = paramsSubspace
	suite.authKeeper = authKeeper
	suite.bankKeeper = bankKeeper
	suite.stakingKeeper = stakingKeeper
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	suite.queryClient = types.NewQueryClient(queryHelper)
	suite.wdistriKeeper = NewKeeper(
		cdc,
		storeKey,
		paramsSubspace,
		suite.authKeeper,
		suite.bankKeeper,
		suite.stakingKeeper,
		wbanktypes.TreasuryPoolName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	suite.msgServer = NewMsgServerImpl(*suite.wdistriKeeper)
}

func (suite *KeeperTestSuite) TestGetAuthority() {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)

	NewKeeperWithAuthority := func(authority string) *Keeper {
		return NewKeeper(
			suite.encCfg.Codec,
			storeKey,
			suite.paramsSubspace,
			suite.authKeeper,
			suite.bankKeeper,
			suite.stakingKeeper,
			wbanktypes.TreasuryPoolName,
			authority,
		)
	}

	tests := map[string]string{
		"some random account": "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
		"gov module account":  authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	}

	for name, expected := range tests {
		suite.T().Run(name, func(t *testing.T) {
			kpr := NewKeeperWithAuthority(expected)
			actual := kpr.GetAuthority()
			suite.Require().Equal(expected, actual)
		})
	}
}

func (suite *KeeperTestSuite) TestGetTreasuryPool() {
	err := suite.wdistriKeeper.Hooks().BeforeDelegationSharesModified(suite.ctx, sdk.AccAddress{}, sdk.ValAddress{})
	assert.NoError(suite.T(), err)
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
			height:              oneDayTotalBlocks,
			regionShares:        []int{1, 1},
			regionWantGetReward: []int{684931507200000, 684931507200000}, // umec
		},
		{
			name:                "one region should get all reward",
			height:              oneDayTotalBlocks * 2,
			regionShares:        []int{1},
			regionWantGetReward: []int{1369863014400000}, // umec
		},
		{
			name:                "one region should get all reward",
			height:              oneDayTotalBlocks * 3,
			regionShares:        []int{1, 2, 2},
			regionWantGetReward: []int{273972602880000, 547945205760000, 547945205760000}, // umec
		},
		{
			name:                "not trigger distribution",
			height:              oneDayTotalBlocks / 2,
			regionShares:        []int{},
			regionWantGetReward: []int{}, // umec
		},
		{
			name:                "second year first day",
			height:              366 * oneDayTotalBlocks,
			regionShares:        []int{1, 1},
			regionWantGetReward: []int{342465753600000, 342465753600000}, // umec
		},
		{
			name:                "second year first half of day",
			height:              366*oneDayTotalBlocks + oneDayTotalBlocks/2,
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
		suite.wdistriKeeper.AllocateBlockRewardEveryday(ctx, abci.RequestEndBlock{Height: ctx.BlockHeight()})
		events := ctx.EventManager().ABCIEvents()
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
func (suite *KeeperTestSuite) SetMockGetBalance(ctx sdk.Context, fee sdkmath.Int) {
	acc := authtypes.NewModuleAddress(suite.wdistriKeeper.feeCollectorName)
	suite.authKeeper.EXPECT().GetModuleAddress(suite.wdistriKeeper.feeCollectorName).Return(acc)
	suite.bankKeeper.EXPECT().GetBalance(ctx, acc, suite.wdistriKeeper.baseDenom).Return(sdk.NewCoin(suite.wdistriKeeper.baseDenom, fee))
}
func (suite *KeeperTestSuite) HelperNewContextWith(height int64) sdk.Context {
	return sdk.NewContext(suite.ctx.MultiStore(), tmproto.Header{Time: tmtime.Now(), Height: height}, false, log.NewNopLogger())
}

type coinAndAddr struct {
	num  int64
	addr string
}

func (suite *KeeperTestSuite) setMockSendCoinsFromModuleToAccountExpect(ctx sdk.Context, want ...coinAndAddr) {
	baseDenom := suite.wdistriKeeper.baseDenom
	for _, w := range want {
		suite.bankKeeper.EXPECT().
			SendCoinsFromModuleToAccount(
				ctx,
				suite.wdistriKeeper.feeCollectorName,
				sdk.MustAccAddressFromBech32(w.addr),
				sdk.NewCoins(sdk.NewCoin(baseDenom, sdk.NewInt(w.num))),
			).Return(nil)
	}
}
