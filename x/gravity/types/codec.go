package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
)

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(Amino)
)

func init() {
	RegisterCodec(ModuleCdc.LegacyAmino)

	// Register all Amino interfaces and concrete types on the authz Amino codec so that this can later be
	// used to properly serialize MsgGrant and MsgExec instances
	RegisterCodec(authzcodec.Amino)
}

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*ExternalClaim)(nil), nil)

	cdc.RegisterConcrete(&MsgBondedRelayer{}, fmt.Sprintf("%s/%s", ModuleName, "MsgBondedRelayer"), nil)
	cdc.RegisterConcrete(&MsgAddDelegate{}, fmt.Sprintf("%s/%s", ModuleName, "MsgAddDelegate"), nil)

	cdc.RegisterConcrete(&MsgRelayerSetConfirm{}, fmt.Sprintf("%s/%s", ModuleName, "MsgRelayerSetConfirm"), nil)
	cdc.RegisterConcrete(&MsgRelayerSetUpdateClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgRelayerSetUpdateClaim"), nil)

	cdc.RegisterConcrete(&MsgBridgeTokenClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgBridgeTokenClaim"), nil)

	cdc.RegisterConcrete(&MsgSendToMeClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgSendToMeClaim"), nil)

	cdc.RegisterConcrete(&MsgSendToExternal{}, fmt.Sprintf("%s/%s", ModuleName, "MsgSendToExternal"), nil)
	cdc.RegisterConcrete(&MsgCancelSendToExternal{}, fmt.Sprintf("%s/%s", ModuleName, "MsgCancelSendToExternal"), nil)
	cdc.RegisterConcrete(&MsgIncreaseBridgeFee{}, fmt.Sprintf("%s/%s", ModuleName, "MsgIncreaseBridgeFee"), nil)
	cdc.RegisterConcrete(&MsgSendToExternalClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgSendToExternalClaim"), nil)

	cdc.RegisterConcrete(&MsgRequestBatch{}, fmt.Sprintf("%s/%s", ModuleName, "MsgRequestBatch"), nil)
	cdc.RegisterConcrete(&MsgConfirmBatch{}, fmt.Sprintf("%s/%s", ModuleName, "MsgConfirmBatch"), nil)

	cdc.RegisterConcrete(&MsgUpdateParams{}, fmt.Sprintf("%s/%s", ModuleName, "MsgUpdateParams"), nil)
	cdc.RegisterConcrete(&MsgProposalRelayers{}, fmt.Sprintf("%s/%s", ModuleName, "MsgProposalRelayers"), nil)
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
