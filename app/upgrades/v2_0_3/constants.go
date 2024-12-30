package v2_0_3

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/st-chain/me-hub/app/upgrades"
	megrouptypes "github.com/st-chain/me-hub/x/megroup/types"
)

const (
	UpgradeName = "v2_0_3"
)

var Upgrade = upgrades.Upgrade{
	Name:          UpgradeName,
	CreateHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added: []string{
			megrouptypes.ModuleName,
		},
		Deleted: []string{},
	},
}
