package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/sample"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/kyc/handler"
	"github.com/stretchr/testify/require"
)

func TestGetDidAndKycActiveAcceptsMinimumLevelTwo(t *testing.T) {
	params.SetAddressPrefixes()
	address := sample.Acc()
	const did = "did:me:test"
	const requiredRegion = "usa"

	tests := []struct {
		name      string
		foundDID  bool
		foundInfo bool
		region    string
		level     didtypes.KycLevel
		status    didtypes.DidStatus
		want      bool
	}{
		{
			name:      "level none rejected",
			foundDID:  true,
			foundInfo: true,
			region:    requiredRegion,
			level:     didtypes.KYC_LEVEL_NONE,
			status:    didtypes.DID_STATUS_ACTIVE,
		},
		{
			name:      "level one rejected",
			foundDID:  true,
			foundInfo: true,
			region:    requiredRegion,
			level:     didtypes.KYC_LEVEL_ONE,
			status:    didtypes.DID_STATUS_ACTIVE,
		},
		{
			name:      "level two accepted",
			foundDID:  true,
			foundInfo: true,
			region:    requiredRegion,
			level:     didtypes.KYC_LEVEL_TWO,
			status:    didtypes.DID_STATUS_ACTIVE,
			want:      true,
		},
		{
			name:      "level three accepted",
			foundDID:  true,
			foundInfo: true,
			region:    requiredRegion,
			level:     didtypes.KYC_LEVEL_THREE,
			status:    didtypes.DID_STATUS_ACTIVE,
			want:      true,
		},
		{
			name:      "level five accepted",
			foundDID:  true,
			foundInfo: true,
			region:    requiredRegion,
			level:     didtypes.KYC_LEVEL_FIVE,
			status:    didtypes.DID_STATUS_ACTIVE,
			want:      true,
		},
		{
			name:      "wrong region rejected",
			foundDID:  true,
			foundInfo: true,
			region:    "eu",
			level:     didtypes.KYC_LEVEL_FIVE,
			status:    didtypes.DID_STATUS_ACTIVE,
		},
		{
			name:      "inactive did rejected",
			foundDID:  true,
			foundInfo: true,
			region:    requiredRegion,
			level:     didtypes.KYC_LEVEL_FIVE,
			status:    didtypes.DID_STATUS_INACTIVE,
		},
		{
			name:      "missing did rejected",
			foundInfo: true,
			region:    requiredRegion,
			level:     didtypes.KYC_LEVEL_FIVE,
			status:    didtypes.DID_STATUS_ACTIVE,
		},
		{
			name:     "missing did info rejected",
			foundDID: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kycKeeper := mockKycKeeper{
				didByAddress: map[string]string{},
				infoByDID:    map[string]didtypes.DidInfo{},
			}
			if tt.foundDID {
				kycKeeper.didByAddress[address.String()] = did
			}
			if tt.foundInfo {
				kycKeeper.infoByDID[did] = didtypes.DidInfo{
					Did:      did,
					Address:  address.String(),
					RegionId: tt.region,
					KycLevel: tt.level,
					Status:   tt.status,
				}
			}

			keeper := Keeper{kycKeeper: kycKeeper}
			gotDID, ok := keeper.GetDidAndKycActive(sdk.Context{}, address, requiredRegion)

			require.Equal(t, tt.want, ok)
			if tt.want {
				require.Equal(t, did, gotDID)
			}
		})
	}
}

type mockKycKeeper struct {
	didByAddress map[string]string
	infoByDID    map[string]didtypes.DidInfo
}

func (m mockKycKeeper) RegisterEventHandler(string, int, string, handler.HandlerFunc) {}

func (m mockKycKeeper) GetDID(_ sdk.Context, addr sdk.AccAddress) (string, bool) {
	did, found := m.didByAddress[addr.String()]
	return did, found
}

func (m mockKycKeeper) GetDidInfo(_ sdk.Context, did string) (didtypes.DidInfo, bool) {
	info, found := m.infoByDID[did]
	return info, found
}
