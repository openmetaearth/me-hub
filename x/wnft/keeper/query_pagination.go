package keeper

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/nft"
)

const maxWNFTQueryLimit uint64 = 100

var (
	wnftNFTOfClassByOwnerKey = []byte{0x03}
	wnftDelimiter            = []byte{0x00}
)

func safeWNFTPageRequest(pageReq *query.PageRequest) *query.PageRequest {
	if pageReq == nil {
		return &query.PageRequest{Limit: maxWNFTQueryLimit}
	}

	safeReq := *pageReq
	if safeReq.Limit == 0 || safeReq.Limit > maxWNFTQueryLimit {
		safeReq.Limit = maxWNFTQueryLimit
	}
	safeReq.CountTotal = false

	return &safeReq
}

func (k Keeper) paginateNFTIDsOfClassByOwner(
	ctx sdk.Context,
	classID string,
	owner sdk.AccAddress,
	pageReq *query.PageRequest,
) ([]string, *query.PageResponse, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), nftOfClassByOwnerStoreKey(owner, classID))

	var tokenIDs []string
	pageRes, err := query.Paginate(store, safeWNFTPageRequest(pageReq), func(key []byte, _ []byte) error {
		tokenIDs = append(tokenIDs, string(key))
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return tokenIDs, pageRes, nil
}

func (k Keeper) paginateNFTsOfClassByOwner(
	ctx sdk.Context,
	classID string,
	owner sdk.AccAddress,
	pageReq *query.PageRequest,
) ([]nft.NFT, *query.PageResponse, error) {
	tokenIDs, pageRes, err := k.paginateNFTIDsOfClassByOwner(ctx, classID, owner, pageReq)
	if err != nil {
		return nil, nil, err
	}

	nfts := make([]nft.NFT, 0, len(tokenIDs))
	for _, tokenID := range tokenIDs {
		nftInfo, found := k.GetNFT(ctx, classID, tokenID)
		if found {
			nfts = append(nfts, nftInfo)
		}
	}

	return nfts, pageRes, nil
}

func (k Keeper) paginateNFTsByOwner(
	ctx sdk.Context,
	owner sdk.AccAddress,
	pageReq *query.PageRequest,
) ([]nft.NFT, *query.PageResponse, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixNftOfClassByOwnerStoreKey(owner))

	var nfts []nft.NFT
	pageRes, err := query.Paginate(store, safeWNFTPageRequest(pageReq), func(key []byte, _ []byte) error {
		classID, tokenID, err := parseNFTOfClassByOwnerKey(key)
		if err != nil {
			return err
		}

		nftInfo, found := k.GetNFT(ctx, classID, tokenID)
		if found {
			nfts = append(nfts, nftInfo)
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return nfts, pageRes, nil
}

func nftOfClassByOwnerStoreKey(owner sdk.AccAddress, classID string) []byte {
	owner = address.MustLengthPrefix(owner)

	key := make([]byte, 0, len(wnftNFTOfClassByOwnerKey)+len(owner)+len(wnftDelimiter)+len(classID)+len(wnftDelimiter))
	key = append(key, wnftNFTOfClassByOwnerKey...)
	key = append(key, owner...)
	key = append(key, wnftDelimiter...)
	key = append(key, []byte(classID)...)
	key = append(key, wnftDelimiter...)
	return key
}

func prefixNftOfClassByOwnerStoreKey(owner sdk.AccAddress) []byte {
	owner = address.MustLengthPrefix(owner)

	key := make([]byte, 0, len(wnftNFTOfClassByOwnerKey)+len(owner)+len(wnftDelimiter))
	key = append(key, wnftNFTOfClassByOwnerKey...)
	key = append(key, owner...)
	key = append(key, wnftDelimiter...)
	return key
}

func parseNFTOfClassByOwnerKey(key []byte) (string, string, error) {
	classID, tokenID, found := bytes.Cut(key, wnftDelimiter)
	if !found || len(classID) == 0 || len(tokenID) == 0 {
		return "", "", fmt.Errorf("invalid nft owner index key")
	}

	return string(classID), string(tokenID), nil
}
