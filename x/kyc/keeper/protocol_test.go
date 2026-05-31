package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	kyctypes "github.com/openmetaearth/me-hub/x/kyc/types"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
	"github.com/stretchr/testify/require"
)

func TestGetAllRegionsDoesNotPrependZeroValueRegions(t *testing.T) {
	stakingKeeper := fakeProtocolStakingKeeper{
		regions: []wstakingtypes.Region{
			{RegionId: "us-east", Name: "US East"},
			{RegionId: "eu-west", Name: "EU West"},
		},
	}
	keeper := Keeper{stkKeeper: stakingKeeper}

	regions := keeper.GetAllRegions(sdk.Context{})

	require.Equal(t, []kyctypes.Region{
		{Id: "us-east", Name: "US East"},
		{Id: "eu-west", Name: "EU West"},
	}, regions)
	for _, region := range regions {
		require.NotEmpty(t, region.Id)
		require.NotEmpty(t, region.Name)
	}
}

type fakeProtocolStakingKeeper struct {
	regions []wstakingtypes.Region
}

func (f fakeProtocolStakingKeeper) GetRegion(_ sdk.Context, regionId string) (wstakingtypes.Region, bool) {
	for _, region := range f.regions {
		if region.RegionId == regionId {
			return region, true
		}
	}
	return wstakingtypes.Region{}, false
}

func (f fakeProtocolStakingKeeper) GetAllRegion(_ sdk.Context) []wstakingtypes.Region {
	return f.regions
}

func (fakeProtocolStakingKeeper) KycReward(sdk.Context, sdk.AccAddress, string, string) error {
	return nil
}

func (fakeProtocolStakingKeeper) RemoveKycReward(sdk.Context, sdk.AccAddress, string) error {
	return nil
}

func (fakeProtocolStakingKeeper) TransferKycRegion(sdk.Context, sdk.AccAddress, string, string, string) error {
	return nil
}

func (fakeProtocolStakingKeeper) SendInviteReward(sdk.Context, string, string, string) error {
	return nil
}
