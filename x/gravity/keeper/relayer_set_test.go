package keeper_test

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/st-chain/me-hub/testutil/helpers"
	"github.com/st-chain/me-hub/x/gravity/types"
)

func (suite *KeeperTestSuite) TestLastPendingGravitySetRequestByAddr() {
	testCases := []struct {
		GravityAddress sdk.AccAddress
		BridgerAddress sdk.AccAddress
		StartHeight    int64

		ExpectGravitySetSize int
	}{
		{
			GravityAddress:       suite.oracleAddrs[0],
			BridgerAddress:       suite.bridgerAddrs[0],
			StartHeight:          1,
			ExpectGravitySetSize: 3,
		},
		{
			GravityAddress:       suite.oracleAddrs[1],
			BridgerAddress:       suite.bridgerAddrs[1],
			StartHeight:          2,
			ExpectGravitySetSize: 2,
		},
		{
			GravityAddress:       suite.oracleAddrs[2],
			BridgerAddress:       suite.bridgerAddrs[2],
			StartHeight:          3,
			ExpectGravitySetSize: 1,
		},
	}

	for i := 1; i <= 3; i++ {
		suite.Keeper().StoreGravitySet(suite.ctx, &types.GravitySet{
			Nonce: uint64(i),
			Members: types.BridgeValidators{{
				Power:           uint64(i),
				ExternalAddress: fmt.Sprintf("0x%d", i),
			}},
			Height: uint64(i),
		})
	}

	wrapSDKContext := sdk.WrapSDKContext(suite.ctx)
	for _, testCase := range testCases {
		oracle := types.Gravity{
			GravityAddress: testCase.GravityAddress.String(),
			BridgerAddress: testCase.BridgerAddress.String(),
			StartHeight:    testCase.StartHeight,
		}
		// save oracle
		suite.Keeper().SetGravity(suite.ctx, oracle)
		suite.Keeper().SetGravityByBridger(suite.ctx, testCase.BridgerAddress, oracle.GetGravity())

		response, err := suite.Keeper().LastPendingGravitySetRequestByAddr(wrapSDKContext,
			&types.QueryLastPendingGravitySetRequestByAddrRequest{
				BridgerAddress: testCase.BridgerAddress.String(),
			})
		require.NoError(suite.T(), err)
		require.EqualValues(suite.T(), testCase.ExpectGravitySetSize, len(response.GravitySets))
	}
}

func (suite *KeeperTestSuite) TestGetUnSlashedGravitySets() {
	height := 100
	index := 10
	for i := 1; i <= index; i++ {
		suite.Keeper().StoreGravitySet(suite.ctx, &types.GravitySet{
			Nonce: uint64(i),
			Members: types.BridgeValidators{{
				Power:           tmrand.Uint64(),
				ExternalAddress: helpers.GenerateAddress().Hex(),
			}},
			Height: uint64(height + i),
		})
	}
	suite.Equal(uint64(0), suite.Keeper().GetLastSlashedGravitySetNonce(suite.ctx))

	sets := suite.Keeper().GetUnSlashedGravitySets(suite.ctx, uint64(height+index))
	require.EqualValues(suite.T(), index-1, sets.Len())

	suite.Keeper().SetLastSlashedGravitySetNonce(suite.ctx, 1)
	sets = suite.Keeper().GetUnSlashedGravitySets(suite.ctx, uint64(height+index))
	require.EqualValues(suite.T(), index-2, sets.Len())

	sets = suite.Keeper().GetUnSlashedGravitySets(suite.ctx, uint64(height+index+1))
	require.EqualValues(suite.T(), index-1, sets.Len())
}

func (suite *KeeperTestSuite) TestKeeper_IterateGravitySetConfirmByNonce() {
	index := tmrand.Intn(20) + 1
	for i := uint64(1); i <= uint64(index); i++ {
		for _, oracle := range suite.oracleAddrs {
			suite.Keeper().SetGravitySetConfirm(suite.ctx, oracle,
				&types.MsgGravitySetConfirm{
					Nonce:           i,
					BridgerAddress:  sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
					ExternalAddress: helpers.GenerateAddress().Hex(),
					Signature:       "",
					ChainName:       suite.chainName,
				},
			)
		}
	}

	index = tmrand.Intn(index) + 1
	var confirms []*types.MsgGravitySetConfirm
	suite.Keeper().IterateGravitySetConfirmByNonce(suite.ctx, uint64(index), func(confirm *types.MsgGravitySetConfirm) bool {
		confirms = append(confirms, confirm)
		return false
	})
	suite.Equal(len(confirms), len(suite.oracleAddrs), index)
}

