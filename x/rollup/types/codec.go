package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgSeqStaking{}, "rollup/MsgSeqStaking", nil)
	cdc.RegisterConcrete(&MsgSeqUnStaking{}, "rollup/MsgSeqUnStaking", nil)
	cdc.RegisterConcrete(&MsgSetRollupParamsRequest{}, "rollup/SetRollupParamsRequest", nil)

}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSeqStaking{}, &MsgSeqUnStaking{}, &MsgSetRollupParamsRequest{})
	//msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
