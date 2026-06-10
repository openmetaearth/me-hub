package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	commontypes "github.com/openmetaearth/me-hub/x/common/types"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
)

// ChannelKeeper defines the expected IBC channel keeper
type ChannelKeeper interface {
	LookupModuleByChannel(ctx sdk.Context, portID, channelID string) (string, *capabilitytypes.Capability, error)
}

type RollappKeeper interface {
	GetParams(ctx sdk.Context) rollapptypes.Params
	GetStateInfo(ctx sdk.Context, rollappId string, index uint64) (val rollapptypes.StateInfo, found bool)
	MustGetStateInfo(ctx sdk.Context, rollappId string, index uint64) rollapptypes.StateInfo
	GetLatestFinalizedStateIndex(ctx sdk.Context, rollappId string) (val rollapptypes.StateInfoIndex, found bool)
	GetAllRollapps(ctx sdk.Context) (list []rollapptypes.Rollapp)
	GetValidTransferFromReceivedPacket(
		ctx sdk.Context,
		packet channeltypes.Packet,
	) (data rollapptypes.TransferData, err error)
	GetValidTransferFromSentPacket(
		ctx sdk.Context,
		packet channeltypes.Packet,
	) (data rollapptypes.TransferData, err error)
	IsSkipDelayRollapp(ctx sdk.Context, rollappId string) bool
}

type EIBCKeeper interface {
	EIBCDemandOrderHandler(ctx sdk.Context, rollappPacket commontypes.RollappPacket, data transfertypes.FungibleTokenPacketData) error
}
