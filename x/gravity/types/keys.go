package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "gravity"

	// RouterKey is the module name router key
	RouterKey = ModuleName

	SlashingModuleAccount = "slashing"
)

var (
	RelayerKey                  = []byte{0x12}
	RelayerAddressByExternalKey = []byte{0x13}

	// RelayerSetRequestKey indexes relayer set requests by nonce
	RelayerSetRequestKey = []byte{0x15}

	// RelayerSetConfirmKey indexes relayer set confirmations by nonce and the validator account address
	RelayerSetConfirmKey = []byte{0x16}

	// RelayerAttestationKey attestation details by nonce and validator address
	// An attestation can be thought of as the 'event to be executed' while
	// the Claims are an individual validator saying that they saw an event
	// occur the Attestation is 'the event' that multiple claims vote on and
	// eventually executes
	RelayerAttestationKey = []byte{0x17}

	// OutgoingTxPoolKey indexes the last nonce for the outgoing tx pool
	OutgoingTxPoolKey = []byte{0x18}

	// OutgoingTxBatchKey indexes outgoing tx batches under a nonce and token address
	OutgoingTxBatchKey = []byte{0x20}

	// OutgoingTxBatchBlockKey indexes outgoing tx batches under a block height and token address
	OutgoingTxBatchBlockKey = []byte{0x21}

	// BatchConfirmKey indexes relayer confirmations by token contract address
	BatchConfirmKey = []byte{0x22}

	// LastEventNonceByRelayerKey indexes latest event nonce by relayer
	LastEventNonceByRelayerKey = []byte{0x23}

	// LastObservedEventNonceKey indexes the latest event nonce
	LastObservedEventNonceKey = []byte{0x24}

	// SequenceKeyPrefix indexes different txIds
	SequenceKeyPrefix = []byte{0x25}

	// KeyLastTxPoolID indexes the lastTxPoolID
	KeyLastTxPoolID = append(append([]byte(nil), SequenceKeyPrefix...), []byte("lastTxPoolId")...)

	// KeyLastOutgoingBatchID indexes the lastBatchID
	KeyLastOutgoingBatchID = append(append([]byte(nil), SequenceKeyPrefix...), []byte("lastBatchId")...)

	// BridgeTokenByContract prefixes the index of asset denom to external token
	BridgeTokenByContractKey = []byte{0x26}

	// BridgeTokenByDenom prefixes the index of assets external token to denom
	BridgeTokenByDenomKey = []byte{0x27}

	// LastSlashedRelayerSetNonce indexes the latest slashed relayerSet nonce
	LastSlashedRelayerSetNonce = []byte{0x28}

	// LatestRelayerSetNonce indexes the latest relayerSet nonce
	LatestRelayerSetNonce = []byte{0x29}

	// LastSlashedBatchBlock indexes the latest slashed batch block height
	LastSlashedBatchBlock = []byte{0x30}

	// Deprecated: LastProposalBlockHeight
	// LastProposalBlockHeight = []byte{0x31}

	// LastObservedBlockHeightKey indexes the latest observed external block height
	LastObservedBlockHeightKey = []byte{0x32}

	// LastObservedRelayerSetKey indexes the latest observed RelayerSet nonce
	LastObservedRelayerSetKey = []byte{0x33}

	// LastEventBlockHeightByRelayerKey indexes latest event blockHeight by relayer
	LastEventBlockHeightByRelayerKey = []byte{0x35}

	// Deprecated: PastExternalSignatureCheckpointKey indexes eth signature checkpoints that have existed
	PastExternalSignatureCheckpointKey = []byte{0x36}

	// LastRelayerSlashBlockHeight indexes the last relayer slash block height
	LastRelayerSlashBlockHeight = []byte{0x37}

	// ProposalRelayerKey -> value ProposalRelayer
	ProposalRelayerKey = []byte{0x38}

	// LastTotalPowerKey relayer set total power
	LastTotalPowerKey = []byte{0x39}

	// ParamsKey is the prefix for params key
	ParamsKey = []byte{0x40}

	// OutgoingTxRelationKey outgoing tx with evm
	OutgoingTxRelationKey = []byte{0x41}
)

// GetRelayerKey returns the following key format
func GetRelayerKey(relayer sdk.AccAddress) []byte {
	return append(RelayerKey, relayer.Bytes()...)
}

