package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/rollapp/types"
)

func TestNewChainIDRejectsKeySeparator(t *testing.T) {
	_, err := types.NewChainID("parent/child")
	require.Error(t, err)
}

func TestMsgCreateRollappValidateBasicRejectsKeySeparator(t *testing.T) {
	msg := types.NewMsgCreateRollapp(sample.AccAddress(), "parent/child", 1, nil)

	err := msg.ValidateBasic()
	require.Error(t, err)
}
