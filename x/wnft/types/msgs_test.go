package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// TestMsgNewClassValidateBasic tests the ValidateBasic method of MsgNewClass
func TestMsgNewClassValidateBasic(t *testing.T) {
	validAddr := sdk.AccAddress("valid_address______").String()

	testCases := []struct {
		name      string
		msg       *MsgNewClass
		expectErr bool
		errMsg    string
	}{
		{
			name: "valid message",
			msg: &MsgNewClass{
				ClassId:     "test-class",
				Sender:      validAddr,
				Name:        "Test Class",
				Symbol:      "TEST",
				Description: "Test Description",
				Uri:         "ipfs://test",
				UriHash:     "hash",
				TotalSupply: 1000,
			},
			expectErr: false,
		},
		{
			name: "empty class id",
			msg: &MsgNewClass{
				ClassId:     "",
				Sender:      validAddr,
				Name:        "Test",
				Symbol:      "TEST",
				TotalSupply: 100,
			},
			expectErr: true,
			errMsg:    "empty class id",
		},
		{
			name: "zero total supply for non-kyc class",
			msg: &MsgNewClass{
				ClassId:     "test-class",
				Sender:      validAddr,
				Name:        "Test",
				Symbol:      "TEST",
				TotalSupply: 0,
			},
			expectErr: true,
			errMsg:    "total supply",
		},
		{
			name: "zero total supply for kyc class - valid",
			msg: &MsgNewClass{
				ClassId:     "kyc",
				Sender:      validAddr,
				Name:        "KYC",
				Symbol:      "KYC",
				TotalSupply: 0,
			},
			expectErr: false,
		},
		{
			name: "empty name",
			msg: &MsgNewClass{
				ClassId:     "test-class",
				Sender:      validAddr,
				Name:        "",
				Symbol:      "TEST",
				TotalSupply: 100,
			},
			expectErr: true,
			errMsg:    "invalid class name",
		},
		{
			name: "empty symbol",
			msg: &MsgNewClass{
				ClassId:     "test-class",
				Sender:      validAddr,
				Name:        "Test",
				Symbol:      "",
				TotalSupply: 100,
			},
			expectErr: true,
			errMsg:    "invalid class symbol",
		},
		{
			name: "invalid sender address",
			msg: &MsgNewClass{
				ClassId:     "test-class",
				Sender:      "invalid",
				Name:        "Test",
				Symbol:      "TEST",
				TotalSupply: 100,
			},
			expectErr: true,
			errMsg:    "invalid sender address",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestNewMsgNewClass tests the constructor
func TestNewMsgNewClass(t *testing.T) {
	classId := "test-class"
	sender := "sender"
	name := "Test"
	symbol := "TST"
	description := "desc"
	uri := "uri"
	uriHash := "hash"
	totalSupply := uint64(1000)

	msg := NewMsgNewClass(classId, sender, name, symbol, description, uri, uriHash, totalSupply)

	require.Equal(t, classId, msg.ClassId)
	require.Equal(t, sender, msg.Sender)
	require.Equal(t, name, msg.Name)
	require.Equal(t, symbol, msg.Symbol)
	require.Equal(t, description, msg.Description)
	require.Equal(t, uri, msg.Uri)
	require.Equal(t, uriHash, msg.UriHash)
	require.Equal(t, totalSupply, msg.TotalSupply)
}

// TestMsgMintNFTValidateBasic tests the ValidateBasic method of MsgMintNFT
func TestMsgMintNFTValidateBasic(t *testing.T) {
	validAddr := sdk.AccAddress("valid_address______").String()

	testCases := []struct {
		name      string
		msg       *MsgMintNFT
		expectErr bool
		errMsg    string
	}{
		{
			name: "valid message",
			msg: &MsgMintNFT{
				ClassId:  "test-class",
				TokenId:  "1",
				Uri:      "ipfs://token",
				UriHash:  "hash",
				Creator:  validAddr,
				Receiver: validAddr,
			},
			expectErr: false,
		},
		{
			name: "empty class id",
			msg: &MsgMintNFT{
				ClassId:  "",
				TokenId:  "1",
				Uri:      "ipfs://token",
				Creator:  validAddr,
				Receiver: validAddr,
			},
			expectErr: true,
			errMsg:    "empty class id",
		},
		{
			name: "empty token id",
			msg: &MsgMintNFT{
				ClassId:  "test-class",
				TokenId:  "",
				Uri:      "ipfs://token",
				Creator:  validAddr,
				Receiver: validAddr,
			},
			expectErr: true,
			errMsg:    "empty token id",
		},
		{
			name: "empty uri",
			msg: &MsgMintNFT{
				ClassId:  "test-class",
				TokenId:  "1",
				Uri:      "",
				Creator:  validAddr,
				Receiver: validAddr,
			},
			expectErr: true,
			errMsg:    "empty uri",
		},
		{
			name: "invalid creator address",
			msg: &MsgMintNFT{
				ClassId:  "test-class",
				TokenId:  "1",
				Uri:      "ipfs://token",
				Creator:  "invalid",
				Receiver: validAddr,
			},
			expectErr: true,
			errMsg:    "invalid mint address",
		},
		{
			name: "invalid receiver address",
			msg: &MsgMintNFT{
				ClassId:  "test-class",
				TokenId:  "1",
				Uri:      "ipfs://token",
				Creator:  validAddr,
				Receiver: "invalid",
			},
			expectErr: true,
			errMsg:    "invalid receiver address",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestNewMsgMintNFT tests the constructor
func TestNewMsgMintNFT(t *testing.T) {
	classId := "test-class"
	tokenId := "1"
	uri := "ipfs://token"
	uriHash := "hash"
	sender := "sender"
	receiver := "receiver"

	msg := NewMsgMintNFT(classId, tokenId, uri, uriHash, sender, receiver)

	require.Equal(t, classId, msg.ClassId)
	require.Equal(t, tokenId, msg.TokenId)
	require.Equal(t, uri, msg.Uri)
	require.Equal(t, uriHash, msg.UriHash)
	require.Equal(t, sender, msg.Creator)
	require.Equal(t, receiver, msg.Receiver)
}

// TestMsgSendValidateBasic tests the ValidateBasic method of MsgSend
func TestMsgSendValidateBasic(t *testing.T) {
	validAddr := sdk.AccAddress("valid_address______").String()

	testCases := []struct {
		name      string
		msg       *MsgSend
		expectErr bool
		errMsg    string
	}{
		{
			name: "valid message",
			msg: &MsgSend{
				ClassId:  "test-class",
				Id:       "1",
				Sender:   validAddr,
				Receiver: validAddr,
			},
			expectErr: false,
		},
		{
			name: "empty class id",
			msg: &MsgSend{
				ClassId:  "",
				Id:       "1",
				Sender:   validAddr,
				Receiver: validAddr,
			},
			expectErr: true,
			errMsg:    "empty class id",
		},
		{
			name: "empty token id",
			msg: &MsgSend{
				ClassId:  "test-class",
				Id:       "",
				Sender:   validAddr,
				Receiver: validAddr,
			},
			expectErr: true,
			errMsg:    "empty nft id",
		},
		{
			name: "invalid sender address",
			msg: &MsgSend{
				ClassId:  "test-class",
				Id:       "1",
				Sender:   "invalid",
				Receiver: validAddr,
			},
			expectErr: true,
			errMsg:    "invalid sender address",
		},
		{
			name: "invalid receiver address",
			msg: &MsgSend{
				ClassId:  "test-class",
				Id:       "1",
				Sender:   validAddr,
				Receiver: "invalid",
			},
			expectErr: true,
			errMsg:    "invalid receiver address",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
