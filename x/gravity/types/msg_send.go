package types

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgSendToExternal       = "send_to_external"
	TypeMsgRequestBatch         = "request_batch"
	TypeMsgConfirmBatch         = "confirm_batch"
	TypeMsgCancelSendToExternal = "cancel_send_to_external"
	TypeMsgIncreaseBridgeFee    = "increase_bridge_fee"
)

// MsgSendToExternal //

// Route should return the name of the module
func (m *MsgSendToExternal) Route() string { return RouterKey }

// Type should return the action
func (m *MsgSendToExternal) Type() string { return TypeMsgSendToExternal }

// ValidateBasic runs stateless checks on the message
// Checks if the Eth address is valid
func (m *MsgSendToExternal) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrapf("unrecognized cross chain name: %s", m.ChainName)
	}
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.Dest); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid dest address: %s", err)
	}
	if !m.Amount.IsValid() || !m.Amount.IsPositive() {
		return errortypes.ErrInvalidRequest.Wrap("invalid amount")
	}
	if m.Amount.Denom != m.BridgeFee.Denom {
		return errortypes.ErrInvalidRequest.Wrap("bridge fee denom not equal amount denom")
	}
	if !m.BridgeFee.IsValid() || !m.BridgeFee.IsPositive() {
		return errortypes.ErrInvalidRequest.Wrap("invalid bridge fee")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgSendToExternal) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgSendToExternal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Sender)}
}

// MsgRequestBatch //

// Route should return the name of the module
func (m *MsgRequestBatch) Route() string { return RouterKey }

// Type should return the action
func (m *MsgRequestBatch) Type() string { return TypeMsgRequestBatch }

// ValidateBasic performs stateless checks
func (m *MsgRequestBatch) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrapf("unrecognized cross chain name: %s", m.ChainName)
	}
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if len(m.Denom) == 0 {
		return errortypes.ErrInvalidRequest.Wrap("empty denom")
	}
	if m.MinimumFee.IsNil() || !m.MinimumFee.IsPositive() {
		return errortypes.ErrInvalidRequest.Wrap("invalid minimum fee")
	}
	if err = ValidateExternalAddr(m.ChainName, m.FeeReceive); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid fee receive address: %s", err)
	}
	if m.BaseFee.IsNil() || m.BaseFee.IsNegative() {
		return errortypes.ErrInvalidRequest.Wrap("invalid base fee")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgRequestBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgRequestBatch) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Sender)}
}

// MsgConfirmBatch //

// Route should return the name of the module
func (m *MsgConfirmBatch) Route() string { return RouterKey }

// Type should return the action
func (m *MsgConfirmBatch) Type() string { return TypeMsgConfirmBatch }

// ValidateBasic performs stateless checks
func (m *MsgConfirmBatch) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrapf("unrecognized cross chain name: %s", m.ChainName)
	}
	if _, err = sdk.AccAddressFromBech32(m.RelayerAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.ExternalAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid external address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.TokenContract); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid token contract: %s", err)
	}
	if len(m.Signature) == 0 {
		return errortypes.ErrInvalidRequest.Wrap("empty signature")
	}
	if _, err = hex.DecodeString(m.Signature); err != nil {
		return errortypes.ErrInvalidRequest.Wrap("could not hex decode signature")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgConfirmBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgConfirmBatch) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.RelayerAddress)}
}

// MsgCancelSendToExternal //

// Route should return the name of the module
func (m *MsgCancelSendToExternal) Route() string { return RouterKey }

// Type should return the action
func (m *MsgCancelSendToExternal) Type() string { return TypeMsgCancelSendToExternal }

// ValidateBasic performs stateless checks
func (m *MsgCancelSendToExternal) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrapf("unrecognized cross chain name: %s", m.ChainName)
	}
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if m.TransactionId == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero transaction id")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgCancelSendToExternal) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgCancelSendToExternal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Sender)}
}

// MsgIncreaseBridgeFee

// Route should return the name of the module
func (m *MsgIncreaseBridgeFee) Route() string { return RouterKey }

// Type should return the action
func (m *MsgIncreaseBridgeFee) Type() string { return TypeMsgIncreaseBridgeFee }

// ValidateBasic performs stateless checks
func (m *MsgIncreaseBridgeFee) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrapf("unrecognized cross chain name: %s", m.ChainName)
	}
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if m.TransactionId == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero transaction id")
	}
	if !m.AddBridgeFee.IsValid() || !m.AddBridgeFee.IsPositive() {
		return errortypes.ErrInvalidRequest.Wrap("invalid bridge fee")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgIncreaseBridgeFee) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgIncreaseBridgeFee) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Sender)}
}
