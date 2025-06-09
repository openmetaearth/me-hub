package mock

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MockTx is a mock implementation of the sdk.Tx interface for testing purposes.
type MockTx struct {
	Msgs []sdk.Msg
}

// GetMsgs returns the messages in the transaction.
func (m *MockTx) GetMsgs() []sdk.Msg {
	return m.Msgs
}

// ValidateBasic is a placeholder implementation for the sdk.Tx interface.
func (m *MockTx) ValidateBasic() error {
	return nil
}
