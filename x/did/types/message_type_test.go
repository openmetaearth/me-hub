package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateStatusMessagesUseOwnLegacyTypes(t *testing.T) {
	require.Equal(t, TypeMsgCreateDid, (&MsgCreateDid{}).Type())
	require.Equal(t, TypeMsgUpdateDidStatus, (&MsgUpdateDidStatus{}).Type())
	require.Equal(t, TypeMsgUpdateServiceStatus, (&MsgUpdateServiceStatus{}).Type())
}
