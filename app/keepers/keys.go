package keepers

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	nftkeeper "github.com/cosmos/cosmos-sdk/x/nft/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	bsctypes "github.com/openmetaearth/me-hub/x/bsc/types"
	daotypes "github.com/openmetaearth/me-hub/x/dao/types"
	delayedacktypes "github.com/openmetaearth/me-hub/x/delayedack/types"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	eibcmoduletypes "github.com/openmetaearth/me-hub/x/eibc/types"
	kyctypes "github.com/openmetaearth/me-hub/x/kyc/types"
	gourpTypes "github.com/openmetaearth/me-hub/x/megroup/types"
	rollappmoduletypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	sequencermoduletypes "github.com/openmetaearth/me-hub/x/sequencer/types"
	trontypes "github.com/openmetaearth/me-hub/x/tron/types"
)

// GenerateKeys generates new keys (KV Store, Transient store, and memory store).
func (a *AppKeepers) GenerateKeys() {
	// Define what keys will be used in the cosmos-sdk key/value store.
	// Cosmos-SDK modules each have a "key" that allows the application to reference what they've stored on the chain.
	a.keys = KVStoreKeys

	// Define transient store keys
	a.tkeys = sdk.NewTransientStoreKeys(paramstypes.TStoreKey, evmtypes.TransientKey, feemarkettypes.TransientKey)

	// MemKeys are for information that is stored only in RAM.
	a.memKeys = sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
}

// GetSubspace gets existing substore from keeper.
func (a *AppKeepers) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := a.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// GetKVStoreKeys gets KV Store keys.
func (a *AppKeepers) GetKVStoreKeys() map[string]*storetypes.KVStoreKey {
	return a.keys
}

// GetTransientStoreKey gets Transient Store keys.
func (a *AppKeepers) GetTransientStoreKey() map[string]*storetypes.TransientStoreKey {
	return a.tkeys
}

// GetMemoryStoreKey get memory Store keys.
func (a *AppKeepers) GetMemoryStoreKey() map[string]*storetypes.MemoryStoreKey {
	return a.memKeys
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (a *AppKeepers) GetKey(storeKey string) *storetypes.KVStoreKey {
	return a.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (a *AppKeepers) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return a.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (a *AppKeepers) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return a.memKeys[storeKey]
}

var KVStoreKeys = sdk.NewKVStoreKeys(
	authtypes.StoreKey,
	authzkeeper.StoreKey,
	banktypes.StoreKey,
	stakingtypes.StoreKey,
	minttypes.StoreKey,
	distrtypes.StoreKey,
	slashingtypes.StoreKey,
	govtypes.StoreKey,
	paramstypes.StoreKey,
	ibcexported.StoreKey,
	upgradetypes.StoreKey,
	feegrant.StoreKey,
	evidencetypes.StoreKey,
	ibctransfertypes.StoreKey,
	capabilitytypes.StoreKey,
	crisistypes.StoreKey,
	consensusparamtypes.StoreKey,
	rollappmoduletypes.StoreKey,
	sequencermoduletypes.StoreKey,
	packetforwardtypes.StoreKey,
	delayedacktypes.StoreKey,
	eibcmoduletypes.StoreKey,
	// ethermint keys
	evmtypes.StoreKey,
	feemarkettypes.StoreKey,
	// did keys
	didtypes.StoreKey,
	kyctypes.StoreKey,
	// me keys
	daotypes.StoreKey,
	nftkeeper.StoreKey,
	wasmtypes.StoreKey,
	//megroup
	gourpTypes.StoreKey,
	// gravity bridge
	bsctypes.StoreKey,
	trontypes.StoreKey,
)
