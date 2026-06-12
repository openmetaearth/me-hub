package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFreeGasAccountUsesOwnLegacyType(t *testing.T) {
	require.Equal(t, TypeMsgUpdateDao, (&MsgUpdateDao{}).Type())
	require.Equal(t, TypeMsgFreeGasAccount, (&MsgFreeGasAccount{}).Type())
	require.NotEqual(t, (&MsgUpdateDao{}).Type(), (&MsgFreeGasAccount{}).Type())
}
