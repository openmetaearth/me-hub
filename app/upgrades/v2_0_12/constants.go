package v2_0_12

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/st-chain/me-hub/app/upgrades"
)

const (
	UpgradeName = "V2.0.12"
)

var Upgrade = upgrades.Upgrade{
	Name:          UpgradeName,
	CreateHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{},
}
