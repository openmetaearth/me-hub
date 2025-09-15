package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/st-chain/me-hub/testutil/helpers"
	"github.com/st-chain/me-hub/x/gravity/types"
)

func (suite *KeeperTestSuite) TestLastPendingBatchRequestByAddr() {
	testCases := []struct {
		Name              string
		GravityAddress    sdk.AccAddress
		BridgerAddress    sdk.AccAddress
		StartHeight       int64
		ExpectStartHeight uint64
	}{
		{
			Name:              "oracle start height with 1, expect oracle set block 3",
			GravityAddress:    suite.oracleAddrs[0],
			BridgerAddress:    suite.bridgerAddrs[0],
			StartHeight:       1,
			ExpectStartHeight: 3,
		},
		{
			Name:              "oracle start height with 2, expect oracle set block 2",
			GravityAddress:    suite.oracleAddrs[1],
			BridgerAddress:    suite.bridgerAddrs[1],
			StartHeight:       2,
			ExpectStartHeight: 3,
		},
		{
			Name:              "oracle start height with 3, expect oracle set block 1",
			GravityAddress:    suite.oracleAddrs[2],
			BridgerAddress:    suite.bridgerAddrs[2],
			StartHeight:       3,
			ExpectStartHeight: 3,
		},
	}
	for i := uint64(1); i <= 3; i++ {
		suite.ctx = suite.ctx.WithBlockHeight(int64(i))
		err := suite.Keeper().StoreBatch(suite.ctx, &types.OutgoingTxBatch{
			Block:      i,
			BatchNonce: i,
			Transactions: types.OutgoingTransferTxs{{
				Id:          i,
				Sender:      sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
				DestAddress: helpers.GenerateAddress().Hex(),
			}},
		})
		require.NoError(suite.T(), err)
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

		response, err := suite.Keeper().LastPendingBatchRequestByAddr(wrapSDKContext,
			&types.QueryLastPendingBatchRequestByAddrRequest{
				BridgerAddress: testCase.BridgerAddress.String(),
			})
		require.NoError(suite.T(), err, testCase.Name)
		require.NotNil(suite.T(), response, testCase.Name)
		require.NotNil(suite.T(), response.Batch, testCase.Name)
		require.EqualValues(suite.T(), testCase.ExpectStartHeight, response.Batch.Block, testCase.Name)
	}
}

func (suite *KeeperTestSuite) TestKeeper_DeleteBatchConfig() {
	tokenContract := helpers.GenerateAddress().Hex()
	batch := &types.OutgoingTxBatch{
		BatchNonce:   1,
		BatchTimeout: 0,
		Transactions: []*types.OutgoingTransferTx{
			{
				Id:          1,
				Sender:      sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
				DestAddress: helpers.GenerateAddress().Hex(),
				Token: types.ERC20Token{
					Contract: tokenContract,
					Amount:   sdkmath.NewInt(1),
				},
				Fee: types.ERC20Token{
					Contract: tokenContract,
					Amount:   sdkmath.NewInt(1),
				},
			},
		},
		TokenContract: tokenContract,
		Block:         100,
		FeeReceive:    helpers.GenerateAddress().Hex(),
	}
	suite.NoError(suite.Keeper().StoreBatch(suite.ctx, batch))

	suite.Equal(uint64(0), suite.Keeper().GetLastSlashedBatchBlock(suite.ctx))
	batches := suite.Keeper().GetUnSlashedBatches(suite.ctx, batch.Block+1)
	suite.Equal(1, len(batches))

	msgConfirmBatch := &types.MsgConfirmBatch{
		Nonce:         batch.BatchNonce,
		TokenContract: tokenContract,
		ChainName:     suite.chainName,
	}
	for i, oracle := range suite.oracleAddrs {
		msgConfirmBatch.BridgerAddress = suite.bridgerAddrs[i].String()
		msgConfirmBatch.ExternalAddress = crypto.PubkeyToAddress(suite.externalPris[i].PublicKey).String()
		suite.Keeper().SetBatchConfirm(suite.ctx, oracle, msgConfirmBatch)
	}
	suite.Keeper().OutgoingTxBatchExecuted(suite.ctx, batch.TokenContract, batch.BatchNonce)

	for _, oracle := range suite.oracleAddrs {
		suite.Nil(suite.Keeper().GetBatchConfirm(suite.ctx, batch.TokenContract, batch.BatchNonce, oracle))
	}
	suite.Nil(suite.Keeper().GetOutgoingTxBatch(suite.ctx, batch.TokenContract, batch.BatchNonce))
}

func (suite *KeeperTestSuite) TestKeeper_IterateBatchBySlashedBatchBlock() {
	index := tmrand.Intn(100)
	for i := 1; i <= index; i++ {
		tokenContract := helpers.GenerateAddress().Hex()
		batch := &types.OutgoingTxBatch{
			BatchNonce:   1,
			BatchTimeout: 0,
			Transactions: []*types.OutgoingTransferTx{
				{
					Id:          1,
					Sender:      sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
					DestAddress: helpers.GenerateAddress().Hex(),
					Token: types.ERC20Token{
						Contract: tokenContract,
						Amount:   sdkmath.NewInt(1),
					},
					Fee: types.ERC20Token{
						Contract: tokenContract,
						Amount:   sdkmath.NewInt(1),
					},
				},
			},
			TokenContract: tokenContract,
			Block:         uint64(100 + i),
			FeeReceive:    helpers.GenerateAddress().Hex(),
		}
		suite.NoError(suite.Keeper().StoreBatch(suite.ctx, batch))
	}
	var batchs []*types.OutgoingTxBatch
	suite.Keeper().IterateBatchByBlockHeight(suite.ctx, 100+1, uint64(100+index+1),
		func(batch *types.OutgoingTxBatch) bool {
			batchs = append(batchs, batch)
			return false
		},
	)
	suite.Equal(len(batchs), index)
}
