package v2_0_18 //nolint:revive

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/openmetaearth/me-hub/app/upgrades"
	regiongovtypes "github.com/openmetaearth/me-hub/x/regiongov/types"
)

const (
	UpgradeName = "v2.0.18"
)

var Upgrade = upgrades.Upgrade{
	Name:          UpgradeName,
	CreateHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added:   []string{regiongovtypes.StoreKey},
		Deleted: []string{},
	},
}
