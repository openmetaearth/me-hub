package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/openmetaearth/me-hub/testutil/helpers"
	"github.com/openmetaearth/me-hub/x/gravity/keeper"
	"github.com/openmetaearth/me-hub/x/gravity/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *KeeperTestSuite) TestQueriesRejectInvalidRelayerAddressWithoutPanic() {
	queryServer := keeper.NewQueryServerImpl(s.Keeper())
	invalidRelayer := "not-a-bech32-address"
	tokenContract := helpers.GenHexAddress().String()

	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "RelayerSetConfirm",
			call: func() error {
				_, err := queryServer.RelayerSetConfirm(s.Ctx, &types.QueryRelayerSetConfirmRequest{
					Nonce:          1,
					RelayerAddress: invalidRelayer,
				})
				return err
			},
		},
		{
			name: "LastPendingRelayerSetRequestByAddr",
			call: func() error {
				_, err := queryServer.LastPendingRelayerSetRequestByAddr(s.Ctx, &types.QueryLastPendingRelayerSetRequestByAddrRequest{
					RelayerAddress: invalidRelayer,
				})
				return err
			},
		},
		{
			name: "LastPendingBatchRequestByAddr",
			call: func() error {
				_, err := queryServer.LastPendingBatchRequestByAddr(s.Ctx, &types.QueryLastPendingBatchRequestByAddrRequest{
					RelayerAddress: invalidRelayer,
				})
				return err
			},
		},
		{
			name: "LastEventNonceByAddr",
			call: func() error {
				_, err := queryServer.LastEventNonceByAddr(s.Ctx, &types.QueryLastEventNonceByAddrRequest{
					RelayerAddress: invalidRelayer,
				})
				return err
			},
		},
		{
			name: "LastEventBlockHeightByAddr",
			call: func() error {
				_, err := queryServer.LastEventBlockHeightByAddr(s.Ctx, &types.QueryLastEventBlockHeightByAddrRequest{
					RelayerAddress: invalidRelayer,
				})
				return err
			},
		},
		{
			name: "BatchConfirm",
			call: func() error {
				_, err := queryServer.BatchConfirm(s.Ctx, &types.QueryBatchConfirmRequest{
					ChainName:      s.chainName,
					TokenContract:  tokenContract,
					Nonce:          1,
					RelayerAddress: invalidRelayer,
				})
				return err
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			var err error
			s.Require().NotPanics(func() {
				err = tc.call()
			})
			s.Require().Error(err)
			s.Require().Equal(codes.InvalidArgument, status.Code(err))
		})
	}
}
