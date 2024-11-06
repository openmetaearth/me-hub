package types

import (
	"bytes"
	"encoding/binary"
	"fmt"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"strconv"
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/types/kv"
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
	RegionKeyPrefix = "Region/value/"
)

const (
	RegionAccountNamePrefix = "Region-Module-Account-"
)

const (
	FixedDepositCfgKeyPrefix        = "FixedDepositCfg/value/"
	FixedDepositCountOfCfgKeyPrefix = "FixedDepositCountOfCfg/value/"
)

var (
	// Keys for store prefixes
	// Last* values are constant during a block.
	LastValidatorPowerKey = []byte{0x11} // prefix for each key to a validator index, for bonded validators
	LastTotalPowerKey     = []byte{0x12} // prefix for the total power

	ValidatorsKey             = []byte{0x21} // prefix for each key to a validator
	ValidatorsByConsAddrKey   = []byte{0x22} // prefix for each key to a validator index, by pubkey
	ValidatorsByPowerIndexKey = []byte{0x23} // prefix for each key to a validator index, sorted by power

	DelegationKey                    = []byte{0x31} // key for a delegation
	UnbondingDelegationKey           = []byte{0x32} // key for an unbonding-delegation
	UnbondingDelegationByValIndexKey = []byte{0x33} // prefix for each key for an unbonding-delegation, by validator operator
	RedelegationKey                  = []byte{0x34} // key for a redelegation
	RedelegationByValSrcIndexKey     = []byte{0x35} // prefix for each key for an redelegation, by source validator operator
	RedelegationByValDstIndexKey     = []byte{0x36} // prefix for each key for an redelegation, by destination validator operator

	UnbondingIDKey    = []byte{0x37} // key for the counter for the incrementing id for UnbondingOperations
	UnbondingIndexKey = []byte{0x38} // prefix for an index for looking up unbonding operations by their IDs
	UnbondingTypeKey  = []byte{0x39} // prefix for an index containing the type of unbonding operations

	UnbondingQueueKey    = []byte{0x41} // prefix for the timestamps in unbonding queue
	RedelegationQueueKey = []byte{0x42} // prefix for the timestamps in redelegations queue
	ValidatorQueueKey    = []byte{0x43} // prefix for the timestamps in validator queue

	HistoricalInfoKey   = []byte{0x50} // prefix for the historical info
	ParamsKey           = []byte{0x51} // prefix for parameters for module x/staking
	ValidatorUpdatesKey = []byte{0x71} // prefix for the end block validator updates key

	StakeKey                    = []byte{0x61} // key for a stake
	UnbondingStakeKey           = []byte{0x72} // key for an unbonding-stake
	UnbondingStakeByValIndexKey = []byte{0x73} // prefix for each key for an unbonding-stake, by validator operator
	UnbondingStakeQueueKey      = []byte{0x74} // prefix for the timestamps in unbonding stake queue
)

// UnbondingType defines the type of unbonding operation
type UnbondingType int

const (
	UnbondingType_Undefined UnbondingType = iota
	UnbondingType_UnbondingDelegation
	UnbondingType_Redelegation
	UnbondingType_ValidatorUnbonding
)

// GetUnbondingTypeKey returns a key for an index containing the type of unbonding operations
func GetUnbondingTypeKey(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return append(UnbondingTypeKey, bz...)
}

// GetUnbondingIndexKey returns a key for the index for looking up UnbondingDelegations by the UnbondingDelegationEntries they contain
func GetUnbondingIndexKey(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return append(UnbondingIndexKey, bz...)
}

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// GetValidatorKey creates the key for the validator with address
// VALUE: staking/Validator
func GetValidatorKey(operatorAddr sdk.ValAddress) []byte {
	return append(ValidatorsKey, address.MustLengthPrefix(operatorAddr)...)
}

// GetValidatorByConsAddrKey creates the key for the validator with pubkey
// VALUE: validator operator address ([]byte)
func GetValidatorByConsAddrKey(addr sdk.ConsAddress) []byte {
	return append(ValidatorsByConsAddrKey, address.MustLengthPrefix(addr)...)
}

// AddressFromValidatorsKey creates the validator operator address from ValidatorsKey
func AddressFromValidatorsKey(key []byte) []byte {
	kv.AssertKeyAtLeastLength(key, 3)
	return key[2:] // remove prefix bytes and address length
}

// AddressFromLastValidatorPowerKey creates the validator operator address from LastValidatorPowerKey
func AddressFromLastValidatorPowerKey(key []byte) []byte {
	kv.AssertKeyAtLeastLength(key, 3)
	return key[2:] // remove prefix bytes and address length
}

