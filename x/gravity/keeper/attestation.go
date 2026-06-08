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
	lastEventNonce := k.GetLastEventNonceByRelayer(ctx, relayerAddr)
	expectedNonce := lastEventNonce + 1

	// Check continuity
	if claim.GetEventNonce() <= lastEventNonce {
		return nil, errorsmod.Wrapf(types.ErrNonContinuousEventNonce, "got %v, expected %v", claim.GetEventNonce(), expectedNonce)
	}
	if claim.GetEventNonce() != expectedNonce && claim.GetEventNonce() <= lastObservedNonce {
		return nil, errorsmod.Wrapf(types.ErrNonContinuousEventNonce, "got %v, expected %v", claim.GetEventNonce(), expectedNonce)
	}
	if claim.GetEventNonce() > lastObservedNonce {
		return nil, errorsmod.Wrapf(types.ErrNonContinuousEventNonce, "got %v, expected %v", claim.GetEventNonce(), expectedNonce)
	}

	gasMeter := ctx.GasMeter()
	ctx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())

	// Tries to get an attestation with the same eventNonce and claim as the claim that was submitted.
	att := k.GetAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash())

	// If it does not exist, create a new one.
	if att == nil {
		att = &types.Attestation{
			Obs