package types

const (
	NotBondedPoolName         = "not_bonded_tokens_pool"
	NotBondedStakePoolName    = "not_bonded_stake_tokens_pool"
	BondedStakePoolName       = "bonded_stake_tokens_pool"
	StakePoolName             = "stake_tokens_pool"
	MeidNFTPoolName           = "meid_nft_pool"
	FixedDepositPrincipalPool = "fixed_deposit_principal_pool"
	GlobalDaoFeePool          = "global_admin_fee_pool"
	BridgeFeePool             = "bridge_fee_pool"
)

// IsAllowedSendToModuleTarget returns true for module accounts the GlobalDao may
// fund through MsgSendToModule. Accounting and reward pools such as
// GlobalDaoFeePool, bonded staking pools, and distribution modules are excluded
// because direct deposits bypass their module-specific accounting flows.
func IsAllowedSendToModuleTarget(moduleName string) bool {
	switch moduleName {
	case StakePoolName, FixedDepositPrincipalPool, BridgeFeePool:
		return true
	default:
		return false
	}
}
