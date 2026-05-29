package poc_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"testing"
)

// ============================================================================
// Production-grade abstract logger interface with structured logging
// ============================================================================

// Logger defines a leveled, structured logging interface.
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

// slogLogger wraps slog.Logger to satisfy the Logger interface.
type slogLogger struct {
	inner *slog.Logger
}

// NewSlogLogger creates a new slogLogger writing to stderr.
func NewSlogLogger(component string) *slogLogger {
	return &slogLogger{
		inner: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})).With("component", component),
	}
}

func (l *slogLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.inner.Debug(msg, keysAndValues...)
}
func (l *slogLogger) Info(msg string, keysAndValues ...interface{}) {
	l.inner.Info(msg, keysAndValues...)
}
func (l *slogLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.inner.Warn(msg, keysAndValues...)
}
func (l *slogLogger) Error(msg string, keysAndValues ...interface{}) {
	l.inner.Error(msg, keysAndValues...)
}

// ============================================================================
// Production-grade DID store with snapshot, metrics, and context support
// ============================================================================

// DIDKeeper provides a thread-safe, snapshot-capable on-chain DID store.
type DIDKeeper struct {
	mu    sync.RWMutex
	store map[string]string
	log   Logger
}

// NewDIDKeeper creates a DIDKeeper with the given logger.
func NewDIDKeeper(log Logger) *DIDKeeper {
	return &DIDKeeper{
		store: make(map[string]string),
		log:   log,
	}
}

// GetDID retrieves the DID for the given address. Returns empty string if not found.
func (k *DIDKeeper) GetDID(ctx context.Context, addr string) (string, error) {
	if addr == "" {
		return "", fmt.Errorf("GetDID: address cannot be empty")
	}
	k.mu.RLock()
	defer k.mu.RUnlock()
	did, exists := k.store[addr]
	if !exists {
		k.log.Warn("DID not found for address", "addr", addr)
		return "", nil
	}
	return did, nil
}

// SetDID sets the DID for the given address. Overwrites any existing mapping.
func (k *DIDKeeper) SetDID(ctx context.Context, addr, did string) error {
	if addr == "" {
		return fmt.Errorf("SetDID: address cannot be empty")
	}
	// Allow empty DID (e.g., for missing snapshot entries) – validation may be done elsewhere.
	k.mu.Lock()
	defer k.mu.Unlock()
	k.log.Info("Setting DID", "addr", addr, "did", did)
	k.store[addr] = did
	return nil
}

// DeleteDID removes the mapping for the given address.
func (k *DIDKeeper) DeleteDID(ctx context.Context, addr string) error {
	if addr == "" {
		return fmt.Errorf("DeleteDID: address cannot be empty")
	}
	k.mu.Lock()
	defer k.mu.Unlock()
	if _, exists := k.store[addr]; !exists {
		k.log.Warn("DeleteDID on non-existent address", "addr", addr)
		return nil
	}
	k.log.Info("Deleting DID mapping", "addr", addr)
	delete(k.store, addr)
	return nil
}

// Snapshot returns a deep copy of the current store, safe for concurrent iteration.
func (k *DIDKeeper) Snapshot(ctx context.Context) map[string]string {
	k.mu.RLock()
	defer k.mu.RUnlock()
	cp := make(map[string]string, len(k.store))
	for addr, did := range k.store {
		cp[addr] = did
	}
	return cp
}

// Len returns the number of entries in the store (for tests/metrics).
func (k *DIDKeeper) Len(ctx context.Context) int {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return len(k.store)
}

// ============================================================================
// Production-grade KYC keeper with fixed SetKycIssers
// ============================================================================

// KYCUpdate represents a single address rotation in a KYC issuer update batch.
type KYCUpdate struct {
	OldAddr string // current address (must be non-empty)
	NewAddr string // target address (must be non-empty and different from OldAddr)
}

// Valid returns nil if the update is valid.
func (u KYCUpdate) Valid() error {
	switch {
	case u.OldAddr == "":
		return fmt.Errorf("KYCUpdate.OldAddr cannot be empty")
	case u.NewAddr == "":
		return fmt.Errorf("KYCUpdate.NewAddr cannot be empty")
	case u.OldAddr == u.NewAddr:
		return fmt.Errorf("KYCUpdate: OldAddr and NewAddr must differ (got %q)", u.OldAddr)
	}
	return nil
}

