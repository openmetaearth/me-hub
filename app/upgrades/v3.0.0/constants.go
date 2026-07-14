package v3_0_0

import (
	storetypes "cosmossdk.io/store/types"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/openmetaearth/me-hub/app/upgrades"
	lightclientmoduletypes "github.com/openmetaearth/me-hub/x/lightclient/types"
)

const (
	UpgradeName = "v3.0.0"
)

var Upgrade = upgrades.Upgrade{
	Name:          UpgradeName,
	CreateHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added: []string{
			consensusparamtypes.StoreKey, // SDK v0.47 → v0.50: consensus params leave x/params
			lightclientmoduletypes.StoreKey,
		},
	},
}
