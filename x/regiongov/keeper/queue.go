package keeper

import (
	"encoding/binary"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

// InsertVotingQueue schedules a proposal for tally at its voting end time.
func (k Keeper) InsertVotingQueue(ctx sdk.Context, end time.Time, id uint64) {
	ctx.KVStore(k.storeKey).Set(types.VotingQueueKey(end, id), types.Uint64ToBytes(id))
}

// DeleteVotingQueue removes a proposal from the tally queue.
func (k Keeper) DeleteVotingQueue(ctx sdk.Context, end time.Time, id uint64) {
	ctx.KVStore(k.storeKey).Delete(types.VotingQueueKey(end, id))
}

// DequeueMatureProposals returns proposal ids whose voting end time is at or before endTime
// and removes those queue entries.
func (k Keeper) DequeueMatureProposals(ctx sdk.Context, endTime time.Time) []uint64 {
	ids, keys := k.matureQueue(ctx, endTime)
	store := ctx.KVStore(k.storeKey)
	for _, key := range keys {
		store.Delete(key)
	}
	return ids
}

func (k Keeper) matureQueue(ctx sdk.Context, endTime time.Time) ([]uint64, [][]byte) {
	store := ctx.KVStore(k.storeKey)
	// PrefixEndBytes covers every proposal id stored under this end time.
	// InclusiveEndBytes only appends 0x00 and would stop before the id suffix.
	iterator := store.Iterator(types.VotingQueuePrefix, sdk.PrefixEndBytes(types.VotingQueueTimePrefix(endTime)))
	defer iterator.Close()

	var ids []uint64
	var keys [][]byte
	for ; iterator.Valid(); iterator.Next() {
		value := iterator.Value()
		if len(value) != 8 {
			continue
		}
		ids = append(ids, binary.BigEndian.Uint64(append([]byte(nil), value...)))
		keys = append(keys, append([]byte(nil), iterator.Key()...))
	}
	return ids, keys
}
