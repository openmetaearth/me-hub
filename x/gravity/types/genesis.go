package types

import "fmt"

// ValidateBasic validates genesis state by looping through the params and
// calling their validation functions
func (m *GenesisState) ValidateBasic() error {
	if err := m.Params.ValidateBasic(); err != nil {
		return err
	}
	if err := m.validateOutgoingTransferLifecycle(); err != nil {
		return err
	}
	return nil
}

func (m *GenesisState) validateOutgoingTransferLifecycle() error {
	seenTxIDs := make(map[uint64]string)

	for _, tx := range m.UnbatchedTransfers {
		if previous, found := seenTxIDs[tx.Id]; found {
			return fmt.Errorf("duplicate outgoing transfer tx id %d in unbatched transfers: previously seen in %s", tx.Id, previous)
		}
		seenTxIDs[tx.Id] = "unbatched transfers"
	}

	for _, batch := range m.Batches {
		location := fmt.Sprintf("batch %d", batch.BatchNonce)
		for _, tx := range batch.Transactions {
			if tx == nil {
				return fmt.Errorf("nil outgoing transfer tx in %s", location)
			}
			if previous, found := seenTxIDs[tx.Id]; found {
				return fmt.Errorf("duplicate outgoing transfer tx id %d in %s: previously seen in %s", tx.Id, location, previous)
			}
			seenTxIDs[tx.Id] = location
		}
	}

	return nil
}