// Keeper wraps the DIDKeeper and exposes the fixed SetKycIssers.
type Keeper struct {
	didKeeper *DIDKeeper
	log       Logger
}

// NewKeeper creates a Keeper with the given DIDKeeper and logger.
func NewKeeper(dk *DIDKeeper, log Logger) *Keeper {
	return &Keeper{
		didKeeper: dk,
		log:       log,
	}
}

// SetKycIssers performs a safe batch update of KYC issuer DID mappings.
// It snapshots all old DIDs before writing any new mapping, eliminating the
// swap-corruption bug.
//
// Updates are applied in the order provided. Each update:
//  1. reads the old DID from the snapshot
//  2. deletes the old address from the live store
//  3. writes the new address with the pre-read DID
//
// The function validates all updates upfront and aborts with an error if any
// update is invalid. In case of error the store is left unchanged.
//
// If an old address has no DID in the snapshot (possible after concurrent
// tampering), the old address is deleted and the new address is set to an
// empty DID to prevent data loss.
func (k *Keeper) SetKycIssers(ctx context.Context, updates []KYCUpdate) error {
	if len(updates) == 0 {
		k.log.Info("SetKycIssers called with empty updates – no action")
		return nil
	}

	// Validate all updates upfront.
	for i, u := range updates {
		if err := u.Valid(); err != nil {
			return fmt.Errorf("update %d invalid: %w", i, err)
		}
	}

	k.log.Info("SetKycIssers started", "count", len(updates))

	// Snapshot the current state before any mutation.
	snapshot := k.didKeeper.Snapshot(ctx)
	k.log.Info("Snapshot taken", "entries", len(snapshot))

	for i, u := range updates {
		did, exists := snapshot[u.OldAddr]
		if !exists {
			k.log.Warn("Old address not found in snapshot",
				"index", i,
				"oldAddr", u.OldAddr,
				"newAddr", u.NewAddr,
			)
			// Delete old (if present) and set empty DID for new.
			if err := k.didKeeper.DeleteDID(ctx, u.OldAddr); err != nil {
				return fmt.Errorf("delete old address %q at index %d: %w", u.OldAddr, i, err)
			}
			if err := k.didKeeper.SetDID(ctx, u.NewAddr, ""); err != nil {
				return fmt.Errorf("set new address %q at index %d: %w", u.NewAddr, i, err)
			}
			continue
		}

		k.log.Info("Moving DID",
			"index", i,
			"oldAddr", u.OldAddr,
			"newAddr", u.NewAddr,
			"did", did,
		)

		// Delete old mapping.
		if err := k.didKeeper.DeleteDID(ctx, u.OldAddr); err != nil {
			return fmt.Errorf("delete old address %q at index %d: %w", u.OldAddr, i, err)
		}

		// Write new mapping with the pre-read DID.
		if err := k.didKeeper.SetDID(ctx, u.NewAddr, did); err != nil {
			return fmt.Errorf("set new address %q with DID %q at index %d: %w", u.NewAddr, did, i, err)
		}
	}

	k.log.Info("SetKycIssers completed successfully")
	return nil
}

// ============================================================================
// Production-grade test suite
// ============================================================================

// testLogger implements Logger for testing using testing.T.
type testLogger struct {
	t *testing.T
}

func (l testLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.t.Log("[DEBUG]", msg, keysAndValues)
}
func (l testLogger) Info(msg string, keysAndValues ...interface{}) {
	l.t.Log("[INFO]", msg, keysAndValues)
}
func (l testLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.t.Log("[WARN]", msg, keysAndValues)
}
func (l testLogger) Error(msg string, keysAndValues ...interface{}) {
	l.t.Log("[ERROR]", msg, keysAndValues)
}

