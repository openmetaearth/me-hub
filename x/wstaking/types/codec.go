package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgUpdateValidator{}, "cosmos-sdk/MsgUpdateValidator", nil)
	cdc.RegisterConcrete(&MsgStake{}, "cosmos-sdk/MsgStake", nil)
	cdc.RegisterConcrete(&MsgUnstake{}, "cosmos-sdk/MsgUnstake", nil)
	cdc.RegisterConcrete(&MsgNewRegion{}, "cosmos-sdk/MsgNewRegion", nil)
	cdc.RegisterConcrete(&MsgWithdrawDelegatorReward{}, "cosmos-sdk/MsgWithdrawDelegatorReward", nil)
	cdc.RegisterConcrete(&MsgRemoveRegion{}, "cosmos-sdk/MsgRemoveRegion", nil)
	cdc.RegisterConcrete(&MsgReplaceConsensusPubKeyRequest{}, "cosmos-sdk/MsgReplaceConsensusPubKeyRequest", nil)
	cdc.RegisterConcrete(&MsgSendToModule{}, "cosmos-sdk/MsgSendToModule", nil)
	cdc.RegisterConcrete(&MsgGrantRegionWithdraw{}, "cosmos-sdk/MsgGrantRegionWithdraw", nil)
	cdc.RegisterConcrete(&MsgRevokeRegionWithdraw{}, "cosmos-sdk/MsgRevokeRegionWithdraw", nil)
}

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
