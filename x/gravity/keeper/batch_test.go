package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/openmetaearth/me-hub/testutil/helpers"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

func (suite *KeeperTestSuite) TestLastPendingBatchRequestByAddr() {
	testCases := []struct {
		Name              string
		RelayerAddress    sdk.AccAddress
		StartHeight       int64
		ExpectStartHeight uint64
	}{
		{
			Name:              "relayer start height with 1, expect relayer set block 3",
			RelayerAddress:    suite.relayerAddrs[0],
			StartHeight:       1,
			ExpectStartHeight: 3,
		},
		{
			Name:              "relayer start height with 2, expect relayer set block 2",
			RelayerAddress:    suite.relayerAddrs[1],
			StartHeight:       2,
			ExpectStartHeight: 3,
		},
		{
			Name:              "relayer start height with 3, expect relayer set block 1",
			RelayerAddress:    suite.relayerAddrs[2],
			StartHeight:       3,
			ExpectStartHeight: 3,
		},
	}
	for i := uint64(1); i <= 3; i++ {
		suite.Ctx = suite.Ctx.WithBlockHeight(int64(i))
		err := suite.Keeper().StoreBatch(suite.Ctx, &types.OutgoingTxBatch{
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

	wrapSDKContext := sdk.WrapSDKContext(suite.Ctx)
	for _, testCase := range testCases {
		relayer := types.Relayer{
			RelayerAddress: testCase.RelayerAddress.String(),
			StartHeight:    testCase.StartHeight,
		}
		suite.Keeper().SetRelayer(suite.Ctx, testCase.RelayerAddress, relayer)

		response, err := suite.QueryClient().LastPendingBatchRequestByAddr(wrapSDKContext,
			&types.QueryLastPendingBatchRequestByAddrRequest{
				RelayerAddress: testCase.RelayerAddress.String(),
			})
		suite.Require().NoError(err, testCase.Name)
		suite.Require().NotNil(response, testCase.Name)
		suite.Require().NotNil(response.Batch, testCase.Name)
		suite.Require().EqualValues(testCase.ExpectStartHeight, response.Batch.Block, testCase.Name)
	}
}

func (suite *KeeperTestSuite) TestKeeper_DeleteBatchConfirm() {
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
	suite.NoError(suite.Keeper().StoreBatch(suite.Ctx, batch))
	suite.Equal(uint64(0), suite.Keeper().GetLastSlashedBatchBlock(suite.Ctx))

	batches := suite.Keeper().GetUnSlashedBatches(suite.Ctx, batch.Block+1)
	suite.Equal(1, len(batches))

	msgConfirmBatch := &types.MsgConfirmBatch{
		Nonce:         batch.BatchNonce,
		TokenContract: tokenContract,
		ChainName:     suite.chainName,
	}
	for i, relayer := range suite.relayerAddrs {
		msgConfirmBatch.RelayerAddress = suite.relayerAddrs[i].String()
		msgConfirmBatch.ExternalAddress = crypto.PubkeyToAddress(suite.externalPris[i].PublicKey).String()
		suite.Keeper().SetBatchConfirm(suite.Ctx, relayer, msgConfirmBatch)
	}
	suite.Keeper().OutgoingTxBatchExecuted(suite.Ctx, batch.TokenContract, batch.BatchNonce)

	for _, relayer := range suite.relayerAddrs {
		suite.Nil(suite.Keeper().GetBatchConfirm(suite.Ctx, batch.TokenContract, batch.BatchNonce, relayer))
	}
	suite.Nil(suite.Keeper().GetOutgoingTxBatch(suite.Ctx, batch.TokenContract, batch.BatchNonce))
}

func (suite *KeeperTestSuite) TestOutgoingTxBatchExecutedCancelsEarlierBatchesAfterIteration() {
	tokenContract := helpers.GenerateAddress().Hex()
	otherTokenContract := helpers.GenerateAddress().Hex()

	olderBatch := newTestOutgoingBatch(tokenContract, 1, 100, 101, 1001)
	middleBatch := newTestOutgoingBatch(tokenContract, 2, 100, 102, 1002)
	executedBatch := newTestOutgoingBatch(tokenContract, 3, 100, 103, 1003)
	otherTokenBatch := newTestOutgoingBatch(otherTokenContract, 1, 100, 104, 2001)

	suite.Require().NoError(suite.Keeper().StoreBatch(suite.Ctx, olderBatch))
	suite.Require().NoError(suite.Keeper().StoreBatch(suite.Ctx, middleBatch))
	suite.Require().NoError(suite.Keeper().StoreBatch(suite.Ctx, executedBatch))
	suite.Require().NoError(suite.Keeper().StoreBatch(suite.Ctx, otherTokenBatch))

	suite.Keeper().OutgoingTxBatchExecuted(suite.Ctx, tokenContract, executedBatch.BatchNonce)

	suite.Nil(suite.Keeper().GetOutgoingTxBatch(suite.Ctx, tokenContract, olderBatch.BatchNonce))
	suite.Nil(suite.Keeper().GetOutgoingTxBatch(suite.Ctx, tokenContract, middleBatch.BatchNonce))
	suite.Nil(suite.Keeper().GetOutgoingTxBatch(suite.Ctx, tokenContract, executedBatch.BatchNonce))
	suite.NotNil(suite.Keeper().GetOutgoingTxBatch(suite.Ctx, otherTokenContract, otherTokenBatch.BatchNonce))

	_, err := suite.Keeper().GetUnbatchedTxById(suite.Ctx, olderBatch.Transactions[0].Id)
	suite.Require().NoError(err)
	_, err = suite.Keeper().GetUnbatchedTxById(suite.Ctx, middleBatch.Transactions[0].Id)
	suite.Require().NoError(err)
	_, err = suite.Keeper().GetUnbatchedTxById(suite.Ctx, executedBatch.Transactions[0].Id)
	suite.Require().Error(err)
}

func (suite *KeeperTestSuite) TestEndBlockerCancelsTimedOutBatchesAfterIteration() {
	tokenContract := helpers.GenerateAddress().Hex()
	expiredBatch := newTestOutgoingBatch(tokenContract, 1, 10, 201, 3001)
	alsoExpiredBatch := newTestOutgoingBatch(tokenContract, 2, 11, 202, 3002)
	activeBatch := newTestOutgoingBatch(tokenContract, 3, 99, 203, 3003)

	suite.Require().NoError(suite.Keeper().StoreBatch(suite.Ctx, expiredBatch))
	suite.Require().NoError(suite.Keeper().StoreBatch(suite.Ctx, alsoExpiredBatch))
	suite.Require().NoError(suite.Keeper().StoreBatch(suite.Ctx, activeBatch))
	suite.Keeper().SetLastObservedBlockHeight(suite.Ctx, 12, uint64(suite.Ctx.BlockHeight()))

	suite.Keeper().EndBlocker(suite.Ctx)

	suite.Nil(suite.Keeper().GetOutgoingTxBatch(suite.Ctx, tokenContract, expiredBatch.BatchNonce))
	suite.Nil(suite.Keeper().GetOutgoingTxBatch(suite.Ctx, tokenContract, alsoExpiredBatch.BatchNonce))
	suite.NotNil(suite.Keeper().GetOutgoingTxBatch(suite.Ctx, tokenContract, activeBatch.BatchNonce))

	_, err := suite.Keeper().GetUnbatchedTxById(suite.Ctx, expiredBatch.Transactions[0].Id)
	suite.Require().NoError(err)
	_, err = suite.Keeper().GetUnbatchedTxById(suite.Ctx, alsoExpiredBatch.Transactions[0].Id)
	suite.Require().NoError(err)
	_, err = suite.Keeper().GetUnbatchedTxById(suite.Ctx, activeBatch.Transactions[0].Id)
	suite.Require().Error(err)
}

func (suite *KeeperTestSuite) TestKeeper_IterateBatch() {
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
		suite.NoError(suite.Keeper().StoreBatch(suite.Ctx, batch))
	}
	var batchs []*types.OutgoingTxBatch
	suite.Keeper().IterateBatchByBlockHeight(suite.Ctx, 100+1, uint64(100+index+1),
		func(batch *types.OutgoingTxBatch) bool {
			batchs = append(batchs, batch)
			return false
		},
	)
	suite.Equal(len(batchs), index)
}

func newTestOutgoingBatch(tokenContract string, nonce, timeout, block, txID uint64) *types.OutgoingTxBatch {
	return &types.OutgoingTxBatch{
		BatchNonce:   nonce,
		BatchTimeout: timeout,
		Transactions: []*types.OutgoingTransferTx{
			{
				Id:          txID,
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
		Block:         block,
		FeeReceive:    helpers.GenerateAddress().Hex(),
	}
}
