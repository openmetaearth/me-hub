package v2_0_15_rc3 //nolint:revive

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/openmetaearth/me-hub/app/upgrades"
)

const (
	UpgradeName = "v2.0.15.rc3"
)

var Upgrade = upgrades.Upgrade{
	Name:          UpgradeName,
	CreateHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}
