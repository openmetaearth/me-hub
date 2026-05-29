package keeper

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/cosmos/cosmos-sdk/types/errors"
)

// SetKycIssers atomically migrates KYC issuer DID mappings from old DAO addresses to new DAO addresses.
// It first snapshots all existing DID-to-address mappings, then deletes old entries and sets new ones.
// This prevents corruption when the old and new address sets overlap (e.g., swapping GlobalDao and MeidDao).
// Parameters:
//   - ctx: SDK context for storage operations.
//   - oldAddrs: slice of current DAO addresses to be rotated out.
//   - newAddrs: slice of replacement DAO addresses, must be same length as oldAddrs.
// Returns:
//   - error: any validation failure or storage write error.
func (k Keeper) SetKycIssers(ctx context.Context, oldAddrs, newAddrs []string) error {
	// Input validation: slices must be non-nil, same length, and contain non-empty addresses.
	if err := validateAddressSlices(oldAddrs, newAddrs); err != nil {
		return err
	}

	// Snapshot all old DID mappings before any write.
	snapshot := make(map[string]string, len(oldAddrs))
	for i, oldAddr := range oldAddrs {
		did, err := k.didKeeper.GetDID(ctx, oldAddr)
		if err != nil {
			return fmt.Errorf("failed to read DID for old address %q at index %d: %w", oldAddr, i, err)
		}
		if did == "" {
			slog.Warn("No DID mapping found for old address, proceeding with empty", "address", oldAddr)
		}
		snapshot[oldAddr] = did
	}

	// Delete all old mappings.
	for _, oldAddr := range oldAddrs {
		if err := k.didKeeper.DeleteDID(ctx, oldAddr); err != nil {
			return fmt.Errorf("failed to delete old DID mapping for address %q: %w", oldAddr, err)
		}
	}

	// Set new mappings from snapshot.
	for i, newAddr := range newAddrs {
		oldAddr := oldAddrs[i]
		did := snapshot[oldAddr]
		if err := k.didKeeper.SetDID(ctx, newAddr, did); err != nil {
			return fmt.Errorf("failed to set DID mapping for new address %q (old: %q): %w", newAddr, oldAddr, err)
		}
	}

	slog.Info("KYC issuer DID mappings updated successfully",
		"old_addresses", formatSlice(oldAddrs),
		"new_addresses", formatSlice(newAddrs),
	)
	return nil
}

// validateAddressSlices checks that both slices are of equal positive length and contain non-empty strings.
func validateAddressSlices(old, new []string) error {
	if len(old) != len(new) {
		return errors.Wrapf(errors.ErrInvalidRequest, "oldAddrs length (%d) does not match newAddrs length (%d)", len(old), len(new))
	}
	if len(old) == 0 {
		return errors.Wrap(errors.ErrInvalidRequest, "address slices must not be empty")
	}
	for i, addr := range old {
		if strings.TrimSpace(addr) == "" {
			return fmt.Errorf("oldAddrs[%d] is empty", i)
		}
	}
	for i, addr := range new {
		if strings.TrimSpace(addr) == "" {
			return fmt.Errorf("newAddrs[%d] is empty", i)
		}
	}
	return nil
}

// formatSlice joins a string slice with commas for logging.
func formatSlice(s []string) string {
	return strings.Join(s, ", ")
}