package types

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

const (
	RollupBlockPrefix = "rollupApp/block/"
	//	RollupAppendBlockPrefix = "rollup/block/append/"
	KeyLastRollupCommit   = "KeyLastRollupCommit"
	KeyBlockWithSubmitter = "KeyBlockWithSubmitter"
	KEY_ELECTION_INTERIM  = "KeyElectionInterim"
	//	KeySubmitterLastSubmitTime = "KeySubmitterLastSubmitTime"
)

const (
	SUBMIT_BLOCK_SUCCESS          int = 0
	SUBMIT_BLOCK_NORMAL_ERR       int = 1
	SUBMIT_BLOCK_DA_VALIDATE_ERR  int = 2
	SUBMIT_BLOCK_DA_VERIFY_ERR    int = 3
	SUBMIT_BLOCK_DA_VERIFY_FAILED int = 4
)

const (
	EventSubmitBlockDA           = "EventSubmitBlockDA"
	EventRegisterRollappInitInfo = "EventRegisterRollappInitInfo"
	EventDAFraudChallenge        = "EventDAFraudChallenge"
)

//type LightBlock tenderminttypes.LightBlock

func GetRollupBlockKeyPrefix(rollappID string) []byte {
	return []byte(fmt.Sprintf("%s%s/", RollupBlockPrefix, rollappID))
}

func GetRollupBlockWithSubmitterKeyPrefix(rollappID string) []byte {
	return []byte(fmt.Sprintf("%s%s/%s/", RollupBlockPrefix, rollappID, KeyBlockWithSubmitter))
}

func GetRollupBlockWithSubmitterKeyByBlockHeight(rollappID string, blockHeight uint64) []byte {
	return []byte(fmt.Sprintf("%s%s/%s/%09d", RollupBlockPrefix, rollappID, KeyBlockWithSubmitter, blockHeight))
}

func ConvertBlockHeightToKey(blockHeight uint64) []byte {
	return []byte(fmt.Sprintf("%09d", blockHeight))
}

func GetSubmitterLastSubmitTimeKey(submitter string) []byte {
	return []byte(fmt.Sprintf("LastSubmitTime_%s", submitter))
}

/*
func GetRollupAppendBlockKeyPrefix(rollappID string) []byte {
	return []byte(fmt.Sprintf("%s/%s", RollupAppendBlockPrefix, rollappID))
}

*/

func ConvertToRecordSubmitVal(submitterAddr string, number uint32) []byte {
	return []byte(fmt.Sprintf("%s-%d", submitterAddr, number))
}

func ParserRecordSubmitVal(val string) (string, uint32, error) {
	res := strings.Split(val, "-")
	if len(res) != 2 {
		return "", 0, fmt.Errorf("Parser error. val = %s,sep is -, len = %d", val, len(res))
	}
	number, err := strconv.Atoi(res[1])
	if err != nil {
		return "", 0, fmt.Errorf("strconv.Atoi error. val = %s,err = %s", res[1], err.Error())
	}

	return res[0], uint32(number), nil

}

func GetRollupBlockKey(startHeight uint64, number uint32) []byte {
	return []byte(fmt.Sprintf("%09d-%d", startHeight, number))
}

func ParserRollupKey(key string) (uint64, uint32, error) {
	res := strings.Split(key, "-")
	if len(res) != 2 {
		return 0, 0, fmt.Errorf("Parser error. key = %s,sep is -, len = %d", key, len(res))
	}
	startHeight, err := strconv.ParseUint(res[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	number, err := strconv.ParseUint(res[1], 10, 32)
	if err != nil {
		return 0, 0, err
	}
	return startHeight, uint32(number), nil

}
func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

/*
// 由于之前的GetRollupBlockKeyPrefix已经加了rollappID，此时已经不需要这里在外面加了
func LastRollupCommitKey(rollappID string) []byte {
	return []byte(fmt.Sprintf("%s/%s/", KEY_LAST_ROLLUP_COMMIT, rollappID))
}

*/
