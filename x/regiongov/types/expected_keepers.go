package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

// StakingKeeper reads region operators and applies passed proposal actions.
// Region.OperatorAddress identifies the voter.
type StakingKeeper interface {
	GetRegion(ctx sdk.Context, regionID string) (wstakingtypes.Region, bool)
	GetAllRegion(ctx sdk.Context) []wstakingtypes.Region
	// WithdrawRegionTreasure pays amount from the region treasury. maxPercent 0
	// skips the equity share check; otherwise amount*100 must be <= balance*maxPercent.
	WithdrawRegionTreasure(ctx sdk.Context, regionID, receiver string, amount sdk.Coins, maxPercent uint32) error
	SetRegionOperator(ctx sdk.Context, regionID, newOperator string) error
}

// DaoKeeper identifies the dao. A dao vote is not part of the 51% count.
type DaoKeeper interface {
	IsGlobalDao(ctx sdk.Context, address string) bool
}
