package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	didtypes "github.com/openmetaearth/me-hub/x/did/types"
)

type MsgCreateSubAccount struct {
	Creator    string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	Did        string `protobuf:"bytes,2,opt,name=did,proto3" json:"did,omitempty"`
	SubAccount string `protobuf:"bytes,3,opt,name=sub_account,json=subAccount,proto3" json:"sub_account,omitempty"`
	Pubkey     string `protobuf:"bytes,4,opt,name=pubkey,proto3" json:"pubkey,omitempty"`
}

func (m *MsgCreateSubAccount) Reset()         { *m = MsgCreateSubAccount{} }
func (m *MsgCreateSubAccount) String() string { return m.Creator }
func (*MsgCreateSubAccount) ProtoMessage()    {}

type MsgCreateSubAccountResponse struct{}

func (m *MsgCreateSubAccountResponse) Reset()         {}
func (m *MsgCreateSubAccountResponse) String() string { return "" }
func (*MsgCreateSubAccountResponse) ProtoMessage()    {}

const TypeMsgCreateSubAccount = "create_sub_account"

func NewMsgCreateSubAccount(creator, did, subAccount, pubkey string) *MsgCreateSubAccount {
	return &MsgCreateSubAccount{
		Creator:    creator,
		Did:        did,
		SubAccount: subAccount,
		Pubkey:     pubkey,
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
	if len(m.Did) != didtypes.DidLength {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "DID length must be %d", didtypes.DidLength)
	}
	if _, err := sdk.AccAddressFromBech32(m.SubAccount); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "invalid sub_account address")
	}
	if m.Pubkey == "" {
		return errors.Wrap(sdkerrors.ErrInvalidPubKey, "pubkey is required")
	}
	return nil
}
