package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(*codec.LegacyAmino) {}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateValidator{},
		&MsgStake{},
		&MsgUnstake{},
		&MsgNewRegion{},
		&MsgWithdrawDelegatorReward{},
		&MsgRemoveRegion{},
		&MsgNewFixedDepositCfg{},
		&MsgSetFixedDepositCfgRate{},
		&MsgSetFixedDepositCfgStatus{},
		&MsgDoFixedDeposit{},
		&MsgWithdrawFixedDeposit{},
		&MsgReplaceConsensusPubKeyRequest{},
		&MsgSendToModule{},
		&MsgGrantRegionWithdraw{},
		&MsgRevokeRegionWithdraw{},
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
