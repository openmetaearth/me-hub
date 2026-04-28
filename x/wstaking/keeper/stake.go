package keeper

import (
	"time"

	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wstaking/types"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Stake performs a stake, set/update everything necessary within the store.
// tokenSrc indicates the bond status of the incoming funds.
func (k Keeper) Stake(ctx sdk.Context, staker sdk.AccAddress, bondAmt math.Int,
	tokenSrc stakingtypes.BondStatus, validator stakingtypes.Validator, subtractAccount bool, tag string) (newShares sdk.Dec, err error) {
	// In some situations, the exchange rate becomes invalid, e.g. if
	// Validator loses all tokens due to slashing. In this case,
	// make all future stakes invalid.
	if validator.InvalidExRate() {
		return sdk.ZeroDec(), types.ErrStakerShareExRateInvalid
	}

	// Get or create the stake object
	stake, found := k.GetStake(ctx, staker, validator.GetOperator())
	if !found {
		stake = types.NewStake(staker, validator.GetOperator(), sdk.ZeroDec())
	}

	// if subtractAccount is true then we are
	// performing a stake and not a restake, thus the source tokens are
	// all non bonded
	if subtractAccount {
		if tokenSrc == stakingtypes.Bonded {
			panic("stake token source cannot be bonded")
		}

		var recipientModule string

		switch {
		case validator.IsBonded():
			recipientModule = types.BondedStakePoolName
		case validator.IsUnbonding(), validator.IsUnbonded():
			recipientModule = types.NotBondedStakePoolName
		default:
			panic("invalid validator status")
		}

		coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), bondAmt))
		if err := k.bankKeeper.Extend().SendCoinsFromModuleToModuleWithTag(ctx, types.StakePoolName, recipientModule, coins, tag); err != nil {
			return sdk.Dec{}, err
		}
	} else {
		// potentially transfer tokens between pools, if
		switch {
		case tokenSrc == stakingtypes.Bonded && validator.IsBonded():
			// do nothing
		case (tokenSrc == stakingtypes.Unbonded || tokenSrc == stakingtypes.Unbonding) && !validator.IsBonded():
			// do nothing
		case (tokenSrc == stakingtypes.Unbonded || tokenSrc == stakingtypes.Unbonding) && validator.IsBonded():
			// transfer pools
			k.NotBondedStakeTokensToBonded(ctx, bondAmt)
		case tokenSrc == stakingtypes.Bonded && !validator.IsBonded():
			// transfer pools
			k.BondedStakeTokensToNotBonded(ctx, bondAmt, validator.Description.RegionID)
		default:
			panic("unknown token source bond status")
		}
	}

	_, newShares = k.AddValidatorTokensAndShares(ctx, validator, bondAmt)

	// Update stake
	stake.Shares = stake.Shares.Add(newShares)
	k.SetStake(ctx, stake)
	k.BondRegion(ctx, validator, stake.Shares.TruncateInt(), true)
	return newShares, nil
}

// GetStake returns a specific stake.
func (k Keeper) GetStake(ctx sdk.Context, stakerAddr sdk.AccAddress, valAddr sdk.ValAddress) (stake types.Stake, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetStakeKey(stakerAddr, valAddr)

	value := store.Get(key)
	if value == nil {
		return stake, false
	}
	k.cdc.MustUnmarshal(value, &stake)
	return stake, true
}

// SetStake sets a stake.
func (k Keeper) SetStake(ctx sdk.Context, stake types.Stake) {
	stakerAddress := sdk.MustAccAddressFromBech32(stake.StakerAddress)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetStakeKey(stakerAddress, stake.GetValidatorAddr()), k.cdc.MustMarshal(&stake))
}

// IterateAllDelegations iterates through all of the delegations.
func (k Keeper) IterateAllStakes(ctx sdk.Context, cb func(stake types.Stake) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.StakeKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var stake types.Stake
		k.cdc.MustUnmarshal(iterator.Value(), &stake)
		if cb(stake) {
			break
		}
	}
}

