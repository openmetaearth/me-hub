package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(Amino)
)

func init() {
	RegisterCodec(ModuleCdc.LegacyAmino)
}

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*ExternalClaim)(nil), nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgBondedRelayer{},
		&MsgAddDelegate{},
		&MsgRelayerSetConfirm{},
		&MsgRelayerSetUpdateClaim{},
		&MsgBridgeTokenClaim{},
		&MsgSendToMeClaim{},
		&MsgSendToExternal{},
		&MsgCancelSendToExternal{},
		&MsgIncreaseBridgeFee{},
		&MsgSendToExternalClaim{},
		&MsgRequestBatch{},
		&MsgConfirmBatch{},
		&MsgUpdateParams{},
		&MsgProposalRelayers{},
	)

	registry.RegisterInterface(
		"gravity.v1beta1.ExternalClaim",
		(*ExternalClaim)(nil),
		&MsgSendToExternalClaim{},
		&MsgSendToMeClaim{},
		&MsgBridgeTokenClaim{},
		&MsgRelayerSetUpdateClaim{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
