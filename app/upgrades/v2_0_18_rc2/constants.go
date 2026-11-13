package v2_0_18_rc2 //nolint:revive

import "github.com/openmetaearth/me-hub/app/upgrades"

const UpgradeName = "v2.0.18.rc2"

// Upgrade reuses the existing stores; regiongov was added by v2.0.18.
var Upgrade = upgrades.Upgrade{
	Name:          UpgradeName,
	CreateHandler: CreateUpgradeHandler,
}
