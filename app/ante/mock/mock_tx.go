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

// GetMsgsV2 returns the messages as protov2 messages.
func (m *MockTx) GetMsgsV2() ([]protov2.Message, error) {
	msgs := make([]protov2.Message, 0, len(m.Msgs))
	for _, msg := range m.Msgs {
		if v2, ok := msg.(protov2.Message); ok {
			msgs = append(msgs, v2)
		}
	}
	return msgs, nil
}

// ValidateBasic is a placeholder implementation for the sdk.Tx interface.
func (m *MockTx) ValidateBasic() error {
	return nil
}
