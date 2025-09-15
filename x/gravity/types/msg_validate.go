package types

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/st-chain/me-hub/utils"
)

var _ MsgValidateBasic = &MsgValidate{}

type MsgValidate struct{}

func (m2 MsgValidate) MsgRelayerSetConfirmValidate(m *MsgRelayerSetConfirm) (err error) {
	if err = utils.ValidateEthereumAddress(m.ExternalAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid external address: %s", err)
	}
	if len(m.Signature) == 0 {
		return errortypes.ErrInvalidRequest.Wrap("empty signature")
	}
	if _, err = hex.DecodeString(m.Signature); err != nil {
		return errortypes.ErrInvalidRequest.Wrap("could not hex decode signature")
	}
	return nil
}

func (m2 MsgValidate) MsgRelayerSetUpdatedClaimValidate(m *MsgRelayerSetUpdatedClaim) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.RelayerAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if len(m.Members) == 0 {
		return errortypes.ErrInvalidRequest.Wrap("empty members")
	}
	for _, member := range m.Members {
		if err = utils.ValidateEthereumAddress(member.ExternalAddress); err != nil {
			return errortypes.ErrInvalidAddress.Wrapf("invalid external address: %s", err)
		}
		if member.Power == 0 {
			return errortypes.ErrInvalidRequest.Wrap("zero power")
		}
	}
	if m.EventNonce == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero event nonce")
	}
	if m.BlockHeight == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero block height")
	}
	return nil
}

func (m2 MsgValidate) MsgBridgeTokenClaimValidate(m *MsgBridgeTokenClaim) (err error) {
	//TODO implement me
	panic("implement me")
}

func (m2 MsgValidate) MsgSendToExternalClaimValidate(m *MsgSendToExternalClaim) (err error) {
	//TODO implement me
	panic("implement me")
}

func (m2 MsgValidate) MsgSendToMeClaimValidate(m *MsgSendToMeClaim) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.RelayerAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if err = utils.ValidateEthereumAddress(m.Sender); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if err = utils.ValidateEthereumAddress(m.TokenContract); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid token contract: %s", err)
	}
	if _, err = sdk.AccAddressFromBech32(m.Receiver); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid receiver address: %s", err)
	}
	if m.Amount.IsNil() || m.Amount.IsNegative() {
		return errortypes.ErrInvalidRequest.Wrap("invalid amount")
	}
	if m.EventNonce == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero event nonce")
	}
	if m.BlockHeight == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero block height")
	}
	return nil
}

func (m2 MsgValidate) MsgBridgeCallClaimValidate(m *MsgBridgeCallClaim) (err error) {
	//TODO implement me
	panic("implement me")
}

func (m2 MsgValidate) MsgSendToExternalValidate(m *MsgSendToExternal) (err error) {
	//TODO implement me
	panic("implement me")
}

func (m2 MsgValidate) MsgRequestBatchValidate(m *MsgRequestBatch) (err error) {
	//TODO implement me
	panic("implement me")
}

func (m2 MsgValidate) MsgConfirmBatchValidate(m *MsgConfirmBatch) (err error) {
	//TODO implement me
	panic("implement me")
}

func (m2 MsgValidate) ValidateAddress(addr string) error {
	//TODO implement me
	panic("implement me")
}

func (m2 MsgValidate) AddressToBytes(addr string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
