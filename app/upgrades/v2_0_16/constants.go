package v2_0_16 //nolint:revive

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/openmetaearth/me-hub/app/upgrades"
)

const (
	UpgradeName = "v2.0.16"
)

var Upgrade = upgrades.Upgrade{
	Name:          UpgradeName,
	CreateHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}
