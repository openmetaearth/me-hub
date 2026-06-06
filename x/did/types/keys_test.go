package types

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterPrefixesUseLengthDelimitedParts(t *testing.T) {
	legacyA := GetLegacyFilterPrefixBySidAndFilter("ab", []byte("cX"))
	legacyB := GetLegacyFilterPrefixBySidAndFilter("abc", []byte("X"))
	require.True(t, bytes.Equal(legacyA, legacyB))

	prefixA := GetFilterPrefixBySidAndFilter("ab", []byte("cX"))
	prefixB := GetFilterPrefixBySidAndFilter("abc", []byte("X"))
	require.False(t, bytes.Equal(prefixA, prefixB))
}
