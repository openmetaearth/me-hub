package keeper_test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/openmetaearth/me-hub/testutil/helpers"
	"github.com/openmetaearth/me-hub/x/gravity/keeper"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

func (s *KeeperTestSuite) TestQueryUnbatchedTxs() {
	ctx := s.Ctx

	initSender := helpers.GenerateAddress().Bytes()
	initToken := sdk.NewCoin("usdt", sdk.NewInt(1000000))
	bridgeToken := s.NewBridgeToken(initSender, initToken)
	queryServer := keeper.NewQueryServerImpl(s.App.BscKeeper)

	// Create some unbatched transactions
	numTxs := 10
	for i := 0; i < numTxs; i++ {
		sender := helpers.GenerateAddress().Bytes()
		dest := fmt.Sprintf("dest%d", i)
		amount := sdk.NewInt(100 + int64(i))
		feeAmount := sdk.NewInt(1 + int64(i))

		err := s.App.BankKeeper.MintCoins(s.Ctx, s.chainName, sdk.NewCoins(initToken))
		s.NoError(err)
		err = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, s.chainName, sender, sdk.NewCoins(initToken))
		s.NoError(err)

		_, err = s.App.BscKeeper.AddToOutgoingPool(
			ctx,
			sender,
			dest,
			sdk.NewCoin(bridgeToken.Denom, amount),
			sdk.NewCoin(bridgeToken.Denom, feeAmount),
		)
		s.Require().NoError(err)
	}

	// Test without pagination
	res, err := queryServer.UnbatchedTxs(ctx, &types.QueryUnbatchedTxsRequest{
		ChainName: s.chainName,
	})
	s.Require().NoError(err)
	s.Require().Len(res.Txs, numTxs)
	s.Require().NotNil(res.Pagination)

	// Test with pagination
	pageLimit := 5
	res, err = queryServer.UnbatchedTxs(ctx, &types.QueryUnbatchedTxsRequest{
		ChainName: s.chainName,
		Pagination: &query.PageRequest{
			Limit: uint64(pageLimit),
		},
	})
	s.Require().NoError(err)
	s.Require().Len(res.Txs, pageLimit)
	s.Require().NotNil(res.Pagination)
	s.Require().NotNil(res.Pagination.NextKey)

	// Test next page
	res, err = queryServer.UnbatchedTxs(ctx, &types.QueryUnbatchedTxsRequest{
		ChainName: s.chainName,
		Pagination: &query.PageRequest{
			Key:   res.Pagination.NextKey,
			Limit: uint64(pageLimit),
		},
	})
	s.Require().NoError(err)
	s.Require().Len(res.Txs, numTxs-pageLimit)
	s.Require().NotNil(res.Pagination)
	s.Require().Nil(res.Pagination.NextKey)
}

func (s *KeeperTestSuite) TestQueryPendingOutgoingTxByAddrPagination() {
	ctx := s.Ctx
	queryServer := keeper.NewQueryServerImpl(s.App.BscKeeper)
	sender := sdk.AccAddress(helpers.GenerateAddress().Bytes())
	otherSender := sdk.AccAddress(helpers.GenerateAddress().Bytes())
	tokenContract := helpers.GenerateAddress().Hex()

	for i := uint64(1); i <= 3; i++ {
		s.Require().NoError(s.Keeper().StoreBatch(ctx, &types.OutgoingTxBatch{
			BatchNonce:    i,
			TokenContract: tokenContract,
			Block:         i,
			Transactions: []*types.OutgoingTransferTx{
				{Id: i, Sender: sender.String(), DestAddress: fmt.Sprintf("dest-batch-%d", i)},
				{Id: i + 100, Sender: otherSender.String(), DestAddress: fmt.Sprintf("dest-other-%d", i)},
			},
		}))
	}

	initToken := sdk.NewCoin("pendingpagination", sdk.NewInt(1000000))
	bridgeToken := s.NewBridgeToken(sender, initToken)
	for i := 0; i < 3; i++ {
		err := s.App.BankKeeper.MintCoins(ctx, s.chainName, sdk.NewCoins(initToken))
		s.Require().NoError(err)
		err = s.App.BankKeeper.SendCoinsFromModuleToAccount(ctx, s.chainName, sender, sdk.NewCoins(initToken))
		s.Require().NoError(err)
		_, err = s.App.BscKeeper.AddToOutgoingPool(
			ctx,
			sender,
			fmt.Sprintf("dest-unbatched-%d", i),
			sdk.NewCoin(bridgeToken.Denom, sdk.NewInt(100+int64(i))),
			sdk.NewCoin(bridgeToken.Denom, sdk.NewInt(1)),
		)
		s.Require().NoError(err)
	}

	res, err := queryServer.PendingOutgoingTxByAddr(ctx, &types.QueryPendingOutgoingTxByAddrRequest{
		ChainName:     s.chainName,
		SenderAddress: sender.String(),
		Pagination:    &query.PageRequest{Limit: 4},
	})
	s.Require().NoError(err)
	s.Require().Len(res.TransfersInBatches, 3)
	s.Require().Len(res.UnbatchedTransfers, 1)
	s.Require().NotNil(res.Pagination)
	s.Require().NotNil(res.Pagination.NextKey)

	res, err = queryServer.PendingOutgoingTxByAddr(ctx, &types.QueryPendingOutgoingTxByAddrRequest{
		ChainName:     s.chainName,
		SenderAddress: sender.String(),
		Pagination: &query.PageRequest{
			Key:   res.Pagination.NextKey,
			Limit: 4,
		},
	})
	s.Require().NoError(err)
	s.Require().Empty(res.TransfersInBatches)
	s.Require().Len(res.UnbatchedTransfers, 2)
	s.Require().Nil(res.Pagination.NextKey)

	res, err = queryServer.PendingOutgoingTxByAddr(ctx, &types.QueryPendingOutgoingTxByAddrRequest{
		ChainName:     s.chainName,
		SenderAddress: sender.String(),
		Pagination: &query.PageRequest{
			Limit:      4,
			CountTotal: true,
		},
	})
	s.Require().NoError(err)
	s.Require().Equal(uint64(6), res.Pagination.Total)
}
