package types

import (
	"encoding/binary"
	"fmt"
)

const (
	MODULE_NAME                = "hubRollUp"
	KEY_LAST_ELECTION_TIME     = "LastElectionTime"
	KEY_FIRST_ELECTION_TIME    = "FirstElectionStartTime"
	KEY_LAST_UNSTAKE_TIME      = "LastUnStakeTime"
	KEY_LAST_ELECTION_INFO     = "LastElectionInfo"
	KEY_PREVIOUS_ELECTION_INFO = "PreviousElectionInfo"
	KEY_ROLLAPP_ID             = "KeyRollappID"
	StoreKey                   = MODULE_NAME
	RouterKey                  = MODULE_NAME
	KeyDaFraudPunishPrefix     = "DaFraudPunish_"
	KeyRollupBlackList         = "RollupBlackList"
	KeyProvedDaFraudPrefix     = "ProvedDaFraud_"
	//	KeyRollappInitInfo         = "KeyRollappInitInfo"
	KeyNeedPunishment         = "KeyNeedPunishment"
	KeyFinishPunishmentPrefix = "FinishPunish_"
	KeyRollappInitInfoPrefix  = "RollappInfo_"
)

const (
	// RollappKeyPrefix is the prefix to retrieve all Rollapp
	RollupKeyPrefix                      = "Rollup/value/"
	RollupStakeKeyPrefix                 = "Rollup/value/Stake/"
	RollupStakeChallengeDaFraudKeyPrefix = RollupStakeKeyPrefix + "ChallengeDaFraud/"
	RollupDaFraudStatics                 = RollupKeyPrefix + "DaFraudStatics/"
	RollupStakeBondNodePrefix            = RollupStakeKeyPrefix + "BondNodeAddr/"
	RollupDelegatorStakeBondNodePrefix   = RollupStakeKeyPrefix + "DelegatorStakeBondNode/"
	//RollupBlackListPrefix                = RollupKeyPrefix + "BlackList"
	//RollupAppIdKeyPrefix = RollupKeyPrefix + KeyRollappIdPrefix
)
const (
	SubmitDaFraudTime int32 = 72
	//用于
	//DaFraudVerifyReserveTime = 12
)

const (
	RESULT_CHG_SUCCESS_SUBMIT_DATA_FAIL   int32 = 0
	RESULT_CHG_SUCCESS_SUBMIT_DATA_SUCESS int32 = 1
	RESULT_CHG_FAIL                       int32 = 2
)

const (
	//	DaySeconds  int64 = 86400
	HourSeconds   int64 = 3600
	MinuteSeconds int64 = 60
)

const (
	NodeNormal    int32 = 0
	NodeSequencer int32 = 1
	NodeBackup    int32 = 2
)

func GetRollupAppKeyPrefix(rollappID string) []byte {
	return []byte(fmt.Sprintf("%s%s/", RollupKeyPrefix, rollappID))
}

func GetRollupAppInitInfKey(rollappID string) []byte {
	return []byte(fmt.Sprintf("%s%s", KeyRollappInitInfoPrefix, rollappID))
}

func GetRollupAppStakeKeyPrefix(rollappID string) []byte {
	return []byte(fmt.Sprintf("%s%s/", RollupStakeKeyPrefix, rollappID))
}

func GetStakeBondNodeAddrPrefix(rollappID string) []byte {
	return []byte(fmt.Sprintf("%s%s/", RollupStakeBondNodePrefix, rollappID))
}

func GetDelegatorStakeNodePrefix(rollappID string) []byte {
	return []byte(fmt.Sprintf("%s%s/", RollupDelegatorStakeBondNodePrefix, rollappID))
}

func GetStakeForChallengeDaFraudPrefix(rollappID string) []byte {
	return []byte(fmt.Sprintf("%s%s/", RollupStakeChallengeDaFraudKeyPrefix, rollappID))
}

func GetDaFraudStaticsPrefix(rollappID string) []byte {
	return []byte(fmt.Sprintf("%s%s/", RollupDaFraudStatics, rollappID))
}

// da fraud
func GetDaFraudPunishmentKey(punishAddr string) []byte {
	return []byte(fmt.Sprintf("%s%s", KeyDaFraudPunishPrefix, punishAddr))
}

func GetProvedDaFraudKeyPrefix(fraudster string) []byte {
	return []byte(fmt.Sprintf("%s%s_", KeyProvedDaFraudPrefix, fraudster))
}

func GetFinishDaFraudPunishKey(fraudster string) []byte {
	return []byte(fmt.Sprintf("%s%s", KeyFinishPunishmentPrefix, fraudster))
}

func GetRollappIdFromInitInfoKey(rollappInitInfoKey []byte) []byte {
	lenPrefix := len(KeyRollappInitInfoPrefix)
	if len(rollappInitInfoKey) <= lenPrefix {
		return nil
	}
	return rollappInitInfoKey[lenPrefix:]
}

//

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func GetRollappIdAndNamespaceVal(rollappID string, namespaceVal []byte) []byte {
	return append([]byte(fmt.Sprintf("%s-")))
}

type ElectionInfo struct {
	StakeAmount uint64
	Address     string
}

type RollappInitExtVal struct {
	IdInDA              []byte `json:"id_in_da"`
	FirstElectBlkHeight uint64 `json:"first_elect_blk_height"`
}

type ElectionsList []ElectionInfo

func (t ElectionsList) Len() int {
	return len(t)
}
func (t ElectionsList) Less(i, j int) bool { //降序排列
	return t[i].StakeAmount > t[j].StakeAmount
}
func (t ElectionsList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

/*
// MsgStake defines a Stake message
type MsgStake struct {
	Delegator sdk.AccAddress `json:"delegator" yaml:"delegator"`
	Amount    sdk.Coin       `json:"amount" yaml:"amount"`
}

// MsgUnstake defines an Unstake message
type MsgUnstake struct {
	Delegator sdk.AccAddress `json:"delegator" yaml:"delegator"`
	Amount    sdk.Coin       `json:"amount" yaml:"amount"`
}



// NewMsgStake creates a new MsgStake instance
func NewMsgStake(delegator sdk.AccAddress, amount sdk.Coin) MsgStake {
	return MsgStake{
		Delegator: delegator,
		Amount:    amount,
	}
}

// NewMsgUnstake creates a new MsgUnstake instance
func NewMsgUnstake(delegator sdk.AccAddress, amount sdk.Coin) MsgUnstake {
	return MsgUnstake{
		Delegator: delegator,
		Amount:    amount,
	}
}

// Route returns the name of the module
func (msg MsgStake) Route() string { return RouterKey }

// Type returns the action
func (msg MsgStake) Type() string { return "Stake" }

// ValidateBasic runs stateless checks on the message
func (msg MsgStake) ValidateBasic() error {
	if msg.Delegator.Empty() {
		return errors.Wrap(errors.ErrInvalidAddress, "missing delegator")
	}
	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return errors.Wrap(errors.ErrInvalidCoins, "invalid amount")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgStake) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgStake) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Delegator}
}

// Route returns the name of the module
func (msg MsgUnstake) Route() string { return RouterKey }

// Type returns the action
func (msg MsgUnstake) Type() string { return "Unstake" }

// ValidateBasic runs stateless checks on the message
func (msg MsgUnstake) ValidateBasic() error {
	if msg.Delegator.Empty() {
		return errors.Wrap(errors.ErrInvalidAddress, "missing delegator")
	}
	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return errors.Wrap(errors.ErrInvalidCoins, "invalid amount")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgUnstake) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgUnstake) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Delegator}
}
*/
