package keeper_test

import (
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestGravityAndBridger() {
	for _, relayer := range suite.relayerAddrs {
		require.True(suite.T(), suite.Keeper().IsProposalRelayer(suite.Ctx, relayer.String()))
	}
}
