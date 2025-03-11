package v2_0_10

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	epochstypes "github.com/osmosis-labs/osmosis/v15/x/epochs/types"
	gammtypes "github.com/osmosis-labs/osmosis/v15/x/gamm/types"
	lockuptypes "github.com/osmosis-labs/osmosis/v15/x/lockup/types"
	poolmanagertypes "github.com/osmosis-labs/osmosis/v15/x/poolmanager/types"
	txfeestypes "github.com/osmosis-labs/osmosis/v15/x/txfees/types"
	"github.com/st-chain/me-hub/app/upgrades"
	daotypes "github.com/st-chain/me-hub/x/dao/types"
	delayedacktypes "github.com/st-chain/me-hub/x/delayedack/types"
	didtypes "github.com/st-chain/me-hub/x/did/types"
	eibctypes "github.com/st-chain/me-hub/x/eibc/types"
	incentivestypes "github.com/st-chain/me-hub/x/incentives/types"
	kyctypes "github.com/st-chain/me-hub/x/kyc/types"
	megrouptypes "github.com/st-chain/me-hub/x/megroup/types"
	rollappmoduletypes "github.com/st-chain/me-hub/x/rollapp/types"
	sequencermoduletypes "github.com/st-chain/me-hub/x/sequencer/types"
	streamermoduletypes "github.com/st-chain/me-hub/x/streamer/types"
)

const (
	UpgradeName       = "v2.0.10"
	didFilePath       = "/config/did.json"
	issuerFilePath    = "/config/issuer.json"
	kycPubkeyFilePath = "/config/kyc_pubkey.json"
	groupFilePath     = "/config/genesis1.3.json" // need export 1.3 genesis
)

var Upgrade = upgrades.Upgrade{
	Name:          UpgradeName,
	CreateHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added: []string{
			rollappmoduletypes.ModuleName,
			sequencermoduletypes.ModuleName,
			streamermoduletypes.ModuleName,
			packetforwardtypes.ModuleName,
			delayedacktypes.ModuleName,
			eibctypes.ModuleName,
			// ethermint keys
			evmtypes.ModuleName,
			feemarkettypes.ModuleName,
			// osmosis keys
			lockuptypes.ModuleName,
			epochstypes.ModuleName,
			gammtypes.ModuleName,
			poolmanagertypes.ModuleName,
			incentivestypes.ModuleName,
			txfeestypes.ModuleName,
			// me keys
			daotypes.ModuleName,
			didtypes.ModuleName,
			kyctypes.ModuleName,
			megrouptypes.ModuleName,
		},
		Deleted: []string{
			"checkin",
		},
	},
	PreUpgradeCmd: PreUpgradeCmd(),
}
