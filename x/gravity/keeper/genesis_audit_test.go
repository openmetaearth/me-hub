package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/openmetaearth/me-hub/testutil/helpers"
	"github.com/openmetaearth/me-hub/x/gravity/keeper"
	"github.com/openmetaearth/me-hub/x/gravity/types"
	"github.com/stretchr/testify/require"
)

func TestValidateBasic_BridgeTokenDuplicates(t *testing.T) {
	contract1 := helpers.GenerateAddress().Hex()
	contract2 := helpers.GenerateAddress().Hex()

	// Test 1: Duplicate denom across different contracts
	state1 := &types.GenesisState{
		Params: types.DefaultParams(),
		BridgeTokens: []types.BridgeToken{
			{
				ContractAddress: contract1,
				Denom:           "udup",
				Supply:          sdkmath.NewInt(100),
			},
			{
				ContractAddress: contract2,
				Denom:           "udup",
				Supply:          sdkmath.NewInt(200),
			},
		},
	}
	err := state1.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "duplicate bridge token denom")

	// Test 2: Duplicate contract address across different denoms
	state2 := &types.GenesisState{
		Params: types.DefaultParams(),
		BridgeTokens: []types.BridgeToken{
			{
				ContractAddress: contract1,
				Denom:           "uone",
				Supply:          sdkmath.NewInt(100),
			},
			{
				ContractAddress: contract1,
				Denom:           "utwo",
				Supply:          sdkmath.NewInt(200),
			},
		},
	}
	err = state2.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "duplicate bridge token contract address")

	// Test 3: Valid bridge tokens
	state3 := &types.GenesisState{
		Params: types.DefaultParams(),
		BridgeTokens: []types.BridgeToken{
			{
				ContractAddress: contract1,
				Denom:           "uone",
				Supply:          sdkmath.NewInt(100),
			},
			{
				ContractAddress: contract2,
				Denom:           "utwo",
				Supply:          sdkmath.NewInt(200),
			},
		},
	}
	err = state3.ValidateBasic()
	require.NoError(t, err)
}

func (suite *KeeperTestSuite) TestInitGenesis_BridgeTokenDuplicates() {
	suite.SetupTest()
	contract1 := helpers.GenerateAddress().Hex()
	contract2 := helpers.GenerateAddress().Hex()

	// Duplicate denom should panic during InitGenesis
	state1 := &types.GenesisState{
		Params: types.DefaultParams(),
		BridgeTokens: []types.BridgeToken{
			{
				ContractAddress: contract1,
				Denom:           "udup",
				Supply:          sdkmath.NewInt(100),
			},
			{
				ContractAddress: contract2,
				Denom:           "udup",
				Supply:          sdkmath.NewInt(200),
			},
		},
	}
	suite.Panics(func() {
		keeper.InitGenesis(suite.Ctx, suite.Keeper(), state1)
	})

	// Duplicate contract address should panic during InitGenesis
	state2 := &types.GenesisState{
		Params: types.DefaultParams(),
		BridgeTokens: []types.BridgeToken{
			{
				ContractAddress: contract1,
				Denom:           "uone",
				Supply:          sdkmath.NewInt(100),
			},
			{
				ContractAddress: contract1,
				Denom:           "utwo",
				Supply:          sdkmath.NewInt(200),
			},
		},
	}
	suite.Panics(func() {
		keeper.InitGenesis(suite.Ctx, suite.Keeper(), state2)
	})
}
