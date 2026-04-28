package keeper_test

import (
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/suite"

	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/x/rollapp/keeper"
	"github.com/openmetaearth/me-hub/x/rollapp/types"

	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

const (
	transferEventCount            = 3 // As emitted by the bank
	createEventCount              = 8
	playEventCountFirst           = 8 // Extra "sender" attribute emitted by the bank
	playEventCountNext            = 7
	rejectEventCount              = 4
	rejectEventCountWithTransfer  = 5 // Extra "sender" attribute emitted by the bank
	forfeitEventCount             = 4
	forfeitEventCountWithTransfer = 5 // Extra "sender" attribute emitted by the bank
	alice                         = "me139mq752delxv78jvtmwxhasyrycufsvr0mue6u"
	bob                           = "me1fxv6zn3my807ps6ph5va48rn0m8zurmp036zlf"
	carol                         = "me1es57vt7wnjyuedepvd4dt7aa7gnvus5hvrmyld"
	balAlice                      = 50000000
	balBob                        = 20000000
	balCarol                      = 10000000
	foreignToken                  = "foreignToken"
	balTokenAlice                 = 5
	balTokenBob                   = 2
	balTokenCarol                 = 1
)

var rollappModuleAddress string

type RollappTestSuite struct {
	apptesting.KeeperTestHelper
	msgServer   types.MsgServer
	queryClient types.QueryClient
}

func (suite *RollappTestSuite) SetupTest(deployerWhitelist ...types.DeployerParams) {
	app := apptesting.Setup(suite.T(), false)
	ctx := app.GetBaseApp().NewContext(false, cometbftproto.Header{})

	err := app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	suite.Require().NoError(err)
	err = app.BankKeeper.SetParams(ctx, banktypes.DefaultParams())
	suite.Require().NoError(err)
	app.RollappKeeper.SetParams(ctx, types.NewParams(true, 2, deployerWhitelist))
	rollappModuleAddress = app.AccountKeeper.GetModuleAddress(types.ModuleName).String()

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.RollappKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.App = app
	suite.msgServer = keeper.NewMsgServerImpl(*app.RollappKeeper)
	suite.Ctx = ctx
	suite.queryClient = queryClient
}

func TestRollappKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(RollappTestSuite))
}

func createNRollapp(keeper *keeper.Keeper, ctx sdk.Context, n int) ([]types.Rollapp, []types.RollappSummary) {
	items := make([]types.Rollapp, n)
	for i := range items {
		items[i].RollappId = strconv.Itoa(i)
		keeper.SetRollapp(ctx, items[i])
	}

	rollappSummaries := []types.RollappSummary{}
	for _, item := range items {
		rollappSummary := types.RollappSummary{
			RollappId: item.RollappId,
		}
		rollappSummaries = append(rollappSummaries, rollappSummary)
	}

	return items, rollappSummaries
}
