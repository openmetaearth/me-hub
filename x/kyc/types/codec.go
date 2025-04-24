package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgApprove{}, "kyc/Approve", nil)
	cdc.RegisterConcrete(&MsgUpdate{}, "kyc/Update", nil)
	cdc.RegisterConcrete(&MsgRemove{}, "kyc/Remove", nil)
	cdc.RegisterConcrete(&MsgCreateSBT{}, "kyc/CreateSBT", nil)
	cdc.RegisterConcrete(&MsgDeleteSBT{}, "kyc/DeleteSBT", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgApprove{},
		&MsgUpdate{},
		&MsgRemove{},
		&MsgCreateSBT{},
		&MsgDeleteSBT{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
