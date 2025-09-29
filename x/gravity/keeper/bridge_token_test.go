package keeper_test

import (
	"github.com/st-chain/me-hub/testutil/helpers"
	"github.com/st-chain/me-hub/x/gravity/types"
)

func (suite *KeeperTestSuite) TestKeeper_BridgeToken() {
	tokenContract := helpers.GenerateAddress().Hex()

	suite.Keeper().AddBridgeToken(suite.Ctx)

	bridgeToken := &types.BridgeToken{Contract: tokenContract, Denom: denom}
	suite.EqualValues(bridgeToken, suite.Keeper().GetBridgeTokenDenom(suite.ctx, tokenContract))

	suite.EqualValues(bridgeToken, suite.Keeper().GetDenomBridgeToken(suite.ctx, denom))

	suite.Keeper().IterateBridgeTokenToDenom(suite.ctx, func(bt *types.BridgeToken) bool {
		suite.Equal(bt.Token, tokenContract)
		suite.Equal(bt.Denom, denom)
		suite.T().Log(bt.Token, bt.Denom)
		return false
	})
}
