package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/gravity/types"
)

// --- PROPOSAL RELAYER --- //
func (k Keeper) SetProposalRelayer(ctx sdk.Context, proposalRelayer *types.ProposalRelayer) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ProposalRelayerKey, k.cdc.MustMarshal(proposalRelayer))
}

func (k Keeper) GetProposalRelayer(ctx sdk.Context) (proposalRelayer types.ProposalRelayer, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ProposalRelayerKey)
	if bz == nil {
		return proposalRelayer, false
	}
	k.cdc.MustUnmarshal(bz, &proposalRelayer)
	return proposalRelayer, true
}

func (k Keeper) IsProposalRelayer(ctx sdk.Context, relayerAddr string) bool {
	proposalRelayer, found := k.GetProposalRelayer(ctx)
	if !found {
		return false
	}
	for _, relayer := range proposalRelayer.Relayers {
		if relayer == relayerAddr {
			return true
		}
	}
	return false
}

// SetRelayerByExternalAddress sets the external address for a given relayer
func (k Keeper) SetRelayerByExternalAddress(ctx sdk.Context, externalAddress string, relayerAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetRelayerAddressByExternalKey(externalAddress), relayerAddr.Bytes())
}

// GetRelayerByExternalAddress returns the external address for a given gravity relayer
func (k Keeper) GetRelayerByExternalAddress(ctx sdk.Context, externalAddress string) (sdk.AccAddress, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetRelayerAddressByExternalKey(externalAddress))
	if bz == nil {
		return nil, false
	}
	return sdk.AccAddress(bz), true
}

// DelRelayerByExternalAddress delete the external address for a give relayer
func (k Keeper) DelRelayerByExternalAddress(ctx sdk.Context, externalAddress string) {
	store := ctx.KVStore(k.storeKey)
	relayerAddr := types.GetRelayerAddressByExternalKey(externalAddress)
	if !store.Has(relayerAddr) {
		return
	}
	store.Delete(relayerAddr)
}

// --- RELAYER TOTAL POWER --- //

// GetLastTotalPower Load the last total relayer power.
func (k Keeper) GetLastTotalPower(ctx sdk.Context) sdkmath.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LastTotalPowerKey)

	if bz == nil {
		return sdkmath.ZeroInt()
	}

	ip := sdk.IntProto{}
	k.cdc.MustUnmarshal(bz, &ip)

	return ip.Int
}

// SetLastTotalPower Set the last total relayer power.
func (k Keeper) SetLastTotalPower(ctx sdk.Context) {
	relayers := k.GetAllRelayers(ctx, true)
	totalPower := sdkmath.ZeroInt()
	for _, relayer := range relayers {
		totalPower = totalPower.Add(relayer.GetPower())
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastTotalPowerKey, k.cdc.MustMarshal(&sdk.IntProto{Int: totalPower}))
}

// --- RELAYERS --- //
func (k Keeper) IterateRelayer(ctx sdk.Context, cb func(relayer types.Relayer) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.RelayerKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		relayer := types.Relayer{}
		k.cdc.MustUnmarshal(iterator.Value(), &relayer)
		if cb(relayer) {
			break
		}
	}
}

// SetRelayer save Relayer data
func (k Keeper) SetRelayer(ctx sdk.Context, relayerAddress sdk.AccAddress, relayer types.Relayer) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetRelayerKey(relayerAddress), k.cdc.MustMarshal(&relayer))
}

func (k Keeper) HasRelayer(ctx sdk.Context, addr sdk.AccAddress) (found bool) {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetRelayerKey(addr))
}

// GetRelayer get Relayer data
func (k Keeper) GetRelayer(ctx sdk.Context, addr sdk.AccAddress) (relayer types.Relayer, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetRelayerKey(addr))
	if value == nil {
		return relayer, false
	}
	k.cdc.MustUnmarshal(value, &relayer)
	return relayer, true
}

func (k Keeper) DelRelayer(ctx sdk.Context, relayer sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetRelayerKey(relayer)
	if !store.Has(key) {
		return
	}
	store.Delete(key)
}

func (k Keeper) GetAllRelayers(ctx sdk.Context, isOnline bool) (relayers types.Relayers) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.RelayerKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var relayer types.Relayer
		k.cdc.MustUnmarshal(iterator.Value(), &relayer)
		if isOnline && !relayer.Online {
			continue
		}
		relayers = append(relayers, relayer)
	}
	return relayers
}

func (k Keeper) GetAllBondedAmount(ctx sdk.Context) sdk.Int {
	relayers := k.GetAllRelayers(ctx, false)
	totalBonded := sdk.ZeroInt()
	for _, relayer := range relayers {
		totalBonded = totalBonded.Add(relayer.DelegateAmount)
	}
	return totalBonded
}

func (k Keeper) SlashRelayer(ctx sdk.Context, relayerAddrStr string) error {
	relayerAddr, err := sdk.AccAddressFromBech32(relayerAddrStr)
	if err != nil {
		return err
	}
	relayer, found := k.GetRelayer(ctx, relayerAddr)
	if !found {
		return types.ErrNotFoundRelayer
	}
	if !relayer.Online {
		return nil
	}
	relayer.SlashTimes += 1
	if uint64(relayer.SlashTimes) >= k.GetParams(ctx).MaxSlashTimes {
		relayer.Online = false
	}
	k.SetRelayer(ctx, relayerAddr, relayer)
	k.SetLastRelayerSlashBlockHeight(ctx, uint64(ctx.BlockHeight()))
	return nil
}

// SetLastRelayerSlashBlockHeight sets the last relayer slash block height.
func (k Keeper) SetLastRelayerSlashBlockHeight(ctx sdk.Context, blockHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastRelayerSlashBlockHeight, sdk.Uint64ToBigEndian(blockHeight))
}

func (s Keeper) checkIsRelayer(ctx sdk.Context, addr sdk.AccAddress) error {
	relayer, found := s.GetRelayer(ctx, addr)
	if !found {
		return types.ErrNotFoundRelayer
	}
	if !relayer.Online {
		return types.ErrRelayerNotOnLine
	}
	return nil
}
