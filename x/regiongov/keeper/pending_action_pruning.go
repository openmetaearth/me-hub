package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

const pendingActionRetention = 7 * 24 * time.Hour

// pruneCompletedPendingActions removes decided requests strictly more than seven
// days after submission, not approval/rejection. Pending requests never expire.
func (k Keeper) pruneCompletedPendingActions(ctx sdk.Context) {
	cutoff := ctx.BlockTime().Add(-pendingActionRetention)
	var expiredIDs []uint64
	k.iterate(ctx, types.PendingActionPrefix, func(_, value []byte) bool {
		var action types.PendingAction
		k.cdc.MustUnmarshal(value, &action)
		completed := action.Status == types.PendingActionStatusApproved || action.Status == types.PendingActionStatusRejected
		if completed && action.SubmitTime.Before(cutoff) {
			expiredIDs = append(expiredIDs, action.Id)
		}
		return false
	})
	// Close the iterator before deleting from its key range. Keep the sequence
	// and the applied equity/freeze state intact; action IDs must not be reused.
	store := ctx.KVStore(k.storeKey)
	for _, id := range expiredIDs {
		store.Delete(types.PendingActionKey(id))
	}
}
