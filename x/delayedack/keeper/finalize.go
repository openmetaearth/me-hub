package keeper

import (
	"bytes"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"

	"github.com/cometbft/cometbft/libs/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	"github.com/osmosis-labs/osmosis/v15/osmoutils"

	commontypes "github.com/openmetaearth/me-hub/x/common/types"
	"github.com/openmetaearth/me-hub/x/delayedack/types"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
)

// FinalizeRollappPackets finalizes the packets for the given rollapp until the given height which is
// the end height of the latest finalized state
func (k Keeper) FinalizeRollappPackets(ctx sdk.Context, ibc porttypes.IBCModule, rollappID string, stateInfo rollapptypes.StateInfo) error {
	stateEndHeight := stateInfo.StartHeight + stateInfo.NumBlocks - 1
	rollappPendingPackets := k.ListRollappPackets(ctx, types.PendingByRollappIDByMaxHeight(rollappID, stateEndHeight))
	if len(rollappPendingPackets) == 0 {
		return nil
	}
	logger := ctx.Logger().With("module", "DelayedAckMiddleware")
	// Get the packets for the rollapp until height
	logger.Debug("finalizing IBC rollapp packets",
		"rollappID", rollappID,
		"state end height", stateEndHeight,
		"num packets", len(rollappPendingPackets))
	for _, rollappPacket := range rollappPendingPackets {
		if err := k.validatePacketProofRoot(ctx, rollappPacket, stateInfo); err != nil {
			return fmt.Errorf("validate rollapp packet proof root: %w", err)
		}
		if err := k.finalizeRollappPacket(ctx, ibc, rollappID, logger, rollappPacket); err != nil {
			return fmt.Errorf("finalize rollapp packet: %w", err)
		}
	}
	return nil
}

type consensusRoot interface {
	GetRoot() exported.Root
}

func (k Keeper) validatePacketProofRoot(
	ctx sdk.Context,
	rollappPacket commontypes.RollappPacket,
	stateInfo rollapptypes.StateInfo,
) error {
	stateRoot, ok := blockDescriptorRootAtHeight(stateInfo, rollappPacket.ProofHeight)
	if !ok {
		return nil
	}

	portID, channelID := packetProofPortChannel(rollappPacket)
	clientID, clientState, err := k.channelKeeper.GetChannelClientState(ctx, portID, channelID)
	if err != nil {
		return errorsmod.Wrapf(err, "get channel client state for %s/%s", portID, channelID)
	}

	consensusHeight := clienttypes.NewHeight(clientState.GetLatestHeight().GetRevisionNumber(), rollappPacket.ProofHeight)
	consensusState, found := k.clientKeeper.GetClientConsensusState(ctx, clientID, consensusHeight)
	if !found {
		return errorsmod.Wrapf(types.ErrUnknownRequest,
			"consensus state not found for client %s at height %s", clientID, consensusHeight)
	}

	rootedConsensusState, ok := consensusState.(consensusRoot)
	if !ok {
		return errorsmod.Wrapf(types.ErrUnknownRequest,
			"consensus state for client %s at height %s has no commitment root", clientID, consensusHeight)
	}

	consensusRoot := rootedConsensusState.GetRoot().GetHash()
	if !bytes.Equal(stateRoot, consensusRoot) {
		return errorsmod.Wrapf(types.ErrUnknownRequest,
			"rollapp state root mismatch at height %d: descriptor root %X != consensus root %X",
			rollappPacket.ProofHeight, stateRoot, consensusRoot)
	}

	return nil
}

func blockDescriptorRootAtHeight(stateInfo rollapptypes.StateInfo, proofHeight uint64) ([]byte, bool) {
	if proofHeight < stateInfo.StartHeight || proofHeight >= stateInfo.StartHeight+stateInfo.NumBlocks {
		return nil, false
	}

	for _, descriptor := range stateInfo.BDs.BD {
		if descriptor.Height == proofHeight {
			if len(descriptor.StateRoot) == 0 {
				return nil, false
			}
			return descriptor.StateRoot, true
		}
	}

	return nil, false
}

