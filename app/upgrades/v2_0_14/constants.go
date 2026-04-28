package v2_0_14

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/openmetaearth/me-hub/app/upgrades"
)

const (
	UpgradeName = "v2.0.14"
)

var Upgrade = upgrades.Upgrade{
	Name:          UpgradeName,
	CreateHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}
