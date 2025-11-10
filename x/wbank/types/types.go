package types

const TreasuryPoolName = "treasury_pool"
const EventTypeFeeToReceivers = "fee_to_receivers"

type FeeReceiverType string

const (
	FeeReceiverDevOperator      FeeReceiverType = "dev_operator"
	FeeReceiverProposerOwner    FeeReceiverType = "block_proposer_owner"
	FeeReceiverKycRegionOwner   FeeReceiverType = "kyc_region_owner"
	FeeReceiverGlobalDaoFeePool FeeReceiverType = "global_dao_fee_pool"
	FeeReceiverContractCreator  FeeReceiverType = "contract_creator"
)
