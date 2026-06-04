package types

const (
	EventTypeStateUpdate  = "state_update"
	EventTypeStatusChange = "status_change"

	AttributeKeyRollappId      = "rollapp_id"
	AttributeKeyStateInfoIndex = "state_info_index"
	AttributeKeyStartHeight    = "start_height"
	AttributeKeyNumBlocks      = "num_blocks"
	AttributeKeyDAPath         = "da_path"
	AttributeKeyStatus         = "status"

	// EventTypeFraud is emitted when a fraud evidence is submitted
	EventTypeFraud             = "fraud_proposal"
	AttributeKeyFraudHeight    = "fraud_height"
	AttributeKeyFraudSequencer = "fraud_sequencer"
	AttributeKeyClientID       = "client_id"

	// EventTypeTransferGenesisTransfersEnabled is when the bridge is enabled
	EventTypeTransferGenesisTransfersEnabled = "transfer_genesis_transfers_enabled"

	// EventTypeDRSViolation is emitted when a rollapp violates the genesis transfer window.
	EventTypeDRSViolation         = "drs_violation"
	AttributeKeyDRSViolationCause = "drs_violation_cause"
)
