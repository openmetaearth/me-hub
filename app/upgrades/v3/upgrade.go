package v3

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/openmetaearth/me-hub/app/upgrades"
	legacydelayedack "github.com/openmetaearth/me-hub/app/upgrades/v3/types/delayedack"
	legacyeibc "github.com/openmetaearth/me-hub/app/upgrades/v3/types/eibc"
	legacyrollapp "github.com/openmetaearth/me-hub/app/upgrades/v3/types/rollapp"
	legacysequencer "github.com/openmetaearth/me-hub/app/upgrades/v3/types/sequencer"
	delayedacktypes "github.com/openmetaearth/me-hub/x/delayedack/types"
	eibctypes "github.com/openmetaearth/me-hub/x/eibc/types"
	lightclientkeeper "github.com/openmetaearth/me-hub/x/lightclient/keeper"
	lightclienttypes "github.com/openmetaearth/me-hub/x/lightclient/types"
	rollappkeeper "github.com/openmetaearth/me-hub/x/rollapp/keeper"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	sequencerkeeper "github.com/openmetaearth/me-hub/x/sequencer/keeper"
	sequencertypes "github.com/openmetaearth/me-hub/x/sequencer/types"
	wgovkeeper "github.com/openmetaearth/me-hub/x/wgov/keeper"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v3.
//
// This upgrade:
//  1. Migrates Cosmos SDK v0.47 → v0.50 (legacy baseapp consensus params if any; gov 4→5 via RunMigrations).
//  2. Aligns settlement modules with Dymension Hub v4 (params → module store, new schemas).
//  3. Migrates rollapps (v2 transfers_enabled → transfer_proof_height; launched/genesis info).
//  4. Adds the lightclient module store key and backfills canonical clients.
//  5. Backfills sequencer dymint-addr / proposer indexes required by lightclient IBC checks.
//
// StoreUpgrades only add lightclient — consensus already existed on med-v2 from genesis.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *upgrades.UpgradeKeepers,
) upgradetypes.UpgradeHandler {
	return func(goCtx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(goCtx)
		logger := ctx.Logger().With("upgrade", UpgradeName)

		// NOTE: DO NOT PUT ANY STATE CHANGES BEFORE RunMigrations().
		migrations, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}

		if err := migrateModuleParams(ctx, keepers); err != nil {
			return nil, fmt.Errorf("migrate consensus params: %w", err)
		}
		logger.Info("migrated consensus params from x/params to x/consensus")

		// RunMigrations already ran gov Migrate4to5 (adds expedited_* from DefaultParams).
		// Fix stake denom + expedited period vs chain VotingPeriod.
		if err := updateGovParams(ctx, keepers.GovKeeper); err != nil {
			panic(fmt.Sprintf("update gov params: %v", err))
		}
		logger.Info("updated gov expedited params after Migrate4to5")

		oldSeqParams := migrateSequencerParams(ctx, keepers)
		logger.Info("migrated sequencer params to module store",
			"notice_period", keepers.SequencerKeeper.GetParams(ctx).NoticePeriod.String(),
		)

		migrateRollappParams(ctx, keepers.RollappKeeper, keepers, oldSeqParams.MinBond)
		logger.Info("migrated rollapp params to module store",
			"dispute_period_in_blocks", keepers.RollappKeeper.DisputePeriodInBlocks(ctx),
			"min_sequencer_bond_global", keepers.RollappKeeper.MinSequencerBondGlobal(ctx).String(),
		)

		migrateDelayedAckParams(ctx, keepers)
		logger.Info("migrated delayedack params to module store")

		migrateEIBCParams(ctx, keepers)
		logger.Info("migrated eibc params to module store")

		// MUST run before migrateRollapps: v2 registeredDenoms live on the rollapp
		// object (field 10), which is dropped when the object is rewritten to v3.
		if err := migrateRollappRegisteredDenoms(ctx, keepers.RollappKeeper); err != nil {
			return nil, fmt.Errorf("migrate rollapp registered denoms: %w", err)
		}
		logger.Info("migrated rollapp registered denoms to keyset")

		// Convert v2 rollapps → v3 (transfers_enabled → transfer_proof_height, launched, etc).
		// MUST run before lightclient migration / before any IBC traffic after upgrade.
		if err := migrateRollapps(ctx, keepers.RollappKeeper); err != nil {
			return nil, fmt.Errorf("migrate rollapps: %w", err)
		}
		logger.Info("migrated rollapps to v3 schema")

		// MUST run before migrateRollappLightClients: canonical clients require
		// SequencerByDymintAddr so MsgUpdateClient from sequencers is accepted.
		if err := backfillSequencerIndexes(ctx, keepers.SequencerKeeper); err != nil {
			return nil, fmt.Errorf("backfill sequencer indexes: %w", err)
		}

		migrateSequencers(ctx, keepers.SequencerKeeper)

		logger.Info("backfilled sequencer dymint-addr and proposer indexes")

		if err := migrateRollappLightClients(ctx, keepers.RollappKeeper, keepers.LightClientKeeper, keepers.IBCKeeper.ChannelKeeper); err != nil {
			return nil, fmt.Errorf("migrate rollapp light clients: %w", err)
		}
		logger.Info("migrated rollapp canonical light clients")

		logger.Info("added lightclient store key")
		logger.Info("upgrade finished successfully.")
		return migrations, nil
	}
}

