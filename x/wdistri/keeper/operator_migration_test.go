package keeper

import (
	"bytes"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/x/wdistri/types"
)

type migrationAccountKeeper struct{ types.AccountKeeper }

func (migrationAccountKeeper) GetModuleAddress(string) sdk.AccAddress {
	return sdk.AccAddress(bytes.Repeat([]byte{9}, 20))
}

func migrationKeeper(t *testing.T) (*Keeper, sdk.Context, *storetypes.KVStoreKey) {
	t.Helper()
	key := storetypes.NewKVStoreKey(distrtypes.StoreKey)
	ctx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient")).Ctx
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	k := NewKeeper(cdc, key, paramtypes.Subspace{}, migrationAccountKeeper{}, nil, nil, "fees", "authority")
	return k, ctx, key
}

func migrationStoreSnapshot(ctx sdk.Context, key *storetypes.KVStoreKey) map[string]string {
	result := make(map[string]string)
	iter := ctx.KVStore(key).Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		result[string(iter.Key())] = string(iter.Value())
	}
	return result
}

func TestMigrateValidatorOperatorPreservesDistribution(t *testing.T) {
	k, ctx, key := migrationKeeper(t)
	oldAddr := sdk.ValAddress(bytes.Repeat([]byte{1}, 20))
	newAddr := sdk.ValAddress(bytes.Repeat([]byte{2}, 20))
	otherAddr := sdk.ValAddress(bytes.Repeat([]byte{3}, 20))
	delegator := sdk.AccAddress(bytes.Repeat([]byte{4}, 20))
	withdraw := sdk.AccAddress(bytes.Repeat([]byte{5}, 20))
	coins := sdk.NewDecCoins(sdk.NewDecCoinFromDec("stake", sdk.MustNewDecFromStr("12.345")))
	outstanding := distrtypes.ValidatorOutstandingRewards{Rewards: coins}
	commission := distrtypes.ValidatorAccumulatedCommission{Commission: coins}
	current := distrtypes.NewValidatorCurrentRewards(coins, 23)
	historical := distrtypes.NewValidatorHistoricalRewards(coins, 2)
	slash := distrtypes.NewValidatorSlashEvent(18, sdk.MustNewDecFromStr("0.125"))
	starting := distrtypes.NewDelegatorStartingInfo(17, sdk.MustNewDecFromStr("19.25"), 50)
	k.SetValidatorOutstandingRewards(ctx, oldAddr, outstanding)
	k.SetValidatorAccumulatedCommission(ctx, oldAddr, commission)
	k.SetValidatorCurrentRewards(ctx, oldAddr, current)
	k.SetValidatorHistoricalRewards(ctx, oldAddr, 17, historical)
	k.SetValidatorHistoricalRewards(ctx, oldAddr, 22, historical)
	k.SetValidatorSlashEvent(ctx, oldAddr, 60, 18, slash)
	secondSlash := distrtypes.NewValidatorSlashEvent(19, sdk.MustNewDecFromStr("0.25"))
	k.SetValidatorSlashEvent(ctx, oldAddr, 60, 19, secondSlash)
	k.SetDelegatorStartingInfo(ctx, oldAddr, delegator, starting)
	k.SetValidatorCurrentRewards(ctx, otherAddr, current)
	k.SetDelegatorWithdrawAddr(ctx, sdk.AccAddress(oldAddr), withdraw)
	k.SetDelegatorWithdrawAddr(ctx, sdk.AccAddress(newAddr), delegator)
	k.SetPreviousProposerConsAddr(ctx, sdk.ConsAddress(oldAddr))
	before := migrationStoreSnapshot(ctx, key)

	cacheCtx, commit := ctx.CacheContext()
	require.NoError(t, k.MigrateValidatorOperator(cacheCtx, oldAddr, newAddr))
	require.Equal(t, before, migrationStoreSnapshot(ctx, key), "migration must respect the caller's cache")
	commit()

	require.Equal(t, outstanding, k.GetValidatorOutstandingRewards(ctx, newAddr))
	require.Equal(t, commission, k.GetValidatorAccumulatedCommission(ctx, newAddr))
	require.Equal(t, current, k.GetValidatorCurrentRewards(ctx, newAddr))
	for _, period := range []uint64{17, 22} {
		require.Equal(t, historical, k.GetValidatorHistoricalRewards(ctx, newAddr, period))
	}
	for _, event := range []distrtypes.ValidatorSlashEvent{slash, secondSlash} {
		actual, found := k.GetValidatorSlashEvent(ctx, newAddr, 60, event.ValidatorPeriod)
		require.True(t, found)
		require.Equal(t, event, actual)
	}
	require.Equal(t, starting, k.GetDelegatorStartingInfo(ctx, newAddr, delegator))
	require.False(t, k.HasDelegatorStartingInfo(ctx, oldAddr, delegator))
	require.Equal(t, current, k.GetValidatorCurrentRewards(ctx, otherAddr))
	require.Equal(t, withdraw, k.GetDelegatorWithdrawAddr(ctx, sdk.AccAddress(oldAddr)))
	require.Equal(t, delegator, k.GetDelegatorWithdrawAddr(ctx, sdk.AccAddress(newAddr)))
	require.Equal(t, sdk.ConsAddress(oldAddr), k.GetPreviousProposerConsAddr(ctx))

	// Migrating back must reproduce every byte, including period/reference counts,
	// and proves no old-address distribution records were left behind.
	require.NoError(t, k.MigrateValidatorOperator(ctx, newAddr, oldAddr))
	require.Equal(t, before, migrationStoreSnapshot(ctx, key))
}

