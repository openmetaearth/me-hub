package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgSubmitDaFraud       = "challengeDaFraud"
	TypeMsgDaFraudVerifyResult = "submitDaFraudVerifyData"
)

//var _ sdk.Msg = &MsgUpdateState{}

func (msg *MsgSubmitDaFraudRequest) Route() string {
	return RouterKey
}

func (msg *MsgSubmitDaFraudRequest) Type() string {
	return TypeMsgSubmitDaFraud
}

func (msg *MsgSubmitDaFraudRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSubmitDaFraudRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSubmitDaFraudRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// an update can't be with no BDs
	if msg.NumBlocks == uint32(0) {
		return errorsmod.Wrap(ErrInputParams, "number of blocks can not be zero")
	}

	if msg.NumBlocks > 100000 {
		return errorsmod.Wrapf(ErrInputParams, "numBlocks(%d)  exceeds max 100000", msg.NumBlocks)
	}

	// check to see that update contains all BDs
	if msg.RollappId == "" {
		return errorsmod.Wrapf(ErrInputParams, "rollappID can not be empty")
	}

	// check to see that startHeight is not zaro
	if msg.StartHeight == 0 {
		return errorsmod.Wrapf(ErrWrongBlockHeight, "StartHeight must be greater than zero")
	}

	if len(msg.Namespace) != 29 {
		return errorsmod.Wrapf(ErrInputParams, "namespace length error.len = %d", msg.NumBlocks)
	}

	if msg.DaBlockHeight < 2 {
		return errorsmod.Wrapf(ErrInputParams, "DaBlockHeight must > 2")
	}

	if nil == msg.DaRoot || nil == msg.Commitment {
		return errorsmod.Wrapf(ErrInputParams, " msg.DaRoot == nil or  nil == msg.Commitment")
	}

	return nil
}

func (msg *MsgDaFraudVerifyResult) Route() string {
	return RouterKey
}

func (msg *MsgDaFraudVerifyResult) Type() string {
	return TypeMsgDaFraudVerifyResult
}

func (msg *MsgDaFraudVerifyResult) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDaFraudVerifyResult) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDaFraudVerifyResult) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// an update can't be with no BDs
	if msg.NumBlocks == uint32(0) {
		return errorsmod.Wrap(ErrInputParams, "number of blocks can not be zero")
	}

	if msg.NumBlocks > 100000 {
		return errorsmod.Wrapf(ErrInputParams, "numBlocks(%d)  exceeds max 100000", msg.NumBlocks)
	}

	// check to see that update contains all BDs
	if msg.RollappId == "" {
		return errorsmod.Wrapf(ErrInputParams, "rollappID can not be empty")
	}

	// check to see that startHeight is not zaro
	if msg.StartHeight == 0 {
		return errorsmod.Wrapf(ErrInputParams, "StartHeight must be greater than zero")
	}

	if msg.Result < 0 {
		return errorsmod.Wrapf(ErrInputParams, "msgResult value error.val= %d length error.len = %d", msg.Result)
	}

	return nil
}
