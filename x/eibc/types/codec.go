package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

func RegisterCodec(*codec.LegacyAmino) {}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgFulfillOrder{},
		&MsgFulfillOrderAuthorized{},
		&MsgUpdateDemandOrder{},
		&MsgCreateOnDemandLP{},
		&MsgDeleteOnDemandLP{},
		&MsgTryFulfillOnDemand{},
		&MsgUpdateParams{},
	)
	registry.RegisterImplementations(
		(*authz.Authorization)(nil),
		&FulfillOrderAuthorization{},
	)
}