func (k Keeper) GetAllStakes(ctx sdk.Context) (stakes []types.Stake) {
	k.IterateAllStakes(ctx, func(stake types.Stake) bool {
		stakes = append(stakes, stake)
		return false
	})
	return stakes
}

func (k Keeper) IterateStakes(ctx sdk.Context, delAddr sdk.AccAddress,
	fn func(index int64, del types.Stake) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)
	stakerPrefixKey := types.GetStakesKey(delAddr)

	iterator := sdk.KVStorePrefixIterator(store, stakerPrefixKey) // smallest to largest
	defer iterator.Close()

	for i := int64(0); iterator.Valid(); iterator.Next() {
		del := types.Stake{}
		k.cdc.MustUnmarshal(iterator.Value(), &del)
		stop := fn(i, del)
		if stop {
			break
		}
		i++
	}
}

// HasMaxUnbondingStakeEntries - check if unbonding stake has maximum number of entries.
func (k Keeper) HasMaxUnbondingStakeEntries(ctx sdk.Context, stakerAddr sdk.AccAddress, validatorAddr sdk.ValAddress) bool {
	ubd, found := k.GetUnbondingStake(ctx, stakerAddr, validatorAddr)
	if !found {
		return false
	}

	return len(ubd.Entries) >= int(k.MaxEntries(ctx))
}

// Unstake unbonds an amount of staker shares from a given validator. It
// will verify that the unbonding entries between the staker and validator
// are not exceeded and unbond the staked tokens (based on shares) by creating
// an unbonding object and inserting it into the unbonding queue which will be
// processed during the staking EndBlocker.
func (k Keeper) Unstake(
	ctx sdk.Context, stakerAddr sdk.AccAddress, valAddr sdk.ValAddress, sharesAmount sdk.Dec,
) (time.Time, error) {
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return time.Time{}, types.ErrNoStakerForAddress
	}

	if k.HasMaxUnbondingStakeEntries(ctx, stakerAddr, valAddr) {
		return time.Time{}, types.ErrMaxUnbondingStakeEntries
	}

	returnAmount, err := k.UnStakeBond(ctx, stakerAddr, valAddr, sharesAmount)
	if err != nil {
		return time.Time{}, err
	}

	// transfer the validator tokens to the not bonded pool
	if validator.IsBonded() {
		k.BondedStakeTokensToNotBonded(ctx, returnAmount, validator.Description.RegionID)
	}

	completionTime := ctx.BlockHeader().Time.Add(time.Second)
	ubs := k.SetUnbondingStakeEntry(ctx, stakerAddr, valAddr, ctx.BlockHeight(), completionTime, returnAmount)
	k.InsertUBSQueue(ctx, ubs, completionTime)
	return completionTime, nil
}

// UnStakeBond unbonds a particular stake and perform associated store operations.
func (k Keeper) UnStakeBond(
	ctx sdk.Context, stakerAddr sdk.AccAddress, valAddr sdk.ValAddress, shares sdk.Dec,
) (amount math.Int, err error) {
	// check if a stake object exists in the store
	stake, found := k.GetStake(ctx, stakerAddr, valAddr)
	if !found {
		return amount, types.ErrNoStakerForAddress
	}

	// ensure that we have enough shares to remove
	if stake.Shares.LT(shares) {
		return amount, sdkerrors.Wrap(types.ErrNotEnoughStakeShares, stake.Shares.String())
	}

	// get validator
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return amount, stakingtypes.ErrNoValidatorFound
	}

	// subtract shares from stake
	stake.Shares = stake.Shares.Sub(shares)

	stakerAddress, err := sdk.AccAddressFromBech32(stake.StakerAddress)
	if err != nil {
		return amount, err
	}

	isValidatorOperator := stakerAddress.Equals(validator.GetOperator())

	// If the stake is the operator of the validator and unstaking will decrease the validator's
	// self-stake below their minimum, we jail the validator.
	if isValidatorOperator && !validator.Jailed &&
		validator.TokensFromShares(stake.Shares).TruncateInt().LT(validator.MinSelfDelegation) {
		k.JailValidator(ctx, validator)
		validator = k.MustGetValidator(ctx, validator.GetOperator())
	}

	if stake.Shares.IsZero() {
		err = k.RemoveStake(ctx, stake)
		if err != nil {
			return amount, err
		}
		k.UnBondRegion(ctx, validator.Description.RegionID)
	} else {
		k.BondRegion(ctx, validator, stake.Shares.TruncateInt(), false)
		k.SetStake(ctx, stake)
		// call the after stake modification hook
		//err = k.AfterDelegationModified(ctx, stakerAddress, stake.GetValidatorAddr())
	}

	// remove the shares and coins from the validator
	// NOTE that the amount is later (in keeper.Stake) moved between staking module pools
	validator, amount = k.RemoveValidatorTokensAndShares(ctx, validator, shares)

	if validator.DelegatorShares.IsZero() && validator.IsUnbonded() {
		// if not unbonded, we must instead remove validator in EndBlocker once it finishes its unbonding period
		k.RemoveValidator(ctx, validator.GetOperator())
		k.RemoveRegion(ctx, validator.Description.RegionID)
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnstake,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, validator.OperatorAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()+params.BaseDenom),
			sdk.NewAttribute(types.AttributeKeyRegionId, validator.Description.RegionID),
		),
	})
	return amount, nil
}

