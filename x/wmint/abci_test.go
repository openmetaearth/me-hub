package wmint

import (
	"fmt"
	"math/big"
	"testing"

	wbanktypes "github.com/openmetaearth/me-hub/x/wbank/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/openmetaearth/me-hub/x/wmint/keeper"
	"github.com/openmetaearth/me-hub/x/wmint/types"
	"github.com/openmetaearth/me-hub/x/wmint/types/mock_types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmtime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"
)

func TestPrintRewardInfo(t *testing.T) {
	calcPerBlockUMEC := func(mul int64) sdk.Int {
		halvingDivisor := sdk.NewDecFromBigInt(new(big.Int).Lsh(big.NewInt(1), uint(mul)))
		amount := sdk.NewDec(int64(types.InitOneYearMintAmount)).
			Quo(sdk.NewDec(int64(types.OneYearTotalBlocks))).
			Quo(halvingDivisor)
		return RoundUpToFourDecimalsDec(amount).MulInt64(100_000_000).TruncateInt()
	}

	firstUmec := calcPerBlockUMEC(0)
	firstMec := sdk.NewDecFromInt(firstUmec).QuoInt64(100_000_000)
	firstDailyUmec := firstUmec.MulRaw(int64(types.OneDayTotalBlocks))
	firstDailyMec := sdk.NewDecFromInt(firstDailyUmec).QuoInt64(100_000_000)
	fmt.Printf("first year per block reward is :%.4f mec %s umec\n", firstMec.MustFloat64(), firstUmec)
	fmt.Printf("first year daily reward is :%.4f mec %s umec\n", firstDailyMec.MustFloat64(), firstDailyUmec)

	secondUmec := calcPerBlockUMEC(1)
	secondMec := sdk.NewDecFromInt(secondUmec).QuoInt64(100_000_000)
	secondDailyUmec := secondUmec.MulRaw(int64(types.OneDayTotalBlocks))
	secondDailyMec := sdk.NewDecFromInt(secondDailyUmec).QuoInt64(100_000_000)
	fmt.Printf("second year per block reward is :%.4f mec %s umec\n", secondMec.MustFloat64(), secondUmec)
	fmt.Printf("second year daily reward is :%.4f mec %s umec\n", secondDailyMec.MustFloat64(), secondDailyUmec)
}

type KeeperTestSuite struct {
	suite.Suite
	ctx            sdk.Context
	wmintKeeper    keeper.Keeper
	bankKeeper     *mock_types.MockBankKeeper
	accKeeper      *mock_types.MockAccountKeeper
	stakingKeeper  *mock_types.MockStakingKeeper
	encCfg         moduletestutil.TestEncodingConfig
	paramsSubspace typesparams.Subspace
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	t := suite.T()
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey("test_key")
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	paramsSubspace := typesparams.NewSubspace(cdc,
		codec.NewLegacyAmino(),
		storeKey,
		memStoreKey,
		"WmintParams",
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{Time: tmtime.Now()}, false, log.NewNopLogger())
	encCfg := moduletestutil.MakeTestEncodingConfig()

	// gomock initializations

	accKeeper := mock_types.NewMockAccountKeeper(t)
	accKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(authtypes.NewModuleAddress(types.ModuleName))
	bankKeeper := mock_types.NewMockBankKeeper(t)
	stakingKeeper := mock_types.NewMockStakingKeeper(t)
	suite.ctx = ctx
	suite.encCfg = encCfg
	suite.paramsSubspace = paramsSubspace
	suite.bankKeeper = bankKeeper
	suite.stakingKeeper = stakingKeeper
	suite.accKeeper = accKeeper
	suite.wmintKeeper = keeper.NewKeeper(encCfg.Codec, storeKey, suite.stakingKeeper, suite.accKeeper, suite.bankKeeper, wbanktypes.TreasuryPoolName, authtypes.NewModuleAddress(govtypes.ModuleName).String())
}

func (suite *KeeperTestSuite) TestBeginBlocker() {
	/*
		first year per block reward is :792.7448 mec 79274480000 umec
		first year daily reward is :13698630.1440 mec 1369863014400000 umec
		second year per block reward is :396.3724 mec 39637240000 umec
		second year daily reward is :6849315.0720 mec 684931507200000 umec
	*/
	testCases := []struct {
		name           string
		targetMinted   int64
		startHeight    int64
		endHeight      int64
		perBlockReward func(height int64) int64
	}{
		{
			name:           "mint at height 1-4",
			targetMinted:   79274480000 * 4,
			perBlockReward: func(height int64) int64 { return 79274480000 },
			startHeight:    1,
			endHeight:      4,
		},
		{
			name:           "mint 1/100 of day",
			targetMinted:   79274480000 * int64(types.OneDayTotalBlocks/100),
			perBlockReward: func(height int64) int64 { return 79274480000 },
			startHeight:    1,
			endHeight:      types.OneDayTotalBlocks / 100,
		},
		{
			name:           "3 blocks at 2nd year",
			targetMinted:   39637240000 * 3,
			perBlockReward: func(height int64) int64 { return 39637240000 },
			startHeight:    types.OneYearTotalBlocks + 1,
			endHeight:      types.OneYearTotalBlocks + 3,
		},
		{
			name:           "100 blocks at 3rd year",
			targetMinted:   19818620000 * 100,
			perBlockReward: func(height int64) int64 { return 19818620000 },
			startHeight:    2*types.OneYearTotalBlocks + 101,
			endHeight:      2*types.OneYearTotalBlocks + 200,
		},
		{
			name:         "100 blocks between 2nd year (30 blocks) and 3rd year (70 blocks)",
			targetMinted: 39637240000*30 + 19818620000*70,
			perBlockReward: func(height int64) int64 {
				if height > 2*types.OneYearTotalBlocks {
					return 19818620000
				}
				return 39637240000
			},
			startHeight: 2*types.OneYearTotalBlocks - 30 + 1,
			endHeight:   2*types.OneYearTotalBlocks + 70,
		},
	}
	for _, testcase := range testCases {
		suite.wmintKeeper.SetMintedCoinAmount(suite.ctx, *big.NewInt(0))
		suite.Run(testcase.name, func() {
			ctx := suite.newContextWith(testcase.startHeight)
			var minted big.Int
			for i := testcase.startHeight; i <= testcase.endHeight; i++ {
				cctx := ctx.WithBlockHeight(i)
				suite.setMockBankKeeper(cctx, testcase.perBlockReward(i))
				BeginBlocker(cctx, suite.wmintKeeper, nil)
				minted = suite.wmintKeeper.GetMintedCoinAmount(cctx)
			}
			assert.Equal(suite.T(), testcase.targetMinted, minted.Int64())
		})
	}
}
func (suite *KeeperTestSuite) newContextWith(height int64) sdk.Context {
	return sdk.NewContext(suite.ctx.MultiStore(), tmproto.Header{Time: tmtime.Now(), Height: height}, false, log.NewNopLogger())
}
func (suite *KeeperTestSuite) setMockBankKeeper(ctx sdk.Context, mintAmount int64) {

	suite.bankKeeper.EXPECT().
		MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin("umec", sdk.NewInt(mintAmount)))).
		Return(nil)

	suite.bankKeeper.EXPECT().
		SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, "treasury_pool", sdk.NewCoins(sdk.NewCoin("umec", sdk.NewInt(mintAmount)))).
		Return(nil)
}
