package types

import (
	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/app/params"
	"math/big"
)

const (
	GlobalRegion     = "ME_EARTH"
	ExperienceRegion = "EXPERIENCE_REGION"

	CaclTotalSupply = 20000000000

	MeidJoinGroupNoReward = 0
	MeidJoinGroupRewarded = 1
)

var (
	Bonus           = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(params.BaseDenomUnit), nil))
	InviteReward    = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(params.BaseDenomUnit-1), nil))
	ValidatorReward = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(params.BaseDenomUnit-2), nil))
	CommitteeReward = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(params.BaseDenomUnit-2), nil))
)

func GetRegionAccountAddr(accountType REGION_ACCOUNT_TYPE, regionId string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(RegionAccountNamePrefix + accountType.String() + regionId)))
}