// GetRelayerAddressByExternalKey returns the following key format
func GetRelayerAddressByExternalKey(externalAddress string) []byte {
	return append(RelayerAddressByExternalKey, []byte(externalAddress)...)
}

// GetRelayerSetKey returns the following key format
func GetRelayerSetKey(nonce uint64) []byte {
	return append(RelayerSetRequestKey, sdk.Uint64ToBigEndian(nonce)...)
}

// GetRelayerSetConfirmKey returns the following key format
func GetRelayerSetConfirmKey(nonce uint64, relayerAddr sdk.AccAddress) []byte {
	return append(RelayerSetConfirmKey, append(sdk.Uint64ToBigEndian(nonce), relayerAddr.Bytes()...)...)
}

// GetAttestationKey returns the following key format
// An attestation is an event multiple people are voting on, this function needs the claim
// details because each Attestation is aggregating all claims of a specific event, lets say
// validator X and validator y where making different claims about the same event nonce
// Note that the claim hash does NOT include the claimer address and only identifies an event
func GetAttestationKey(eventNonce uint64, claimHash []byte) []byte {
	return append(RelayerAttestationKey, append(sdk.Uint64ToBigEndian(eventNonce), claimHash...)...)
}

func GetAttestationKeyByNonce(eventNonce uint64) []byte {
	return append(RelayerAttestationKey, append(sdk.Uint64ToBigEndian(eventNonce))...)
}

// GetOutgoingTxPoolContractPrefix returns the following key format
// This prefix is used for iterating over unbatched transactions for a given contract
func GetOutgoingTxPoolContractPrefix(tokenContract string) []byte {
	return append(OutgoingTxPoolKey, []byte(tokenContract)...)
}

// GetOutgoingTxPoolKey returns the following key format
func GetOutgoingTxPoolKey(fee ERC20Token, id uint64) []byte {
	amount := make([]byte, 32)
	amount = fee.Amount.BigInt().FillBytes(amount)
	return append(OutgoingTxPoolKey, append([]byte(fee.Contract), append(amount, sdk.Uint64ToBigEndian(id)...)...)...)
}

// GetOutgoingTxBatchKey returns the following key format
func GetOutgoingTxBatchKey(tokenContract string, batchNonce uint64) []byte {
	return append(append(OutgoingTxBatchKey, []byte(tokenContract)...), sdk.Uint64ToBigEndian(batchNonce)...)
}

// GetOutgoingTxBatchBlockKey returns the following key format
func GetOutgoingTxBatchBlockKey(blockHeight uint64) []byte {
	return append(OutgoingTxBatchBlockKey, sdk.Uint64ToBigEndian(blockHeight)...)
}

// GetBatchConfirmKey returns the following key format
func GetBatchConfirmKey(tokenContract string, batchNonce uint64, relayerAddr sdk.AccAddress) []byte {
	return append(BatchConfirmKey, append([]byte(tokenContract), append(sdk.Uint64ToBigEndian(batchNonce), relayerAddr.Bytes()...)...)...)
}

// GetLastEventNonceByRelayerKey returns the following key format
func GetLastEventNonceByRelayerKey(relayerAddr sdk.AccAddress) []byte {
	return append(LastEventNonceByRelayerKey, relayerAddr.Bytes()...)
}

// GetLastEventBlockHeightByRelayerKey returns the following key format
func GetLastEventBlockHeightByRelayerKey(relayerAddr sdk.AccAddress) []byte {
	return append(LastEventBlockHeightByRelayerKey, relayerAddr.Bytes()...)
}

// GetBridgeTokenByContract returns the following key format
func GetBridgeTokenByContractKey(tokenContract string) []byte {
	return append(BridgeTokenByContractKey, []byte(tokenContract)...)
}

// GetBridgeTokenByDenom returns the following key format
func GetBridgeTokenByDenomKey(denom string) []byte {
	return append(BridgeTokenByDenomKey, []byte(denom)...)
}

// GetPastExternalSignatureCheckpointKey returns the following key format
func GetPastExternalSignatureCheckpointKey(blockHeight uint64, checkpoint []byte) []byte {
	return append(PastExternalSignatureCheckpointKey, append(sdk.Uint64ToBigEndian(blockHeight), checkpoint...)...)
}

func GetOutgoingTxRelationKey(txID uint64) []byte {
	return append(OutgoingTxRelationKey, sdk.Uint64ToBigEndian(txID)...)
}
