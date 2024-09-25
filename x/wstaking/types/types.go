package types

import (
	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/app/params"
	"math/big"
)

const (
	MeEarthRegionName    = "ME_EARTH"
	ExperienceRegionName = "EXPERIENCE_REGION"

	CaclTotalSupply = 20000000000

	MeidJoinGroupNoReward = 0
	MeidJoinGroupRewarded = 1
	MeEarthRegionId       = "me_earth"
	ExperienceRegionId    = "experience_region"
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

func GetClassId(regionName string) string {
	return regionName + "-NFT-CLASS-ID"
}

func GetClassName(regionName string) string {
	return regionName + "-NFT-CLASS-NAME"
}

func GetClassSymbol(regionName string) string {
	return regionName + "-NFT-CLASS-SYMBOL"
}

func GetClassDescription(regionId string) string {
	return "nft class for region " + regionId
}
