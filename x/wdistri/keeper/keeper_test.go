package keeper

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/st-chain/me-hub/x/wdistri/types"
	"github.com/st-chain/me-hub/x/wdistri/types/mock_types"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/baseapp"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

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
	ctx           sdk.Context
	wdistriKeeper *Keeper
	authKeeper    *mock_types.MockAccountKeeper
	bankKeeper    *mock_types.MockBankKeeper
	stakingKeeper *mock_types.MockStakingKeeper
	queryClient   types.QueryClient
	msgServer     types.MsgServer
	encCfg        moduletestutil.TestEncodingConfig
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	t := suite.T()
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
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

	suite.ctx = ctx
	suite.encCfg = encCfg
	suite.authKeeper = mock_types.NewMockAccountKeeper(t)
	suite.bankKeeper = mock_types.NewMockBankKeeper(t)
	suite.stakingKeeper = mock_types.NewMockStakingKeeper(t)
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	suite.queryClient = types.NewQueryClient(queryHelper)
	suite.wdistriKeeper = NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		suite.authKeeper,
		suite.bankKeeper,
		suite.stakingKeeper,
		"treasury_pool",
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	suite.msgServer = NewMsgServerImpl(*suite.wdistriKeeper)
}
