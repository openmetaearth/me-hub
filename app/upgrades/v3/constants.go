package v3

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/st-chain/me-hub/app/upgrades"
	eibctypes "github.com/st-chain/me-hub/x/eibc/types"
)

const (
	UpgradeName = "v3"
)

var Upgrade = upgrades.Upgrade{
	Name:          UpgradeName,
	CreateHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added: []string{eibctypes.ModuleName},
	},
}
