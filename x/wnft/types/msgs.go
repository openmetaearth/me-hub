package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft"
)

const (
	TypeMsgNewClass = "new_class"
	TypeMsgMintNFT  = "mint_nft"
	// TypeMsgSend nft message types
	TypeMsgSend = "send"
)

var (
	_ sdk.Msg = &MsgNewClass{}
	_ sdk.Msg = &MsgMintNFT{}
	_ sdk.Msg = &MsgSend{}
)

// Route implements the sdk.Msg interface.
func (msg MsgNewClass) Route() string { return nft.RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgNewClass) Type() string { return TypeMsgNewClass }

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgNewClass) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}
func (msg MsgNewClass) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{signer}
}
func (msg MsgNewClass) ValidateBasic() error {
	if len(msg.ClassId) == 0 {
		return nft.ErrEmptyClassID
	}

	if msg.TotalSupply == 0 && msg.ClassId != "kyc" {
		return ErrEmptyTotalSupply
	}

	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", msg.Sender)
	}

	return nil
}

func NewMsgNewClass(classId, sender, name, symbol, description, uri, uriHash string, totalSupply uint64) *MsgNewClass {
	return &MsgNewClass{
		ClassId:     classId,
		Sender:      sender,
		Name:        name,
		Symbol:      symbol,
		Description: description,
		Uri:         uri,
		UriHash:     uriHash,
		TotalSupply: totalSupply,
	}
}

// ValidateBasic implements the Msg.ValidateBasic method.
func (m MsgMintNFT) ValidateBasic() error {
	if len(m.ClassId) == 0 {
		return nft.ErrEmptyClassID
	}

	if len(m.TokenId) == 0 {
		return ErrEmptyTokenId
	}

	if len(m.Uri) == 0 {
		return ErrEmptyUri
	}

	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", m.Sender)
	}

	return nil
}

// GetSigners returns the expected signers for MsgMintNFT.
func (m MsgMintNFT) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{signer}
}

// Route implements the sdk.Msg interface.
func (msg MsgMintNFT) Route() string { return nft.RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgMintNFT) Type() string { return TypeMsgMintNFT }

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgMintNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgMintNFT(class_id, token_id, uri, uriHash, sender string) *MsgMintNFT {
	return &MsgMintNFT{
		ClassId: class_id,
		TokenId: token_id,
		Uri:     uri,
		UriHash: uriHash,
		Sender:  sender,
	}
}

// ValidateBasic implements the Msg.ValidateBasic method.
func (m MsgSend) ValidateBasic() error {
	if len(m.ClassId) == 0 {
		return nft.ErrEmptyClassID
	}

	if len(m.Id) == 0 {
		return nft.ErrEmptyNFTID
	}

	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", m.Sender)
	}

	_, err = sdk.AccAddressFromBech32(m.Receiver)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid receiver address (%s)", m.Receiver)
	}
	return nil
}

// GetSigners returns the expected signers for MsgSend.
func (m MsgSend) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{signer}
}

// GetSignBytes get the bytes for the message signer to sign on
func (msg MsgSend) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// Route implements the LegacyMsg interface.
func (msg MsgSend) Route() string { return nft.RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgSend) Type() string { return TypeMsgSend }
