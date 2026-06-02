package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisStateValidateBasicRejectsDuplicateOutgoingTransferLifecycle(t *testing.T) {
	for _, tc := range []struct {
		name        string
		genesis     GenesisState
		errContains string
	}{
		{
			name: "valid distinct outgoing transfers",
			genesis: GenesisState{
				Params:             DefaultParams(),
				UnbatchedTransfers: []OutgoingTransferTx{outgoingTransferTx(1)},
				Batches:            []OutgoingTxBatch{outgoingTxBatch(2, 3)},
			},
		},
		{
			name: "duplicate unbatched outgoing transfer",
			genesis: GenesisState{
				Params:             DefaultParams(),
				UnbatchedTransfers: []OutgoingTransferTx{outgoingTransferTx(7), outgoingTransferTx(7)},
			},
			errContains: "duplicate outgoing transfer tx id 7",
		},
		{
			name: "duplicate batched outgoing transfer",
			genesis: GenesisState{
				Params:  DefaultParams(),
				Batches: []OutgoingTxBatch{outgoingTxBatch(3, 7, 7)},
			},
			errContains: "duplicate outgoing transfer tx id 7",
		},
		{
			name: "outgoing transfer in unbatched pool and batch",
			genesis: GenesisState{
				Params:             DefaultParams(),
				UnbatchedTransfers: []OutgoingTransferTx{outgoingTransferTx(7)},
				Batches:            []OutgoingTxBatch{outgoingTxBatch(3, 7)},
			},
			errContains: "duplicate outgoing transfer tx id 7",
		},
		{
			name: "nil outgoing transfer in batch",
			genesis: GenesisState{
				Params: DefaultParams(),
				Batches: []OutgoingTxBatch{
					{
						BatchNonce:   3,
						Transactions: []*OutgoingTransferTx{nil},
					},
				},
			},
			errContains: "nil outgoing transfer tx in batch 3",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.genesis.ValidateBasic()
			if tc.errContains == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.errContains)
			}
		})
	}
}

func outgoingTxBatch(nonce uint64, txIDs ...uint64) OutgoingTxBatch {
	txs := make([]*OutgoingTransferTx, 0, len(txIDs))
	for _, txID := range txIDs {
		tx := outgoingTransferTx(txID)
		txs = append(txs, &tx)
	}

	return OutgoingTxBatch{
		BatchNonce:    nonce,
		TokenContract: "0x1111111111111111111111111111111111111111",
		Transactions:  txs,
	}
}

func outgoingTransferTx(id uint64) OutgoingTransferTx {
	return OutgoingTransferTx{
		Id:          id,
		Sender:      "me1p8u5377smm8zkfq9dmcg6prwflap0ndj4ht34z",
		DestAddress: "0x2222222222222222222222222222222222222222",
		Token: ERC20Token{
			Contract: "0x1111111111111111111111111111111111111111",
			Amount:   sdk.NewInt(100),
		},
		Fee: ERC20Token{
			Contract: "0x1111111111111111111111111111111111111111",
			Amount:   sdk.NewInt(1),
		},
	}
}