func packetProofPortChannel(rollappPacket commontypes.RollappPacket) (string, string) {
	switch rollappPacket.Type {
	case commontypes.RollappPacket_ON_RECV:
		return rollappPacket.Packet.DestinationPort, rollappPacket.Packet.DestinationChannel
	case commontypes.RollappPacket_ON_ACK, commontypes.RollappPacket_ON_TIMEOUT:
		return rollappPacket.Packet.SourcePort, rollappPacket.Packet.SourceChannel
	default:
		return rollappPacket.Packet.SourcePort, rollappPacket.Packet.SourceChannel
	}
}

type wrappedFunc func(ctx sdk.Context) error

func (k Keeper) finalizeRollappPacket(
	ctx sdk.Context,
	ibc porttypes.IBCModule,
	rollappID string,
	logger log.Logger,
	rollappPacket commontypes.RollappPacket,
) error {
	logContext := []interface{}{
		"rollappID", rollappID,
		"sequence", rollappPacket.Packet.Sequence,
		"source channel", rollappPacket.Packet.SourceChannel,
		"destination channel", rollappPacket.Packet.DestinationChannel,
		"type", rollappPacket.Type,
	}

	var packetErr error
	switch rollappPacket.Type {
	case commontypes.RollappPacket_ON_RECV:
		ack := ibc.OnRecvPacket(ctx, *rollappPacket.Packet, rollappPacket.Relayer)
		if ack != nil { // NOTE: in practice ack should not be nil, since ibc transfer core module always returns something
			packetErr = osmoutils.ApplyFuncIfNoError(ctx, k.writeRecvAck(rollappPacket, ack))
		}
	case commontypes.RollappPacket_ON_ACK:
		packetErr = osmoutils.ApplyFuncIfNoError(ctx, k.onAckPacket(rollappPacket, ibc))
	case commontypes.RollappPacket_ON_TIMEOUT:
		packetErr = osmoutils.ApplyFuncIfNoError(ctx, k.onTimeoutPacket(rollappPacket, ibc))
	default:
		logger.Error("Unknown rollapp packet type", logContext...)
	}
	// Update the packet with the error
	if packetErr != nil {
		rollappPacket.Error = packetErr.Error()
	}
	// Update status to finalized
	_, err := k.UpdateRollappPacketWithStatus(ctx, rollappPacket, commontypes.Status_FINALIZED)
	if err != nil {
		// If we failed finalizing the packet we return an error to abort the end blocker otherwise it's
		// invariant breaking
		return err
	}

	logger.Debug("finalized IBC rollapp packet", logContext...)
	return nil
}

func (k Keeper) writeRecvAck(rollappPacket commontypes.RollappPacket, ack exported.Acknowledgement) wrappedFunc {
	return func(ctx sdk.Context) (err error) {
		var chanCap *capabilitytypes.Capability
		_, chanCap, err = k.LookupModuleByChannel(
			ctx,
			rollappPacket.Packet.DestinationPort,
			rollappPacket.Packet.DestinationChannel,
		)
		if err != nil {
			return
		}
		/*
			Here, we do the inverse of what we did when we updated the packet transfer address, when we fulfilled the order
			to ensure the ack matches what the rollapp expects.
		*/
		rollappPacket = rollappPacket.RestoreOriginalTransferTarget()
		return k.WriteAcknowledgement(ctx, chanCap, rollappPacket.Packet, ack)
	}
}

func (k Keeper) onAckPacket(rollappPacket commontypes.RollappPacket, ibc porttypes.IBCModule) wrappedFunc {
	return func(ctx sdk.Context) (err error) {
		return ibc.OnAcknowledgementPacket(
			ctx,
			*rollappPacket.Packet,
			rollappPacket.Acknowledgement,
			rollappPacket.Relayer,
		)
	}
}

func (k Keeper) onTimeoutPacket(rollappPacket commontypes.RollappPacket, ibc porttypes.IBCModule) wrappedFunc {
	return func(ctx sdk.Context) (err error) {
		return ibc.OnTimeoutPacket(ctx, *rollappPacket.Packet, rollappPacket.Relayer)
	}
}
