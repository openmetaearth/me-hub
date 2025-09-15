package keeper_test

import (
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestGravityAndBridger() {
	for _, oracle := range suite.oracleAddrs {
		require.True(suite.T(), suite.Keeper().IsProposalGravity(suite.ctx, oracle.String()))
	}

	for _, bridger := range suite.bridgerAddrs {
		oracle, found := suite.Keeper().GetGravityAddressByBridgerKey(suite.ctx, bridger)
		require.False(suite.T(), found)
		require.Equal(suite.T(), "", oracle.String())
	}
}
