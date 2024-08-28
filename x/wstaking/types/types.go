package types

import (
	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	GlobalRegion     = "ME_EARTH"
	ExperienceRegion = "EXPERIENCE_REGION"

	CaclTotalSupply = 20000000000

	MeidJoinGroupNoReward = 0
	MeidJoinGroupRewarded = 1
)

func GetRegionAccountAddr(accountType REGION_ACCOUNT_TYPE, regionId string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(RegionAccountNamePrefix + accountType.String() + regionId)))
}
