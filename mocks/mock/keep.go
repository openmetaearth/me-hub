package mock

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	wstaking "github.com/openmetaearth/me-hub/x/wstaking/keeper"

	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
)

type MockMeid struct {
	Account    string `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
	Creator    string `protobuf:"bytes,2,opt,name=creator,proto3" json:"creator,omitempty"`
	RegionId   string `protobuf:"bytes,3,opt,name=regionId,proto3" json:"regionId,omitempty"`
	RegionName string `protobuf:"bytes,4,opt,name=regionName,proto3" json:"regionName,omitempty"`
	RewardType int32  `protobuf:"varint,5,opt,name=RewardType,proto3" json:"RewardType,omitempty"`
}
type MockRegion struct {
	RegionId            string `protobuf:"bytes,1,opt,name=regionId,proto3" json:"regionId,omitempty"`
	Name                string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Creator             string `protobuf:"bytes,3,opt,name=creator,proto3" json:"creator,omitempty"`
	OperatorAddress     string `protobuf:"bytes,4,opt,name=operator_address,json=operatorAddress,proto3" json:"operator_address,omitempty"`
	NftClassId          string `protobuf:"bytes,5,opt,name=nft_class_id,json=nftClassId,proto3" json:"nft_class_id,omitempty"`
	RegionTreasureAddr  string `protobuf:"bytes,6,opt,name=region_treasure_addr,json=regionTreasureAddr,proto3" json:"region_treasure_addr,omitempty"`
	DepositInterestAddr string `protobuf:"bytes,8,opt,name=deposit_interest_addr,json=depositInterestAddr,proto3" json:"deposit_interest_addr,omitempty"`
	// tokens define the region tokens share
	RegionShare      github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,9,opt,name=region_share,json=regionShare,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"region_share"`
	DelegateInterest github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,10,opt,name=delegate_interest,json=delegateInterest,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"delegate_interest"`
	DelegateAmount   github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,11,opt,name=delegate_amount,json=delegateAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"delegate_amount"`
}
type MockBankKeeper struct {
	bank.Keeper
}

func (mbk MockBankKeeper) GetTreasuryPoolName() string {
	return "treasury"
}

type MockStakingKeeper struct {
	*wstaking.Keeper
	Meid   map[string]*MockMeid
	Region map[string]*MockRegion
}

func NewMockStakingKeeper(sk *wstaking.Keeper) *MockStakingKeeper {
	ms := &MockStakingKeeper{
		Keeper: sk,
		Meid:   make(map[string]*MockMeid),
		Region: make(map[string]*MockRegion),
	}
	ms.Meid["cosmos1lugrmnrk3ngky85n3hsrxumr3ca7m643h59t72"] = &MockMeid{}
	return ms
}

func (msk MockStakingKeeper) CheckRegionName(name string) (string, error) {
	return "oh yeah", nil
}

func (msk MockStakingKeeper) GetMeid(ctx sdk.Context, account string) (val MockMeid, found bool) {
	v, ok := msk.Meid[account]
	if !ok {
		return MockMeid{}, false
	}
	return *v, ok
}

func (msk *MockStakingKeeper) SetMeid(ctx sdk.Context, meid MockMeid) {
	msk.Meid[meid.Account] = &meid
}

func (msk MockStakingKeeper) GetRegion(ctx sdk.Context, regionId string) (val MockRegion, found bool) {
	v, ok := msk.Region[regionId]
	return *v, ok
}
