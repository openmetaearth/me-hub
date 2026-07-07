package keeper_test

import (
	"testing"

	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/stretchr/testify/suite"
)

const (
	delayedAckEventType = "delayedack"
)

type DelayedAckTestSuite struct {
	apptesting.KeeperTestHelper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(DelayedAckTestSuite))
}

func (s *DelayedAckTestSuite) SetupTest() {
	app := apptesting.Setup(s.T())
	ctx := app.GetBaseApp().NewContext(false)

	s.App = app
	s.Ctx = ctx
}

func (s *DelayedAckTestSuite) CreateRollappWithName(name string) {
	s.CreateRollappByName(name)
}
