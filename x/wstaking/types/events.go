package types

const (
	EventTypeUnstake                      = "unstake"
	EventTypeStake                        = "stake"
	EventTypeCompleteUnStakeBonding       = "complete_stake_unbonding"
	EventTypeCompleteUnDelBonding         = "complete_del_unbonding"
	EventTypeSettleDelRewardsForKyc       = "settle_del_rewards_for_kyc"
	EventTypeWithdrawFromRegion           = "withdraw_from_region"
	EventTypeWithdrawFromGlobalDaoFeePool = "withdraw_from_global_dao_fee_pool"

	EventTypeMeidNew    = "meid_new"
	EventTypeMeidRemove = "meid_remove"
)

const (
	AttributeKeyValidator                = "validator"
	AttributeKeyRegionId                 = "regionId"
	AttributeKeyNewShares                = "new_shares"
	AttributeKeyCompletionTime           = "completion_time"
	AttributeKeyDelegator                = "delegator"
	AttributeKeyTotalAmountDelegate      = "total_amount_delegate"
	AttributeKeyAmountDelegateInterest   = "amount_of_delegate_interest"
	AttributeKeyRegionTreasure           = "region_treasure"
	AttributeKeyDelegatorAddress         = "delegator_address"
	AttributeKeyPersonalDelegateInterest = "personal_delegate_interest"
	AttributeKeyIsMeid                   = "is_meid"

	AttributeKeyExpired      = "fixed_expired"
	AttributeKeyAccount      = "account"
	AttributeKeyCreator      = "creator"
	AttributeKeyExecTime     = "exec_time"
	AttributeKeyInterest     = "fixed_deposit_interest"
	AttributeKeyInterestAddr = "fixed_deposit_interest_address"
	AttributeKeyTreasureAddr = "fixed_deposit_treasure_address"

	EventTransferRegion         = "transfer_region"
	AttributeKeyFromAddress     = "from_address"
	AttributeKeyTransferAddress = "transfer_address"
	AttributeKeyFromRegion      = "from_region"
	AttributeKeyToRegion        = "to_region"

	AttributeKeyReceiver  = "receiver"
	AttributeKeyReceiver2 = "receiver2"
	AttributeKeyReceiver3 = "receiver3"
	AttributeKeyReceiver4 = "receiver4"

	AttributeKeyAmount  = "amount"
	AttributeKeyAmount2 = "amount2"
	AttributeKeyAmount3 = "amount3"
	AttributeKeyAmount4 = "amount4"

	AttributeKeyMeidInviteAddress                = "meid_invite_address"
	AttributeKeyMeidInviteReward                 = "meid_invite_reward"
	AttributeKeySendMeidInviteAddress            = "send_meid_invite_reward_address"
	AttributeKeyReceiveMeidInviteAddress_Society = "society_receive_meid_invite_reward_address"
	AttributeKeyReceiveMeidInviteAddress_Node    = "node_receive_meid_invite_reward_address"
	AttributeKeyMeidNumAddReward                 = "meid_unmber_add_reward"

	AttributeKeyNewOwnerAddress = "new_owner"

	AttributeKeyUnKycLockTime = "unkyc_lock_time"
)
