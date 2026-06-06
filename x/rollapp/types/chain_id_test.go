package types_test

import (
	"testing"

	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/stretchr/testify/require"
)

func TestNewChainIDRejectsKeySeparator(t *testing.T) {
	_, err := types.NewChainID("parent/child")
	require.Error(t, err)
}

func TestRollappValidateBasicRejectsKeySeparator(t *testing.T) {
	rollapp := types.NewRollapp(sample.AccAddress(), "parent/child", 1, nil, true)

	err := rollapp.ValidateBasic()
	require.Error(t, err)
}
