package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

// validateOperatorReplacement checks the destination and overlapping key rotation.
func (k Keeper) validateOperatorReplacement(ctx sdk.Context, newOwner sdk.AccAddress, validator stakingtypes.Validator) error {
	if err := sdk.VerifyAddressFormat(newOwner); err != nil {
		return err
	}
	oldAddr, newAddr := validator.GetOperator(), sdk.ValAddress(newOwner)
	if _, found := k.GetValidator(ctx, newAddr); found {
		return types.ErrValidatorExist
	}
	if _, module := k.authKeeper.GetAccount(ctx, newOwner).(authtypes.ModuleAccountI); module {
		return types.ErrValidatorAddress
	}
	if k.bankKeeper.Extend().BlockedAddr(newOwner) {
		return types.ErrValidatorAddress
	}
	// A key rotation retains the old consensus-address alias for two blocks. Do not
	// rename its operator until both the queued update and that alias are removed.
	rotation, err := k.GetReplaceConsensusPubKeyInfo(ctx)
	if err != nil {
		return err
	}
	if rotation != nil && strings.EqualFold(rotation.OperatorAddress, oldAddr.String()) {
		return sdkerrors.Wrap(types.ErrExistingReplaceValidator, "finish consensus public key replacement before changing the operator")
	}
	if k.distributionKeeper == nil {
		return sdkerrors.Wrap(types.ErrResetValidator, "distribution keeper is not configured")
	}
	return nil
}

// resetValidator renames the same validator without invoking removal/creation hooks.
func (k Keeper) resetValidator(ctx sdk.Context, newOwner sdk.AccAddress, validator stakingtypes.Validator) error {
	if err := k.validateOperatorReplacement(ctx, newOwner, validator); err != nil {
		return err
	}
	oldAddr, newAddr := validator.GetOperator(), sdk.ValAddress(newOwner)
	cacheCtx, write := ctx.CacheContext()
	if err := k.migrateValidatorPositions(cacheCtx, oldAddr, newAddr); err != nil {
		return err
	}
	if err := k.distributionKeeper.MigrateValidatorOperator(cacheCtx, oldAddr, newAddr); err != nil {
		return err
	}
	store := cacheCtx.KVStore(k.storeKey)
	if store.Has(stakingtypes.GetLastValidatorPowerKey(newAddr)) {
		return sdkerrors.Wrap(types.ErrValidatorExist, "new operator has last-validator power")
	}
	// Remove only the old keys, NOT RemoveValidator: its hooks pay commission,
	// discard reward history, and remove the still-valid slashing pubkey mapping.
	k.DeleteValidatorByPowerIndex(cacheCtx, validator)
	store.Delete(stakingtypes.GetValidatorKey(oldAddr))
	if validator.IsUnbonding() {
		k.DeleteValidatorQueue(cacheCtx, validator)
	}
	if power := store.Get(stakingtypes.GetLastValidatorPowerKey(oldAddr)); power != nil {
		store.Set(stakingtypes.GetLastValidatorPowerKey(newAddr), power)
		k.DeleteLastValidatorPower(cacheCtx, oldAddr)
	}
	validator.OperatorAddress = newAddr.String()
	validator.OwnerAddress = newOwner.String()
	k.SetValidator(cacheCtx, validator)
	if err := k.SetValidatorByConsAddr(cacheCtx, validator); err != nil {
		return err
	}
	k.SetValidatorByPowerIndex(cacheCtx, validator)
	if validator.IsUnbonding() {
		k.InsertUnbondingValidatorQueue(cacheCtx, validator)
	}
	for _, id := range validator.UnbondingIds {
		k.SetValidatorByUnbondingID(cacheCtx, validator, id)
	}

	// Changing ownership must not change the region's shares or reward accounting.
	shares := validator.Tokens
	if region, found := k.GetRegion(cacheCtx, validator.Description.RegionID); found {
		shares = region.RegionShare
	}
	k.BondRegion(cacheCtx, validator, shares, true)
	cacheCtx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeResetValidator,
		sdk.NewAttribute(types.AttributeKeyValidator, oldAddr.String()),
		sdk.NewAttribute(types.AttributeKeyNewValidator, newAddr.String()),
		sdk.NewAttribute(types.AttributeKeyNewOwnerAddress, validator.OwnerAddress),
	))
	write()
	return nil
}

// migrateValidatorPositions preserves amounts, owners, heights, completion times,
// and unbonding IDs. Collect before modifying stores; reject reuse of an address
// with outstanding positions rather than overwriting or merging those positions.
func (k Keeper) migrateValidatorPositions(ctx sdk.Context, oldAddr, newAddr sdk.ValAddress) error {
	oldOperator, newOperator := oldAddr.String(), newAddr.String()
	positions, err := k.getValidatorPositions(ctx, oldAddr)
	if err != nil {
		return err
	}
	target, err := k.getValidatorPositions(ctx, newAddr)
	if err != nil {
		return err
	}
	if !target.empty() {
		return sdkerrors.Wrap(types.ErrValidatorExist, "new operator has outstanding staking positions")
	}

	for _, stake := range positions.stakes {
		if err := k.RemoveStake(ctx, *stake); err != nil {
			return err
		}
		stake.ValidatorAddress = newOperator
		k.SetStake(ctx, *stake)
	}
	// Do not defer this to the KYC-only ChangeDelegationValidator queue: experience
	// region delegators may have no KYC, and callers need consistent state now.
	for _, del := range positions.delegations {
		del.ValidatorAddress = newOperator
		k.SetDelegation(ctx, del)
	}
	k.migrateUnbondingStakes(ctx, positions.unbondingStakes, oldOperator, newOperator)
	k.migrateUnbondingDelegations(ctx, positions.unbondings, oldOperator, newOperator)
	k.migrateRedelegations(ctx, positions.redelegations, oldOperator, newOperator)
	return nil
}

