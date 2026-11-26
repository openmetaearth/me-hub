package keeper

import (
	"sort"
	"strings"
	"time"

	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (k *Keeper) validateOperatorReplacement(ctx sdk.Context, newOwner sdk.AccAddress, validator stakingtypes.Validator) error {
	if err := sdk.VerifyAddressFormat(newOwner); err != nil {
		return err
	}
	newAddr := sdk.ValAddress(newOwner)
	if _, err := k.GetValidator(ctx, newAddr); err == nil {
		return types.ErrValidatorExist
	} else if !errorsmod.IsOf(err, stakingtypes.ErrNoValidatorFound) {
		return err
	}
	if _, module := k.authKeeper.GetAccount(ctx, newOwner).(authtypes.ModuleAccountI); module {
		return types.ErrValidatorAddress
	}
	if k.bankKeeper.Extend().BlockedAddr(newOwner) {
		return types.ErrValidatorAddress
	}
	store := ctx.KVStore(k.storeKey)
	if store.Has(stakingtypes.GetLastValidatorPowerKey(newAddr)) {
		return errorsmod.Wrap(types.ErrValidatorExist, "new operator has last-validator power")
	}
	iterator := storetypes.KVStorePrefixIterator(store, stakingtypes.GetUBDsByValIndexKey(newAddr))
	occupied := iterator.Valid()
	iterator.Close()
	if occupied {
		return errorsmod.Wrap(types.ErrValidatorExist, "new operator has unbonding delegations")
	}
	// A pending consensus-key rotation still refers to the old operator.
	rotation, err := k.GetReplaceConsensusPubKeyInfo(ctx)
	if err != nil {
		return err
	}
	if rotation != nil && strings.EqualFold(rotation.OperatorAddress, validator.OperatorAddress) {
		return errorsmod.Wrap(types.ErrExistingReplaceValidator, "finish consensus public key replacement before changing the operator")
	}
	return nil
}

// resetValidator follows the remove/create lifecycle, keeping the consensus key
// and completing all staking address changes atomically in the current block.
func (k *Keeper) resetValidator(ctx sdk.Context, newOwner sdk.AccAddress, validator stakingtypes.Validator) error {
	if err := k.validateOperatorReplacement(ctx, newOwner, validator); err != nil {
		return err
	}
	oldAddr, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
	if err != nil {
		return err
	}
	newAddr := sdk.ValAddress(newOwner)
	logger := k.Logger(ctx).With("height", ctx.BlockHeight(), "region_id", validator.Description.RegionID,
		"old_operator", oldAddr.String(), "new_operator", newAddr.String())
	started := time.Now()
	logger.Info("operator replacement started")
	ctx, write := ctx.CacheContext()
	// UpdateValidator may also edit RegionID; select delegators using the stored
	// source region, not the validator's proposed destination region.
	storedValidator, err := k.GetValidator(ctx, oldAddr)
	if err != nil {
		return err
	}
	if err := k.migrateValidatorPositions(ctx, storedValidator.Description.RegionID, oldAddr, newAddr); err != nil {
		return err
	}

	// Same cleanup as RemoveValidator, but propagate hook errors instead of only
	// logging them: a failed distribution settlement must roll back the proposal.
	consAddr, err := validator.GetConsAddr()
	if err != nil {
		return err
	}
	if err := k.Hooks().AfterValidatorRemoved(ctx, consAddr, oldAddr); err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	store.Delete(stakingtypes.GetValidatorKey(oldAddr))
	store.Delete(stakingtypes.GetValidatorByConsAddrKey(consAddr))
	if err := k.DeleteValidatorByPowerIndex(ctx, validator); err != nil {
		return err
	}
	if validator.IsUnbonding() {
		if err := k.DeleteValidatorQueue(ctx, validator); err != nil {
			return err
		}
	}
	// Keep last-set tracking, including a jailed validator awaiting a zero update.
	if power := store.Get(stakingtypes.GetLastValidatorPowerKey(oldAddr)); power != nil {
		store.Set(stakingtypes.GetLastValidatorPowerKey(newAddr), power)
		if err := k.DeleteLastValidatorPower(ctx, oldAddr); err != nil {
			return err
		}
	}

	validator.OperatorAddress = newAddr.String()
	validator.OwnerAddress = newOwner.String()
	if err := k.SetValidator(ctx, validator); err != nil {
		return err
	}
	if err := k.SetValidatorByConsAddr(ctx, validator); err != nil {
		return err
	}
	if err := k.SetValidatorByPowerIndex(ctx, validator); err != nil {
		return err
	}
	if validator.IsUnbonding() {
		if err := k.InsertUnbondingValidatorQueue(ctx, validator); err != nil {
			return err
		}
	}
	for _, id := range validator.UnbondingIds {
		if err := k.SetValidatorByUnbondingID(ctx, validator, id); err != nil {
			return err
		}
	}
	shares := validator.Tokens
	if region, found := k.GetRegion(ctx, validator.Description.RegionID); found {
		shares = region.RegionShare
	}
	k.BondRegion(ctx, validator, shares, true)
	if err := k.Hooks().AfterValidatorCreated(ctx, newAddr); err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeResetValidator,
		sdk.NewAttribute(types.AttributeKeyValidator, oldAddr.String()),
		sdk.NewAttribute(types.AttributeKeyNewValidator, newAddr.String()),
		sdk.NewAttribute(types.AttributeKeyNewOwnerAddress, validator.OwnerAddress),
	))
	write()
	logger.Info("operator replacement completed", "elapsed", time.Since(started))
	return nil
}

