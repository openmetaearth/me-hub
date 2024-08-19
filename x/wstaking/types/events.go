package types

const (
	EventTypeUnstake                = "unstake"
	EventTypeStake                  = "stake"
	EventTypeCompleteUnStakeBonding = "complete_stake_unbonding"
	EventTypeCompleteUnDelBonding   = "complete_del_unbonding"
)

const (
	AttributeKeyValidator      = "validator"
	AttributeKeyRegionId       = "regionId"
	AttributeKeyNewShares      = "new_shares"
	AttributeKeyCompletionTime = "completion_time"
	AttributeKeyDelegator      = "delegator"
	AttributeKeyFromAddress    = "from_address"
)
