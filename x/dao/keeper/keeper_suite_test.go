package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/stretchr/testify/suite"

	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/x/dao/keeper"
	"github.com/openmetaearth/me-hub/x/dao/types"
)

type DaoKeeperTestSuite struct {
	apptesting.KeeperTestHelper
	msgServer   types.MsgServer
	queryClient types.QueryClient
}

func (suite *DaoKeeperTestSuite) SetupTest() {
	app := apptesting.Setup(suite.T())
	ctx := app.GetBaseApp().NewContext(false)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.DaoKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.App = app
	suite.msgServer = keeper.NewMsgServerImpl(app.DaoKeeper)
	suite.Ctx = ctx
	suite.queryClient = queryClient
}

func TestDaoKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(DaoKeeperTestSuite))
}
