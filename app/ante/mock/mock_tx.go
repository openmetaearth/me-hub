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

// MockFeeTx is a mock implementation of sdk.FeeTx for testing purposes.
type MockFeeTx struct {
	Msgs      []sdk.Msg
	FeeAmount sdk.Coins
	GasLimit  uint64
	Payer     sdk.AccAddress
	Granter   sdk.AccAddress
}

func (m *MockFeeTx) GetMsgs() []sdk.Msg         { return m.Msgs }
func (m *MockFeeTx) ValidateBasic() error       { return nil }
func (m *MockFeeTx) GetGas() uint64             { return m.GasLimit }
func (m *MockFeeTx) GetFee() sdk.Coins          { return m.FeeAmount }
func (m *MockFeeTx) FeePayer() sdk.AccAddress   { return m.Payer }
func (m *MockFeeTx) FeeGranter() sdk.AccAddress { return m.Granter }
