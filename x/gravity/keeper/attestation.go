package keeper

import (
	"encoding/hex"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/gravity/types"
)

func (k Keeper) Attest(ctx sdk.Context, relayerAddr sdk.AccAddress, claim types.ExternalClaim) (*types.Attestation, error) {
	anyClaim, err := codectypes.NewAnyWithValue(claim)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrUnknown, "msg to any")
	}

	lastObservedNonce := k.GetLastObservedEventNonce(ctx)
	// Check that the nonce of this event is exactly one higher than the last nonce stored by this relayer.
	// We check the event nonce in processAttestation as well, but checking it here gives individual eth signers a chance to retry,
	// and prevents validators from submitting two claims with the same nonce.
	// This prevents there being two attestations with the same nonce that get 2/3s of the votes
	// in the endBlocker.
	lastEventNonce := k.GetLastEventNonceByRelayer(ctx, relayerAddr)
	expectedNonce := lastEventNonce + 1

	// fist check continuity
	if claim.GetEventNonce() <= lastEventNonce {
		return nil, errorsmod.Wrapf(types.ErrNonContinuousEventNonce, "got %v, expected %v", claim.GetEventNonce(), expectedNonce)
	}
	if claim.GetEventNonce() != expectedNonce && claim.GetEventNonce() > lastObservedNonce {
		// second: if not continuous, event nonce must greater than last observed nonce.
		return nil, errorsmod.Wrapf(types.ErrNonContinuousEventNonce, "got %v, expected %v", claim.GetEventNonce(), expectedNonce)
	}

	gasMeter := ctx.GasMeter()
	ctx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())

	// Tries to get an attestation with the same eventNonce and claim as the claim that was submitted.
	att := k.GetAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash())

	// If it does not exist, create a new one.
	if att == nil {
		att = &types.Attestation{
			Observed: false,
			Height:   uint64(ctx.BlockHeight()),
			Claim:    anyClaim,
		}
	}

	// Check if relayer already voted
	for _, existingVote := range att.Votes {
		if existingVote == relayerAddr.String() {
			return nil, errorsmod.Wrap(types.ErrDuplicate, "relayer already voted on this attestation")
		}
	}

	// Add the relayer's vote to this attestation
	att.Votes = append(att.Votes, relayerAddr.String())
	k.SetAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash(), att)

	if !att.Observed && claim.GetEventNonce() == lastObservedNonce+1 {
		k.TryAttestation(ctx, att, claim)
	}

	ctx = ctx.WithGasMeter(gasMeter)
	k.SetLastEventNonceByRelayer(ctx, relayerAddr, claim.GetEventNonce())
	k.SetLastEventBlockHeightByRelayer(ctx, relayerAddr, claim.GetBlockHeight())
	return att, nil
}

