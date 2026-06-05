package types

import (
	errorsmod "cosmossdk.io/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/decred/dcrd/dcrec/edwards"
)

const (
	TypeMsgCreateSequencer        = "create_sequencer"
	TypeMsgUnbond                 = "unbond"
	TypeMsgReplaceRollappPorposer = "replace_rollapp_proposer"
)

var (
	_ sdk.Msg                            = &MsgCreateSequencer{}
	_ sdk.Msg                            = &MsgUnbond{}
	_ codectypes.UnpackInterfacesMessage = (*MsgCreateSequencer)(nil)
)

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgCreateSequencer) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(msg.DymintPubKey, &pubKey)
}

/* --------------------------- MsgCreateSequencer --------------------------- */
func NewMsgCreateSequencer(creator string, pubkey cryptotypes.PubKey, rollappId string, description *Description, bond sdk.Coin) (*MsgCreateSequencer, error) {
	var pkAny *codectypes.Any
	if pubkey != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(pubkey); err != nil {
			return nil, err
		}
	}

	return &MsgCreateSequencer{
		Creator:      creator,
		DymintPubKey: pkAny,
		RollappId:    rollappId,
		Description:  *description,
		Bond:         bond,
	}, nil
}

func (msg *MsgCreateSequencer) Route() string {
	return RouterKey
}

func (msg *MsgCreateSequencer) Type() string {
	return TypeMsgCreateSequencer
}

func (msg *MsgCreateSequencer) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateSequencer) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateSequencer) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// public key also checked by the application logic
	if msg.DymintPubKey != nil {
		// check it is a pubkey
		if _, err = codectypes.NewAnyWithValue(msg.DymintPubKey); err != nil {
			return errorsmod.Wrapf(ErrInvalidPubKey, "invalid sequencer pubkey(%s)", err)
		}

		// cast to cryptotypes.PubKey type
		pk, ok := msg.DymintPubKey.GetCachedValue().(cryptotypes.PubKey)
		if !ok {
			return errorsmod.Wrapf(ErrInvalidType, "Expecting cryptotypes.PubKey, got %T", pk)
		}

		_, err = edwards.ParsePubKey(edwards.Edwards(), pk.Bytes())
		// err means the pubkey validation failed
		if err != nil {
			return errorsmod.Wrapf(ErrInvalidPubKey, "%s", err)
		}

	}

	if _, err := msg.Description.EnsureLength(); err != nil {
		return err
	}

	if !msg.Bond.IsValid() || !msg.Bond.IsPositive() {
		return errorsmod.Wrapf(ErrInvalidCoins, "invalid bond amount: %s", msg.Bond.String())
	}

	return nil
}

/* -------------------------------- MsgUnbond ------------------------------- */
func NewMsgUnbond(creator string) *MsgUnbond {
	return &MsgUnbond{
		Creator: creator,
	}
}

func (msg *MsgUnbond) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

func (msg *MsgUnbond) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

//

func NewMsgReplaceProposerRequest(creator, rollappId, oldProposer, newProposer string, blockHeight int64) (*MsgReplaceProposerRequest, error) {
	return &MsgReplaceProposerRequest{
		Creator: creator,
		ReplaceProposer: &MsgRepalceProposer{
			RollappId:   rollappId,
			OldProposer: oldProposer,
			NewProposer: newProposer,
			BlockHeight: blockHeight,
		},
	}, nil
}
func (msg *MsgReplaceProposerRequest) Route() string {
	return RouterKey
}

func (msg *MsgReplaceProposerRequest) Type() string {
	return TypeMsgReplaceRollappPorposer
}

func (msg *MsgReplaceProposerRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgReplaceProposerRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgReplaceProposerRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.ReplaceProposer == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "ReplaceProposer can not  be nil")
	}
	if msg.ReplaceProposer.RollappId == "" {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "rollapp id cannot be empty")
	}
	_, err = sdk.AccAddressFromBech32(msg.ReplaceProposer.OldProposer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid OldProposer address.addr = %s,err = %s",
			msg.ReplaceProposer.OldProposer, err.Error())
	}

	_, err = sdk.AccAddressFromBech32(msg.ReplaceProposer.NewProposer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid NewProposer address.addr = %s,err = %s",
			msg.ReplaceProposer.NewProposer, err.Error())
	}

	if msg.ReplaceProposer.BlockHeight < 1 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid block number (%d)", msg.ReplaceProposer.BlockHeight)
	}

	return nil
}