// GetValidatorsByPowerIndexKey creates the validator by power index.
// Power index is the key used in the power-store, and represents the relative
// power ranking of the validator.
// VALUE: validator operator address ([]byte)
func GetValidatorsByPowerIndexKey(validator stakingtypes.Validator, powerReduction math.Int) []byte {
	// NOTE the address doesn't need to be stored because counter bytes must always be different
	// NOTE the larger values are of higher value

	consensusPower := sdk.TokensToConsensusPower(validator.Tokens, powerReduction)
	consensusPowerBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(consensusPowerBytes, uint64(consensusPower))

	powerBytes := consensusPowerBytes
	powerBytesLen := len(powerBytes) // 8

	addr, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
	if err != nil {
		panic(err)
	}
	operAddrInvr := sdk.CopyBytes(addr)
	addrLen := len(operAddrInvr)

	for i, b := range operAddrInvr {
		operAddrInvr[i] = ^b
	}

	// key is of format prefix || powerbytes || addrLen (1byte) || addrBytes
	key := make([]byte, 1+powerBytesLen+1+addrLen)

	key[0] = ValidatorsByPowerIndexKey[0]
	copy(key[1:powerBytesLen+1], powerBytes)
	key[powerBytesLen+1] = byte(addrLen)
	copy(key[powerBytesLen+2:], operAddrInvr)

	return key
}

// GetLastValidatorPowerKey creates the bonded validator index key for an operator address
func GetLastValidatorPowerKey(operator sdk.ValAddress) []byte {
	return append(LastValidatorPowerKey, address.MustLengthPrefix(operator)...)
}

// ParseValidatorPowerRankKey parses the validators operator address from power rank key
func ParseValidatorPowerRankKey(key []byte) (operAddr []byte) {
	powerBytesLen := 8

	// key is of format prefix (1 byte) || powerbytes || addrLen (1byte) || addrBytes
	operAddr = sdk.CopyBytes(key[powerBytesLen+2:])

	for i, b := range operAddr {
		operAddr[i] = ^b
	}

	return operAddr
}

// GetValidatorQueueKey returns the prefix key used for getting a set of unbonding
// validators whose unbonding completion occurs at the given time and height.
func GetValidatorQueueKey(timestamp time.Time, height int64) []byte {
	heightBz := sdk.Uint64ToBigEndian(uint64(height))
	timeBz := sdk.FormatTimeBytes(timestamp)
	timeBzL := len(timeBz)
	prefixL := len(ValidatorQueueKey)

	bz := make([]byte, prefixL+8+timeBzL+8)

	// copy the prefix
	copy(bz[:prefixL], ValidatorQueueKey)

	// copy the encoded time bytes length
	copy(bz[prefixL:prefixL+8], sdk.Uint64ToBigEndian(uint64(timeBzL)))

	// copy the encoded time bytes
	copy(bz[prefixL+8:prefixL+8+timeBzL], timeBz)

	// copy the encoded height
	copy(bz[prefixL+8+timeBzL:], heightBz)

	return bz
}

// ParseValidatorQueueKey returns the encoded time and height from a key created
// from GetValidatorQueueKey.
func ParseValidatorQueueKey(bz []byte) (time.Time, int64, error) {
	prefixL := len(ValidatorQueueKey)
	if prefix := bz[:prefixL]; !bytes.Equal(prefix, ValidatorQueueKey) {
		return time.Time{}, 0, fmt.Errorf("invalid prefix; expected: %X, got: %X", ValidatorQueueKey, prefix)
	}

	timeBzL := sdk.BigEndianToUint64(bz[prefixL : prefixL+8])
	ts, err := sdk.ParseTimeBytes(bz[prefixL+8 : prefixL+8+int(timeBzL)])
	if err != nil {
		return time.Time{}, 0, err
	}

	height := sdk.BigEndianToUint64(bz[prefixL+8+int(timeBzL):])

	return ts, int64(height), nil
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

// GetDelegationsKey creates the prefix for a delegator for all validators
func GetDelegationsKey(delAddr sdk.AccAddress) []byte {
	return append(DelegationKey, address.MustLengthPrefix(delAddr)...)
}

// GetUBDKey creates the key for an unbonding delegation by delegator and validator addr
// VALUE: staking/UnbondingDelegation
func GetUBDKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetUBDsKey(delAddr.Bytes()), address.MustLengthPrefix(valAddr)...)
}

// GetUBDByValIndexKey creates the index-key for an unbonding delegation, stored by validator-index
// VALUE: none (key rearrangement used)
func GetUBDByValIndexKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetUBDsByValIndexKey(valAddr), address.MustLengthPrefix(delAddr)...)
}

// GetUBDsKey creates the prefix for all unbonding delegations from a delegator
func GetUBDsKey(delAddr sdk.AccAddress) []byte {
	return append(UnbondingDelegationKey, address.MustLengthPrefix(delAddr)...)
}

// GetUBDsByValIndexKey creates the prefix keyspace for the indexes of unbonding delegations for a validator
func GetUBDsByValIndexKey(valAddr sdk.ValAddress) []byte {
	return append(UnbondingDelegationByValIndexKey, address.MustLengthPrefix(valAddr)...)
}

// GetHistoricalInfoKey returns a key prefix for indexing HistoricalInfo objects.
func GetHistoricalInfoKey(height int64) []byte {
	return append(HistoricalInfoKey, []byte(strconv.FormatInt(height, 10))...)
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
