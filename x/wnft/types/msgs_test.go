package types_test

import (
	"testing"

	_ "github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wnft/types"
	"github.com/stretchr/testify/require"
)

const wnftMsgSender = "me139mq752delxv78jvtmwxhasyrycufsvr0mue6u"

func TestMsgNewClassValidateBasicRejectsReservedClassID(t *testing.T) {
	msg := types.NewMsgNewClass(
		types.ReservedKycClassID,
		wnftMsgSender,
		"KYC",
		"KYC",
		"",
		"",
		"",
		0,
	)

	err := msg.ValidateBasic()

	require.Error(t, err)
	require.True(t, types.ErrReservedClassId.Is(err))
}

func TestMsgNewClassValidateBasicAllowsRegularClass(t *testing.T) {
	msg := types.NewMsgNewClass(
		"public-class",
		wnftMsgSender,
		"Public Class",
		"PUB",
		"",
		"",
		"",
		1,
	)

	require.NoError(t, msg.ValidateBasic())
}
