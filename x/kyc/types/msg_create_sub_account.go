package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCreateSubAccount = "create_sub_account"

func NewMsgCreateSubAccount(creator, subAccount, subAccountPubkey string) *MsgCreateSubAccount {
	return &MsgCreateSubAccount{
		Creator:    creator,
		SubAccount: subAccount,
		SubAccountPubkey:     subAccountPubkey,
	}
}

func (m *MsgCreateSubAccount) Route() string { return RouterKey }
func (m *MsgCreateSubAccount) Type() string  { return TypeMsgCreateSubAccount }

func (m *MsgCreateSubAccount) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (m *MsgCreateSubAccount) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgCreateSubAccount) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "invalid creator address")
	}
	if _, err := sdk.AccAddressFromBech32(m.SubAccount); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "invalid sub_account address")
	}
	if m.SubAccountPubkey == "" {
		return errors.Wrap(sdkerrors.ErrInvalidPubKey, "pubkey is required")
	}
	return nil
}
