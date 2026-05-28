package types

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgSendToExternal //

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

// MsgRequestBatch //

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

// MsgConfirmBatch //

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

// MsgCancelSendToExternal //

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

// MsgIncreaseBridgeFee

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

