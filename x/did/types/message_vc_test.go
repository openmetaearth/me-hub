package types

import (
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

const (
	testIssuer = "me1kjnt3ypezt3yf58w8upujvejdtt5xsvkq5dpk4"
	testDid    = "did:me:holder"
	testSid    = "kyc"
)

func TestMain(m *testing.M) {
	sdk.GetConfig().SetBech32PrefixForAccount("me", "mepub")
	os.Exit(m.Run())
}

func TestMsgCreateVCValidateBasicDataSize(t *testing.T) {
	require.NoError(t, newValidCreateVC(make([]byte, maxCredentialDataLength)).ValidateBasic())

	err := newValidCreateVC(make([]byte, maxCredentialDataLength+1)).ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "data length exceeds")
}

func TestMsgUpdateVCValidateBasicDataSize(t *testing.T) {
	require.NoError(t, newValidUpdateVC(make([]byte, maxCredentialDataLength)).ValidateBasic())

	err := newValidUpdateVC(make([]byte, maxCredentialDataLength+1)).ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "data length exceeds")
}

func newValidCreateVC(data []byte) *MsgCreateVC {
	return NewMsgCreateVC(testIssuer, testDid, testSid, "hash", "https://example.com/vc", data, nil)
}

func newValidUpdateVC(data []byte) *MsgUpdateVC {
	return NewMsgUpdateVC(testIssuer, testDid, testSid, "hash", "https://example.com/vc", data, nil)
}
