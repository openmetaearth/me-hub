package types

import (
	"math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestMsgNewClass_ValidateBasic_TotalSupplyCap(t *testing.T) {
	// Generate a valid sender address
	sender := sdk.AccAddress([]byte("sender___________________")).String()

	tests := []struct {
		name        string
		msg         *MsgNewClass
		expectError bool
	}{
		{
			name: "valid total supply within limit",
			msg: &MsgNewClass{
				ClassId:     "test-class",
				Sender:      sender,
				Name:        "Test Class",
				Symbol:      "TC",
				TotalSupply: 100,
			},
			expectError: false,
		},
		{
			name: "total supply at uint32 max - valid",
			msg: &MsgNewClass{
				ClassId:     "test-class",
				Sender:      sender,
				Name:        "Test Class",
				Symbol:      "TC",
				TotalSupply: math.MaxUint32,
			},
			expectError: false,
		},
		{
			name: "total supply exceeds uint32 max - invalid",
			msg: &MsgNewClass{
				ClassId:     "test-class",
				Sender:      sender,
				Name:        "Test Class",
				Symbol:      "TC",
				TotalSupply: math.MaxUint32 + 1,
			},
			expectError: true,
		},
		{
			name: "total supply at uint64 max - invalid",
			msg: &MsgNewClass{
				ClassId:     "test-class",
				Sender:      sender,
				Name:        "Test Class",
				Symbol:      "TC",
				TotalSupply: math.MaxUint64,
			},
			expectError: true,
		},
		{
			name: "zero total supply (non-kyc) - invalid",
			msg: &MsgNewClass{
				ClassId:     "test-class",
				Sender:      sender,
				Name:        "Test Class",
				Symbol:      "TC",
				TotalSupply: 0,
			},
			expectError: true,
		},
		{
			name: "zero total supply for kyc class - valid",
			msg: &MsgNewClass{
				ClassId:     "kyc",
				Sender:      sender,
				Name:        "KYC Class",
				Symbol:      "KYC",
				TotalSupply: 0,
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectError && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}