type validatorPositions struct {
	stakes          []*types.Stake
	delegations     []stakingtypes.Delegation
	unbondingStakes []types.UnbondingStake
	unbondings      []stakingtypes.UnbondingDelegation
	redelegations   []stakingtypes.Redelegation
}

func (p validatorPositions) empty() bool {
	return len(p.stakes)+len(p.delegations)+len(p.unbondingStakes)+len(p.unbondings)+len(p.redelegations) == 0
}

func (k Keeper) getValidatorPositions(ctx sdk.Context, addr sdk.ValAddress) (validatorPositions, error) {
	var positions validatorPositions
	stakes, err := k.GetStakesByValidator(ctx, addr)
	if err != nil {
		return positions, err
	}
	positions.stakes = stakes
	positions.unbondings = k.GetUnbondingDelegationsFromValidator(ctx, addr)
	operator := addr.String()
	k.IterateAllDelegations(ctx, func(del stakingtypes.Delegation) bool {
		if del.ValidatorAddress == operator {
			positions.delegations = append(positions.delegations, del)
		}
		return false
	})
	k.IterateUnbondingStakes(ctx, func(ubs types.UnbondingStake) bool {
		if ubs.ValidatorAddress == operator {
			positions.unbondingStakes = append(positions.unbondingStakes, ubs)
		}
		return false
	})
	k.IterateRedelegations(ctx, func(_ int64, red stakingtypes.Redelegation) bool {
		if red.ValidatorSrcAddress == operator || red.ValidatorDstAddress == operator {
			positions.redelegations = append(positions.redelegations, red)
		}
		return false
	})
	return positions, nil
}

func (k Keeper) migrateUnbondingStakes(ctx sdk.Context, stakes []types.UnbondingStake, oldOperator, newOperator string) {
	for _, ubs := range stakes {
		k.RemoveUnbondingStake(ctx, ubs)
		ubs.ValidatorAddress = newOperator
		k.SetUnbondingStake(ctx, ubs)
		for _, entry := range ubs.Entries {
			pairs := k.GetUBSQueueTimeSlice(ctx, entry.CompletionTime)
			changed := false
			for i := range pairs {
				if pairs[i].StakerAddress == ubs.StakerAddress && pairs[i].ValidatorAddress == oldOperator {
					pairs[i].ValidatorAddress = newOperator
					changed = true
				}
			}
			if changed {
				k.SetUBSQueueTimeSlice(ctx, entry.CompletionTime, pairs)
			}
		}
	}
}

func (k Keeper) migrateUnbondingDelegations(ctx sdk.Context, unbondings []stakingtypes.UnbondingDelegation, oldOperator, newOperator string) {
	for _, ubd := range unbondings {
		k.RemoveUnbondingDelegation(ctx, ubd)
		ubd.ValidatorAddress = newOperator
		k.SetUnbondingDelegation(ctx, ubd)
		for _, entry := range ubd.Entries {
			if entry.UnbondingId != 0 {
				k.SetUnbondingDelegationByUnbondingID(ctx, ubd, entry.UnbondingId)
			}
			pairs := k.GetUBDQueueTimeSlice(ctx, entry.CompletionTime)
			changed := false
			for i := range pairs {
				if pairs[i].DelegatorAddress == ubd.DelegatorAddress && pairs[i].ValidatorAddress == oldOperator {
					pairs[i].ValidatorAddress = newOperator
					changed = true
				}
			}
			if changed {
				k.SetUBDQueueTimeSlice(ctx, entry.CompletionTime, pairs)
			}
		}
	}
}

func (k Keeper) migrateRedelegations(ctx sdk.Context, redelegations []stakingtypes.Redelegation, oldOperator, newOperator string) {
	for _, red := range redelegations {
		oldSrc, oldDst := red.ValidatorSrcAddress, red.ValidatorDstAddress
		k.RemoveRedelegation(ctx, red)
		if oldSrc == oldOperator {
			red.ValidatorSrcAddress = newOperator
		}
		if oldDst == oldOperator {
			red.ValidatorDstAddress = newOperator
		}
		k.SetRedelegation(ctx, red)
		for _, entry := range red.Entries {
			if entry.UnbondingId != 0 {
				k.SetRedelegationByUnbondingID(ctx, red, entry.UnbondingId)
			}
			triplets := k.GetRedelegationQueueTimeSlice(ctx, entry.CompletionTime)
			changed := false
			for i := range triplets {
				if triplets[i].DelegatorAddress == red.DelegatorAddress && triplets[i].ValidatorSrcAddress == oldSrc && triplets[i].ValidatorDstAddress == oldDst {
					triplets[i].ValidatorSrcAddress = red.ValidatorSrcAddress
					triplets[i].ValidatorDstAddress = red.ValidatorDstAddress
					changed = true
				}
			}
			if changed {
				k.SetRedelegationQueueTimeSlice(ctx, entry.CompletionTime, triplets)
			}
		}
	}
}
