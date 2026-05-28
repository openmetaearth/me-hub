package types

import (
	gomath "math"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/app/params"
)

var (
	_ sdk.Msg = &MsgStake{}
	_ sdk.Msg = &MsgUnstake{}
	_ sdk.Msg = &MsgNewRegion{}
	_ sdk.Msg = &MsgRemoveRegion{}
	_ sdk.Msg = &MsgWithdrawDelegatorReward{}
	_ sdk.Msg = &MsgWithdrawFromRegion{}
	_ sdk.Msg = &MsgWithdrawFromGlobalDaoFeePool{}
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

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgStake) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.StakerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid staker address: %s", err)
	}

	if _, err := sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return errorsmod.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid stake amount",
		)
	}

	minSelfStake := math.NewInt(int64(gomath.Pow10(params.BaseDenomUnit)))
	if !msg.Amount.Amount.Mod(minSelfStake).Equal(math.NewInt(0)) {
		return errorsmod.Wrapf(
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

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUnstake) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.StakerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid staker address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return errorsmod.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid shares amount",
		)
	}

	minSelfStake := math.NewInt(int64(gomath.Pow10(params.BaseDenomUnit)))
	if !msg.Amount.Amount.Mod(minSelfStake).Equal(math.NewInt(0)) {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"invalid unstake amount: got %s, expected %s integer multiple", msg.Amount.Amount, minSelfStake)
	}

	return nil
}

func NewMsgNewRegion(creator string, name string, validator string) *MsgNewRegion {
	return &MsgNewRegion{
		Creator:         creator,
		Name:            name,
		OperatorAddress: validator,
	}
}

func (msg *MsgNewRegion) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.ValAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
	}

	return nil
}

func NewMsgRemoveRegion(creator string, regionId string) *MsgRemoveRegion {
	return &MsgRemoveRegion{
		Creator:  creator,
		RegionId: regionId,
	}
}

func (msg *MsgRemoveRegion) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

func NewMsgWithdrawFromRegion(withdrawer string, regionId string, receiver string, amount sdk.Coins) *MsgWithdrawFromRegion {
	return &MsgWithdrawFromRegion{
		Withdrawer: withdrawer,
		RegionId:   regionId,
		Receiver:   receiver,
		Amount:     amount,
	}
}

func (msg *MsgWithdrawFromRegion) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Withdrawer)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
	}

	if !msg.Amount.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if !msg.Amount.IsAllPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	return nil
}

// NewMsgRecord creates a new MsgNewRecord instance.
func NewMsgRecord(actionNum, url, addr string) *MsgNewRecord {
	return &MsgNewRecord{
		ActionNumber: actionNum,
		ActionUrl:    url,
		From:         addr,
	}
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgNewRecord) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid staker address: %s", err)
	}
	return nil
}

// NewMsgReviewRecord creates a new MsgReviewRecord instance.
func NewMsgReviewRecord(hash, result, address, id, reviewedAddress string) *MsgReviewRecord {
	return &MsgReviewRecord{
		RecordHash:      hash,
		ReviewResult:    result,
		From:            address,
		ActionNumber:    id,
		ReviewedAddress: reviewedAddress,
	}
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgReviewRecord) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid staker address: %s", err)
	}
	if msg.ReviewResult == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid review result")
	}
	if msg.RecordHash == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid record hash")
	}
	if msg.ActionNumber == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid action number")
	}
	return nil
}

func NewMsgWithdrawFromGlobalDaoFeePool(withdrawer string, amount sdk.Coins) *MsgWithdrawFromGlobalDaoFeePool {
	return &MsgWithdrawFromGlobalDaoFeePool{
		Withdrawer: withdrawer,
		Amount:     amount,
	}
}

func (msg *MsgWithdrawFromGlobalDaoFeePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Withdrawer)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

// NewMsgDelegate creates a new MsgDelegate instance.
//
//nolint:interfacer
func NewMsgDelegate(delAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin, valStr string) *types.MsgDelegate {
	valAddrStr := valAddr.String()
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
func NewMsgUndelegate(delAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin) *types.MsgUndelegate {
	return &types.MsgUndelegate{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: valAddr.String(),
		Amount:           amount,
	}
}

func NewMsgWithdrawDelegatorReward(delAddr sdk.AccAddress, valAddr sdk.ValAddress) *MsgWithdrawDelegatorReward {
	return &MsgWithdrawDelegatorReward{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: valAddr.String(),
	}
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

func NewMsgTransferRegion(from, to, creatorAddr string, address []string) *MsgTransferRegion {
	return &MsgTransferRegion{FromRegion: from, ToRegion: to, Address: address, Creator: creatorAddr}
}

func (msg *MsgTransferRegion) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

func NewMsgReplaceConsensusPubKeyRequest(creator, operator string, pubkey cryptotypes.PubKey, blockNumber int64) (*MsgReplaceConsensusPubKeyRequest, error) {
	codecPubKey, err := codectypes.NewAnyWithValue(pubkey)
	if err != nil {
		return nil, err
	}

	return &MsgReplaceConsensusPubKeyRequest{
		Creator: creator,
		ReplacePubKey: &MsgReplaceConsensusPubKey{
			OperatorAddress: operator,
			PubKey:          codecPubKey,
			BlockNumber:     blockNumber,
		},
	}, nil
}

func (msg *MsgReplaceConsensusPubKeyRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.ReplacePubKey == nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "replace pubkey cannot be nil")
	}
	_, err = sdk.ValAddressFromBech32(msg.ReplacePubKey.OperatorAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
	}
	if msg.ReplacePubKey == nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "replace pubkey cannot be nil")
	}
	if msg.ReplacePubKey.BlockNumber < 1 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid block number (%d)", msg.ReplacePubKey.BlockNumber)
	}

	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgReplaceConsensusPubKeyRequest) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if msg.ReplacePubKey == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "replace pubkey cannot be nil")
	}

	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(msg.ReplacePubKey.PubKey, &pubKey)
}
