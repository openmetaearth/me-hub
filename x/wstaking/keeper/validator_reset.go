package keeper

import (
	"sort"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (k Keeper) validateOperatorReplacement(ctx sdk.Context, newOwner sdk.AccAddress, validator stakingtypes.Validator) error {
	if err := sdk.VerifyAddressFormat(newOwner); err != nil {
		return err
	}
	newAddr := sdk.ValAddress(newOwner)
	if _, found := k.GetValidator(ctx, newAddr); found {
		return types.ErrValidatorExist
	}
	if _, module := k.authKeeper.GetAccount(ctx, newOwner).(authtypes.ModuleAccountI); module {
		return types.ErrValidatorAddress
	}
	if k.bankKeeper.Extend().BlockedAddr(newOwner) {
		return types.ErrValidatorAddress
	}
	store := ctx.KVStore(k.storeKey)
	if store.Has(stakingtypes.GetLastValidatorPowerKey(newAddr)) {
		return sdkerrors.Wrap(types.ErrValidatorExist, "new operator has last-validator power")
	}
	iterator := sdk.KVStorePrefixIterator(store, stakingtypes.GetUBDsByValIndexKey(newAddr))
	occupied := iterator.Valid()
	iterator.Close()
	if occupied {
		return sdkerrors.Wrap(types.ErrValidatorExist, "new operator has unbonding delegations")
	}
	// A pending consensus-key rotation still refers to the old operator.
	rotation, err := k.GetReplaceConsensusPubKeyInfo(ctx)
	if err != nil {
		return err
	}
	if rotation != nil && strings.EqualFold(rotation.OperatorAddress, validator.OperatorAddress) {
		return sdkerrors.Wrap(types.ErrExistingReplaceValidator, "finish consensus public key replacement before changing the operator")
	}
	return nil
}

// resetValidator follows the remove/create lifecycle, keeping the consensus key
// and completing all staking address changes atomically in the current block.
func (k Keeper) resetValidator(ctx sdk.Context, newOwner sdk.AccAddress, validator stakingtypes.Validator) error {
	if err := k.validateOperatorReplacement(ctx, newOwner, validator); err != nil {
		return err
	}
	oldAddr, newAddr := validator.GetOperator(), sdk.ValAddress(newOwner)
	logger := k.Logger(ctx).With("height", ctx.BlockHeight(), "region_id", validator.Description.RegionID,
		"old_operator", oldAddr.String(), "new_operator", newAddr.String())
	started := time.Now()
	logger.Info("operator replacement started")
	ctx, write := ctx.CacheContext()
	// UpdateValidator may also edit RegionID; select delegators using the stored
	// source region, not the validator's proposed destination region.
	oldRegionID := k.MustGetValidator(ctx, oldAddr).Description.RegionID
	if err := k.migrateValidatorPositions(ctx, oldRegionID, oldAddr, newAddr); err != nil {
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
	k.DeleteValidatorByPowerIndex(ctx, validator)
	if validator.IsUnbonding() {
		k.DeleteValidatorQueue(ctx, validator)
	}
	// Keep last-set tracking, including a jailed validator awaiting a zero update.
	if power := store.Get(stakingtypes.GetLastValidatorPowerKey(oldAddr)); power != nil {
		store.Set(stakingtypes.GetLastValidatorPowerKey(newAddr), power)
		k.DeleteLastValidatorPower(ctx, oldAddr)
	}

	validator.OperatorAddress = newAddr.String()
	validator.OwnerAddress = newOwner.String()
	k.SetValidator(ctx, validator)
	if err := k.SetValidatorByConsAddr(ctx, validator); err != nil {
		return err
	}
	k.SetValidatorByPowerIndex(ctx, validator)
	if validator.IsUnbonding() {
		k.InsertUnbondingValidatorQueue(ctx, validator)
	}
	for _, id := range validator.UnbondingIds {
		k.SetValidatorByUnbondingID(ctx, validator, id)
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

func (k Keeper) migrateValidatorPositions(ctx sdk.Context, regionID string, oldAddr, newAddr sdk.ValAddress) error {
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
		k.IterateAllDelegations(ctx, func(del stakingtypes.Delegation) bool {
			if del.ValidatorAddress == oldOperator {
				delegations = append(delegations, del)
			}
			return false
		})
		for _, del := range delegations {
			del.ValidatorAddress = newOperator
			k.SetDelegation(ctx, del)
		}
	} else {
		// Other regions require KYC/DID. Reuse ChangeDelegationValidator's
		// indexed path now rather than queueing it for the following block.
		k.changeRegionDelegationValidator(ctx, regionID, oldOperator, newOperator)
	}
	k.migrateUnbondingDelegations(ctx, k.GetUnbondingDelegationsFromValidator(ctx, oldAddr), oldOperator, newOperator)
	return nil
}

type unbondingQueueUpdate struct {
	completionTime time.Time
	delegators     map[string]struct{}
}

// Rewrite each affected timeslice once, not once per unbonding entry.
func (k Keeper) migrateUnbondingDelegations(ctx sdk.Context, unbondings []stakingtypes.UnbondingDelegation, oldOperator, newOperator string) {
	updates := make(map[string]*unbondingQueueUpdate)
	for _, ubd := range unbondings {
		k.RemoveUnbondingDelegation(ctx, ubd)
		ubd.ValidatorAddress = newOperator
		k.SetUnbondingDelegation(ctx, ubd)
		for _, entry := range ubd.Entries {
			if entry.UnbondingId != 0 {
				k.SetUnbondingDelegationByUnbondingID(ctx, ubd, entry.UnbondingId)
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
		pairs := k.GetUBDQueueTimeSlice(ctx, update.completionTime)
		changed := false
		for i := range pairs {
			if _, ok := update.delegators[pairs[i].DelegatorAddress]; ok && pairs[i].ValidatorAddress == oldOperator {
				pairs[i].ValidatorAddress = newOperator
				changed = true
			}
		}
		if changed {
			k.SetUBDQueueTimeSlice(ctx, update.completionTime, pairs)
		}
	}
}
