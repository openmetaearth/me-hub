package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferRegionUsesOwnLegacyType(t *testing.T) {
	require.Equal(t, TypeMsgTransferRegion, (&MsgTransferRegion{}).Type())
}
