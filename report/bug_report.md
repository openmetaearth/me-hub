package keeper

import (
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	ErrInvalidAddressLength = errors.New("old and new address lists must have equal length")
	ErrEmptyAddress         = errors.New("address cannot be empty")
	ErrDuplicateAddress     = errors.New("duplicate address detected within a list")
	ErrInvalidAddressFormat = errors.New("address is not a valid bech32 address")
	ErrNilDIDKeeper         = errors.New("didKeeper is nil")
	ErrDIDOperation         = errors.New("DID operation failed")
	ErrNilAddressList       = errors.New("address list must not be nil")
)

func (k Keeper) SetKycIssuers(ctx sdk.Context, oldAddrs, newAddrs []string) error {
	logger := k.Logger(ctx)
	logger.Debug("SetKycIssuers called", "old_addrs", strings.Join(oldAddrs, ","), "new_addrs", strings.Join(newAddrs, ","))

	if k.didKeeper == nil {
		return ErrNilDIDKeeper
	}
	if err := validateAddressPairs(oldAddrs, newAddrs); err != nil {
		return fmt.Errorf("input validation: %w", err)
	}
	if len(oldAddrs) == 0 {
		logger.Warn("SetKycIssuers called with empty address lists; no action taken")
		return nil
	}

	// Snapshot all existing DID mappings before any write to avoid corruption during swap.
	snapshot := make(map[string]string, len(oldAddrs))
	hasDID := make(map[string]bool, len(oldAddrs))
	for _, oldAddr := range oldAddrs {
		did, found := k.didKeeper.GetDID(ctx, oldAddr)
		if !found {
			logger.Warn("no existing DID mapping for old DAO address; skipping", "old_address", oldAddr)
			hasDID[oldAddr] = false
		} else {
			snapshot[oldAddr] = did
			hasDID[oldAddr] = true
		}
	}

	for i, oldAddr := range oldAddrs {
		newAddr := newAddrs[i]
		// Skip if addresses are identical to avoid unnecessary delete/set
		if oldAddr == newAddr {
			logger.Debug("old and new addresses are identical; skipping", "address", oldAddr)
			continue
		}
		if !hasDID[oldAddr] {
			logger.Debug("skipping address without a DID mapping", "old_address", oldAddr, "new_address", newAddr)
			continue
		}
		did := snapshot[oldAddr]
		if err := k.didKeeper.DeleteDID(ctx, oldAddr); err != nil {
			return fmt.Errorf("%w: delete old mapping for address %q (did %q): %w", ErrDIDOperation, oldAddr, did, err)
		}
		if err := k.didKeeper.SetDID(ctx, newAddr, did); err != nil {
			return fmt.Errorf("%w: set new mapping for address %q (did %q): %w", ErrDIDOperation, newAddr, did, err)
		}
	}

	logger.Info("KYC issuer DID mappings migrated successfully", "count", len(oldAddrs))
	return nil
}

func validateAddressPairs(oldAddrs, newAddrs []string) error {
	if oldAddrs == nil || newAddrs == nil {
		return fmt.Errorf("%w: old nil=%v, new nil=%v", ErrNilAddressList, oldAddrs == nil, newAddrs == nil)
	}
	if len(oldAddrs) != len(newAddrs) {
		return fmt.Errorf("%w: old length %d, new length %d", ErrInvalidAddressLength, len(oldAddrs), len(newAddrs))
	}

	seenOld := make(map[string]int, len(oldAddrs))
	seenNew := make(map[string]int, len(newAddrs))

	for i, addr := range oldAddrs {
		if addr == "" {
			return fmt.Errorf("%w: old address at index %d", ErrEmptyAddress, i)
		}
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return fmt.Errorf("%w: old address %q at index %d: %w", ErrInvalidAddressFormat, addr, i, err)
		}
		if _, exists := seenOld[addr]; exists {
			return fmt.Errorf("%w: old address %q at index %d", ErrDuplicateAddress, addr, i)
		}
		seenOld[addr] = i
	}

	for i, addr := range newAddrs {
		if addr == "" {
			return fmt.Errorf("%w: new address at index %d", ErrEmptyAddress, i)
		}
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return fmt.Errorf("%w: new address %q at index %d: %w", ErrInvalidAddressFormat, addr, i, err)
		}
		if _, exists := seenNew[addr]; exists {
			return fmt.Errorf("%w: new address %q at index %d", ErrDuplicateAddress, addr, i)
		}
		seenNew[addr] = i
	}

	return nil
}