func (k *Keeper) migrateValidatorPositions(ctx sdk.Context, regionID string, oldAddr, newAddr sdk.ValAddress) error {
	oldOperator, newOperator := oldAddr.String(), newAddr.String()
	stakes, err := k.GetStakesByValidator(ctx, oldAddr)
	if err != nil {
		return err
	}
	for _, stake := range stakes {
		if err := k.RemoveStake(ctx, *stake); err != nil {
			return err
		}
		stake.ValidatorAddress = newOperator
		k.SetStake(ctx, *stake)
	}
	if regionID == types.ExperienceRegionId {
		// Experience delegators have no KYC/DID index. Collect before writing.
		var delegations []stakingtypes.Delegation
		if err := k.IterateAllDelegations(ctx, func(del stakingtypes.Delegation) bool {
			if del.ValidatorAddress == oldOperator {
				delegations = append(delegations, del)
			}
			return false
		}); err != nil {
			return err
		}
		for _, del := range delegations {
			del.ValidatorAddress = newOperator
			k.SetDelegation(ctx, del)
		}
	} else {
		// Other regions require KYC/DID. Reuse ChangeDelegationValidator's
		// indexed path now rather than queueing it for the following block.
		k.changeRegionDelegationValidator(ctx, regionID, oldOperator, newOperator)
	}
	unbondings, err := k.GetUnbondingDelegationsFromValidator(ctx, oldAddr)
	if err != nil {
		return err
	}
	return k.migrateUnbondingDelegations(ctx, unbondings, oldOperator, newOperator)
}

type unbondingQueueUpdate struct {
	completionTime time.Time
	delegators     map[string]struct{}
}

// Rewrite each affected timeslice once, not once per unbonding entry.
func (k *Keeper) migrateUnbondingDelegations(ctx sdk.Context, unbondings []stakingtypes.UnbondingDelegation, oldOperator, newOperator string) error {
	updates := make(map[string]*unbondingQueueUpdate)
	for _, ubd := range unbondings {
		if err := k.RemoveUnbondingDelegation(ctx, ubd); err != nil {
			return err
		}
		ubd.ValidatorAddress = newOperator
		if err := k.SetUnbondingDelegation(ctx, ubd); err != nil {
			return err
		}
		for _, entry := range ubd.Entries {
			if entry.UnbondingId != 0 {
				if err := k.SetUnbondingDelegationByUnbondingID(ctx, ubd, entry.UnbondingId); err != nil {
					return err
				}
			}
			key := string(stakingtypes.GetUnbondingDelegationTimeKey(entry.CompletionTime))
			update := updates[key]
			if update == nil {
				update = &unbondingQueueUpdate{completionTime: entry.CompletionTime, delegators: make(map[string]struct{})}
				updates[key] = update
			}
			update.delegators[ubd.DelegatorAddress] = struct{}{}
		}
	}
	keys := make([]string, 0, len(updates))
	for key := range updates {
		keys = append(keys, key)
	}
	sort.Strings(keys) // State-write order must not depend on Go map iteration.
	for _, key := range keys {
		update := updates[key]
		pairs, err := k.GetUBDQueueTimeSlice(ctx, update.completionTime)
		if err != nil {
			return err
		}
		changed := false
		for i := range pairs {
			if _, ok := update.delegators[pairs[i].DelegatorAddress]; ok && pairs[i].ValidatorAddress == oldOperator {
				pairs[i].ValidatorAddress = newOperator
				changed = true
			}
		}
		if changed {
			if err := k.SetUBDQueueTimeSlice(ctx, update.completionTime, pairs); err != nil {
				return err
			}
		}
	}
	return nil
}