func legacySubspace(keepers *upgrades.UpgradeKeepers, name string, kt paramstypes.KeyTable) paramstypes.Subspace {
	ss := keepers.ParamsKeeper.Subspace(name)
	if !ss.HasKeyTable() {
		ss = ss.WithKeyTable(kt)
	}
	return ss
}

// migrateModuleParams moves CometBFT consensus parameters from the legacy
// x/params "baseapp" subspace into x/consensus. Required for SDK v0.47 → v0.50.
//
//nolint:staticcheck // ConsensusParamsKeyTable / Paramspace are intentionally used for upgrade migration.
func migrateModuleParams(ctx sdk.Context, keepers *upgrades.UpgradeKeepers) error {
	if keepers.ConsensusKeeper == nil {
		return errorsmod.Wrap(fmt.Errorf("nil ConsensusKeeper"), "migrate consensus params")
	}

	baseAppLegacySS := keepers.ParamsKeeper.Subspace(baseapp.Paramspace).
		WithKeyTable(paramstypes.ConsensusParamsKeyTable()) //nolint:staticcheck

	return baseapp.MigrateParams(ctx, baseAppLegacySS, keepers.ConsensusKeeper.ParamsStore)
}

func migrateSequencerParams(ctx sdk.Context, keepers *upgrades.UpgradeKeepers) legacysequencer.Params {
	ss := legacySubspace(keepers, legacysequencer.ModuleName, legacysequencer.ParamKeyTable())

	var old legacysequencer.Params
	if ss.Has(ctx, legacysequencer.KeyMinBond) {
		ss.GetParamSet(ctx, &old)
	}

	// Proto schema changed completely (min_bond/unbonding_time → notice_period + liveness/dishonor).
	// Start from defaults, then preserve unbonding_time as notice_period when present.
	p := sequencertypes.DefaultParams()
	if old.UnbondingTime > 0 {
		p.NoticePeriod = old.UnbondingTime
	}
	keepers.SequencerKeeper.SetParams(ctx, p)
	return old
}

func migrateRollappParams(
	ctx sdk.Context,
	rk *rollappkeeper.Keeper,
	keepers *upgrades.UpgradeKeepers,
	oldMinBond sdk.Coin,
) {
	ss := legacySubspace(keepers, legacyrollapp.ModuleName, legacyrollapp.ParamKeyTable())

	params := rollapptypes.DefaultParams()

	if ss.Has(ctx, legacyrollapp.KeyDisputePeriodInBlocks) {
		var disputePeriod uint64
		ss.Get(ctx, legacyrollapp.KeyDisputePeriodInBlocks, &disputePeriod)
		if disputePeriod > 0 {
			params.DisputePeriodInBlocks = disputePeriod
		}
	}

	// min_bond moved from sequencer params → rollapp.params.min_sequencer_bond_global
	if oldMinBond.IsValid() && !oldMinBond.IsZero() {
		params.MinSequencerBondGlobal = oldMinBond
	}

	rk.SetParams(ctx, params)
}

func migrateDelayedAckParams(ctx sdk.Context, keepers *upgrades.UpgradeKeepers) {
	ss := legacySubspace(keepers, legacydelayedack.ModuleName, legacydelayedack.ParamKeyTable())

	params := delayedacktypes.DefaultParams()
	if ss.Has(ctx, legacydelayedack.KeyEpochIdentifier) {
		var old legacydelayedack.Params
		ss.GetParamSet(ctx, &old)
		if old.EpochIdentifier != "" {
			params.EpochIdentifier = old.EpochIdentifier
		}
		if !old.BridgingFee.IsNil() {
			params.BridgingFee = old.BridgingFee
		}
		if old.DeletePacketsEpochLimit >= 0 {
			params.DeletePacketsEpochLimit = old.DeletePacketsEpochLimit
		}
	}
	keepers.DelayedAckKeeper.SetParams(ctx, params)
}

func migrateEIBCParams(ctx sdk.Context, keepers *upgrades.UpgradeKeepers) {
	ss := legacySubspace(keepers, legacyeibc.ModuleName, legacyeibc.ParamKeyTable())

	params := eibctypes.DefaultParams()
	if ss.Has(ctx, legacyeibc.KeyEpochIdentifier) {
		var old legacyeibc.Params
		ss.GetParamSet(ctx, &old)
		if old.EpochIdentifier != "" {
			params.EpochIdentifier = old.EpochIdentifier
		}
		if !old.TimeoutFee.IsNil() {
			params.TimeoutFee = old.TimeoutFee
		}
		if !old.ErrackFee.IsNil() {
			params.ErrackFee = old.ErrackFee
		}
	}
	keepers.EIBCKeeper.SetParams(ctx, params)
}

