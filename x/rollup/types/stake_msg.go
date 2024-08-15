package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeStakeForSequencer = "StakeForSequencer"
	TypeUnStake           = "UnStake"
)

const MaxAllowedSequencers = 100

func NewMsgSeqStaking(creator string, rollappId string, version, amount uint64) *MsgSeqStaking {
	return &MsgSeqStaking{
		Creator:   creator,
		RollappId: rollappId,
		Version:   version,
		Amount:    amount,
	}
}

func (msg *MsgSeqStaking) Route() string {
	return RouterKey
}

func (msg *MsgSeqStaking) Type() string {
	return TypeStakeForSequencer
}

func (msg *MsgSeqStaking) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSeqStaking) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSeqStaking) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return err
	}
	if msg.Amount < 1 {
		return fmt.Errorf("stake amount can not be 0")
	}
	if msg.RollappId == "" {
		return fmt.Errorf("stake RollappId can not be empty")
	}
	return nil
}

func NewMsgSeqUnStaking(creator string, rollappId string, version, amount uint64) *MsgSeqUnStaking {
	return &MsgSeqUnStaking{
		Creator:   creator,
		RollappId: rollappId,
		Version:   version,
		Amount:    amount,
	}
}

func (msg *MsgSeqUnStaking) Route() string {
	return RouterKey
}

func (msg *MsgSeqUnStaking) Type() string {
	return TypeUnStake
}

func (msg *MsgSeqUnStaking) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSeqUnStaking) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSeqUnStaking) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return err
	}
	if msg.Amount < 1 {
		return fmt.Errorf("stake amount can not be 0")
	}
	if msg.RollappId == "" {
		return fmt.Errorf("stake RollappId can not be empty")
	}
	return nil
}
