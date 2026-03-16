package v2_0_13_patch_4

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/st-chain/me-hub/app/upgrades"
)

const (
	UpgradeName = "v2.0.13.patch.4"
)

var Upgrade = upgrades.Upgrade{
	Name:          UpgradeName,
	CreateHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}
