package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// EpochHooks is a minimal interface matching osmosis epochs hooks used by delayedack.
type EpochHooks interface {
	BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error
	AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error
}
