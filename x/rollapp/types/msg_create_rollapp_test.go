package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestMsgCreateRollapp_ValidateBasic_WhitespaceDetection ensures that
// RollappId values containing leading or trailing whitespace are rejected
// by ValidateBasic, preventing the namespace squatting vulnerability
// described in ME-Hub bug bounty Medium severity report.
func TestMsgCreateRollapp_ValidateBasic_WhitespaceDetection(t *testing.T) {
	tests := []struct {
		name      string
		rollappID string
		wantErr   bool
	}{
		{
			name:      "non-EIP RollappId with surrounding spaces",
			rollappID: "  dup  ",
			wantErr:   true,
		},
		{
			name:      "EIP155 RollappId with surrounding spaces",
			rollappID: "  evil_1-1  ",
			wantErr:   true,
		},
		{
			name:      "EIP155 RollappId with leading space",
			rollappID: " good_1-1",
			wantErr:   true,
		},
		{
			name:      "EIP155 RollappId with trailing space",
			rollappID: "good_1-1 ",
			wantErr:   true,
		},
		{
			name:      "valid non-EIP RollappId, no whitespace",
			rollappID: "dup",
			wantErr:   false,
		},
		{
			name:      "valid EIP155 RollappId, no whitespace",
			rollappID: "good_1-1",
			wantErr:   false,
		},
		{
			name:      "empty string after trimming (only spaces)",
			rollappID: "   ",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &MsgCreateRollapp{
				RollappId:             tt.rollappID,
				MaxSequencers:         1,
				PermissionedAddresses: []string{},
				Creator:               "cosmos1testcreatoraddress",
			}
			err := msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err, "expected error for RollappId %q", tt.rollappID)
			} else {
				require.NoError(t, err, "unexpected error for RollappId %q", tt.rollappID)
			}
		})
	}
}