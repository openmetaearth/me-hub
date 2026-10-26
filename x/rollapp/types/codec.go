package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(*codec.LegacyAmino) {}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateRollapp{},
		&MsgUpdateRollappInformation{},
		&MsgTransferOwnership{},
		&MsgUpdateState{},
		&MsgAddApp{},
		&MsgUpdateApp{},
		&MsgRemoveApp{},
		&MsgRollappFraudProposal{},
		&MsgMarkObsoleteRollapps{},
		&MsgForceGenesisInfoChange{},
		&MsgUpdateParams{},
		&MsgToggleTEE{},
		&MsgSkipDelayRollapp{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
