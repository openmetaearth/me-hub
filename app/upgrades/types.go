package upgrades

import (
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	feemarketkeeper "github.com/evmos/ethermint/x/feemarket/keeper"
	"github.com/spf13/cobra"

	delayedackkeeper "github.com/openmetaearth/me-hub/x/delayedack/keeper"
	eibckeeper "github.com/openmetaearth/me-hub/x/eibc/keeper"
	evmkeeper "github.com/openmetaearth/me-hub/x/evm/keeper"
	gravitykeeper "github.com/openmetaearth/me-hub/x/gravity/keeper"
	lightclientkeeper "github.com/openmetaearth/me-hub/x/lightclient/keeper"
	rollappkeeper "github.com/openmetaearth/me-hub/x/rollapp/keeper"
	sequencerkeeper "github.com/openmetaearth/me-hub/x/sequencer/keeper"
	wbankkeeper "github.com/openmetaearth/me-hub/x/wbank/keeper"
	wgovkeeper "github.com/openmetaearth/me-hub/x/wgov/keeper"
	wmintkeeper "github.com/openmetaearth/me-hub/x/wmint/keeper"
	wstakingkeeper "github.com/openmetaearth/me-hub/x/wstaking/keeper"
)

// Keepers contains only the keepers required by legacy upgrade handlers. It
// keeps the upgrade package independent from package app and avoids an import
// cycle now that AppKeepers lives in package app.
type Keepers struct {
	BankKeeper      wbankkeeper.BaseKeeperWrapper
	GovKeeper       *wgovkeeper.Keeper
	StakingKeeper   *wstakingkeeper.Keeper
	EvmKeeper       *evmkeeper.Keeper
	FeeMarketKeeper feemarketkeeper.Keeper
	BscKeeper       gravitykeeper.Keeper
	TronKeeper      gravitykeeper.Keeper
}

// UpgradeKeepers contains the keepers required by the v3 migration handler.
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

// BaseAppParamManager defines an interface that BaseApp is expected to fulfill
// that allows upgrade handlers to modify BaseApp parameters.
type BaseAppParamManager interface {
	GetConsensusParams(ctx sdk.Context) cometbftproto.ConsensusParams
	StoreConsensusParams(ctx sdk.Context, cp cometbftproto.ConsensusParams) error
}

// Upgrade defines a struct containing necessary fields that a SoftwareUpgradeProposal
// must have written, in order for the state migration to go smoothly.
// An upgrade must implement this struct, and then set it in the app.go.
// The app.go will then define the handler.
type Upgrade struct {
	// Upgrade version name, for the upgrade handler, e.g. `v4`
	Name string

	// CreateHandler defines the function that creates an upgrade handler
	CreateHandler func(*module.Manager, module.Configurator, BaseAppParamManager, *Keepers) upgradetypes.UpgradeHandler

	// CreateHandlerV3 defines a handler using the keeper set required by v3.
	CreateHandlerV3 func(*module.Manager, module.Configurator, *UpgradeKeepers) upgradetypes.UpgradeHandler

	// Store upgrades, should be used for any new modules introduced, new modules deleted, or store names renamed.
	StoreUpgrades storetypes.StoreUpgrades

	PreUpgradeCmd *cobra.Command
}
