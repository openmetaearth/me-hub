package types

const (
	EventTypeBondedRelayer          = "bonded_relayer"
	EventTypeAddDelegate            = "add_delegate"
	EventTypeUnBondedRelayer        = "unbonded_relayer"
	EventTypeRelayerSetConfirm      = "relayer_set_confirm"
	EventTypeSendToExternal         = "send_to_external"
	EventTypeSendToExternalCanceled = "send_to_external_canceled"
	EventTypeRelayerSetUpdate       = "relayer_set_update"
	EventTypeContractEvent          = "contract_event_claim"
	EventTypeOutgoingBatch          = "outgoing_batch"
	EventTypeOutgoingBatchConfirm   = "outgoing_batch_confirm"
	EventTypeOutgoingBatchCanceled  = "outgoing_batch_canceled"
	EventTypeIncreaseBridgeFee      = "increase_bridge_fee"
)

const (
	AttributeKeyReceiver             = "receiver"
	AttributeKeySlashAmount          = "slash_amount"
	AttributeKeyUnbondAmount         = "unbond_amount"
	AttributeKeyExternalAddress      = "external_address"
	AttributeKeyOutgoingTxID         = "outgoing_tx_id"
	AttributeKeyOutgoingTxIds        = "outgoing_tx_ids"
	AttributeKeyOutgoingBatchNonce   = "batch_nonce"
	AttributeKeyTokenContract        = "token_contract"
	AttributeKeyOutgoingBatchTimeout = "outgoing_batch_timeout"
	AttributeKeyRelayerSetNonce      = "relayer_set_nonce"
	AttributeKeyRelayerSetLen        = "relayer_set_len"
	AttributeKeyClaimType            = "claim_type"
	AttributeKeyEventNonce           = "event_nonce"
	AttributeKeyClaimHash            = "claim_hash"
	AttributeKeyBlockHeight          = "block_height"
	AttributeKeyStateSuccess         = "state_success"
	AttributeKeyIncreaseFee          = "increase_fee"
	AttributeKeyBridgeFee            = "bridge_fee"
	AttributeKeyRefundAmount         = "refund_amount"
)