// TryAttestation checks if an attestation has enough votes to be applied to the consensus state
// and has not already been marked Observed, then calls processAttestation to actually apply it to the state,
// and then marks it Observed and emits an event.
func (k Keeper) TryAttestation(ctx sdk.Context, att *types.Attestation, claim types.ExternalClaim) {
	// If the attestation has not yet been Observed, sum up the votes and see if it is ready to apply to the state.
	// This conditional stops the attestation from accidentally being applied twice.
	// Sum the current powers of all validators who have voted and see if it passes the current threshold
	totalPower := k.GetLastTotalPower(ctx)
	requiredPower := types.AttestationVotesPowerThreshold.Mul(totalPower).Quo(sdk.NewIntFromUint64(types.PowerBase))
	attestationPower := sdkmath.NewInt(0)

	for _, relayerStr := range att.Votes {
		relayerAddr, err := sdk.AccAddressFromBech32(relayerStr)
		if err != nil {
			k.Logger(ctx).Error("TryAttestation", "invalid relayer address", relayerStr, "error", err)
			continue
		}
		relayer, found := k.GetRelayer(ctx, relayerAddr)
		if !found {
			k.Logger(ctx).Error("TryAttestation", "not found relayer", relayerAddr.String(), "claimEventNonce",
				claim.GetEventNonce(), "claimType", claim.GetType(), "claimHeight", claim.GetBlockHeight())
			continue
		}
		relayerPower := relayer.GetPower()
		// Add it to the attestation power's sum
		attestationPower = attestationPower.Add(relayerPower)
		if attestationPower.LT(requiredPower) {
			continue
		}

		k.SetLastObservedEventNonce(ctx, claim.GetEventNonce())

		// in case of web3 event is long time ago, we set the last observed me block height need long enough.
		k.SetLastObservedBlockHeight(ctx, claim.GetBlockHeight(), uint64(ctx.BlockHeight()))

		att.Observed = true
		k.SetAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash(), att)

		err = k.processAttestation(ctx, claim)
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeContractEvent,
			sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
			sdk.NewAttribute(types.AttributeKeyClaimType, claim.GetType().String()),
			sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(claim.GetEventNonce())),
			sdk.NewAttribute(types.AttributeKeyClaimHash, fmt.Sprint(hex.EncodeToString(claim.ClaimHash()))),
			sdk.NewAttribute(types.AttributeKeyBlockHeight, fmt.Sprint(claim.GetBlockHeight())),
			sdk.NewAttribute(types.AttributeKeyStateSuccess, fmt.Sprint(err == nil)),
		))
		// execute the timeout logic
		//k.cleanupTimedOutBatches(ctx)
		k.PruneAttestations(ctx)
		break
	}
}

// processAttestation actually applies the attestation to the consensus state
func (k Keeper) processAttestation(ctx sdk.Context, claim types.ExternalClaim) error {
	// then execute in a new Tx so that we can store state on failure
	xCtx, commit := ctx.CacheContext()
	if err := k.AttestationHandler(xCtx, claim); err != nil {
		// execute with a transient storage
		// If the attestation fails, something has gone wrong and we can't recover it. Log and move on
		// The attestation will still be marked "Observed", and validators can still be slashed for not
		// having voted for it.
		k.Logger(ctx).Error("attestation failed", "cause", err.Error(), "claim type", claim.GetType(),
			"id", hex.EncodeToString(types.GetAttestationKey(claim.GetEventNonce(), claim.ClaimHash())),
			"nonce", fmt.Sprint(claim.GetEventNonce()),
		)
		return err
	}
	commit() // persist transient storage
	return nil
}

// SetAttestation sets the attestation in the store
func (k Keeper) SetAttestation(ctx sdk.Context, eventNonce uint64, claimHash []byte, att *types.Attestation) {
	store := ctx.KVStore(k.storeKey)
	aKey := types.GetAttestationKey(eventNonce, claimHash)
	store.Set(aKey, k.cdc.MustMarshal(att))
}

// GetAttestation return an attestation given a nonce
func (k Keeper) GetAttestation(ctx sdk.Context, eventNonce uint64, claimHash []byte) *types.Attestation {
	store := ctx.KVStore(k.storeKey)
	aKey := types.GetAttestationKey(eventNonce, claimHash)
	bz := store.Get(aKey)
	if len(bz) == 0 {
		return nil
	}
	var att types.Attestation
	k.cdc.MustUnmarshal(bz, &att)
	return &att
}

// DeleteAttestation deletes an attestation given an event nonce and claim
func (k Keeper) DeleteAttestation(ctx sdk.Context, claim types.ExternalClaim) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetAttestationKey(claim.GetEventNonce(), claim.ClaimHash()))
}

// IterateAttestationAndClaim iterates through all attestations
func (k Keeper) IterateAttestationAndClaim(ctx sdk.Context, cb func(*types.Attestation, types.ExternalClaim) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.RelayerAttestationKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		att := new(types.Attestation)
		k.cdc.MustUnmarshal(iter.Value(), att)
		claim, err := types.UnpackAttestationClaim(k.cdc, att)
		if err != nil {
			k.Logger(ctx).Error("failed to unpack attestation claim", "error", err)
			continue
		}
		// cb returns true to stop early
		if cb(att, claim) {
			return
		}
	}
}