func TestMigrateValidatorOperatorRejectsDestinationState(t *testing.T) {
	setters := map[string]func(*Keeper, sdk.Context, sdk.ValAddress){
		"outstanding": func(k *Keeper, ctx sdk.Context, addr sdk.ValAddress) {
			k.SetValidatorOutstandingRewards(ctx, addr, distrtypes.ValidatorOutstandingRewards{})
		},
		"commission": func(k *Keeper, ctx sdk.Context, addr sdk.ValAddress) {
			k.SetValidatorAccumulatedCommission(ctx, addr, distrtypes.ValidatorAccumulatedCommission{})
		},
		"current": func(k *Keeper, ctx sdk.Context, addr sdk.ValAddress) {
			k.SetValidatorCurrentRewards(ctx, addr, distrtypes.NewValidatorCurrentRewards(sdk.DecCoins{}, 1))
		},
		"historical": func(k *Keeper, ctx sdk.Context, addr sdk.ValAddress) {
			k.SetValidatorHistoricalRewards(ctx, addr, 0, distrtypes.NewValidatorHistoricalRewards(sdk.DecCoins{}, 1))
		},
		"slash": func(k *Keeper, ctx sdk.Context, addr sdk.ValAddress) {
			k.SetValidatorSlashEvent(ctx, addr, 10, 2, distrtypes.NewValidatorSlashEvent(2, sdk.ZeroDec()))
		},
		"starting": func(k *Keeper, ctx sdk.Context, addr sdk.ValAddress) {
			k.SetDelegatorStartingInfo(ctx, addr, sdk.AccAddress(addr), distrtypes.NewDelegatorStartingInfo(0, sdk.ZeroDec(), 0))
		},
	}
	for name, setter := range setters {
		t.Run(name, func(t *testing.T) {
			k, ctx, key := migrationKeeper(t)
			oldAddr := sdk.ValAddress(bytes.Repeat([]byte{1}, 20))
			newAddr := sdk.ValAddress(bytes.Repeat([]byte{2}, 20))
			k.SetValidatorCurrentRewards(ctx, oldAddr, distrtypes.NewValidatorCurrentRewards(sdk.DecCoins{}, 9))
			setter(k, ctx, newAddr)
			before := migrationStoreSnapshot(ctx, key)
			require.ErrorContains(t, k.MigrateValidatorOperator(ctx, oldAddr, newAddr), "already exists")
			require.Equal(t, before, migrationStoreSnapshot(ctx, key))
		})
	}
}

func TestMigrateValidatorOperatorPreservesRawRecords(t *testing.T) {
	k, ctx, key := migrationKeeper(t)
	oldAddr := sdk.ValAddress(bytes.Repeat([]byte{1}, 20))
	newAddr := sdk.ValAddress(bytes.Repeat([]byte{2}, 32))
	store := ctx.KVStore(key)
	// Preserve exact key suffixes and opaque protobuf bytes without decoding or
	// normalizing them, and handle different length-prefixed address lengths.
	oldKey := distrtypes.GetValidatorHistoricalRewardsKey(oldAddr, 71)
	newKey := distrtypes.GetValidatorHistoricalRewardsKey(newAddr, 71)
	value := []byte{0x98, 0x06, 0x01} // unknown protobuf field
	store.Set(oldKey, value)
	store.Set(distrtypes.GetValidatorOutstandingRewardsKey(oldAddr), []byte{})
	require.NoError(t, k.MigrateValidatorOperator(ctx, oldAddr, newAddr))
	require.Equal(t, value, store.Get(newKey))
	require.False(t, store.Has(oldKey))
	require.True(t, store.Has(distrtypes.GetValidatorOutstandingRewardsKey(newAddr)))
	require.Empty(t, store.Get(distrtypes.GetValidatorOutstandingRewardsKey(newAddr)))
	require.False(t, store.Has(distrtypes.GetValidatorOutstandingRewardsKey(oldAddr)))
}

func TestMigrateValidatorOperatorValidatesAddresses(t *testing.T) {
	k, ctx, key := migrationKeeper(t)
	addr := sdk.ValAddress(bytes.Repeat([]byte{1}, 20))
	before := migrationStoreSnapshot(ctx, key)
	require.Error(t, k.MigrateValidatorOperator(ctx, addr, addr))
	require.Error(t, k.MigrateValidatorOperator(ctx, nil, addr))
	require.Error(t, k.MigrateValidatorOperator(ctx, addr, nil))
	require.Equal(t, before, migrationStoreSnapshot(ctx, key))
	require.NoError(t, k.MigrateValidatorOperator(ctx, addr, sdk.ValAddress(bytes.Repeat([]byte{2}, 20))))
	require.Equal(t, before, migrationStoreSnapshot(ctx, key), "absent source state should not invent records")
}
