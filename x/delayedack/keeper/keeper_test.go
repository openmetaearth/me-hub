package keeper_test

import (
	"testing"

	"github.com/openmetaearth/me-hub/app/apptesting"

	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
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
	app := apptesting.Setup(s.T(), false)
	ctx := app.GetBaseApp().NewContext(false, cometbftproto.Header{})

	s.App = app
	s.Ctx = ctx
}
