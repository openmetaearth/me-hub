package types

import (
	"strings"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/openmetaearth/me-hub/testutil/sample"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
	"github.com/stretchr/testify/require"
)

func TestValidateApproveLevel(t *testing.T) {
	tests := []struct {
		name  string
		level didtypes.KycLevel
		err   error
	}{
		{
			name:  "rejects none",
			level: didtypes.KYC_LEVEL_NONE,
			err:   sdkerrors.ErrInvalidType,
		},
		{
			name:  "rejects unknown enum",
			level: didtypes.KycLevel(999),
			err:   sdkerrors.ErrInvalidType,
		},
		{
			name:  "accepts level one",
			level: didtypes.KYC_LEVEL_ONE,
		},
		{
			name:  "accepts level two",
			level: didtypes.KYC_LEVEL_TWO,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateApproveLevel(tt.level)

			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgApproveValidateBasicRejectsKycLevelNone(t *testing.T) {
	msg := MsgApprove{
		Issuer:   sample.AccAddress(),
		Did:      strings.Repeat("1", didtypes.DidLength),
		RegionId: wstakingtypes.MeEarthRegionName,
		Address:  sample.AccAddress(),
		Pubkey:   "pubkey",
		Level:    didtypes.KYC_LEVEL_NONE,
	}

	require.ErrorIs(t, msg.ValidateBasic(), sdkerrors.ErrInvalidType)
}
