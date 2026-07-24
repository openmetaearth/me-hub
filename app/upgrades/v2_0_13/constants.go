package v2_0_13

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/openmetaearth/me-hub/app/upgrades"
)

const (
	UpgradeName = "v2.0.13"
)

var Upgrade = upgrades.Upgrade{
	Name:          UpgradeName,
	CreateHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}
