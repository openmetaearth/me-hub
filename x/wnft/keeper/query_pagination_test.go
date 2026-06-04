package keeper

import (
	"bytes"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
)

func TestSafeWNFTPageRequestCapsLimitAndDisablesTotal(t *testing.T) {
	require.Equal(t, maxWNFTQueryLimit, safeWNFTPageRequest(nil).Limit)

	req := &query.PageRequest{
		Limit:      maxWNFTQueryLimit + 1,
		CountTotal: true,
	}
	got := safeWNFTPageRequest(req)

	require.Equal(t, maxWNFTQueryLimit, got.Limit)
	require.False(t, got.CountTotal)
	require.Equal(t, maxWNFTQueryLimit+1, req.Limit)
	require.True(t, req.CountTotal)
}

func TestParseNFTOfClassByOwnerKey(t *testing.T) {
	classID := "album"
	tokenID := "nft-1"
	gotClassID, gotTokenID, err := parseNFTOfClassByOwnerKey([]byte(classID + string(wnftDelimiter) + tokenID))

	require.NoError(t, err)
	require.Equal(t, classID, gotClassID)
	require.Equal(t, tokenID, gotTokenID)
}

func TestNFTOfClassByOwnerStoreKeyUsesOwnerClassPrefix(t *testing.T) {
	owner := sdk.AccAddress(bytes.Repeat([]byte{0x01}, 20))
	key := nftOfClassByOwnerStoreKey(owner, "album")

	require.True(t, bytes.HasPrefix(key, wnftNFTOfClassByOwnerKey))
	require.True(t, bytes.HasSuffix(key, append([]byte("album"), wnftDelimiter...)))
	require.True(t, bytes.HasPrefix(key, prefixNftOfClassByOwnerStoreKey(owner)))
}
