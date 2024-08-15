package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

const TypeMsgNewRegion = "new-region"
const TypeMsgRetrieveCoinFromRegion = "retrieve-coin-from-region"
const TypeMsgRemoveRegion = "remove-region"
const TypeMsgRetrieveFeeFromGlobalAdminFeePool = "retrieve-fee-from-global-admin-fee-pool"

var _ sdk.Msg = &MsgNewRegion{}
var _ sdk.Msg = &MsgRemoveRegion{}
var _ sdk.Msg = &MsgRetrieveCoinsFromRegion{}
var _ sdk.Msg = &MsgRetrieveFeeFromGlobalAdminFeePool{}

func NewMsgNewRegion(creator string, regionId string, name string, validator string) *MsgNewRegion {
	return &MsgNewRegion{
		Creator:         creator,
		RegionId:        regionId,
		Name:            name,
		OperatorAddress: validator,
	}
}

func (msg *MsgNewRegion) Route() string {
	return RouterKey
}

func (msg *MsgNewRegion) Type() string {
	return TypeMsgNewRegion
}

func (msg *MsgNewRegion) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgNewRegion) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgNewRegion) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.ValAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
	}

	return nil
}

func NewMsgRemoveRegion(creator string, regionId string) *MsgRemoveRegion {
	return &MsgRemoveRegion{
		Creator:  creator,
		RegionId: regionId,
	}
}

func (msg *MsgRemoveRegion) Route() string {
	return RouterKey
}

func (msg *MsgRemoveRegion) Type() string {
	return TypeMsgRemoveRegion
}

func (msg *MsgRemoveRegion) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRemoveRegion) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveRegion) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

func NewMsgRetrieveCoinsFromRegion(admin string, regionId string, receiver string, amount sdk.Coins) *MsgRetrieveCoinsFromRegion {
	return &MsgRetrieveCoinsFromRegion{
		Admin:    admin,
		RegionId: regionId,
		Receiver: receiver,
		Amount:   amount,
	}
}

func (msg *MsgRetrieveCoinsFromRegion) Route() string {
	return RouterKey
}

func (msg *MsgRetrieveCoinsFromRegion) Type() string {
	return TypeMsgRetrieveCoinFromRegion
}

func (msg *MsgRetrieveCoinsFromRegion) GetSigners() []sdk.AccAddress {
	admin, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{admin}
}

func (msg *MsgRetrieveCoinsFromRegion) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRetrieveCoinsFromRegion) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
	}

	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if !msg.Amount.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	return nil
}

func NewMsgRetrieveFeeFromGlobalAdminFeePool(admin string, amount sdk.Coins) *MsgRetrieveFeeFromGlobalAdminFeePool {
	return &MsgRetrieveFeeFromGlobalAdminFeePool{
		Admin:  admin,
		Amount: amount,
	}
}

func (msg *MsgRetrieveFeeFromGlobalAdminFeePool) Route() string {
	return RouterKey
}

func (msg *MsgRetrieveFeeFromGlobalAdminFeePool) Type() string {
	return TypeMsgRetrieveFeeFromGlobalAdminFeePool
}

func (msg *MsgRetrieveFeeFromGlobalAdminFeePool) GetSigners() []sdk.AccAddress {
	admin, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{admin}
}

func (msg *MsgRetrieveFeeFromGlobalAdminFeePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRetrieveFeeFromGlobalAdminFeePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

// NewMsgDelegate creates a new MsgDelegate instance.
//
//nolint:interfacer
func NewMsgDelegate(delAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin, valStr string) *types.MsgDelegate {
	valAddrStr := valAddr.String()
	fmt.Println(valAddrStr)
	//if valStr == NotBondedPoolName && valAddr.Empty() {
	//	valAddrStr = valStr
	//}
	return &types.MsgDelegate{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: valAddrStr,
		Amount:           amount,
	}
}
