package mock

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	protov2 "google.golang.org/protobuf/proto"
)

// MockTx is a mock implementation of the sdk.Tx interface for testing purposes.
type MockTx struct {
	Msgs []sdk.Msg
}

// GetMsgs returns the messages in the transaction.
func (m *MockTx) GetMsgs() []sdk.Msg {
	return m.Msgs
}

// GetMsgsV2 returns the messages using the protobuf v2 interface.
func (m *MockTx) GetMsgsV2() ([]protov2.Message, error) {
	return nil, nil
}

// ValidateBasic is a placeholder implementation for the sdk.Tx interface.
func (m *MockTx) ValidateBasic() error {
	return nil
}