// RemoveStake removes a stake
func (k Keeper) RemoveStake(ctx sdk.Context, stake types.Stake) error {
	stakerAddress := sdk.MustAccAddressFromBech32(stake.StakerAddress)

	// TODO: Consider calling hooks outside of the store wrapper functions, it's unobvious.
	//if err := k.BeforeDelegationRemoved(ctx, stakerAddress, stake.GetValidatorAddr()); err != nil {
	//	return err
	//}

	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetStakeKey(stakerAddress, stake.GetValidatorAddr()))
	return nil
}

// SetUnbondingStakeEntry adds an entry to the unbonding stake at
// the given addresses. It creates the unbonding stake if it does not exist.
func (k Keeper) SetUnbondingStakeEntry(
	ctx sdk.Context, stakerAddr sdk.AccAddress, validatorAddr sdk.ValAddress,
	creationHeight int64, minTime time.Time, balance math.Int,
) types.UnbondingStake {
	ubs, found := k.GetUnbondingStake(ctx, stakerAddr, validatorAddr)
	if found {
		ubs.AddEntry(creationHeight, minTime, balance)
	} else {
		ubs = types.NewUnbondingStake(stakerAddr, validatorAddr, creationHeight, minTime, balance)
	}
	k.SetUnbondingStake(ctx, ubs)
	return ubs
}

// ValidateUnbondAmount validates that a given unbond amount is valied
// based on upon the converted shares. If the amount is valid, the total
// amount of respective shares is returned, otherwise an error is returned.
func (k Keeper) ValidateUnbondAmount(
	ctx sdk.Context, stakerAddr sdk.AccAddress, valAddr sdk.ValAddress, amt math.Int,
) (shares sdk.Dec, err error) {
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return shares, stakingtypes.ErrNoValidatorFound
	}

	valTokens := math.ZeroInt()

	// ensure validator's tokens can not less than meid amount or delegate amount
	if validator.MeidAmount.GTE(validator.DelegationAmount) {
		valTokens = validator.Tokens.Sub(amt)
		if valTokens.LT(validator.MeidAmount) {
			return shares, types.ErrValidatorTokensAmount
		}
	} else {
		valTokens = validator.Tokens.Sub(amt)
		if valTokens.LT(validator.DelegationAmount) {
			return shares, types.ErrValidatorTokensAmount
		}
	}

	sta, found := k.GetStake(ctx, stakerAddr, valAddr)
	if !found {
		return shares, types.ErrNoStake
	}

	shares, err = validator.SharesFromTokens(amt)
	if err != nil {
		return shares, err
	}

	sharesTruncated, err := validator.SharesFromTokensTruncated(amt)
	if err != nil {
		return shares, err
	}

	staShares := sta.GetShares()
	if sharesTruncated.GT(staShares) {
		return shares, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid shares amount")
	}

	// Cap the shares at the stake's shares. Shares being greater could occur
	// due to rounding, however we don't want to truncate the shares or take the
	// minimum because we want to allow for the full withdraw of shares from a
	// stake.
	if shares.GT(staShares) {
		shares = staShares
	}

	return shares, nil
}

