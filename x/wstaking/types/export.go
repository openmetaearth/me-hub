package types

import github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"

var _ RegionI = &Region{}

type RegionI interface {
	GetRegionId() string
	GetCreator() string
	GetName() string
	GetOperatorAddress() string
	GetNftClassId() string
	GetRegionTreasureAddr() string
	GetDepositInterestAddr() string
	GetRegionShare() github_com_cosmos_cosmos_sdk_types.Int
	GetDelegateInterest() github_com_cosmos_cosmos_sdk_types.Dec
	GetDelegateAmount() github_com_cosmos_cosmos_sdk_types.Int
}

func (m *Region) GetRegionShare() github_com_cosmos_cosmos_sdk_types.Int {
	return m.RegionShare
}
func (m *Region) GetDelegateInterest() github_com_cosmos_cosmos_sdk_types.Dec {
	return m.DelegateInterest
}
func (m *Region) GetDelegateAmount() github_com_cosmos_cosmos_sdk_types.Int {
	return m.DelegateAmount
}
