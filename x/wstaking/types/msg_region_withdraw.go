package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewMsgGrantRegionWithdraw(creator, regionId, address string) *MsgGrantRegionWithdraw {
	return &MsgGrantRegionWithdraw{
		Creator:  creator,
		RegionId: regionId,
		Address:  address,
	}
}

func (msg *MsgGrantRegionWithdraw) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address: %s", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.Address); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid granted address: %s", err)
	}
	if msg.RegionId == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "region_id cannot be empty")
	}
	return nil
}

func NewMsgRevokeRegionWithdraw(creator, regionId string) *MsgRevokeRegionWithdraw {
	return &MsgRevokeRegionWithdraw{
		Creator:  creator,
		RegionId: regionId,
	}
}

func (msg *MsgRevokeRegionWithdraw) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address: %s", err)
	}
	if msg.RegionId == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "region_id cannot be empty")
	}
	return nil
}
