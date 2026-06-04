package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/kyc/handler"
	"github.com/stretchr/testify/require"
)

type fakeKycKeeper struct {
	dids     map[string]string
	didInfos map[string]didtypes.DidInfo
}

func (f fakeKycKeeper) RegisterEventHandler(_ string, _ int, _ string, _ handler.HandlerFunc) {}

func (f fakeKycKeeper) GetDID(_ sdk.Context, addr sdk.AccAddress) (string, bool) {
	did, found := f.dids[string(addr)]
	return did, found
}

func (f fakeKycKeeper) GetDidInfo(_ sdk.Context, did string) (didtypes.DidInfo, bool) {
	info, found := f.didInfos[did]
	return info, found
}

func TestGetDidAndKycActiveAcceptsLevelTwoAndHigher(t *testing.T) {
	regionID := "me_earth"

	testCases := []struct {
		name   string
		level  didtypes.KycLevel
		active bool
	}{
		{name: "none", level: didtypes.KYC_LEVEL_NONE, active: false},
		{name: "one", level: didtypes.KYC_LEVEL_ONE, active: false},
		{name: "two", level: didtypes.KYC_LEVEL_TWO, active: true},
		{name: "three", level: didtypes.KYC_LEVEL_THREE, active: true},
		{name: "four", level: didtypes.KYC_LEVEL_FOUR, active: true},
		{name: "five", level: didtypes.KYC_LEVEL_FIVE, active: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			addr := sdk.AccAddress([]byte(tc.name + "_address________"))
			did := "did:" + tc.name
			k := Keeper{
				kycKeeper: fakeKycKeeper{
					dids: map[string]string{
						string(addr): did,
					},
					didInfos: map[string]didtypes.DidInfo{
						did: {
							RegionId: regionID,
							KycLevel: tc.level,
							Status:   didtypes.DID_STATUS_ACTIVE,
						},
					},
				},
			}

			gotDid, ok := k.GetDidAndKycActive(sdk.Context{}, addr, regionID)

			require.Equal(t, tc.active, ok)
			if tc.active {
				require.Equal(t, did, gotDid)
			} else {
				require.Empty(t, gotDid)
			}
		})
	}
}
