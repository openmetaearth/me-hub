package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/stretchr/testify/require"
)

func validMsgMintNFT(tokenID string) MsgMintNFT {
	return MsgMintNFT{
		ClassId:  "class1",
		TokenId:  tokenID,
		Uri:      "ipfs://token",
		UriHash:  "hash",
		Creator:  sample.AccAddress(),
		Receiver: sample.AccAddress(),
	}
}

func TestMsgMintNFTValidateBasicRejectsNonCanonicalTokenIDs(t *testing.T) {
	tests := []struct {
		name    string
		tokenID string
		err     error
	}{
		{
			name:    "valid",
			tokenID: "1",
		},
		{
			name:    "leading zero",
			tokenID: "01",
			err:     sdkerrors.ErrInvalidRequest,
		},
		{
			name:    "multiple leading zeros",
			tokenID: "001",
			err:     sdkerrors.ErrInvalidRequest,
		},
		{
			name:    "zero",
			tokenID: "0",
			err:     sdkerrors.ErrInvalidRequest,
		},
		{
			name:    "not numeric",
			tokenID: "one",
			err:     sdkerrors.ErrInvalidRequest,
		},
		{
			name:    "empty",
			tokenID: "",
			err:     ErrEmptyTokenId,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validMsgMintNFT(tt.tokenID).ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
