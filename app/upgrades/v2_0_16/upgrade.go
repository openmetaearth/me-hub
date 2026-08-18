package v2_0_16 //nolint:revive

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	appkeepers "github.com/openmetaearth/me-hub/app/keepers"
	"github.com/openmetaearth/me-hub/app/upgrades"
	gravitykeeper "github.com/openmetaearth/me-hub/x/gravity/keeper"
	gravitytypes "github.com/openmetaearth/me-hub/x/gravity/types"
)

// bscSkipToEventNonce skips fraudulent BSC TransactionBatchExecutedEvent
// nonces 210-213 (submitBatch without a corresponding ME batch) so relayers
// can resume from event nonce 214 without panicking / stalling the bridge.
const bscSkipToEventNonce uint64 = 213

// CreateUpgradeHandler creates an SDK upgrade handler for v2.0.15.rc3
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ upgrades.BaseAppParamManager,
	keepers *appkeepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)
		logger.Info("upgrade starting...")

		maxSendToExternalUsdAmount := sdkmath.NewInt(10_000_000_000)
		bscParams := keepers.BscKeeper.GetParams(ctx)
		bscParams.MaxDelegate = sdkmath.NewInt(500_000_000)
		bscParams.MaxRelayers = 5
		bscParams.MaxSlashTimes = 100
		bscParams.MaxSendToExternalUsdAmount = maxSendToExternalUsdAmount
		bscParams.SlashFraction = sdk.MustNewDecFromStr("0.01")
		if err := keepers.BscKeeper.SetParams(ctx, &bscParams); err != nil {
			return nil, fmt.Errorf("failed to set BSC max send to external amount: %w", err)
		}

		tronParams := keepers.TronKeeper.GetParams(ctx)
		tronParams.MaxDelegate = sdkmath.NewInt(500_000_000)
		tronParams.MaxRelayers = 5
		tronParams.MaxSlashTimes = 100
		tronParams.MaxSendToExternalUsdAmount = maxSendToExternalUsdAmount
		tronParams.SlashFraction = sdk.MustNewDecFromStr("0.01")
		if err := keepers.TronKeeper.SetParams(ctx, &tronParams); err != nil {
			return nil, fmt.Errorf("failed to set Tron max send to external amount: %w", err)
		}

		// skipBscFraudulentBatchEvents(ctx, keepers.BscKeeper)
		logger.Info("upgrade finished successfully.")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

// skipBscFraudulentBatchEvents advances BSC lastObservedEventNonce past the
// attacker's submitBatch events (210-213) and aligns every relayer's personal
// event nonce so Attest continuity accepts the next legitimate claim (214+).
func skipBscFraudulentBatchEvents(ctx sdk.Context, k gravitykeeper.Keeper) {
	logger := ctx.Logger().With("upgrade", UpgradeName, "migration", "skip_bsc_batch_events")

	current := k.GetLastObservedEventNonce(ctx)
	if current >= bscSkipToEventNonce {
		logger.Info("BSC lastObservedEventNonce already at or past skip target; nothing to do",
			"current", current, "target", bscSkipToEventNonce)
		return
	}

	logger.Info("skipping fraudulent BSC TransactionBatchExecutedEvent nonces",
		"from", current, "to", bscSkipToEventNonce)

	k.SetLastObservedEventNonce(ctx, bscSkipToEventNonce)

	// Drop any unobserved attestations for the skipped nonces so they cannot
	// be retried after the upgrade.
	k.IterateAttestationAndClaim(ctx, func(_ *gravitytypes.Attestation, claim gravitytypes.ExternalClaim) bool {
		nonce := claim.GetEventNonce()
		if nonce > current && nonce <= bscSkipToEventNonce {
			k.DeleteAttestation(ctx, claim)
		}
		return false
	})

	for _, relayer := range k.GetAllRelayers(ctx, false) {
		addr := sdk.MustAccAddressFromBech32(relayer.RelayerAddress)
		last := k.GetLastEventNonceByRelayer(ctx, addr)
		if last < bscSkipToEventNonce {
			k.SetLastEventNonceByRelayer(ctx, addr, bscSkipToEventNonce)
			logger.Info("updated relayer lastEventNonce",
				"relayer", relayer.RelayerAddress, "from", last, "to", bscSkipToEventNonce)
		}
	}

	logger.Info("BSC event nonce skip complete",
		"last_observed_event_nonce", k.GetLastObservedEventNonce(ctx))
}
