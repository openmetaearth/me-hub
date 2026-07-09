package v4_0_0

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/openmetaearth/me-hub/app/upgrades"
	lightclientmoduletypes "github.com/openmetaearth/me-hub/x/lightclient/types"
)

const (
	UpgradeName = "v4.0.0"
)

var Upgrade = upgrades.Upgrade{
	Name:          UpgradeName,
	CreateHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added: []string{
			lightclientmoduletypes.StoreKey,
		},
	},
}
