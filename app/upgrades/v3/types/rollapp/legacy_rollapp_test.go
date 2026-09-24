package rollapp_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	legacyrollapp "github.com/openmetaearth/me-hub/app/upgrades/v3/types/rollapp"
)

// encodeStringField encodes a protobuf length-delimited string field.
func encodeStringField(fieldNum int, s string) []byte {
	tag := byte(fieldNum<<3 | 2)
	b := []byte{tag, byte(len(s))}
	return append(b, s...)
}

func TestLegacyRollappUnmarshal_RegisteredDenoms(t *testing.T) {
	// field 1: rollappId, field 10: registeredDenoms (repeated)
	bz := append(encodeStringField(1, "mecheckin_100-1"), encodeStringField(10, "mea")...)
	bz = append(bz, encodeStringField(10, "ibc/ABC")...)
	// unrelated field 2 (owner) should be skipped
	bz = append(bz, encodeStringField(2, "me1owner")...)

	var legacy legacyrollapp.LegacyRollapp
	require.NoError(t, legacy.Unmarshal(bz))
	require.Equal(t, "mecheckin_100-1", legacy.RollappId)
	require.Equal(t, []string{"mea", "ibc/ABC"}, legacy.RegisteredDenoms)
}
