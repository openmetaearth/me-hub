package types

import (
	gomath "math"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/app/params"
)

const (
	TypeMsgNewRegion                       = "new-region"
	TypeMsgRetrieveCoinFromRegion          = "retrieve-coin-from-region"
	TypeMsgWithdrawDelegatorReward         = "withdraw_delegator_reward"
	TypeMsgRemoveRegion                    = "remove-region"
	TypeMsgRetrieveFeeFromGlobalDaoFeePool = "retrieve-fee-from-global-dao-fee-pool"
	TypeMsgRecord                          = "new_record"
	TypeReviewRecord                       = "review_record"
	TypeMsgStake                           = "stake"
	TypeMsgUnstake                         = "unstake"
	TypeMsgWithdrawFromRegion              = "withdraw_from_region"
	TypeMsgWithdrawFromGlobalDaoFeePool    = "withdraw_from_global_dao_fee_pool"
	TypeMsgResetValidator                  = "create_validator"
	TypeMsgNewMeid                         = "new_meid"
	TypeMsgRemoveMeid                      = "remove_meid"
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

func NewMsgNewRegion(creator string, name string, validator string) *MsgNewRegion {
	return &MsgNewRegion{
		Creator:         creator,
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

func NewMsgWithdrawFromRegion(withdrawer string, regionId string, receiver string, amount sdk.Coins) *MsgWithdrawFromRegion {
	return &MsgWithdrawFromRegion{
		Withdrawer: withdrawer,
		RegionId:   regionId,
		Receiver:   receiver,
		Amount:     amount,
	}
}

func (msg *MsgWithdrawFromRegion) Route() string {
	return RouterKey
}

func (msg *MsgWithdrawFromRegion) Type() string {
	return TypeMsgWithdrawFromRegion
}

func (msg *MsgWithdrawFromRegion) GetSigners() []sdk.AccAddress {
	admin, err := sdk.AccAddressFromBech32(msg.Withdrawer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{admin}
}

func (msg *MsgWithdrawFromRegion) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdrawFromRegion) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Withdrawer)
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

// NewMsgRecord creates a new MsgNewRecord instance.
func NewMsgRecord(actionNum, url, addr string) *MsgNewRecord {
	return &MsgNewRecord{
		ActionNumber: actionNum,
		ActionUrl:    url,
		From:         addr,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgNewRecord) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgNewRecord) Type() string { return TypeMsgRecord }

// GetSigners implements the sdk.Msg interface.
func (msg MsgNewRecord) GetSigners() []sdk.AccAddress {
	staker, _ := sdk.AccAddressFromBech32(msg.From)
	return []sdk.AccAddress{staker}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgNewRecord) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
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

// Route implements the sdk.Msg interface.
func (msg MsgReviewRecord) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgReviewRecord) Type() string { return TypeReviewRecord }

// GetSigners implements the sdk.Msg interface.
func (msg MsgReviewRecord) GetSigners() []sdk.AccAddress {
	staker, _ := sdk.AccAddressFromBech32(msg.From)
	return []sdk.AccAddress{staker}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgReviewRecord) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
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

func (msg *MsgWithdrawFromGlobalDaoFeePool) Route() string {
	return RouterKey
}

func (msg *MsgWithdrawFromGlobalDaoFeePool) Type() string {
	return TypeMsgWithdrawFromGlobalDaoFeePool
}

func (msg *MsgWithdrawFromGlobalDaoFeePool) GetSigners() []sdk.AccAddress {
	admin, err := sdk.AccAddressFromBech32(msg.Withdrawer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{admin}
}

func (msg *MsgWithdrawFromGlobalDaoFeePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdrawFromGlobalDaoFeePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Withdrawer)
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

func NewMsgTransferRegion(from, to, creatorAddr string, address []string) *MsgTransferRegion {
	return &MsgTransferRegion{FromRegion: from, ToRegion: to, Address: address, Creator: creatorAddr}
}
func (msg *MsgTransferRegion) Route() string {
	return RouterKey
}

func (msg *MsgTransferRegion) Type() string {
	return "msg_transfer_meid"
}

func (msg *MsgTransferRegion) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgTransferRegion) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgTransferRegion) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
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
func (msg *MsgReplaceConsensusPubKeyRequest) Route() string {
	return RouterKey
}

func (msg *MsgReplaceConsensusPubKeyRequest) Type() string {
	return TypeMsgReplaceConsensusPubKey
}

func (msg *MsgReplaceConsensusPubKeyRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgReplaceConsensusPubKeyRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgReplaceConsensusPubKeyRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.ReplacePubKey == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "replace pubkey cannot be nil")
	}
	_, err = sdk.ValAddressFromBech32(msg.ReplacePubKey.OperatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
	}
	if msg.ReplacePubKey.BlockNumber < 1 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid block number (%d)", msg.ReplacePubKey.BlockNumber)
	}

	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgReplaceConsensusPubKeyRequest) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if msg.ReplacePubKey == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "replace pubkey cannot be nil")
	}

	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(msg.ReplacePubKey.PubKey, &pubKey)
}
