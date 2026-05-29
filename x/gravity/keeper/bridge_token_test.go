package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/testutil/helpers"
	"github.com/openmetaearth/me-hub/x/gravity/keeper"
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

func (suite *KeeperTestSuite) TestGenesisValidationAndImport() {
	// 1. Test same denom with two contracts in ValidateBasic
	genesis1 := &types.GenesisState{
		Params: types.DefaultParams(),
		BridgeTokens: []types.BridgeToken{
			{
				ContractAddress: helpers.GenerateAddress().Hex(),
				Denom:           "testdenom",
				Name:            "Test 1",
				Symbol:          "T1",
				Decimal:         18,
				Supply:          sdk.NewInt(100),
			},
			{
				ContractAddress: helpers.GenerateAddress().Hex(),
				Denom:           "testdenom",
				Name:            "Test 2",
				Symbol:          "T2",
				Decimal:         18,
				Supply:          sdk.NewInt(200),
			},
		},
	}
	suite.Error(genesis1.ValidateBasic())

	// 2. Test same contract with two denoms in ValidateBasic
	contractAddr := helpers.GenerateAddress().Hex()
	genesis2 := &types.GenesisState{
		Params: types.DefaultParams(),
		BridgeTokens: []types.BridgeToken{
			{
				ContractAddress: contractAddr,
				Denom:           "denom1",
				Name:            "Test 1",
				Symbol:          "T1",
				Decimal:         18,
				Supply:          sdk.NewInt(100),
			},
			{
				ContractAddress: contractAddr,
				Denom:           "denom2",
				Name:            "Test 2",
				Symbol:          "T2",
				Decimal:         18,
				Supply:          sdk.NewInt(200),
			},
		},
	}
	suite.Error(genesis2.ValidateBasic())

	// 3. Test defensive check in InitGenesis
	suite.Panics(func() {
		keeper.InitGenesis(suite.Ctx, suite.Keeper(), genesis1)
	})
}

func (suite *KeeperTestSuite) TestGenesisExportImport() {
	genesis := &types.GenesisState{
		Params: types.DefaultParams(),
		BridgeTokens: []types.BridgeToken{
			{
				ContractAddress: helpers.GenerateAddress().Hex(),
				Denom:           "canonicaldenom",
				Name:            "Canonical",
				Symbol:          "CAN",
				Decimal:         18,
				Supply:          sdk.NewInt(500),
			},
		},
	}

	// 1. Import
	keeper.InitGenesis(suite.Ctx, suite.Keeper(), genesis)

	// 2. Verify keeper state
	token, err := suite.Keeper().GetBridgeTokenByDenom(suite.Ctx, "canonicaldenom")
	suite.NoError(err)
	suite.Equal("canonicaldenom", token.Denom)

	// 3. Export
	exported := keeper.ExportGenesis(suite.Ctx, suite.Keeper())
	suite.Len(exported.BridgeTokens, 1)
	suite.Equal("canonicaldenom", exported.BridgeTokens[0].Denom)

	// 4. Re-import after resetting context
	suite.SetupTest()
	keeper.InitGenesis(suite.Ctx, suite.Keeper(), exported)
	tokenReimported, err := suite.Keeper().GetBridgeTokenByDenom(suite.Ctx, "canonicaldenom")
	suite.NoError(err)
	suite.Equal("canonicaldenom", tokenReimported.Denom)
}


