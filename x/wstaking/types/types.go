package types

import (
	"math/big"

	"github.com/cometbft/cometbft/crypto"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/params"
)

type MEID_REWARD_TYPE int32

const (
	MeEarthRegionName    = "ME_EARTH"
	ExperienceRegionName = "EXPERIENCE_REGION"

	CaclTotalSupply = 20000000000

	MeidJoinGroupNoReward MEID_REWARD_TYPE = 0
	MeidJoinGroupRewarded MEID_REWARD_TYPE = 1
	MeEarthRegionId                        = "me_earth"
	ExperienceRegionId                     = "experience_region"
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

func (m MEID_REWARD_TYPE) Int32() int32 {
	return int32(m)
}

type UpdatePubKeyInfo struct {
	OperatorAddress string `json:"operator_address"`
	OldConsAddress  []byte `json:"old_cons_address"`
	PubKey          []byte `json:"pub_key"`
	UpdateAtHeight  int64  `json:"update_at_height"`
}

type ReplaceNodePubKey struct {
	OperatorAddress string
	OldPubKey       cryptotypes.PubKey
	NewPubKey       cryptotypes.PubKey
}
