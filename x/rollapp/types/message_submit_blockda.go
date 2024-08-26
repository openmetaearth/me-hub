package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgSubmitBlockDA     = "submit_block_da"
	TypeGetLastSubmitBlockDA = "get_last_submit_blocksInfo"
)

//var _ sdk.Msg = &MsgUpdateState{}

func NewMsgSubmitBlkDA(creator string, rollappId string, startHeight uint64, numBlocks uint32, dAPath string, version uint64, blocks *MsgLightBlkInfos,
	daCommit []byte, daRoot []byte) *MsgBlkDAInfo {
	return &MsgBlkDAInfo{
		Creator:         creator,
		RollappId:       rollappId,
		StartHeight:     startHeight,
		NumBlocks:       numBlocks,
		DAPath:          dAPath,
		Version:         version,
		Blocks:          *blocks,
		CommitmentProof: daCommit,
		DaRoot:          daRoot,
	}
}

func (msg *MsgBlkDAInfo) Route() string {
	return RouterKey
}

func (msg *MsgBlkDAInfo) Type() string {
	return TypeMsgSubmitBlockDA
}

func (msg *MsgBlkDAInfo) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgBlkDAInfo) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBlkDAInfo) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// an update can't be with no BDs
	if msg.NumBlocks == uint32(0) {
		return errorsmod.Wrap(ErrInvalidNumBlocks, "number of blocks can not be zero")
	}

	if msg.NumBlocks > 100000 {
		return errorsmod.Wrapf(ErrInvalidNumBlocks, "numBlocks(%d)  exceeds max 100000", msg.NumBlocks)
	}

	// check to see that update contains all BDs
	if len(msg.Blocks.LightBlocks) != int(msg.NumBlocks) {
		return errorsmod.Wrapf(ErrInvalidNumBlocks, "number of blocks (%d) != number of light block(%d)", msg.NumBlocks, len(msg.Blocks.LightBlocks))
	}

	// check to see that startHeight is not zaro
	if msg.StartHeight == 0 {
		return errorsmod.Wrapf(ErrWrongBlockHeight, "StartHeight must be greater than zero")
	}

	// check that the blocks are sequential by height
	for i := uint32(0); i < msg.NumBlocks; i++ {
		if msg.Blocks.LightBlocks[i].SignedHeader.Header.Height != int64(msg.StartHeight+uint64(i)) {
			return ErrInvalidBlockSequence
		}
	}

	if msg.DaRoot == nil || msg.CommitmentProof == nil {
		return errorsmod.Wrapf(ErrValidateSubmitBlock, " msg.DaRoot == nil or msg.CommitmentProof == nil")
	}

	return nil
}

func (msg *MsgLastSubmitBlkRequest) Route() string {
	return RouterKey
}

func (msg *MsgLastSubmitBlkRequest) Type() string {
	return TypeGetLastSubmitBlockDA
}

func (msg *MsgLastSubmitBlkRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgLastSubmitBlkRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgLastSubmitBlkRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// an update can't be with no BDs
	if msg.NumBlocks == uint32(0) {
		return errorsmod.Wrap(ErrInvalidNumBlocks, "number of blocks can not be zero")
	}

	if msg.NumBlocks > 100000 {
		return errorsmod.Wrapf(ErrInvalidNumBlocks, "numBlocks(%d)  exceeds max 100000", msg.NumBlocks)
	}

	// check to see that update contains all BDs
	if len(msg.Blocks.LightBlocks) != int(msg.NumBlocks) {
		return errorsmod.Wrapf(ErrInvalidNumBlocks, "number of blocks (%d) != number of light block(%d)", msg.NumBlocks, len(msg.Blocks.LightBlocks))
	}

	// check to see that startHeight is not zaro
	if msg.StartHeight == 0 {
		return errorsmod.Wrapf(ErrWrongBlockHeight, "StartHeight must be greater than zero")
	}

	// check that the blocks are sequential by height
	for i := uint32(0); i < msg.NumBlocks; i++ {
		if msg.Blocks.LightBlocks[i].SignedHeader.Header.Height != int64(msg.StartHeight+uint64(i)) {
			return ErrInvalidBlockSequence
		}
	}

	if msg.DaRoot == nil || msg.CommitmentProof == nil {
		return errorsmod.Wrapf(ErrValidateSubmitBlock, " msg.DaRoot == nil or msg.CommitmentProof == nil")
	}

	return nil
}
