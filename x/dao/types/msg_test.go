package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	_ "github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateDaoValidateBasic(t *testing.T) {
	validAddresses := DaoAddresses{
		GlobalDao:      sample.AccAddress(),
		MeidDao:        sample.AccAddress(),
		DevOperator:    sample.AccAddress(),
		AirdropAddress: sample.AccAddress(),
	}

	tests := []struct {
		name      string
		creator   string
		addresses DaoAddresses
		wantError string
	}{
		{
			name:      "valid dao address update",
			creator:   sample.AccAddress(),
			addresses: validAddresses,
		},
		{
			name:      "invalid creator",
			creator:   "not-a-bech32-address",
			addresses: validAddresses,
			wantError: "not-a-bech32-address",
		},
		{
			name:      "all dao addresses empty",
			creator:   sample.AccAddress(),
			addresses: DaoAddresses{},
			wantError: "global dao address is empty",
		},
		{
			name:    "empty global dao",
			creator: sample.AccAddress(),
			addresses: DaoAddresses{
				MeidDao:        validAddresses.MeidDao,
				DevOperator:    validAddresses.DevOperator,
				AirdropAddress: validAddresses.AirdropAddress,
			},
			wantError: "global dao address is empty",
		},
		{
			name:    "empty meid dao",
			creator: sample.AccAddress(),
			addresses: DaoAddresses{
				GlobalDao:      validAddresses.GlobalDao,
				DevOperator:    validAddresses.DevOperator,
				AirdropAddress: validAddresses.AirdropAddress,
			},
			wantError: "meid dao address is empty",
		},
		{
			name:    "empty dev operator",
			creator: sample.AccAddress(),
			addresses: DaoAddresses{
				GlobalDao:      validAddresses.GlobalDao,
				MeidDao:        validAddresses.MeidDao,
				AirdropAddress: validAddresses.AirdropAddress,
			},
			wantError: "dev operator address is empty",
		},
		{
			name:    "empty airdrop address",
			creator: sample.AccAddress(),
			addresses: DaoAddresses{
				GlobalDao:   validAddresses.GlobalDao,
				MeidDao:     validAddresses.MeidDao,
				DevOperator: validAddresses.DevOperator,
			},
			wantError: "airdrop address is empty",
		},
		{
			name:    "invalid dao address",
			creator: sample.AccAddress(),
			addresses: DaoAddresses{
				GlobalDao:      validAddresses.GlobalDao,
				MeidDao:        "not-a-bech32-address",
				DevOperator:    validAddresses.DevOperator,
				AirdropAddress: validAddresses.AirdropAddress,
			},
			wantError: "meid dao address not-a-bech32-address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgUpdateDao{
				Creator:      tt.creator,
				DaoAddresses: tt.addresses,
			}

			err := msg.ValidateBasic()
			if tt.wantError != "" {
				require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
				require.ErrorContains(t, err, tt.wantError)
				return
			}

			require.NoError(t, err)
		})
	}
}
