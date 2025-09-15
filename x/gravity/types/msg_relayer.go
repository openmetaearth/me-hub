package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/st-chain/me-hub/utils"
)

const (
	TypeMsgBondedRelayer = "bonded_relayer"
)

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
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "the creator is not a valid bech32 address")
	}
	if err := utils.ValidateEthereumAddress(m.ExternalAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid external address: %s", err)
	}
	if !m.DelegateAmount.IsValid() || m.DelegateAmount.IsNegative() {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid delegation amount")
	}
	return nil
}

// Route implements the sdk.Msg interface.
func (m *MsgAddDelegate) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (m *MsgAddDelegate) Type() string { return TypeMsgBondedRelayer }

func (m *MsgAddDelegate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.RelayerAddress)}
}

// GetSignBytes returns the message bytes to sign over.
func (m *MsgAddDelegate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgAddDelegate) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.RelayerAddress); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "the creator is not a valid bech32 address")
	}
	if !m.Amount.IsValid() || m.Amount.IsNegative() {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid delegation amount")
	}
	return nil
}
