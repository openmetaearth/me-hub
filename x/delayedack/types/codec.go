package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(*codec.LegacyAmino) {}

// RegisterInterfaces registers interfaces types with the interface registry.
func RegisterInterfaces(reg types.InterfaceRegistry) {
	reg.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgFinalizePacket{},
		&MsgFinalizePacketByPacketKey{},
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(reg, &_Msg_serviceDesc)
}
