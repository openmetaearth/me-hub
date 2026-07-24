package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSkipDelayRollapp = "skip_delay_rollapp"

var _ sdk.Msg = &MsgSkipDelayRollapp{}

func NewMsgSkipDelayRollapp(creator string, rollappId string, skip bool) *MsgSkipDelayRollapp {
	return &MsgSkipDelayRollapp{
		Creator:   creator,
		RollappId: rollappId,
		Skip:      skip,
	}
}

func (msg *MsgSkipDelayRollapp) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return err
	}
	if len(msg.RollappId) == 0 {
		return ErrInvalidRollappID
	}
	return nil
}
