package types

import "fmt"

// ValidateBasic validates genesis state by checking params and
// cross-field accounting invariants between UnbatchedTransfers and Batches.
// It ensures no outgoing transfer id is duplicated within or across the two
// mutually exclusive lifecycle collections, preventing the double-spend vector
// where a genesis import could place the same tx in both UnbatchedTransfers
// (allowing cancel/refund) and a Batch (allowing external execution).
func (m *GenesisState) ValidateBasic() error {
	if err := m.Params.ValidateBasic(); err != nil {
		return err
	}

	unbatchedIds := make(map[uint64]bool)

	// Check for duplicates within UnbatchedTransfers
	for _, tx := range m.UnbatchedTransfers {
		if unbatchedIds[tx.Id] {
			return fmt.Errorf("duplicate outgoing tx id %d in UnbatchedTransfers", tx.Id)
		}
		unbatchedIds[tx.Id] = true
	}

	// Check for duplicates within Batches and overlap with UnbatchedTransfers.
	// Also validate that each batch transaction's token contract matches its
	// parent batch token contract.
	batchedIds := make(map[uint64]uint64) // tx.Id -> batch nonce
	for _, batch := range m.Batches {
		for _, tx := range batch.Transactions {
			if tx == nil {
				continue
			}

			// Duplicate within Batches (same tx in multiple batches)
			if prevNonce, ok := batchedIds[tx.Id]; ok {
				return fmt.Errorf(
					"duplicate outgoing tx id %d in Batches (batch nonce %d and %d)",
					tx.Id, prevNonce, batch.BatchNonce,
				)
			}
			batchedIds[tx.Id] = batch.BatchNonce

			// Overlap between UnbatchedTransfers and Batches —
			// a tx must not be both cancelable (unbatched) and executable (batched)
			if unbatchedIds[tx.Id] {
				return fmt.Errorf(
					"outgoing tx id %d exists in both UnbatchedTransfers and Batches (batch nonce %d)",
					tx.Id, batch.BatchNonce,
				)
			}

			// Token contract consistency within a batch
			if batch.TokenContract != "" && tx.Token.Contract != batch.TokenContract {
				return fmt.Errorf(
					"outgoing tx id %d token contract %s does not match batch nonce %d token contract %s",
					tx.Id, tx.Token.Contract, batch.BatchNonce, batch.TokenContract,
				)
			}
		}
	}

	return nil
}
