package types

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	RollupBlockPrefix = "rollupApp/block/"
	//	RollupAppendBlockPrefix = "rollup/block/append/"
	KeyLastRollupCommit = "KeyLastRollupCommit"
	//KEY_APPEND_BLOCK       = "KeyAppendBlock"
	KEY_ELECTION_INTERIM = "KeyElectionInterim"
)

const (
	SUBMIT_BLOCK_SUCCESS          int = 0
	SUBMIT_BLOCK_NORMAL_ERR       int = 1
	SUBMIT_BLOCK_DA_VALIDATE_ERR  int = 2
	SUBMIT_BLOCK_DA_VERIFY_ERR    int = 3
	SUBMIT_BLOCK_DA_VERIFY_FAILED int = 4
)

const (
	EventSubmitBlockDA = "EventSubmitBlockDA"
)

//type LightBlock tenderminttypes.LightBlock

func GetRollupBlockKeyPrefix(rollappID string) []byte {
	return []byte(fmt.Sprintf("%s%s/", RollupBlockPrefix, rollappID))
}

/*
func GetRollupAppendBlockKeyPrefix(rollappID string) []byte {
	return []byte(fmt.Sprintf("%s/%s", RollupAppendBlockPrefix, rollappID))
}

*/

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

/*
// 由于之前的GetRollupBlockKeyPrefix已经加了rollappID，此时已经不需要这里在外面加了
func LastRollupCommitKey(rollappID string) []byte {
	return []byte(fmt.Sprintf("%s/%s/", KEY_LAST_ROLLUP_COMMIT, rollappID))
}

*/
