package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// MemberJoinedKeyPrefix is the prefix to retrieve all MemberJoined
	MemberJoinedKeyPrefix = "MemberJoined/value/"
)

// MemberJoinedKey returns the store key to retrieve a MemberJoined from the index fields
func MemberJoinedKey(
	address string,
) []byte {
	var key []byte

	addressBytes := []byte(address)
	key = append(key, addressBytes...)
	key = append(key, []byte("/")...)

	return key
}
