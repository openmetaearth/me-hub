package v2_0_13

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	lockuptypes "github.com/osmosis-labs/osmosis/v15/x/lockup/types"
	"github.com/st-chain/me-hub/app/upgrades"
	bsctypes "github.com/st-chain/me-hub/x/bsc/types"
)

const (
	UpgradeName = "v2.0.13"
)

var Upgrade = upgrades.Upgrade{
	Name:          UpgradeName,
	CreateHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added: []string{
			bsctypes.ModuleName,
		},
		Deleted: []string{
			lockuptypes.ModuleName,
		},
	},
}
