package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateDid{}, "did/MsgCreateDid", nil)
	cdc.RegisterConcrete(&MsgUpdateDidStatus{}, "did/MsgUpdateDidStatus", nil)
	//cdc.RegisterConcrete(&MsgRemoveDid{}, "did/MsgRemoveDid", nil)
	cdc.RegisterConcrete(&MsgCreateService{}, "did/MsgCreateService", nil)
	cdc.RegisterConcrete(&MsgUpdateServiceStatus{}, "did/MsgUpdateServiceStatus", nil)
	//cdc.RegisterConcrete(&MsgRemoveService{}, "did/MsgRemoveService", nil)
	cdc.RegisterConcrete(&MsgCreateVC{}, "did/MsgCreateVC", nil)
	cdc.RegisterConcrete(&MsgUpdateVC{}, "did/MsgUpdateVC", nil)
	cdc.RegisterConcrete(&MsgRemoveVC{}, "did/MsgRemoveVC", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateDid{},
		&MsgUpdateDidStatus{},
		//&MsgRemoveDid{},
		&MsgCreateService{},
		&MsgUpdateServiceStatus{},
		//&MsgRemoveService{},
		&MsgCreateVC{},
		&MsgUpdateVC{},
		&MsgRemoveVC{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
