package types

import sdkmath "cosmossdk.io/math"

var _ RegionI = &Region{}

type RegionI interface {
	GetRegionId() string
	GetCreator() string
	GetName() string
	GetOperatorAddress() string
	GetNftClassId() string
	GetRegionTreasureAddr() string
	GetDepositInterestAddr() string
	GetRegionShare() sdkmath.Int
	GetDelegateInterest() sdkmath.LegacyDec
	GetDelegateAmount() sdkmath.Int
}

func (m *Region) GetRegionShare() sdkmath.Int {
	return m.RegionShare
}
func (m *Region) GetDelegateInterest() sdkmath.LegacyDec {
	return m.DelegateInterest
}
func (m *Region) GetDelegateAmount() sdkmath.Int {
	return m.DelegateAmount
}