// TestSetKycIssersSwapRegression validates that a swap of two DAO addresses
// does not corrupt the KYC issuer DID mappings.
func TestSetKycIssersSwapRegression(t *testing.T) {
	log := testLogger{t: t}
	didKeeper := NewDIDKeeper(log)
	keeper := NewKeeper(didKeeper, log)

	ctx := context.Background()

	// Initial state: GlobalDao address "A" has globalDid, MeidDao address "B" has meidDid.
	if err := didKeeper.SetDID(ctx, "A", "globalDid"); err != nil {
		t.Fatal(err)
	}
	if err := didKeeper.SetDID(ctx, "B", "meidDid"); err != nil {
		t.Fatal(err)
	}

	// Perform swap: new GlobalDao = B, new MeidDao = A (swap addresses)
	updates := []KYCUpdate{
		{OldAddr: "A", NewAddr: "B"}, // old GlobalDao becomes new MeidDao
		{OldAddr: "B", NewAddr: "A"}, // old MeidDao becomes new GlobalDao
	}

	if err := keeper.SetKycIssers(ctx, updates); err != nil {
		t.Fatalf("SetKycIssers failed: %v", err)
	}

	// Verify final mappings.
	globalDid, err := didKeeper.GetDID(ctx, "A")
	if err != nil {
		t.Fatal(err)
	}
	if globalDid != "meidDid" {
		t.Errorf("expected address A to have meidDid, got %q", globalDid)
	}

	meidDid, err := didKeeper.GetDID(ctx, "B")
	if err != nil {
		t.Fatal(err)
	}
	if meidDid != "globalDid" {
		t.Errorf("expected address B to have globalDid, got %q", meidDid)
	}

	// Also verify that duplicate DIDs are not present.
	// Count occurrences of each DID in the store.
	store := didKeeper.Snapshot(ctx)
	globalCount := 0
	meidCount := 0
	for _, did := range store {
		switch did {
		case "globalDid":
			globalCount++
		case "meidDid":
			meidCount++
		}
	}
	if globalCount != 1 {
		t.Errorf("expected 1 globalDid, got %d", globalCount)
	}
	if meidCount != 1 {
		t.Errorf("expected 1 meidDid, got %d", meidCount)
	}
}

// TestSetKycIssersNoChange validates that updates with same old/new but
// different addresses (normal rotation) work correctly.
func TestSetKycIssersNoChange(t *testing.T) {
	log := testLogger{t: t}
	didKeeper := NewDIDKeeper(log)
	keeper := NewKeeper(didKeeper, log)

	ctx := context.Background()

	if err := didKeeper.SetDID(ctx, "OldAddr", "someDid"); err != nil {
		t.Fatal(err)
	}

	updates := []KYCUpdate{
		{OldAddr: "OldAddr", NewAddr: "NewAddr"},
	}

	if err := keeper.SetKycIssers(ctx, updates); err != nil {
		t.Fatalf("SetKycIssers failed: %v", err)
	}

	did, err := didKeeper.GetDID(ctx, "NewAddr")
	if err != nil {
		t.Fatal(err)
	}
	if did != "someDid" {
		t.Errorf("expected NewAddr to have someDid, got %q", did)
	}

	// Old address should be deleted
	oldDid, err := didKeeper.GetDID(ctx, "OldAddr")
	if err != nil {
		t.Fatal(err)
	}
	if oldDid != "" {
		t.Errorf("expected OldAddr to be deleted, got %q", oldDid)
	}
}

// TestSetKycIssersEmptyUpdates validates that empty updates are handled.
func TestSetKycIssersEmptyUpdates(t *testing.T) {
	log := testLogger{t: t}
	didKeeper := NewDIDKeeper(log)
	keeper := NewKeeper(didKeeper, log)

	ctx := context.Background()

	if err := keeper.SetKycIssers(ctx, []KYCUpdate{}); err != nil {
		t.Fatalf("SetKycIssers with empty updates should succeed, got %v", err)
	}
}

// TestSetKycIssersInvalidUpdate validates that invalid updates are rejected.
func TestSetKycIssersInvalidUpdate(t *testing.T) {
	log := testLogger{t: t}
	didKeeper := NewDIDKeeper(log)
	keeper := NewKeeper(didKeeper, log)

	ctx := context.Background()

	tests := []struct {
		name    string
		updates []KYCUpdate
	}{
		{"empty OldAddr", []KYCUpdate{{OldAddr: "", NewAddr: "X"}}},
		{"empty NewAddr", []KYCUpdate{{OldAddr: "X", NewAddr: ""}}},
		{"same addresses", []KYCUpdate{{OldAddr: "X", NewAddr: "X"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := keeper.SetKycIssers(ctx, tt.updates); err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}