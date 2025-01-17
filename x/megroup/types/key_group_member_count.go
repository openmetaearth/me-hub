package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// GroupMemberCountKeyPrefix is the prefix to retrieve all GroupMemberCount
	GroupMemberCountKeyPrefix = "GroupMemberCount/value/"
)

// GroupMemberCountKey returns the store key to retrieve a GroupMemberCount from the index fields
func GroupMemberCountKey(
	groupId uint64,
) []byte {
	var key []byte

	groupIdBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(groupIdBytes, groupId)
	key = append(key, groupIdBytes...)
	key = append(key, []byte("/")...)

	return key
}