func (suite *KeeperTestSuite) TestKeeper_DeleteGravitySetConfirm() {
	member := make([]types.BridgeValidator, 0, len(suite.externalPris))
	for i, external := range suite.externalPris {
		externalAddr := suite.PubKeyToExternalAddr(external.PublicKey)
		member = append(member, types.BridgeValidator{
			Power:           uint64(i),
			ExternalAddress: externalAddr,
		})
	}
	oracleSet := &types.GravitySet{
		Nonce:   1,
		Members: member,
		Height:  uint64(suite.ctx.BlockHeight()),
	}
	suite.Keeper().StoreGravitySet(suite.ctx, oracleSet)

	for i, external := range suite.externalPris {
		externalAddress, signature := suite.SignGravitySetConfirm(external, oracleSet)
		suite.Keeper().SetGravitySetConfirm(suite.ctx, suite.oracleAddrs[i],
			&types.MsgGravitySetConfirm{
				Nonce:           oracleSet.Nonce,
				BridgerAddress:  suite.bridgerAddrs[i].String(),
				ExternalAddress: externalAddress,
				Signature:       hex.EncodeToString(signature),
				ChainName:       suite.chainName,
			},
		)
	}
	suite.Keeper().SetLastObservedGravitySet(suite.ctx,
		&types.GravitySet{
			Nonce: oracleSet.Nonce + 1,
		},
	)

	params := suite.Keeper().GetParams(suite.ctx)
	params.SignedWindow = 10
	err := suite.Keeper().SetParams(suite.ctx, &params)
	suite.Require().NoError(err)
	suite.Commit()
	for _, oracle := range suite.oracleAddrs {
		suite.NotNil(suite.Keeper().GetGravitySetConfirm(suite.ctx, oracleSet.Nonce, oracle))
	}

	suite.Commit(int64(params.SignedWindow + 1))
	for _, oracle := range suite.oracleAddrs {
		suite.Nil(suite.Keeper().GetGravitySetConfirm(suite.ctx, oracleSet.Nonce, oracle))
	}
}

func (suite *KeeperTestSuite) TestKeeper_IterateGravitySet() {
	member := make([]types.BridgeValidator, 0, len(suite.externalPris))
	for i, external := range suite.externalPris {
		member = append(member, types.BridgeValidator{
			Power:           uint64(i),
			ExternalAddress: crypto.PubkeyToAddress(external.PublicKey).String(),
		})
	}
	for i := 1; i <= 10; i++ {
		suite.Keeper().StoreGravitySet(suite.ctx, &types.GravitySet{
			Nonce:   uint64(i),
			Members: member,
			Height:  uint64(i + 100),
		})
	}
	i := uint64(0)
	oracleSets := types.GravitySets{}
	suite.Keeper().IterateGravitySetByNonce(suite.ctx, 0, func(oracleSet *types.GravitySet) bool {
		i = i + 1
		suite.Equal(i, oracleSet.Nonce)
		oracleSets = append(oracleSets, oracleSet)
		return false
	})
	suite.Equal(len(oracleSets), 10)

	oracleSets = types.GravitySets{}
	suite.Keeper().IterateGravitySetByNonce(suite.ctx, 1, func(oracleSet *types.GravitySet) bool {
		oracleSets = append(oracleSets, oracleSet)
		return false
	})
	suite.Equal(len(oracleSets), 10)

	oracleSets = types.GravitySets{}
	suite.Keeper().IterateGravitySetByNonce(suite.ctx, 2, func(oracleSet *types.GravitySet) bool {
		oracleSets = append(oracleSets, oracleSet)
		return false
	})
	suite.Equal(len(oracleSets), 9)

	suite.Keeper().IterateGravitySets(suite.ctx, true, func(oracleSet *types.GravitySet) bool {
		suite.Equal(i, oracleSet.Nonce, oracleSet.Nonce)
		i = i - 1
		return false
	})

	suite.Keeper().IterateGravitySets(suite.ctx, false, func(oracleSet *types.GravitySet) bool {
		i = i + 1
		suite.Equal(i, oracleSet.Nonce)
		return false
	})
}
