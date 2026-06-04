package types

import (
	"testing"

	errorsmod "cosmossdk.io/errors"
	"github.com/stretchr/testify/require"
)

func TestGenesisValidateBasicRejectsDuplicateOutgoingTransferTxIDs(t *testing.T) {
	testCases := []struct {
		name    string
		genesis GenesisState
	}{
		{
			name: "duplicate unbatched transfer ids",
			genesis: GenesisState{
				Params: DefaultParams(),
				UnbatchedTransfers: []OutgoingTransferTx{
					{Id: 7},
					{Id: 7},
				},
			},
		},
		{
			name: "duplicate ids across unbatched transfer and batch",
			genesis: GenesisState{
				Params: DefaultParams(),
				UnbatchedTransfers: []OutgoingTransferTx{
					{Id: 7},
				},
				Batches: []OutgoingTxBatch{
					{
						BatchNonce: 3,
						Transactions: []*OutgoingTransferTx{
							{Id: 7},
						},
					},
				},
			},
		},
		{
			name: "duplicate ids across batches",
			genesis: GenesisState{
				Params: DefaultParams(),
				Batches: []OutgoingTxBatch{
					{
						BatchNonce: 3,
						Transactions: []*OutgoingTransferTx{
							{Id: 7},
						},
					},
					{
						BatchNonce: 4,
						Transactions: []*OutgoingTransferTx{
							{Id: 7},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.genesis.ValidateBasic()

			require.Error(t, err)
			require.True(t, errorsmod.IsOf(err, ErrDuplicate), "expected ErrDuplicate, got %v", err)
		})
	}
}

func TestGenesisValidateBasicAcceptsUniqueOutgoingTransferTxIDs(t *testing.T) {
	genesis := GenesisState{
		Params: DefaultParams(),
		UnbatchedTransfers: []OutgoingTransferTx{
			{Id: 7},
		},
		Batches: []OutgoingTxBatch{
			{
				BatchNonce: 3,
				Transactions: []*OutgoingTransferTx{
					{Id: 8},
				},
			},
		},
	}

	require.NoError(t, genesis.ValidateBasic())
}

func TestGenesisValidateBasicRejectsNilBatchTransfer(t *testing.T) {
	genesis := GenesisState{
		Params: DefaultParams(),
		Batches: []OutgoingTxBatch{
			{
				BatchNonce:   3,
				Transactions: []*OutgoingTransferTx{nil},
			},
		},
	}

	err := genesis.ValidateBasic()

	require.Error(t, err)
	require.True(t, errorsmod.IsOf(err, ErrInvalid), "expected ErrInvalid, got %v", err)
}