// IterateAttestations iterates through all attestations
func (k Keeper) IterateAttestationsByNonce(ctx sdk.Context, nonce uint64, cb func(*types.Attestation) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetAttestationKeyByNonce(nonce))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		att := new(types.Attestation)
		k.cdc.MustUnmarshal(iter.Value(), att)
		// cb returns true to stop early
		if cb(att) {
			return
		}
	}
}

// IterateAttestations iterates through all attestations
func (k Keeper) IterateAttestations(ctx sdk.Context, cb func(*types.Attestation) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.RelayerAttestationKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		att := new(types.Attestation)
		k.cdc.MustUnmarshal(iter.Value(), att)
		// cb returns true to stop early
		if cb(att) {
			return
		}
	}
}

// GetLastObservedEventNonce returns the latest observed event nonce
func (k Keeper) GetLastObservedEventNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastObservedEventNonceKey)
	if len(bytes) == 0 {
		return 0
	}
	return sdk.BigEndianToUint64(bytes)
}

// SetLastObservedEventNonce sets the latest observed event nonce
func (k Keeper) SetLastObservedEventNonce(ctx sdk.Context, eventNonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastObservedEventNonceKey, sdk.Uint64ToBigEndian(eventNonce))
}

// GetLastObservedBlockHeight height gets the block height to of the last observed attestation from
// the store
func (k Keeper) GetLastObservedBlockHeight(ctx sdk.Context) types.LastObservedBlockHeight {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastObservedBlockHeightKey)
	if len(bytes) == 0 {
		return types.LastObservedBlockHeight{
			ExternalBlockHeight: 0,
			BlockHeight:         0,
		}
	}
	height := types.LastObservedBlockHeight{}
	k.cdc.MustUnmarshal(bytes, &height)
	return height
}

// SetLastObservedBlockHeight sets the block height in the store.
func (k Keeper) SetLastObservedBlockHeight(ctx sdk.Context, externalBlockHeight, blockHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	height := types.LastObservedBlockHeight{
		ExternalBlockHeight: externalBlockHeight,
		BlockHeight:         blockHeight,
	}
	store.Set(types.LastObservedBlockHeightKey, k.cdc.MustMarshal(&height))
}

// GetLastEventNonceByGravity returns the latest event nonce for a given relayer
func (k Keeper) GetLastEventNonceByRelayer(ctx sdk.Context, relayerAddr sdk.AccAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.GetLastEventNonceByRelayerKey(relayerAddr))

	if len(bytes) == 0 {
		// in the case that we have no existing value this is the first
		// time a relayerAddr is submitting a claim. Since we don't want to force
		// them to replay the entire history of all events ever we can't start
		// at zero
		lastEventNonce := k.GetLastObservedEventNonce(ctx)
		if lastEventNonce >= 1 {
			return lastEventNonce - 1
		} else {
			return 0
		}
	}
	return sdk.BigEndianToUint64(bytes)
}

// DelLastEventNonceByRelayer delete the latest event nonce for a given relayer
func (k Keeper) DelLastEventNonceByRelayer(ctx sdk.Context, relayerAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLastEventNonceByRelayerKey(relayerAddr)
	if !store.Has(key) {
		return
	}
	store.Delete(key)
}

// SetLastEventNonceByRelayer sets the latest event nonce for a give relayer
func (k Keeper) SetLastEventNonceByRelayer(ctx sdk.Context, relay sdk.AccAddress, eventNonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetLastEventNonceByRelayerKey(relay), sdk.Uint64ToBigEndian(eventNonce))
}

func (k Keeper) SetLastEventBlockHeightByRelayer(ctx sdk.Context, relay sdk.AccAddress, blockHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetLastEventBlockHeightByRelayerKey(relay), sdk.Uint64ToBigEndian(blockHeight))
}

func (k Keeper) GetLastEventBlockHeightByRelayer(ctx sdk.Context, relay sdk.AccAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLastEventBlockHeightByRelayerKey(relay)
	if !store.Has(key) {
		return 0
	}
	data := store.Get(key)
	return sdk.BigEndianToUint64(data)
}
