package v3

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/openmetaearth/me-hub/app/upgrades"
	lightclientmoduletypes "github.com/openmetaearth/me-hub/x/lightclient/types"
)

const (
	UpgradeName = "v3.0.0"
)

var Upgrade = upgrades.Upgrade{
	Name:            UpgradeName,
	CreateHandlerV3: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added: []string{
			// NOTE: consensus store must NOT be listed here — med-v2 already mounted
			// consensusparamtypes.StoreKey from genesis (height 1). Marking it Added
			// sets IAVL initialVersion to upgradeHeight and fails with:
			// "initial version set to N, but found earlier version 1".
			lightclientmoduletypes.StoreKey,
		},
	},
}
