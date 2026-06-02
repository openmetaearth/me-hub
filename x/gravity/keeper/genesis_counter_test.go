package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/openmetaearth/me-hub/x/gravity/keeper"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

func (s *KeeperTestSuite) TestGenesisImportRebuildsOutgoingAutoIncrementCounters() {
	k := s.Keeper()
	contract := "0x1111111111111111111111111111111111111111"
	outgoingTx := func(id uint64) *types.OutgoingTransferTx {
		return &types.OutgoingTransferTx{
			Id:          id,
			Sender:      s.relayerAddrs[0].String(),
			DestAddress: "0x2222222222222222222222222222222222222222",
			Token:       types.NewERC20Token(sdkmath.NewInt(100), contract),
			Fee:         types.NewERC20Token(sdkmath.NewInt(1), contract),
		}
	}

	state := &types.GenesisState{
		Params: types.DefaultParams(),
		UnbatchedTransfers: []types.OutgoingTransferTx{
			*outgoingTx(7),
		},
		Batches: []types.OutgoingTxBatch{
			{
				BatchNonce:    11,
				BatchTimeout:  100,
				Transactions:  []*types.OutgoingTransferTx{outgoingTx(13)},
				TokenContract: contract,
				Block:         1,
			},
		},
	}

	keeper.InitGenesis(s.Ctx, k, state)

	s.Require().EqualValues(14, k.AutoIncrementID(s.Ctx, types.KeyLastTxPoolID))
	s.Require().EqualValues(12, k.AutoIncrementID(s.Ctx, types.KeyLastOutgoingBatchID))
}
