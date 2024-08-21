package types

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/st-chain/me-hub/app/params"
	gomath "math"
)

const (
	TypeMsgNewRegion                         = "new-region"
	TypeMsgRetrieveCoinFromRegion            = "retrieve-coin-from-region"
	TypeMsgWithdrawDelegatorReward           = "withdraw_delegator_reward"
	TypeMsgRemoveRegion                      = "remove-region"
	TypeMsgRetrieveFeeFromGlobalAdminFeePool = "retrieve-fee-from-global-admin-fee-pool"
	TypeMsgStake                             = "stake"
	TypeMsgUnstake                           = "begin_unstaking"
)

var (
	_ sdk.Msg = &MsgStake{}
	_ sdk.Msg = &MsgUnstake{}
	_ sdk.Msg = &MsgNewRegion{}
	_ sdk.Msg = &MsgRemoveRegion{}
	_ sdk.Msg = &MsgWithdrawDelegatorReward{}
	_ sdk.Msg = &MsgRetrieveCoinsFromRegion{}
	_ sdk.Msg = &MsgRetrieveFeeFromGlobalAdminFeePool{}
)

// NewMsgStake creates a new MsgStake instance.
//
//nolint:interfacer
func NewMsgStake(stakerAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin) *MsgStake {
	return &MsgStake{
		StakerAddress:    stakerAddr.String(),
		ValidatorAddress: valAddr.String(),
		Amount:           amount,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgStake) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgStake) Type() string { return TypeMsgStake }

// GetSigners implements the sdk.Msg interface.
func (msg MsgStake) GetSigners() []sdk.AccAddress {
	staker, _ := sdk.AccAddressFromBech32(msg.StakerAddress)
	return []sdk.AccAddress{staker}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgStake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgStake) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.StakerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid staker address: %s", err)
	}

	if _, err := sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid stake amount",
		)
	}

	minSelfStake := math.NewInt(int64(gomath.Pow10(params.BaseDenomUnit)))
	if !msg.Amount.Amount.Mod(minSelfStake).Equal(math.NewInt(0)) {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"invalid stake amount: got %s, expected %s integer multiple", msg.Amount.Amount, minSelfStake)
	}
	return nil
}

// NewMsgUnstake creates a new MsgUnstake instance.
//
//nolint:interfacer
func NewMsgUnstake(stakerAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin) *MsgUnstake {
	return &MsgUnstake{
		StakerAddress:    stakerAddr.String(),
		ValidatorAddress: valAddr.String(),
		Amount:           amount,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUnstake) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUnstake) Type() string { return TypeMsgUnstake }

// GetSigners implements the sdk.Msg interface.
func (msg MsgUnstake) GetSigners() []sdk.AccAddress {
	staker, _ := sdk.AccAddressFromBech32(msg.StakerAddress)
	return []sdk.AccAddress{staker}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgUnstake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUnstake) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.StakerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid staker address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid shares amount",
		)
	}

	minSelfStake := math.NewInt(int64(gomath.Pow10(params.BaseDenomUnit)))
	if !msg.Amount.Amount.Mod(minSelfStake).Equal(math.NewInt(0)) {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"invalid unstake amount: got %s, expected %s integer multiple", msg.Amount.Amount, minSelfStake)
	}

	return nil
}

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

// NewMsgUndelegate creates a new MsgUndelegate instance.
//
//nolint:interfacer
func NewMsgUndelegate(delAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin, isMeid bool) *types.MsgUndelegate {
	return &types.MsgUndelegate{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: valAddr.String(),
		Amount:           amount,
		IsMeid:           isMeid,
	}
}

func NewMsgWithdrawDelegatorReward(delAddr sdk.AccAddress, valAddr sdk.ValAddress) *MsgWithdrawDelegatorReward {
	return &MsgWithdrawDelegatorReward{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: valAddr.String(),
	}
}

func (msg MsgWithdrawDelegatorReward) Route() string { return ModuleName }
func (msg MsgWithdrawDelegatorReward) Type() string  { return TypeMsgWithdrawDelegatorReward }

// GetSigners Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawDelegatorReward) GetSigners() []sdk.AccAddress {
	delegator, _ := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	return []sdk.AccAddress{delegator}
}

// GetSignBytes get the bytes for the message signer to sign on
func (msg MsgWithdrawDelegatorReward) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic quick validity check
func (msg MsgWithdrawDelegatorReward) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DelegatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}
	if len(msg.ValidatorAddress) > 0 {
		if _, err := sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
			if _, err := sdk.AccAddressFromBech32(msg.ValidatorAddress); err != nil {
				return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s,err=%s", msg.ValidatorAddress, err)
			}
		}
	}

	return nil
}