// SetUnbondingStake sets the unbonding stake and associated index.
func (k Keeper) SetUnbondingStake(ctx sdk.Context, ubs types.UnbondingStake) {
	stakerAddress := sdk.MustAccAddressFromBech32(ubs.StakerAddress)

	store := ctx.KVStore(k.storeKey)
	addr, err := sdk.ValAddressFromBech32(ubs.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	key := types.GetUBSKey(stakerAddress, addr)
	store.Set(key, k.cdc.MustMarshal(&ubs))
	store.Set(types.GetUBSByValIndexKey(stakerAddress, addr), []byte{}) // index, store empty bytes
}

// GetUnbondingStake returns a unbonding stake.
func (k Keeper) GetUnbondingStake(ctx sdk.Context, stakerAddr sdk.AccAddress, valAddr sdk.ValAddress) (ubs types.UnbondingStake, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUBSKey(stakerAddr, valAddr)
	value := store.Get(key)

	if value == nil {
		return ubs, false
	}

	k.cdc.MustUnmarshal(value, &ubs)
	return ubs, true
}

// RemoveUnbondingStake removes the unbonding stake object and associated index.
func (k Keeper) RemoveUnbondingStake(ctx sdk.Context, ubd types.UnbondingStake) {
	stakerAddress := sdk.MustAccAddressFromBech32(ubd.StakerAddress)

	store := ctx.KVStore(k.storeKey)
	addr, err := sdk.ValAddressFromBech32(ubd.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	key := types.GetUBSKey(stakerAddress, addr)
	store.Delete(key)
	store.Delete(types.GetUBSByValIndexKey(stakerAddress, addr))
}

// IterateUnbondingStakes iterates through all of the unbonding stakes.
func (k Keeper) IterateUnbondingStakes(ctx sdk.Context, cb func(ubs types.UnbondingStake) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.UnbondingStakeKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		ubs := types.UnbondingStake{}
		k.cdc.MustUnmarshal(iterator.Value(), &ubs)
		if cb(ubs) {
			break
		}
	}
}

// UBSQueueIterator returns all the unbonding queue timeslices from time 0 until endTime.
func (k Keeper) UBSQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.UnbondingStakeQueueKey,
		sdk.InclusiveEndBytes(types.GetUnbondingStakeTimeKey(endTime)))
}

// SequeueAllMatureUBSQueue returns a concatenated list of all the timeslices inclusively previous to
// currTime, and deletes the timeslices from the queue.
func (k Keeper) SequeueAllMatureUBSQueue(ctx sdk.Context, currTime time.Time) (matureUnbonds []types.SVPair) {
	store := ctx.KVStore(k.storeKey)
	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	unbondingTimesliceIterator := k.UBSQueueIterator(ctx, currTime)
	defer unbondingTimesliceIterator.Close()
	for ; unbondingTimesliceIterator.Valid(); unbondingTimesliceIterator.Next() {
		timeslice := types.SVPairs{}
		value := unbondingTimesliceIterator.Value()
		k.cdc.MustUnmarshal(value, &timeslice)
		matureUnbonds = append(matureUnbonds, timeslice.Pairs...)
		store.Delete(unbondingTimesliceIterator.Key())
	}
	return matureUnbonds
}

