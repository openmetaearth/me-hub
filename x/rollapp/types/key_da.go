package types

import (
	"fmt"
)

const (
	KeyCelestia            = "celestia"
	NameSpacePrefix        = "Namespace_"
	DaFraudChallengePrefix = "DaFraudChallenge"
)
const (
	LenNamespace int = 29
)

const (
	STATUS_CHG_DA_FRAUD_ING     int32 = 0
	STATUS_CHG_DA_FRAUD_SUCCESS int32 = 1
	STATUS_CHG_DA_FRAUD_FAIL    int32 = 2
)

func GetRollappWithCelestiaKey() []byte {
	return []byte(fmt.Sprintf("%s%s/", RollappKeyPrefix, KeyCelestia))
}

func GetDaChallengeKeyPrefix(rollappId string) []byte {
	return []byte(fmt.Sprintf("%s%s/%s/", RollappKeyPrefix, rollappId, DaFraudChallengePrefix))
}

func GetNamespaceKey(rollappID string) []byte {
	return []byte(fmt.Sprintf("%s%s", NameSpacePrefix, rollappID))
}

func ParserRollappIdFrNamespaceIdKey(key []byte) (string, error) {
	nsIdPrefixLen := len(NameSpacePrefix)
	if len(key) <= nsIdPrefixLen {
		return "", fmt.Errorf("key length error. key = %s,keyLen = %d", string(key), len(key))
	}
	return string(key[nsIdPrefixLen:]), nil
}
