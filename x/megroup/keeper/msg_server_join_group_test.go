package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShouldSendAdminJoinReward(t *testing.T) {
	tests := []struct {
		name      string
		applicant string
		admin     string
		want      bool
	}{
		{
			name:      "different applicant and admin",
			applicant: "applicant",
			admin:     "admin",
			want:      true,
		},
		{
			name:      "same applicant and admin",
			applicant: "admin",
			admin:     "admin",
			want:      false,
		},
		{
			name:      "empty admin",
			applicant: "applicant",
			admin:     "",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, shouldSendAdminJoinReward(tt.applicant, tt.admin))
		})
	}
}
