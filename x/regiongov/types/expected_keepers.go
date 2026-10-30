package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

// StakingKeeper reads region operators. Region.OperatorAddress identifies the voter.
type StakingKeeper interface {
	GetRegion(ctx sdk.Context, regionID string) (wstakingtypes.Region, bool)
	GetAllRegion(ctx sdk.Context) []wstakingtypes.Region
}

// DaoKeeper identifies the dao. A dao vote is not part of the 51% count.
type DaoKeeper interface {
	IsGlobalDao(ctx sdk.Context, address string) bool
}
