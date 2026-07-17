package v3

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	rollappkeeper "github.com/openmetaearth/me-hub/x/rollapp/keeper"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
)

// migrateRollapps converts v2 rollapp objects into the v3 schema.
//
// Critical for IBC: v2 stored genesis_state.transfers_enabled (proto field 2).
// v3 reserved that field and uses transfer_proof_height (field 3) instead.
// IsTransferEnabled() is now TransferProofHeight != 0. Without this migration,
// existing IBC-connected rollapps look "pre-genesis-bridge", so normal transfer
// packets are validated as GenesisBridgeData and fail with
// "missing fields in genesis bridge info".
func migrateRollapps(ctx sdk.Context, rk *rollappkeeper.Keeper) error {
	if rk == nil {
		return fmt.Errorf("nil RollappKeeper")
	}

	defaultMinBond := rk.MinSequencerBondGlobal(ctx)
	for _, old := range rk.GetAllRollapps(ctx) {
		latestHeight, _ := rk.GetLatestHeight(ctx, old.RollappId)
		newRollapp := ConvertOldRollappToNew(old, latestHeight, defaultMinBond)
		if err := newRollapp.ValidateBasic(); err != nil {
			return fmt.Errorf("validate migrated rollapp %s: %w", old.RollappId, err)
		}
		rk.SetRollapp(ctx, newRollapp)
	}
	return nil
}

// ConvertOldRollappToNew maps a rollapp unmarshaled from v2 store bytes into v3.
func ConvertOldRollappToNew(old rollapptypes.Rollapp, latestHeight uint64, defaultMinBond sdk.Coin) rollapptypes.Rollapp {
	genesisState := old.GenesisState

	// v2 transfers_enabled is dropped on unmarshal into v3. Rollapps that already
	// completed the canonical-channel / genesis-event flow have channel_id set.
	if genesisState.TransferProofHeight == 0 && old.ChannelId != "" {
		if latestHeight > 0 {
			genesisState.TransferProofHeight = latestHeight
		} else {
			genesisState.TransferProofHeight = 1
		}
	}

	launched := old.Launched || old.ChannelId != "" || genesisState.TransferProofHeight != 0

	genesisInfo := old.GenesisInfo
	if launched {
		if !genesisInfo.Launchable() {
			// v2 had no GenesisInfo on the rollapp object; fill a sealed stub so
			// Launched && Sealed invariant holds. Native denom may be empty.
			checksum := sha256.Sum256([]byte(old.RollappId))
			prefix := bech32PrefixFromRollappID(old.RollappId)
			genesisInfo = rollapptypes.GenesisInfo{
				GenesisChecksum: hex.EncodeToString(checksum[:]),
				Bech32Prefix:    prefix,
				InitialSupply:   math.ZeroInt(),
				Sealed:          true,
			}
		} else {
			genesisInfo.Sealed = true
		}
	}

	minBond := old.MinSequencerBond
	if len(minBond) == 0 {
		if defaultMinBond.IsValid() && !defaultMinBond.IsZero() {
			minBond = sdk.NewCoins(defaultMinBond)
		} else {
			minBond = sdk.NewCoins(rollapptypes.DefaultMinSequencerBondGlobalCoin)
		}
	}

	initialSequencer := old.InitialSequencer
	if initialSequencer == "" {
		initialSequencer = "*"
	}

	vmType := old.VmType
	if vmType == rollapptypes.Rollapp_Unspecified {
		vmType = rollapptypes.Rollapp_EVM
	}

	revisions := old.Revisions
	if len(revisions) == 0 {
		revisions = []rollapptypes.Revision{{
			Number:      0,
			StartHeight: 0,
		}}
	}

	return rollapptypes.Rollapp{
		RollappId:                    old.RollappId,
		Owner:                        old.Owner,
		GenesisState:                 genesisState,
		ChannelId:                    old.ChannelId,
		Metadata:                     old.Metadata,
		GenesisInfo:                  genesisInfo,
		InitialSequencer:             initialSequencer,
		MinSequencerBond:             minBond,
		VmType:                       vmType,
		Launched:                     launched,
		PreLaunchTime:                old.PreLaunchTime,
		LivenessEventHeight:          old.LivenessEventHeight,
		LivenessCountdownStartHeight: old.LivenessCountdownStartHeight,
		Revisions:                    revisions,
		EnableTee:                    old.EnableTee,
	}
}

func bech32PrefixFromRollappID(rollappID string) string {
	return "me"
}