// backfillSequencerIndexes fills v3-only indexes that InitGenesis would set for
// new chains, but which are missing after an in-place upgrade from med-v2:
//   - dymintProposerAddr → sequencer account (required by lightclient MsgUpdateClient)
//   - opted_in=true for bonded sequencers (v2 had no opted_in; default false blocks proposers)
//   - proposer-by-rollapp key (v2 stored proposer as a bool on the sequencer object)
func backfillSequencerIndexes(ctx sdk.Context, k *sequencerkeeper.Keeper) error {
	if k == nil {
		return fmt.Errorf("nil SequencerKeeper")
	}

	// rollappID → best bonded sequencer to assign as proposer if missing
	proposers := map[string]sequencertypes.Sequencer{}

	for _, seq := range k.AllSequencers(ctx) {
		if seq.DymintPubKey != nil {
			proposerAddr, err := seq.ProposerAddr()
			if err != nil {
				return fmt.Errorf("sequencer %s proposer addr: %w", seq.Address, err)
			}
			if err := k.SetSequencerByDymintAddr(ctx, proposerAddr, seq.Address); err != nil {
				return fmt.Errorf("set sequencer by dymint addr %s: %w", seq.Address, err)
			}
		}

		if seq.Bonded() && !seq.OptedIn {
			seq.OptedIn = true
			k.SetSequencer(ctx, seq)
		}

		if !seq.Bonded() {
			continue
		}
		cur, ok := proposers[seq.RollappId]
		if !ok || seq.Tokens.IsAllGT(cur.Tokens) {
			proposers[seq.RollappId] = seq
		}
	}

	for rollappID, seq := range proposers {
		if k.GetProposer(ctx, rollappID).Sentinel() {
			k.SetProposer(ctx, rollappID, seq.Address)
		}
	}

	return nil
}

// migrateRollappLightClients backfills canonical IBC client IDs for rollapps that
// already have a ChannelId from me-hub v2 (transfergenesis / canonical channel hack).
func migrateRollappLightClients(
	ctx sdk.Context,
	rk *rollappkeeper.Keeper,
	lightClientKeeper *lightclientkeeper.Keeper,
	ibcChannelKeeper lightclienttypes.IBCChannelKeeperExpected,
) error {
	if lightClientKeeper == nil || ibcChannelKeeper == nil {
		return errorsmod.Wrap(fmt.Errorf("nil keeper"), "lightclient migration")
	}

	for _, rollapp := range rk.GetAllRollapps(ctx) {
		if rollapp.ChannelId == "" {
			continue
		}

		_, connection, err := ibcChannelKeeper.GetChannelConnection(ctx, ibctransfertypes.PortID, rollapp.ChannelId)
		if err != nil {
			// Match Dymension: skip if connection cannot be resolved for this channel.
			ctx.Logger().With("upgrade", UpgradeName).Error(
				"skip canonical client migration: channel connection not found",
				"rollapp_id", rollapp.RollappId,
				"channel_id", rollapp.ChannelId,
				"err", err,
			)
			continue
		}

		lightClientKeeper.SetCanonicalClient(ctx, rollapp.RollappId, connection.GetClientID())
	}

	return nil
}

// updateGovParams fixes gov params after SDK v0.47 → v0.50 (gov consensus 4 → 5).
// Migrate4to5 fills new expedited_* fields from DefaultParams:
//   - ExpeditedMinDeposit uses denom "stake" → rewrite to chain MinDeposit denom
//   - ExpeditedVotingPeriod defaults to 1 day, which is invalid when VotingPeriod
//     is shorter (common on testnets, e.g. 300s) → clamp to VotingPeriod / 2
func updateGovParams(ctx sdk.Context, k *wgovkeeper.Keeper) error {
	if k == nil {
		return fmt.Errorf("nil GovKeeper")
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return fmt.Errorf("get gov params: %w", err)
	}
	if len(params.MinDeposit) == 0 {
		return fmt.Errorf("gov MinDeposit is empty after migration")
	}

	// Expedited min deposit = 5 × min deposit (same denom as chain bond denom).
	params.ExpeditedMinDeposit = sdk.NewCoins(
		sdk.NewCoin(params.MinDeposit[0].Denom, params.MinDeposit[0].Amount.MulRaw(5)),
	)

	if params.VotingPeriod != nil {
		if params.ExpeditedVotingPeriod == nil ||
			params.ExpeditedVotingPeriod.Seconds() >= params.VotingPeriod.Seconds() {
			expedited := *params.VotingPeriod / 2
			params.ExpeditedVotingPeriod = &expedited
		}
	}

	if err := params.ValidateBasic(); err != nil {
		return fmt.Errorf("gov params invalid after migration fix: %w", err)
	}

	if err := k.Params.Set(ctx, params); err != nil {
		return fmt.Errorf("set gov params: %w", err)
	}
	return nil
}