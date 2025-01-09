package v2_0_6

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/st-chain/me-hub/app/upgrades"
)

const (
	UpgradeName = "v2.0.6"
)

var Upgrade = upgrades.Upgrade{
	Name:          UpgradeName,
	CreateHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{},
}
