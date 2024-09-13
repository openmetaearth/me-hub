package keeper_test

import (
	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/st-chain/me-hub/app/apptesting"
	"github.com/st-chain/me-hub/app/params"
	testutilstypes "github.com/st-chain/me-hub/testutil/types"
	"github.com/st-chain/me-hub/x/wstaking/keeper"
	wstakingkeeper "github.com/st-chain/me-hub/x/wstaking/keeper"
	"github.com/st-chain/me-hub/x/wstaking/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	msgServer           keeper.MsgServer
	queryClient         types.QueryClient
	meEarthValidator    stakingtypes.Validator
	experienceValidator stakingtypes.Validator
	usaValidator        stakingtypes.Validator
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) Keeper() *wstakingkeeper.Keeper {
	return suite.App.StakingKeeper
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
	nativeQuerier := keeper.Querier{Keeper: app.StakingKeeper}
	types.RegisterQueryServer(queryHelper, nativeQuerier)
	queryClient := types.NewQueryClient(queryHelper)

	suite.App = app

	stakingKeeperMsgSrv := stakingkeeper.NewMsgServerImpl(app.StakingKeeper.Keeper)
	suite.msgServer = keeper.NewMsgServerImpl(app.StakingKeeper, stakingKeeperMsgSrv)
	suite.Ctx = ctx
	suite.queryClient = queryClient

	suite.InitializeDao()

	validators := suite.Keeper().GetValidators(suite.Ctx, 10)
	suite.Require().True(len(validators) >= 3)
	suite.meEarthValidator = validators[0]
	suite.experienceValidator = validators[1]
	suite.usaValidator = validators[2]

	newRegion := types.MsgNewRegion{
		Creator:         suite.Dao.GlobalDao,
		Name:            types.ExperienceRegionName,
		OperatorAddress: suite.experienceValidator.OperatorAddress,
	}
	_, err = suite.msgServer.NewRegion(suite.Ctx, &newRegion)
	suite.Require().NoError(err)
}

func SetValidatorV1(ctx sdk.Context, k *keeper.Keeper, validator testutilstypes.ValidatorV1) {
	store := ctx.KVStore(k.GetStoreKey())
	bz := k.GetCdc().MustMarshal(&validator)
	addr, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
	if err != nil {
		panic(err)
	}
	store.Set(stakingtypes.GetValidatorKey(addr), bz)
}

func GetValidatorV2(ctx sdk.Context, k *keeper.Keeper, addr sdk.ValAddress) (validator testutilstypes.ValidatorV2, found bool) {
	store := ctx.KVStore(k.GetStoreKey())
	value := store.Get(stakingtypes.GetValidatorKey(addr))
	if value == nil {
		return validator, false
	}
	err := k.GetCdc().Unmarshal(value, &validator)
	if err != nil {
		panic(err)
	}
	return validator, true
}

func (suite *KeeperTestSuite) TestMigrateValidator() {
	val1 := testutilstypes.ValidatorV1{
		OperatorAddress: "mevaloper139mq752delxv78jvtmwxhasyrycufsvr707ate",
		ConsensusPubkey: nil,
		Jailed:          false,
		Status:          stakingtypes.Bonded,
		Tokens:          sdk.NewInt(100),
		StakerShares:    sdk.NewDec(100),
		Description: stakingtypes.Description{
			Moniker:         "node1",
			Identity:        "",
			Website:         "",
			SecurityContact: "",
			Details:         "",
			RegionId:        "usa",
		},
		UnbondingHeight:         0,
		UnbondingTime:           time.Time{},
		Commission:              stakingtypes.Commission{},
		MinSelfStake:            sdk.Int{},
		DelegationAmount:        sdk.Int{},
		MeidAmount:              sdk.Int{},
		OwnerAddress:            "me139mq752delxv78jvtmwxhasyrycufsvr0mue6u",
		UnbondingIds:            nil,
		UnbondingOnHoldRefCount: 0,
	}
	SetValidatorV1(suite.Ctx, suite.App.StakingKeeper, val1)
	suite.T().Log(val1.String())

	addr, err := sdk.ValAddressFromBech32(val1.OperatorAddress)
	if err != nil {
		panic(err)
	}
	//test panicked: proto: wrong wireType = 2 for field UnbondingOnHoldRefCount
	validator, found := GetValidatorV2(suite.Ctx, suite.App.StakingKeeper, addr)
	require.True(suite.T(), found)

	validators := suite.App.StakingKeeper.GetAllValidators(suite.Ctx)
	require.Equal(suite.T(), len(validators), 4)
	for _, v := range validators {
		if v.OperatorAddress == validator.OperatorAddress {
			suite.T().Log(validator.String())
			require.Equal(suite.T(), validator.String(), v.String())
		}
	}
}
