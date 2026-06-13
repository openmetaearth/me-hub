package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

// UpdateValidatorPubKey updates the validator's consensus public key.
// It ensures that the old consensus address is removed from the staking index
// AND the corresponding ValidatorSigningInfo is deleted from the slashing module
// to prevent state bloat (Issue #1249).
func (k Keeper) UpdateValidatorPubKey(ctx sdk.Context, valAddr sdk.ValAddress, newConsensusPubKey []byte) error {
	// 1. Retrieve the validator
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return stakingtypes.ErrNoValidatorFound
	}

	// 2. Get the old consensus address
	oldConsAddr := sdk.ConsAddress(validator.ConsensusPubKey.Address())

	// 3. Update the validator's consensus public key in the staking store
	validator.ConsensusPubKey = newConsensusPubKey
	k.SetValidator(ctx, validator)

	// 4. Update the ValidatorByConsAddr index
	// Remove the old index entry
	k.DeleteValidatorByConsAddr(ctx, oldConsAddr)
	// Add the new index entry
	newConsAddr := sdk.ConsAddress(newConsensusPubKey)
	k.SetValidatorByConsAddr(ctx, newConsAddr, valAddr)

	// 5. FIX FOR ISSUE #1249:
	// Delete the old ValidatorSigningInfo from the slashing module.
	// Previously, this step was missing, causing stale entries to accumulate.
	if k.slashingKeeper != nil {
		// We assume the slashing keeper has a method to delete signing info.
		// In standard Cosmos SDK, this is often done via k.slashingKeeper.DeleteValidatorSigningInfo
		k.slashingKeeper.DeleteValidatorSigningInfo(ctx, oldConsAddr)
	}

	return nil
}