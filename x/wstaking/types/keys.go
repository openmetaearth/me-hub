package types

import (
	"encoding/binary"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName is the name of the staking module
	ModuleName = "staking"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// RouterKey is the msg router key for the staking module
	RouterKey = ModuleName
)

const (
	FixedDepositKey            = "FixedDeposit-value-"
	FixedDepositCountKey       = "FixedDeposit-count-"
	FixedDepositKeyAcct        = "FixedDeposit-index-acct-"
	FixedDepositTotalAmountKey = "fixed-deposit-total-amount/"
)

const (
	MeidKeyPrefix       = "Meid/value/"
	MeidRegionKeyPrefix = "Meid-region"
)

const (
	MeidNFTKeyPrefix        = "MeidNFT/value/"
	MeidNFTAccountKeyPrefix = "MeidNFT-account"
)

const (
	RegionKeyPrefix         = "Region/value/"
	RegionWithdrawKeyPrefix = "RegionWithdraw/value/"
)

const (
	RegionAccountNamePrefix = "Region-Module-Account-"
)

const (
	FixedDepositCfgKeyPrefix        = "FixedDepositCfg/value/"
	FixedDepositCountOfCfgKeyPrefix = "FixedDepositCountOfCfg/value/"
)

const (
	TypeMsgReplaceConsensusPubKey   = "replace_consensus_pubkey"
	ReplaceConsensusPubKey          = "ReplaceNodePubKey"
	EventTypeReplacePubKey          = "replace_pubkey"
	EventTypeStartReplacePubKey     = "start_replace_pubkey"
	EventTypeDelayRemoveOldConsAddr = "delay_remove_old_cons_addr"
	EventTypeReplacePubKeyFailed    = "replace_pubkey_failed"

	AttributeKeyOperatorAddress = "operator_address"
	AttributeKeyPubKey          = "pub_key"
	AttributeKeyOldConsAddr     = "old_cons_addr"
	AttributeKeyNowConsAddr     = "now_cons_addr"
	AttributeKeyUpdateAtHeight  = "update_at_height"
	AttributeKeyFailedReason    = "failed_reason"
)

var (
	StakeKey                    = []byte{0x61} // key for a stake
	UnbondingStakeKey           = []byte{0x72} // key for an unbonding-stake
	UnbondingStakeByValIndexKey = []byte{0x73} // prefix for each key for an unbonding-stake, by validator operator
	UnbondingStakeQueueKey      = []byte{0x74} // prefix for the timestamps in unbonding stake queue

	NewRecordKey                 = []byte{0x88} //key for new record
	ReviewRecordKey              = []byte{0x89} // key for new review record
	InviteKey                    = []byte{0x90}
	ChangeDelegationValidatorKey = []byte{0x91}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// GetStakeKey creates the key for staker bond with validator
// VALUE: staking/Stake
func GetStakeKey(stakeAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetStakesKey(stakeAddr), address.MustLengthPrefix(valAddr)...)
}

// GetStakesKey creates the prefix for a staker for all validators
func GetStakesKey(stakeAddr sdk.AccAddress) []byte {
	return append(StakeKey, address.MustLengthPrefix(stakeAddr)...)
}

// GetUBSKey creates the key for an unbonding stake by staker and validator addr
// VALUE: staking/UnbondingStake
func GetUBSKey(stakerAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetUBSsKey(stakerAddr.Bytes()), address.MustLengthPrefix(valAddr)...)
}

// GetUBDsKey creates the prefix for all unbonding stakes from a staker
func GetUBSsKey(delAddr sdk.AccAddress) []byte {
	return append(UnbondingStakeKey, address.MustLengthPrefix(delAddr)...)
}

// GetUBSByValIndexKey creates the index-key for an unbonding stake, stored by validator-index
// VALUE: none (key rearrangement used)
func GetUBSByValIndexKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetUBSsByValIndexKey(valAddr), address.MustLengthPrefix(delAddr)...)
}

// GetUBDsByValIndexKey creates the prefix keyspace for the indexes of unbonding stakes for a validator
func GetUBSsByValIndexKey(valAddr sdk.ValAddress) []byte {
	return append(UnbondingStakeByValIndexKey, address.MustLengthPrefix(valAddr)...)
}

// GetUnbondingStakeTimeKey creates the prefix for all unbonding stakes from a staker
func GetUnbondingStakeTimeKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)
	return append(UnbondingStakeQueueKey, bz...)
}

// GetRecordKey creates the prefix for a record
func GetRecordKey(acc sdk.AccAddress) []byte {
	return append(NewRecordKey, address.MustLengthPrefix(acc)...)
}

// GetReviewRecordKey creates the prefix for a review
func GetReviewRecordKey(recordNum string) []byte {
	b := []byte(recordNum)
	return append(ReviewRecordKey, b...)
}

// MeidKey returns the store key to retrieve a Meid from the index fields
func MeidKey(account string) []byte {
	var key []byte

	accountBytes := []byte(account)
	key = append(key, accountBytes...)
	key = append(key, []byte("/")...)

	return key
}

// RegionKey returns the store key to retrieve a Region from the index fields
func RegionKey(regionId string) []byte {
	var key []byte

	regionIdBytes := []byte(regionId)
	key = append(key, regionIdBytes...)
	key = append(key, []byte("/")...)

	return key
}

func MeidNFTKey(umeid string) []byte {
	var key []byte
	meidBytes := []byte(umeid)
	key = append(key, meidBytes...)
	key = append(key, []byte("/")...)

	return key
}

func FixedDepositCfgKey(term int64) []byte {
	var key = make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(term))
	key = append(key, []byte("/")...)

	return key
}
