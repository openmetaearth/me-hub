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
		}
	}

	// Process attestation
	err = k.processAttestation(ctx, claim)
	if err != nil {
		return nil, err
	}

	// Only after successful processing, update the attestation and last observed nonce
	att.Observed = true
	k.SetAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash(), att)
	k.SetLastObservedEventNonce(ctx, claim.GetEventNonce())

	return att, nil
}