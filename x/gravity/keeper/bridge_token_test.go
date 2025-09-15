package keeper_test

import (
	"fmt"

	"github.com/st-chain/me-hub/testutil/helpers"
	"github.com/st-chain/me-hub/x/gravity/types"
)

func (suite *KeeperTestSuite) TestKeeper_BridgeToken() {
	tokenContract := helpers.GenerateAddress().Hex()
	denom, err := suite.Keeper().SetIbcDenomTrace(suite.ctx, tokenContract, "")
	suite.NoError(err)
	suite.Equal(fmt.Sprintf("%s%s", suite.chainName, tokenContract), denom)

	suite.Keeper().AddBridgeToken(suite.ctx, tokenContract, denom)

	bridgeToken := &types.BridgeToken{Token: tokenContract, Denom: denom}
	suite.EqualValues(bridgeToken, suite.Keeper().GetBridgeTokenDenom(suite.ctx, tokenContract))

	suite.EqualValues(bridgeToken, suite.Keeper().GetDenomBridgeToken(suite.ctx, denom))

	suite.Keeper().IterateBridgeTokenToDenom(suite.ctx, func(bt *types.BridgeToken) bool {
		suite.Equal(bt.Token, tokenContract)
		suite.Equal(bt.Denom, denom)
		suite.T().Log(bt.Token, bt.Denom)
		return false
	})
}
