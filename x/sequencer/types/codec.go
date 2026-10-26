package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

var ModuleCdc codec.Codec

func init() {
	registry := cdctypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(registry)
	ModuleCdc = codec.NewProtoCodec(registry)
}

func RegisterCodec(*codec.LegacyAmino) {}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateSequencer{},
		&MsgUnbond{},
		&MsgIncreaseBond{},
		&MsgDecreaseBond{},
		&MsgUpdateRewardAddress{},
		&MsgUpdateWhitelistedRelayers{},
		&MsgKickProposer{},
		&MsgUpdateOptInStatus{},
		&MsgUpdateParams{},
		&MsgPunishSequencer{},
		&MsgUpdateSequencerInformation{},
		&MsgReplaceProposerRequest{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
