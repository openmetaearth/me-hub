package v2_0_13

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	epochstypes "github.com/osmosis-labs/osmosis/v15/x/epochs/types"
	gammtypes "github.com/osmosis-labs/osmosis/v15/x/gamm/types"
	lockuptypes "github.com/osmosis-labs/osmosis/v15/x/lockup/types"
	poolmanagertypes "github.com/osmosis-labs/osmosis/v15/x/poolmanager/types"
	txfeestypes "github.com/osmosis-labs/osmosis/v15/x/txfees/types"
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
			gammtypes.ModuleName,
			poolmanagertypes.ModuleName,
			txfeestypes.ModuleName,
			epochstypes.ModuleName,
		},
	},
}
