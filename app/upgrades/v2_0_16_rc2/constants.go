package v2_0_16_rc2

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/openmetaearth/me-hub/app/upgrades"
)

const (
	UpgradeName = "v2.0.16.uat.rc2"
)

var Upgrade = upgrades.Upgrade{
	Name:          UpgradeName,
	CreateHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}
