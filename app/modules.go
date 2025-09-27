package app

import (
	"fmt"
	"github.com/CosmWasm/wasmd/x/wasm"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	"github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcclientclient "github.com/cosmos/ibc-go/v7/modules/core/02-client/client"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	"github.com/evmos/ethermint/x/evm"
	evmclient "github.com/evmos/ethermint/x/evm/client"
	"github.com/evmos/ethermint/x/feemarket"
	"github.com/osmosis-labs/osmosis/v15/x/epochs"
	"github.com/osmosis-labs/osmosis/v15/x/gamm"
	"github.com/osmosis-labs/osmosis/v15/x/lockup"
	"github.com/osmosis-labs/osmosis/v15/x/poolmanager"
	"github.com/osmosis-labs/osmosis/v15/x/txfees"
	"github.com/st-chain/me-hub/x/bsc"
	"github.com/st-chain/me-hub/x/dao"
	"github.com/st-chain/me-hub/x/delayedack"
	"github.com/st-chain/me-hub/x/denommetadata"
	denommetadatamoduleclient "github.com/st-chain/me-hub/x/denommetadata/client"
	did "github.com/st-chain/me-hub/x/did"
	"github.com/st-chain/me-hub/x/eibc"
	kyc "github.com/st-chain/me-hub/x/kyc"
	groupmodule "github.com/st-chain/me-hub/x/megroup"
	"github.com/st-chain/me-hub/x/rollapp"
	"github.com/st-chain/me-hub/x/wbank"
	"github.com/st-chain/me-hub/x/wdistri"
	"github.com/st-chain/me-hub/x/wgov"
	"github.com/st-chain/me-hub/x/wmint"
	"github.com/st-chain/me-hub/x/wnft"
	"github.com/st-chain/me-hub/x/wstaking"

	rollappmoduleclient "github.com/st-chain/me-hub/x/rollapp/client"
	"github.com/st-chain/me-hub/x/sequencer"
)

// ModuleBasics defines the module BasicManager is in charge of setting up basic,
// non-dependant module elements, such as codec registration
// and genesis verification.
var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	authzmodule.AppModuleBasic{},
	genutil.NewAppModuleBasic(GenTxMessageValidator),
	wbank.AppModuleBasic{},
	capability.AppModuleBasic{},
	consensus.AppModuleBasic{},
	wstaking.AppModuleBasic{},
	wmint.AppModuleBasic{},
	wdistri.AppModuleBasic{},
	wgov.NewAppModuleBasic([]client.ProposalHandler{
		paramsclient.ProposalHandler,
		upgradeclient.LegacyProposalHandler,
		upgradeclient.LegacyCancelProposalHandler,
		ibcclientclient.UpdateClientProposalHandler,
		ibcclientclient.UpgradeProposalHandler,
		rollappmoduleclient.SubmitFraudHandler,
		denommetadatamoduleclient.CreateDenomMetadataHandler,
		denommetadatamoduleclient.UpdateDenomMetadataHandler,
		evmclient.UpdateVirtualFrontierBankContractProposalHandler,
	}),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	feegrantmodule.AppModuleBasic{},
	ibc.AppModuleBasic{},
	ibctm.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	transfer.AppModuleBasic{},
	vesting.AppModuleBasic{},
	rollapp.AppModuleBasic{},
	sequencer.AppModuleBasic{},
	denommetadata.AppModuleBasic{},
	packetforward.AppModuleBasic{},
	delayedack.AppModuleBasic{},
	eibc.AppModuleBasic{},

	// Ethermint modules
	evm.AppModuleBasic{},
	feemarket.AppModuleBasic{},

	// did modules
	did.AppModuleBasic{},
	kyc.AppModuleBasic{},

	// Osmosis modules
	lockup.AppModuleBasic{},
	epochs.AppModuleBasic{},
	gamm.AppModuleBasic{},
	poolmanager.AppModuleBasic{},
	txfees.AppModuleBasic{},
	dao.AppModuleBasic{},
	wnft.AppModuleBasic{},
	wasm.AppModuleBasic{},
	groupmodule.AppModuleBasic{},
	bsc.AppModuleBasic{},
)

func GenTxMessageValidator(msgs []sdk.Msg) error {
	if len(msgs) == 0 {
		return fmt.Errorf("unexpected number of GenTx messages; got: %d, expected great than 0", len(msgs))
	}
	if _, ok := msgs[0].(*stakingtypes.MsgCreateValidator); !ok {
		return fmt.Errorf("unexpected GenTx message type; expected: MsgCreateValidator, got: %T", msgs[0])
	}
	return nil
}
