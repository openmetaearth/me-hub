package types

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/utils/gerrc"
)

const TypeMsgReplaceRollappPorposer = "replace_rollapp_proposer"

var _ sdk.Msg = &MsgReplaceProposerRequest{}

func NewMsgReplaceProposerRequest(creator, rollappId, oldProposer, newProposer string, blockHeight int64) (*MsgReplaceProposerRequest, error) {
	return &MsgReplaceProposerRequest{
		Creator: creator,
		ReplaceProposer: &MsgRepalceProposer{
			RollappId:   rollappId,
			OldProposer: oldProposer,
			NewProposer: newProposer,
			BlockHeight: blockHeight,
		},
	}, nil
}

func (msg *MsgReplaceProposerRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "creator")
	}
	if msg.ReplaceProposer == nil {
		return errorsmod.Wrap(gerrc.ErrInvalidArgument, "ReplaceProposer can not be nil")
	}
	if msg.ReplaceProposer.RollappId == "" {
		return errorsmod.Wrap(gerrc.ErrInvalidArgument, "rollapp id cannot be empty")
	}
	if _, err = sdk.AccAddressFromBech32(msg.ReplaceProposer.OldProposer); err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "old proposer")
	}
	if _, err = sdk.AccAddressFromBech32(msg.ReplaceProposer.NewProposer); err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "new proposer")
	}
	if msg.ReplaceProposer.BlockHeight < 1 {
		return errorsmod.Wrapf(gerrc.ErrInvalidArgument, "invalid block number (%d)", msg.ReplaceProposer.BlockHeight)
	}
	return nil
}
