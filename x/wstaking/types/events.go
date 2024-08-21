package types

const (
	EventTypeUnstake                  = "unstake"
	EventTypeStake                    = "stake"
	EventTypeCompleteUnStakeBonding   = "complete_stake_unbonding"
	EventTypeCompleteUnDelBonding     = "complete_del_unbonding"
	EventTypeAddFixedDepositCfg       = "add_fixed_deposit_cfg"
	EventTypeRemoveFixedDepositCfg    = "remove_fixed_deposit_cfg"
	EventTypeSetFixedDepositCfgStatus = "set_fixed_deposit_cfg_status"
	EventTypeSetFixedDepositCfgRate   = "set_fixed_deposit_cfg_rate"
)

const (
	AttributeKeyCompletionTime = "completion_time"
	AttributeKeyFromAddress    = "from_address"
	AttributeKeyAccount        = "account"
	AttributeKeyTerm           = "fixed_deposit_cofig_term"
	AttributeKeyRate           = "fixed_deposit_config_rate"
	AttributeKeyStatus         = "fixed_deposit_config_status"
	EventTypeAdminUpdated      = "admin_updated"
	EventTypeUnMeidDelegate    = "unmeid_delegate"
	EventTypeUnDelegate        = "undelegate"

	AttributeKeyLastAdmin                = "last_admin"
	AttributeKeyCurrentAdmin             = "current_admin"
	AttributeKeyValidator                = "validator"
	AttributeKeyRewards                  = "rewards"
	AttributeKeyRegionTreasure           = "region_treasure"
	AttributeKeyNewShares                = "new_shares"
	AttributeKeyTotalAmountDelegate      = "total_amount_delegate"
	AttributeKeyRegionId                 = "regionId"
	AttributeKeyAmountDelegateInterest   = "amount_of_delegate_interest"
	AttributeKeyDelegatorAddress         = "delegator_address"
	AttributeKeyPersonalDelegateInterest = "personal_delegate_interest"
	AttributeKeyIsMeid                   = "is_meid"

	EventTypeSetWithdrawAddress = "set_withdraw_address"
	EventTypeRewards            = "rewards"
	EventTypeCommission         = "commission"
	EventTypeWithdrawRewards    = "withdraw_rewards"
	EventTypeWithdrawCommission = "withdraw_commission"
	EventTypeProposerReward     = "proposer_reward"

	EventTypeWithdrawDelegatorReward = "withdraw_delegator_reward"
	EventTypeRegionTreasuryReword    = "region_treasury_reword"

	AttributeKeyWithdrawAddress       = "withdraw_address"
	AttributeKeyDelegator             = "delegator"
	AttributeValueCategory            = ModuleName
	AttributeKeyRegionTreasuryAddress = "region_treasury_address"
)
