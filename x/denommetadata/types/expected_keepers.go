package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
)

// BankKeeper defines the expected interface needed
type BankKeeper interface {
	GetDenomMetaData(ctx sdk.Context, denom string) (types.Metadata, bool)
	SetDenomMetaData(ctx sdk.Context, denomMetaData types.Metadata)
}

type DenomMetadataKeeper interface {
	CreateDenomMetadata(ctx sdk.Context, metadata types.Metadata) error
}

type RollappKeeper interface {
	SetRollapp(ctx sdk.Context, rollapp rollapptypes.Rollapp)
	GetValidTransferFromReceivedPacket(
		ctx sdk.Context,
		packet channeltypes.Packet,
	) (data rollapptypes.TransferData, err error)
	GetValidTransferFromSentPacket(
		ctx sdk.Context,
		packet channeltypes.Packet,
	) (data rollapptypes.TransferData, err error)
	GetValidTransferFromSendPacket(
		ctx sdk.Context,
		packetData []byte,
		sourcePort, sourceChannel string,
	) (data rollapptypes.TransferData, err error)
}
