package v3_0_0

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/openmetaearth/me-hub/app/upgrades"
)

const (
	UpgradeName = "v3.0.0"
)

var Upgrade = upgrades.Upgrade{
	Name:          UpgradeName,
	CreateHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{},
}
