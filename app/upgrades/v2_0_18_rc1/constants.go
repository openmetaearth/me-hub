package v2_0_18_rc1 //nolint:revive

import "github.com/openmetaearth/me-hub/app/upgrades"

const UpgradeName = "v2.0.18.rc1"

// Upgrade follows v2.0.18, which already added the regiongov store.
var Upgrade = upgrades.Upgrade{
	Name:          UpgradeName,
	CreateHandler: CreateUpgradeHandler,
}
