package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/rollapp/types"
)

const drsViolationCauseGenesisTransferAfterOpen = "genesis_transfer_after_bridge_open"

// HandleDRSViolation freezes a rollapp that sends genesis transfers after the bridge is open.
func (k Keeper) HandleDRSViolation(ctx sdk.Context, rollappID string) error {
	rollapp, found := k.GetRollapp(ctx, rollappID)
	if !found {
		return types.ErrUnknownRollappID
	}

	rollapp.Frozen = true
	k.SetRollapp(ctx, rollapp)

	stateInfo, hasLatestState := k.latestStateInfo(ctx, rollappID)
	if hasLatestState && stateInfo.Sequencer != "" {
		fraudHeight := stateInfo.StartHeight + stateInfo.NumBlocks - 1
		if err := k.hooks.FraudSubmitted(ctx, rollappID, fraudHeight, stateInfo.Sequencer); err != nil {
			return err
		}
	}

	k.RevertPendingStates(ctx, rollappID)

	attrs := []sdk.Attribute{
		sdk.NewAttribute(types.AttributeKeyRollappId, rollappID),
		sdk.NewAttribute(types.AttributeKeyDRSViolationCause, drsViolationCauseGenesisTransferAfterOpen),
	}
	if hasLatestState {
		attrs = append(attrs,
			sdk.NewAttribute(types.AttributeKeyFraudHeight, fmt.Sprint(stateInfo.StartHeight+stateInfo.NumBlocks-1)),
			sdk.NewAttribute(types.AttributeKeyFraudSequencer, stateInfo.Sequencer),
		)
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeDRSViolation, attrs...))

	return nil
}

func (k Keeper) latestStateInfo(ctx sdk.Context, rollappID string) (types.StateInfo, bool) {
	stateInfoIndex, found := k.GetLatestStateInfoIndex(ctx, rollappID)
	if !found {
		return types.StateInfo{}, false
	}

	stateInfo, found := k.GetStateInfo(ctx, rollappID, stateInfoIndex.Index)
	if !found {
		return types.StateInfo{}, false
	}

	return stateInfo, true
}