// CompleteStakeUnBonding completes the unbonding of all mature entries in the
// retrieved unbonding stake object and returns the total unbonding balance
// or an error upon failure.
func (k Keeper) CompleteStakeUnBonding(ctx sdk.Context, stakerAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error) {
	ubs, found := k.GetUnbondingStake(ctx, stakerAddr, valAddr)
	if !found {
		return nil, types.ErrNoUnbondingStake
	}

	bondDenom := k.GetParams(ctx).BondDenom
	balances := sdk.NewCoins()
	ctxTime := ctx.BlockHeader().Time

	// loop through all the entries and complete unbonding mature entries
	for i := 0; i < len(ubs.Entries); i++ {
		entry := ubs.Entries[i]
		if entry.IsMature(ctxTime) {
			ubs.RemoveEntry(int64(i))
			i--

			// track unstake only when remaining or truncated shares are non-zero
			if !entry.Balance.IsZero() {
				amt := sdk.NewCoin(bondDenom, entry.Balance)
				if err := k.bankKeeper.UnstakeCoinsFromModuleToModule(
					ctx, types.NotBondedStakePoolName, types.StakePoolName, sdk.NewCoins(amt),
				); err != nil {
					return nil, err
				}

				balances = balances.Add(amt)
			}
		}
	}

	// set the unbonding stake or remove it if there are no more entries
	if len(ubs.Entries) == 0 {
		k.RemoveUnbondingStake(ctx, ubs)
	} else {
		k.SetUnbondingStake(ctx, ubs)
	}

	return balances, nil
}

// InsertUBSQueue inserts an unbonding stake to the appropriate timeslice
// in the unbonding queue.
func (k Keeper) InsertUBSQueue(ctx sdk.Context, ubs types.UnbondingStake, completionTime time.Time) {
	svPair := types.SVPair{StakerAddress: ubs.StakerAddress, ValidatorAddress: ubs.ValidatorAddress}

	timeSlice := k.GetUBSQueueTimeSlice(ctx, completionTime)
	if len(timeSlice) == 0 {
		k.SetUBSQueueTimeSlice(ctx, completionTime, []types.SVPair{svPair})
	} else {
		timeSlice = append(timeSlice, svPair)
		k.SetUBSQueueTimeSlice(ctx, completionTime, timeSlice)
	}
}

// GetUBSQueueTimeSlice gets a specific unbonding queue timeslice. A timeslice
// is a slice of SVPair corresponding to unbonding stakes that expire at a
// certain time.
func (k Keeper) GetUBSQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (svPairs []types.SVPair) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetUnbondingStakeTimeKey(timestamp))
	if bz == nil {
		return []types.SVPair{}
	}

	pairs := types.SVPairs{}
	k.cdc.MustUnmarshal(bz, &pairs)

	return pairs.Pairs
}

// SetUBSQueueTimeSlice sets a specific unbonding queue timeslice.
func (k Keeper) SetUBSQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys []types.SVPair) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&types.SVPairs{Pairs: keys})
	store.Set(types.GetUnbondingStakeTimeKey(timestamp), bz)
}

func (k Keeper) ParserStakeKey(key []byte) (stakerAddr sdk.AccAddress, valAddr sdk.ValAddress, err error) {
	totalKeyLen := len(key)
	if totalKeyLen < 3 {
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "invalid stake key length: %d", totalKeyLen)
	}
	if key[0] != types.StakeKey[0] {
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "invalid stake key prefix: %X", key[0])
	}

	stakeAddrLen := int(key[1])
	if stakeAddrLen+2 >= totalKeyLen {
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "invalid stake key. length: %d,stakerAddrlength:%d",
			totalKeyLen, stakeAddrLen)
	}
	stakerAddr = key[2 : 2+stakeAddrLen]

	valAddrLen := int(key[2+stakeAddrLen])
	if 3+stakeAddrLen+valAddrLen != totalKeyLen {
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "invalid stake key. length: %d,stakerAddrLen:%d,valAddrLen:%d",
			totalKeyLen, stakeAddrLen, valAddrLen)
	}
	valAddr = key[2+stakeAddrLen+1:]
	return stakerAddr, valAddr, nil
}

func (k Keeper) GetStakesByValidator(ctx sdk.Context, valAddr sdk.ValAddress) ([]*types.Stake, error) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.StakeKey)
	defer iterator.Close()
	var stakes []*types.Stake
	for ; iterator.Valid(); iterator.Next() {
		_, vAddr, err := k.ParserStakeKey(iterator.Key())
		if err != nil {
			return nil, err
		}
		if vAddr.Equals(valAddr) {
			var stakeInfo types.Stake
			k.cdc.MustUnmarshal(iterator.Value(), &stakeInfo)
			stakes = append(stakes, &stakeInfo)
		}

	}
	return stakes, nil
}
