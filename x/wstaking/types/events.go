package types

// sudo module event types
const (
	EventTypeAdminUpdated   = "admin_updated"
	EventTypeUnMeidDelegate = "unmeid_delegate"
	EventTypeUnDelegate     = "undelegate"

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
	AttributeKeyAccount                  = "account"
	AttributeKeyTerm                     = "fixed_deposit_cofig_term"
	AttributeKeyRate                     = "fixed_deposit_config_rate"
	AttributeKeyStatus                   = "fixed_deposit_config_status"

	EventTypeAddFixedDepositCfg       = "add_fixed_deposit_cfg"
	EventTypeRemoveFixedDepositCfg    = "remove_fixed_deposit_cfg"
	EventTypeSetFixedDepositCfgStatus = "set_fixed_deposit_cfg_status"
	EventTypeSetFixedDepositCfgRate   = "set_fixed_deposit_cfg_rate"
)
