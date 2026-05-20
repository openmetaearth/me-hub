package keeper

import (
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	testutilstypes "github.com/openmetaearth/me-hub/testutil/types"
)

// MigrateValidatorsFromV1 migrates all validators stored in V1 format (me-hub v47 era)
// to the current V2 validator format (me-hub v50). The key schema changes are:
//
//	V1 field 6  = staker_shares → V2 field 6  = delegator_shares
//	V1 field 11 = min_self_stake → V2 field 11 = min_self_delegation
//	V1 field 12 = delegation_amount (bytes) → V2 field 15 = delegation_amount
//	V1 field 13 = meid_amount → V2 field 16 = meid_amount
//	V1 field 14 = owner_address → V2 field 14 = owner_address
//	V1 field 15 = unbonding_ids → V2 field 13 = unbonding_ids
//	V1 field 16 = unbonding_on_hold_ref_count → V2 field 12 = unbonding_on_hold_ref_count
func (k *Keeper) MigrateValidatorsFromV1(ctx sdk.Context) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), stakingtypes.ValidatorsKey)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var v1 testutilstypes.ValidatorV1
		if err := k.cdc.Unmarshal(iterator.Value(), &v1); err != nil {
			// Already migrated or different format — skip
			continue
		}

		// Map V1 fields to current stakingtypes.Validator
		v2 := stakingtypes.Validator{
			OperatorAddress:         v1.OperatorAddress,
			ConsensusPubkey:         v1.ConsensusPubkey,
			Jailed:                  v1.Jailed,
			Status:                  v1.Status,
			Tokens:                  v1.Tokens,
			DelegatorShares:         v1.StakerShares,
			Description:             v1.Description,
			UnbondingHeight:         v1.UnbondingHeight,
			UnbondingTime:           v1.UnbondingTime,
			Commission:              v1.Commission,
			MinSelfDelegation:       v1.MinSelfStake,
			UnbondingOnHoldRefCount: v1.UnbondingOnHoldRefCount,
			UnbondingIds:            v1.UnbondingIds,
			OwnerAddress:            v1.OwnerAddress,
			DelegationAmount:        v1.DelegationAmount,
			MeidAmount:              v1.MeidAmount,
		}

		bz := k.cdc.MustMarshal(&v2)
		store.Set(iterator.Key(), bz)
	}
	return nil
}
