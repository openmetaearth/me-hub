package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateGroup{}, "megroup/CreateGroup", nil)
	cdc.RegisterConcrete(&MsgUpdateGroup{}, "megroup/UpdateGroup", nil)
	cdc.RegisterConcrete(&MsgDeleteGroup{}, "megroup/DeleteGroup", nil)
	cdc.RegisterConcrete(&MsgJoinGroup{}, "megroup/JoinGroup", nil)
	cdc.RegisterConcrete(&MsgLeaveGroupRequest{}, "megroup/MsgLeaveGroupRequest", nil)

	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateGroup{},
		&MsgUpdateGroup{},
		&MsgDeleteGroup{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgJoinGroup{},
		&MsgLeaveGroupRequest{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
