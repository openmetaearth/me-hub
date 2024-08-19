package types

import (
	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	StakePoolName = "stake_tokens_pool"

	GlobalRegion     = "ME_EARTH"
	ExperienceRegion = "EXPERIENCE_REGION"
)

func GetRegionAccountAddr(accountType REGION_ACCOUNT_TYPE, regionId string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(RegionAccountNamePrefix + accountType.String() + regionId)))
}
