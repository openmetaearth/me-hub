package upgrades

import (
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	delayedackkeeper "github.com/st-chain/me-hub/x/delayedack/keeper"
	eibckeeper "github.com/st-chain/me-hub/x/eibc/keeper"
	rollappkeeper "github.com/st-chain/me-hub/x/rollapp/keeper"
	sequencerkeeper "github.com/st-chain/me-hub/x/sequencer/keeper"
	wgovkeeper "github.com/st-chain/me-hub/x/wgov/keeper"
	wmintkeeper "github.com/st-chain/me-hub/x/wmint/keeper"
)

// Upgrade defines a struct containing necessary fields that a SoftwareUpgradeProposal
// must have written, in order for the state migration to go smoothly.
// An upgrade must implement this struct, and then set it in the app.go.
// The app.go will then define the handler.
type Upgrade struct {
	// Upgrade version name, for the upgrade handler, e.g. `v4`
	Name string

	// CreateHandler defines the function that creates an upgrade handler
	CreateHandler func(
		mm *module.Manager,
		configurator module.Configurator,
		appKeepers *UpgradeKeepers,
	) upgradetypes.UpgradeHandler

	// Store upgrades, should be used for any new modules introduced, new modules deleted, or store names renamed.
	StoreUpgrades storetypes.StoreUpgrades
}

type UpgradeKeepers struct {
	AccountKeeper    *authkeeper.AccountKeeper
	GovKeeper        *wgovkeeper.Keeper
	RollappKeeper    *rollappkeeper.Keeper
	SequencerKeeper  *sequencerkeeper.Keeper
	ParamsKeeper     *paramskeeper.Keeper
	DelayedAckKeeper *delayedackkeeper.Keeper
	EIBCKeeper       *eibckeeper.Keeper
	MintKeeper       *wmintkeeper.Keeper
	SlashingKeeper   *slashingkeeper.Keeper
	ConsensusKeeper  *consensusparamkeeper.Keeper
}
