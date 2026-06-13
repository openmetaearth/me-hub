package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"

	commontypes "github.com/openmetaearth/me-hub/x/common/types"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
)

// ChannelKeeper defines the expected IBC channel keeper
type ChannelKeeper interface {
	LookupModuleByChannel(ctx sdk.Context, portID, channelID string) (string, *capabilitytypes.Capability, error)
	GetChannelClientState(ctx sdk.Context, portID, channelID string) (string, exported.ClientState, error)
}

type IBCClientKeeper interface {
	GetClientConsensusState(ctx sdk.Context, clientID string, height exported.Height) (exported.ConsensusState, bool)
}

type RollappKeeper interface {
	GetParams(ctx sdk.Context) rollapptypes.Params
	GetStateInfo(ctx sdk.Context, rollappId string, index uint64) (val rollapptypes.StateInfo, found bool)
	MustGetStateInfo(ctx sdk.Context, rollappId string, index uint64) rollapptypes.StateInfo
	GetLatestFinalizedStateIndex(ctx sdk.Context, rollappId string) (val rollapptypes.StateInfoIndex, found bool)
	GetAllRollapps(ctx sdk.Context) (list []rollapptypes.Rollapp)
	GetValidTransfer(
		ctx sdk.Context,
		packetData []byte,
		raPortOnHub, raChanOnHub string,
	) (data rollapptypes.TransferData, err error)
	IsSkipDelayRollapp(ctx sdk.Context, rollappId string) bool
}

type EIBCKeeper interface {
	EIBCDemandOrderHandler(ctx sdk.Context, rollappPacket commontypes.RollappPacket, data transfertypes.FungibleTokenPacketData) error
}
