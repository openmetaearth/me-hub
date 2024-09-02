package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgWithdrawTreasury{}, "cosmos-sdk/MsgWithdrawTreasury", nil)
	cdc.RegisterConcrete(&MsgSendToTreasury{}, "cosmos-sdk/MsgSendToTreasury", nil)
	cdc.RegisterConcrete(&MsgSendToAirdrop{}, "cosmos-sdk/MsgSendToAirdrop", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgWithdrawTreasury{},
		&MsgSendToTreasury{},
		&MsgSendToAirdrop{},
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
