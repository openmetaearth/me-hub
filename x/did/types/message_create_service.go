package types

import (
	"cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateService = "create_service"
)

func NewMsgCreateService(creator, sid, name, description string, issuers []string) *MsgCreateService {
	return &MsgCreateService{
		Creator:     creator,
		Sid:         sid,
		Name:        name,
		Description: description,
		Issuers:     issuers,
	}
}

// Route implements the sdk.Msg interface.
func (m *MsgCreateService) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (m *MsgCreateService) Type() string { return TypeMsgCreateService }

func (m *MsgCreateService) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

// GetSignBytes returns the message bytes to sign over.
func (m *MsgCreateService) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgCreateService) GetService() Service {
	return NewService(m.Sid, m.Name, m.Description, SERVICE_STATUS_ACTIVE, m.Issuers)
}

func (m *MsgCreateService) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "the creator is not a valid bech32 address")
	}

	if len(m.Sid) < 2 || len(m.Sid) > 8 {
		return errors.Wrap(sdkerrors.ErrInvalidType, "sid length must be between 2 and 8")
	}
	if len(m.Name) == 0 || len(m.Name) > 20 {
		return errors.Wrap(sdkerrors.ErrInvalidType, "name length exceeds 8")
	}
	if len(m.Description) > 1024 {
		return errors.Wrap(sdkerrors.ErrInvalidType, "description length exceeds 1024")
	}
	for _, issuer := range m.Issuers {
		if len(issuer) != DidLength {
			return errors.Wrap(sdkerrors.ErrInvalidType, fmt.Sprintf("issuer length must be equal to %d", DidLength))
		}
	}

	return nil
}
