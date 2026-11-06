package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgSubmitProposal{}, "regiongov/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(&MsgVote{}, "regiongov/MsgVote", nil)
	cdc.RegisterConcrete(&MsgSubmitEquityUpdate{}, "regiongov/MsgSubmitEquityUpdate", nil)
	cdc.RegisterConcrete(&MsgSubmitFreezeEquity{}, "regiongov/MsgSubmitFreezeEquity", nil)
	cdc.RegisterConcrete(&MsgSubmitUnfreezeEquity{}, "regiongov/MsgSubmitUnfreezeEquity", nil)
	cdc.RegisterConcrete(&MsgApprovePendingAction{}, "regiongov/MsgApprovePendingAction", nil)
	cdc.RegisterConcrete(&MsgRejectPendingAction{}, "regiongov/MsgRejectPendingAction", nil)
	cdc.RegisterConcrete(&MsgFreezeRegion{}, "regiongov/MsgFreezeRegion", nil)
	cdc.RegisterConcrete(&MsgUnfreezeRegion{}, "regiongov/MsgUnfreezeRegion", nil)
	cdc.RegisterConcrete(&MsgFreezeEquity{}, "regiongov/MsgFreezeEquity", nil)
	cdc.RegisterConcrete(&MsgUnfreezeEquity{}, "regiongov/MsgUnfreezeEquity", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSubmitProposal{},
		&MsgVote{},
		&MsgSubmitEquityUpdate{},
		&MsgSubmitFreezeEquity{},
		&MsgSubmitUnfreezeEquity{},
		&MsgApprovePendingAction{},
		&MsgRejectPendingAction{},
		&MsgFreezeRegion{},
		&MsgUnfreezeRegion{},
		&MsgFreezeEquity{},
		&MsgUnfreezeEquity{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(Amino)
)

func init() {
	RegisterCodec(Amino)
	cryptocodec.RegisterCrypto(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)
	Amino.Seal()
}
