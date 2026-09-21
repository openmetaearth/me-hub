package types

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func validTestCoreID() []byte {
	coreID := make([]byte, ZkCoreIDLength)
	coreID[0] = 0x01
	coreID[1] = 0xc1
	for i := 2; i < ZkCoreIDLength-2; i++ {
		coreID[i] = byte(i)
	}
	var checksum uint16
	for _, b := range coreID[:ZkCoreIDLength-2] {
		checksum += uint16(b)
	}
	binary.LittleEndian.PutUint16(coreID[ZkCoreIDLength-2:], checksum)
	return coreID
}

func TestValidateAndConvertZkCoreID(t *testing.T) {
	coreID := validTestCoreID()
	require.NoError(t, ValidateZkCoreID(coreID))

	coreIDInt, err := ZkCoreIDToInt(coreID)
	require.NoError(t, err)
	require.Positive(t, coreIDInt.Sign())

	invalid := append([]byte(nil), coreID...)
	invalid[len(invalid)-1] ^= 0xff
	require.Error(t, ValidateZkCoreID(invalid))
	require.Error(t, ValidateZkCoreID(coreID[:len(coreID)-1]))
}

func TestCalculateZkCoreIDBindingChallenge(t *testing.T) {
	coreID := validTestCoreID()
	challenge, err := CalculateZkCoreIDBindingChallenge(
		"me-hub-2401",
		"me1kjnt3ypezt3yf58w8upujvejdtt5xsvkq5dpk4",
		"1234567890123",
		coreID,
		7,
		1_900_000_000,
	)
	require.NoError(t, err)
	require.Positive(t, challenge.Sign())
	require.Less(t, challenge.BitLen(), 253)
}
