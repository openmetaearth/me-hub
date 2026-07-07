package hooks_test

import (
	"testing"

	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/stretchr/testify/suite"
)

type HooksTestSuite struct {
	apptesting.KeeperTestHelper
}

func TestHooksTestSuite(t *testing.T) {
	suite.Run(t, new(HooksTestSuite))
}

func (suite *HooksTestSuite) SetupTest() {
	app := apptesting.Setup(suite.T())
	ctx := app.GetBaseApp().NewContext(false)

	suite.App = app
	suite.Ctx = ctx
}
