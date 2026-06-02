package keeper_test

import (
	sdkmath "cosmossdk.io/math"

	"github.com/openmetaearth/me-hub/x/gravity/keeper"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

const genesisCounterTestContract = "0x1111111111111111111111111111111111111111"

func (s *KeeperTestSuite) TestGenesisImportRebuildsOutgoingAutoIncrementCounters() {
	k := s.Keeper()
	state := &types.GenesisState{
		Params: types.DefaultParams(),
		UnbatchedTransfers: []types.OutgoingTransferTx{
			*s.genesisCounterOutgoingTx(7),
		},
		Batches: []types.OutgoingTxBatch{
			{
				BatchNonce:    11,
				BatchTimeout:  100,
				Transactions:  []*types.OutgoingTransferTx{s.genesisCounterOutgoingTx(13)},
				TokenContract: genesisCounterTestContract,
				Block:         1,
			},
		},
	}

	keeper.InitGenesis(s.Ctx, k, state)

	batch := k.GetOutgoingTxBatch(s.Ctx, genesisCounterTestContract, 11)
	s.Require().NotNil(batch)
	s.Require().Len(batch.Transactions, 1)
	s.Require().EqualValues(13, batch.Transactions[0].Id)
	s.Require().EqualValues(14, k.AutoIncrementID(s.Ctx, types.KeyLastTxPoolID))
	s.Require().EqualValues(12, k.AutoIncrementID(s.Ctx, types.KeyLastOutgoingBatchID))
}

func (s *KeeperTestSuite) TestGenesisImportKeepsEmptyOutgoingCountersAtDefault() {
	k := s.Keeper()

	keeper.InitGenesis(s.Ctx, k, &types.GenesisState{Params: types.DefaultParams()})

	s.Require().EqualValues(1, k.AutoIncrementID(s.Ctx, types.KeyLastTxPoolID))
	s.Require().EqualValues(1, k.AutoIncrementID(s.Ctx, types.KeyLastOutgoingBatchID))
}

func (s *KeeperTestSuite) TestGenesisImportSkipsNilBatchedTransfersWhenRebuildingCounters() {
	k := s.Keeper()
	state := &types.GenesisState{
		Params: types.DefaultParams(),
		Batches: []types.OutgoingTxBatch{
			{
				BatchNonce:    11,
				BatchTimeout:  100,
				Transactions:  []*types.OutgoingTransferTx{nil, s.genesisCounterOutgoingTx(13)},
				TokenContract: genesisCounterTestContract,
				Block:         1,
			},
		},
	}

	keeper.InitGenesis(s.Ctx, k, state)

	s.Require().EqualValues(14, k.AutoIncrementID(s.Ctx, types.KeyLastTxPoolID))
	s.Require().EqualValues(12, k.AutoIncrementID(s.Ctx, types.KeyLastOutgoingBatchID))
}

func (s *KeeperTestSuite) genesisCounterOutgoingTx(id uint64) *types.OutgoingTransferTx {
	return &types.OutgoingTransferTx{
		Id:          id,
		Sender:      s.relayerAddrs[0].String(),
		DestAddress: "0x2222222222222222222222222222222222222222",
		Token:       types.NewERC20Token(sdkmath.NewInt(100), genesisCounterTestContract),
		Fee:         types.NewERC20Token(sdkmath.NewInt(1), genesisCounterTestContract),
	}
}
