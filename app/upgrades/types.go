package upgrades

import (
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	delayedackkeeper "github.com/openmetaearth/me-hub/x/delayedack/keeper"
	eibckeeper "github.com/openmetaearth/me-hub/x/eibc/keeper"
	lightclientkeeper "github.com/openmetaearth/me-hub/x/lightclient/keeper"
	rollappkeeper "github.com/openmetaearth/me-hub/x/rollapp/keeper"
	sequencerkeeper "github.com/openmetaearth/me-hub/x/sequencer/keeper"
	wgovkeeper "github.com/openmetaearth/me-hub/x/wgov/keeper"
	wmintkeeper "github.com/openmetaearth/me-hub/x/wmint/keeper"
	wstakingkeeper "github.com/openmetaearth/me-hub/x/wstaking/keeper"
)

// UpgradeKeepers contains the keepers required by upgrade handlers.
type UpgradeKeepers struct {
	AccountKeeper     *authkeeper.AccountKeeper
	GovKeeper         *wgovkeeper.Keeper
	RollappKeeper     *rollappkeeper.Keeper
	SequencerKeeper   *sequencerkeeper.Keeper
	ParamsKeeper      *paramskeeper.Keeper
	DelayedAckKeeper  *delayedackkeeper.Keeper
	EIBCKeeper        *eibckeeper.Keeper
	LightClientKeeper *lightclientkeeper.Keeper
	IBCKeeper         *ibckeeper.Keeper
	MintKeeper        *wmintkeeper.Keeper
	SlashingKeeper    *slashingkeeper.Keeper
	ConsensusKeeper   *consensusparamkeeper.Keeper
	StakingKeeper     *wstakingkeeper.Keeper
}

// Upgrade defines a struct containing necessary fields that a SoftwareUpgradeProposal
// must have written, in order for the state migration to go smoothly.
// An upgrade must implement this struct, and then set it in the app.go.
// The app.go will then define the handler.
type Upgrade struct {
	// Upgrade version name, for the upgrade handler, e.g. `v4`
	Name string

	// CreateHandler defines the function that creates an upgrade handler.
	CreateHandler func(*module.Manager, module.Configurator, *UpgradeKeepers) upgradetypes.UpgradeHandler

	// Store upgrades, should be used for any new modules introduced, new modules deleted, or store names renamed.
	StoreUpgrades storetypes.StoreUpgrades
}
