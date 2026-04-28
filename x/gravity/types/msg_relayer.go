package types

import (
	"cosmossdk.io/errors"
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/openmetaearth/me-hub/app/params"
)

const (
	TypeMsgBondedRelayer     = "bonded_relayer"
	TypeMsgAddDelegate       = "add_delegate"
	TypeMsgUnbondedRelayer   = "unbonded_relayer"
	TypeMsgRelayerSetConfirm = "relayer_set_confirm"
	TypeMsgProposalRelayers  = "update_relayers"
	TypeMsgUpdateParams      = "update_params"
)

// MsgBondedRelayer //

// Route implements the sdk.Msg interface.
func (m *MsgBondedRelayer) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (m *MsgBondedRelayer) Type() string { return TypeMsgBondedRelayer }

func (m *MsgBondedRelayer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.RelayerAddress)}
}

// GetSignBytes returns the message bytes to sign over.
func (m *MsgBondedRelayer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgBondedRelayer) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.RelayerAddress); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "relayer address is not a valid bech32 address")
	}
	if err := ValidateExternalAddr(m.ChainName, m.ExternalAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid external address: %s", err)
	}
	if !m.DelegateAmount.IsValid() || !m.DelegateAmount.IsPositive() {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid delegation amount")
	}
	if m.DelegateAmount.Denom != params.BaseDenom {
		return sdkerrors.ErrInvalidRequest.Wrapf("delegate denom got %s, expected %s", m.DelegateAmount.Denom, params.BaseDenom)
	}
	return nil
}

// MsgAddDelegate //

// Route implements the sdk.Msg interface.
func (m *MsgAddDelegate) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (m *MsgAddDelegate) Type() string { return TypeMsgAddDelegate }

func (m *MsgAddDelegate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.RelayerAddress)}
}

// GetSignBytes returns the message bytes to sign over.
func (m *MsgAddDelegate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgAddDelegate) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.RelayerAddress); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "relayer address is not a valid bech32 address")
	}
	if !m.Amount.IsValid() || !m.Amount.IsPositive() {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid delegation amount")
	}
	if m.Amount.Denom != params.BaseDenom {
		return sdkerrors.ErrInvalidRequest.Wrapf("delegate denom got %s, expected %s", m.Amount.Denom, params.BaseDenom)
	}
	return nil
}

// MsgUnbondedRelayer //

func (m *MsgUnbondedRelayer) Route() string { return RouterKey }

func (m *MsgUnbondedRelayer) Type() string { return TypeMsgUnbondedRelayer }

func (m *MsgUnbondedRelayer) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.RelayerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid relayer address: %s", err)
	}
	return nil
}

func (m *MsgUnbondedRelayer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgUnbondedRelayer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.RelayerAddress)}
}

// MsgProposalRelayers //

// Route implements the sdk.Msg interface.
func (m *MsgProposalRelayers) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (m *MsgProposalRelayers) Type() string { return TypeMsgProposalRelayers }

func (m *MsgProposalRelayers) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Authority)}
}

// GetSignBytes returns the message bytes to sign over.
func (m *MsgProposalRelayers) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgProposalRelayers) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "authority is not a valid bech32 address")
	}
	if len(m.Relayers) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("relayers list cannot be empty")
	}
	for _, relayer := range m.Relayers {
		if _, err := sdk.AccAddressFromBech32(relayer); err != nil {
			return errors.Wrap(sdkerrors.ErrInvalidAddress, "relayer address is not a valid bech32 address")
		}
	}
	return nil
}

// MsgRelayerSetConfirm //

// Route should return the name of the module
func (m *MsgRelayerSetConfirm) Route() string { return RouterKey }

// Type should return the action
func (m *MsgRelayerSetConfirm) Type() string { return TypeMsgRelayerSetConfirm }

// ValidateBasic performs stateless checks
func (m *MsgRelayerSetConfirm) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.RelayerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid relayer address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.ExternalAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid external address: %s", err)
	}
	if len(m.Signature) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("empty signature")
	}
	if _, err = hex.DecodeString(m.Signature); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrap("could not hex decode signature")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgRelayerSetConfirm) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgRelayerSetConfirm) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.RelayerAddress)}
}

// MsgUpdateParams //

// Route returns the MsgUpdateParams message route.
func (m *MsgUpdateParams) Route() string { return ModuleName }

// Type returns the MsgUpdateParams message type.
func (m *MsgUpdateParams) Type() string { return TypeMsgUpdateParams }

// GetSignBytes returns the raw bytes for a MsgUpdateParams message that
// the expected signer needs to sign.
func (m *MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Authority)}
}

func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.Wrap(err, "authority")
	}
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if err := m.Params.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "params")
	}
	return nil
}
