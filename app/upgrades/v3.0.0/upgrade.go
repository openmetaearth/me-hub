package v3_0_0

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/openmetaearth/me-hub/app/upgrades"
	legacydelayedack "github.com/openmetaearth/me-hub/app/upgrades/v3.0.0/types/delayedack"
	legacyeibc "github.com/openmetaearth/me-hub/app/upgrades/v3.0.0/types/eibc"
	legacyrollapp "github.com/openmetaearth/me-hub/app/upgrades/v3.0.0/types/rollapp"
	legacysequencer "github.com/openmetaearth/me-hub/app/upgrades/v3.0.0/types/sequencer"
	delayedacktypes "github.com/openmetaearth/me-hub/x/delayedack/types"
	eibctypes "github.com/openmetaearth/me-hub/x/eibc/types"
	rollappkeeper "github.com/openmetaearth/me-hub/x/rollapp/keeper"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	sequencertypes "github.com/openmetaearth/me-hub/x/sequencer/types"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v3.0.0.
//
// This upgrade:
//  1. Migrates the chain from Cosmos SDK v0.47 to v0.50 (incl. wstaking validators).
//  2. Aligns settlement modules with Dymension Hub v4 (params → module store, new schemas).
//  3. Adds the lightclient module store key.
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

		// Existing rollapps need MinSequencerBond filled (moved out of sequencer params).
		backfillRollappMinSequencerBond(ctx, keepers.RollappKeeper)
		logger.Info("backfilled rollapp min_sequencer_bond")

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

func backfillRollappMinSequencerBond(ctx sdk.Context, rk *rollappkeeper.Keeper) {
	minBond := rk.MinSequencerBondGlobal(ctx)
	for _, ra := range rk.GetAllRollapps(ctx) {
		if len(ra.MinSequencerBond) == 0 && minBond.IsValid() && !minBond.IsZero() {
			ra.MinSequencerBond = sdk.NewCoins(minBond)
			rk.SetRollapp(ctx, ra)
		}
	}
}
