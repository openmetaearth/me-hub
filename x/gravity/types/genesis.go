package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
)

// ValidateBasic validates genesis state by looping through the params and
// calling their validation functions
func (m *GenesisState) ValidateBasic() error {
	if err := m.Params.ValidateBasic(); err != nil {
		return err
	}
	if err := m.validateUniqueOutgoingTransferTxIDs(); err != nil {
		return err
	}
	return nil
}

func (m *GenesisState) validateUniqueOutgoingTransferTxIDs() error {
	seen := make(map[uint64]string, len(m.UnbatchedTransfers))

	for i, tx := range m.UnbatchedTransfers {
		location := fmt.Sprintf("unbatched transfer index %d", i)
		if previous, found := seen[tx.Id]; found {
			return errorsmod.Wrapf(ErrDuplicate, "outgoing transfer tx id %d appears in %s and %s", tx.Id, previous, location)
		}
		seen[tx.Id] = location
	}

	for batchIndex, batch := range m.Batches {
		for txIndex, tx := range batch.Transactions {
			if tx == nil {
				return errorsmod.Wrapf(ErrInvalid, "nil outgoing transfer tx in batch index %d transaction index %d", batchIndex, txIndex)
			}
			location := fmt.Sprintf("batch index %d transaction index %d", batchIndex, txIndex)
			if previous, found := seen[tx.Id]; found {
				return errorsmod.Wrapf(ErrDuplicate, "outgoing transfer tx id %d appears in %s and %s", tx.Id, previous, location)
			}
			seen[tx.Id] = location
		}
	}

	return nil
}
