package keeper_test

import (
	"github.com/openmetaearth/me-hub/testutil/helpers"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

func (suite *KeeperTestSuite) TestKeeper_BridgeToken() {
	tokenContract := helpers.GenerateAddress().Hex()
	tokenContract2 := helpers.GenerateAddress().Hex()
	denom := "test"

	bridgeToken := &types.BridgeToken{ContractAddress: tokenContract, Denom: denom}
	suite.Keeper().SetBridgeToken(suite.Ctx, bridgeToken)

	suite.Keeper().SetBridgeToken(suite.Ctx, &types.BridgeToken{ContractAddress: tokenContract2, Denom: "test2"})

	b1, err := suite.Keeper().GetBridgeTokenByDenom(suite.Ctx, denom)
	suite.NoError(err)
	suite.EqualValues(bridgeToken, b1)

	b2, err := suite.Keeper().GetBridgeTokenByContract(suite.Ctx, tokenContract)
	suite.NoError(err)
	suite.EqualValues(bridgeToken, b2)

	suite.Keeper().IterateBridgeTokenByDenom(suite.Ctx, func(bt *types.BridgeToken) bool {
		if bt.Denom == denom {
			suite.Equal(bt.ContractAddress, tokenContract)
		}
		if bt.Denom == "test2" {
			suite.Equal(bt.ContractAddress, tokenContract2)
		}
		return false
	})
}
