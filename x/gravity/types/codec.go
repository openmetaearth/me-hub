package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*ExternalClaim)(nil), nil)

	cdc.RegisterConcrete(&MsgBondedRelayer{}, fmt.Sprintf("%s/%s", ModuleName, "MsgBondedRelayer"), nil)
	cdc.RegisterConcrete(&MsgAddDelegate{}, fmt.Sprintf("%s/%s", ModuleName, "MsgAddDelegate"), nil)
	cdc.RegisterConcrete(&MsgWithdrawReward{}, fmt.Sprintf("%s/%s", ModuleName, "MsgWithdrawReward"), nil)

	cdc.RegisterConcrete(&MsgRelayerSetConfirm{}, fmt.Sprintf("%s/%s", ModuleName, "MsgRelayerSetConfirm"), nil)
	cdc.RegisterConcrete(&MsgRelayerSetUpdatedClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgRelayerSetUpdatedClaim"), nil)

	cdc.RegisterConcrete(&MsgBridgeTokenClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgBridgeTokenClaim"), nil)

	cdc.RegisterConcrete(&MsgSendToMeClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgSendToMeClaim"), nil)
	cdc.RegisterConcrete(&MsgBridgeCallClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgBridgeCallClaim"), nil)

	cdc.RegisterConcrete(&MsgSendToExternal{}, fmt.Sprintf("%s/%s", ModuleName, "MsgSendToExternal"), nil)
	cdc.RegisterConcrete(&MsgCancelSendToExternal{}, fmt.Sprintf("%s/%s", ModuleName, "MsgCancelSendToExternal"), nil)
	cdc.RegisterConcrete(&MsgIncreaseBridgeFee{}, fmt.Sprintf("%s/%s", ModuleName, "MsgIncreaseBridgeFee"), nil)
	cdc.RegisterConcrete(&MsgSendToExternalClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgSendToExternalClaim"), nil)

	cdc.RegisterConcrete(&MsgRequestBatch{}, fmt.Sprintf("%s/%s", ModuleName, "MsgRequestBatch"), nil)
	cdc.RegisterConcrete(&MsgConfirmBatch{}, fmt.Sprintf("%s/%s", ModuleName, "MsgConfirmBatch"), nil)

	// register Proposal
	cdc.RegisterConcrete(&MsgUpdateParams{}, fmt.Sprintf("%s/%s", ModuleName, "MsgUpdateParams"), nil)
	cdc.RegisterConcrete(&MsgUpdateChainRelayers{}, fmt.Sprintf("%s/%s", ModuleName, "MsgUpdateChainRelayers"), nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgBondedRelayer{},
		&MsgAddDelegate{},
		&MsgWithdrawReward{},

		&MsgRelayerSetConfirm{},
		&MsgRelayerSetUpdatedClaim{},

		&MsgBridgeTokenClaim{},

		&MsgSendToMeClaim{},
		&MsgBridgeCallClaim{},

		&MsgSendToExternal{},
		&MsgCancelSendToExternal{},
		&MsgIncreaseBridgeFee{},
		&MsgSendToExternalClaim{},

		&MsgRequestBatch{},
		&MsgConfirmBatch{},

		&MsgUpdateParams{},
		&MsgUpdateChainRelayers{},
	)

	registry.RegisterInterface(
		"gravity.v1beta1.ExternalClaim",
		(*ExternalClaim)(nil),
		&MsgSendToExternalClaim{},
		&MsgSendToMeClaim{},
		&MsgBridgeCallClaim{},
		&MsgBridgeTokenClaim{},
		&MsgRelayerSetUpdatedClaim{},
	)

	//registry.RegisterImplementations(
	//	(*govv1betal.Content)(nil),
	//	&InitCrossChainParamsProposal{},
	//	&UpdateChainRelayersProposal{},
	//)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
