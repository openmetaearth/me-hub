package types

import (
	"fmt"
	"github.com/cometbft/cometbft/crypto/tmhash"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gogo/protobuf/proto"
)

// ExternalClaim represents a claim on ethereum state
type ExternalClaim interface {
	proto.Message
	// GetEventNonce All Ethereum claims that we relay from the Gravity contract and into the module
	// have a nonce that is monotonically increasing and unique, since this nonce is
	// issued by the Ethereum contract it is immutable and must be agreed on by all validators
	// any disagreement on what claim goes to what nonce means someone is lying.
	GetEventNonce() uint64
	// GetBlockHeight The block height that the claimed event occurred on. This EventNonce provides sufficient
	// ordering for the execution of all claims. The block height is used only for batchTimeouts + logicTimeouts
	// when we go to create a new batch we set the timeout some number of batches out from the last
	// known height plus projected block progress since then.
	GetBlockHeight() uint64
	// GetClaimer the delegate address of the claimer, for MsgSendToExternalClaim and MsgSendToMeClaim
	// this is sent in as the sdk.AccAddress of the delegated key. it is up to the user
	// to disambiguate this into a sdk.ValAddress
	GetClaimer() sdk.AccAddress
	// GetType Which type of claim this is
	GetType() ClaimType
	ValidateBasic() error
	ClaimHash() []byte
}

var (
	_ ExternalClaim = &MsgSendToMeClaim{}
	_ ExternalClaim = &MsgBridgeTokenClaim{}
	_ ExternalClaim = &MsgSendToExternalClaim{}
	_ ExternalClaim = &MsgRelayerSetUpdateClaim{}
)

func UnpackAttestationClaim(cdc codectypes.AnyUnpacker, att *Attestation) (ExternalClaim, error) {
	var msg ExternalClaim
	err := cdc.UnpackAny(att.Claim, &msg)
	return msg, err
}

// MsgSendToMeClaim

// GetType returns the type of the claim
func (m *MsgSendToMeClaim) GetType() ClaimType {
	return CLAIM_TYPE_SEND_TO_ME
}

func (m *MsgSendToMeClaim) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrapf("unrecognized cross chain name: %s", m.ChainName)
	}
	if _, err = sdk.AccAddressFromBech32(m.RelayerAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid relayer address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.Sender); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.TokenContract); err != nil {
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


func (m *MsgSendToMeClaim) GetClaimer() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.RelayerAddress)
}

// ClaimHash Hash implements BridgeSendToExternal.Hash
func (m *MsgSendToMeClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%d/%s/%s/%s/%s", m.BlockHeight, m.EventNonce, m.TokenContract, m.Sender, m.Amount.String(), m.Receiver)
	return tmhash.Sum([]byte(path))
}

// MsgSendToExternalClaim //

// GetType returns the claim type
func (m *MsgSendToExternalClaim) GetType() ClaimType {
	return CLAIM_TYPE_SEND_TO_EXTERNAL
}

// ValidateBasic performs stateless checks
func (m *MsgSendToExternalClaim) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrapf("unrecognized cross chain name: %s", m.ChainName)
	}
	if _, err = sdk.AccAddressFromBech32(m.RelayerAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid relayer address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.TokenContract); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid token contract: %s", err)
	}
	if m.EventNonce == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero event nonce")
	}
	if m.BlockHeight == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero block height")
	}
	if m.BatchNonce == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero batch nonce")
	}
	return nil
}

// ClaimHash Hash implements SendToFxBatch.Hash
func (m *MsgSendToExternalClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%d/%s/%d", m.BlockHeight, m.EventNonce, m.TokenContract, m.BatchNonce)
	return tmhash.Sum([]byte(path))
}

func (m *MsgSendToExternalClaim) GetClaimer() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.RelayerAddress)
}

// MsgBridgeTokenClaim //

func (m *MsgBridgeTokenClaim) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrapf("unrecognized cross chain name: %s", m.ChainName)
	}
	if _, err = sdk.AccAddressFromBech32(m.RelayerAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.TokenContract); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid token contract: %s", err)
	}
	if len(m.Name) == 0 {
		return errortypes.ErrInvalidRequest.Wrap("empty token name")
	}
	if len(m.Symbol) == 0 {
		return errortypes.ErrInvalidRequest.Wrap("empty token symbol")
	}
	if m.EventNonce == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero event nonce")
	}
	if m.BlockHeight == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero block height")
	}
	return nil
}

func (m *MsgBridgeTokenClaim) GetClaimer() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.RelayerAddress)
}

func (m *MsgBridgeTokenClaim) GetType() ClaimType {
	return CLAIM_TYPE_BRIDGE_TOKEN
}

func (m *MsgBridgeTokenClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%d/%s/%s/%s/%d", m.BlockHeight, m.EventNonce, m.TokenContract, m.Name, m.Symbol, m.Decimals)
	return tmhash.Sum([]byte(path))
}

// MsgRelayerSetUpdateClaim //

// GetType returns the type of the claim
func (m *MsgRelayerSetUpdateClaim) GetType() ClaimType {
	return CLAIM_TYPE_RELAYER_SET_UPDATED
}

// ValidateBasic performs stateless checks
func (m *MsgRelayerSetUpdateClaim) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrapf("unrecognized cross chain name: %s", m.ChainName)
	}
	if _, err = sdk.AccAddressFromBech32(m.RelayerAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if len(m.Members) == 0 {
		return errortypes.ErrInvalidRequest.Wrap("empty members")
	}
	for _, member := range m.Members {
		if err = ValidateExternalAddr(m.ChainName, member.ExternalAddress); err != nil {
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

func (m *MsgRelayerSetUpdateClaim) GetClaimer() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.RelayerAddress)
}

// ClaimHash Hash implements BridgeSendToExternal.Hash
func (m *MsgRelayerSetUpdateClaim) ClaimHash() []byte {
	var membersStr string
	for _, member := range m.Members {
		membersStr += member.String()
	}
	path := fmt.Sprintf("%d/%d/%d/%s", m.BlockHeight, m.RelayerSetNonce, m.EventNonce, membersStr)
	return tmhash.Sum([]byte(path))
}
