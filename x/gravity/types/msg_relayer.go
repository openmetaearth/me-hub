package types

import (
	"encoding/hex"

	"cosmossdk.io/errors"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/openmetaearth/me-hub/app/params"
)

// MsgBondedRelayer //

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

func (m *MsgUnbondedRelayer) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.RelayerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid relayer address: %s", err)
	}
	return nil
}

// MsgProposalRelayers //

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

// MsgUpdateParams //

func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "authority")
	}
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return sdkerrors.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if err := m.Params.ValidateBasic(); err != nil {
		return errorsmod.Wrap(err, "params")
	}
	return nil
}
