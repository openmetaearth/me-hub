package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	types "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/cosmos/gogoproto/proto"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgNewClass{}, "cosmos-sdk/MsgNewClass", nil)
	cdc.RegisterConcrete(&MsgMintNFT{}, "cosmos-sdk/MsgMintNFT", nil)
	cdc.RegisterConcrete(&ClassMetadata{}, "metaearth.wnft.ClassMetadata", nil)
	cdc.RegisterConcrete(&Extension{}, "metaearth.wnft.Extension", nil)
}

// RegisterInterfaces registers the interfaces types with the interface registry.
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgNewClass{},
		&MsgMintNFT{},
		&MsgSend{},
	)

	registry.RegisterImplementations(
		(*proto.Message)(nil),
		&Extension{},
		&ClassMetadata{},
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
