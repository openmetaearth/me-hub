package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var _ RegionI = &Region{}

type RegionI interface {
	GetRegionId() string
	GetCreator() string
	GetName() string
	GetOperatorAddress() string
	GetNftClassId() string
	GetRegionTreasureAddr() string
	GetDepositInterestAddr() string
	GetRegionShare() sdk.Int
	GetDelegateInterest() sdk.Dec
	GetDelegateAmount() sdk.Int
}

func (m *Region) GetRegionShare() sdk.Int {
	return m.RegionShare
}

func (m *Region) GetDelegateInterest() sdk.Dec {
	return m.DelegateInterest
}

func (m *Region) GetDelegateAmount() sdk.Int {
	return m.DelegateAmount
}
