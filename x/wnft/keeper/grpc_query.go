package keeper

import (
	"context"

	"cosmossdk.io/x/nft"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/openmetaearth/me-hub/x/wnft/types"
)

type Querier struct {
	*Keeper
}

var _ types.QueryServer = Keeper{}

// Classes return all NFT classes
func (k Keeper) ClassAddress(goCtx context.Context, r *types.QueryClassAddressRequest) (*types.QueryClassAddressResponse, error) {
	if r == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	class, ok := k.GetClass(ctx, r.ClassId)
	if !ok {
		return &types.QueryClassAddressResponse{Exists: false}, nil
	}

	nftResp, err := k.Keeper.NFTs(goCtx, &nft.QueryNFTsRequest{
		ClassId:    r.ClassId,
		Owner:      r.Address,
		Pagination: r.Pagination,
	})
	if err != nil {
		return nil, err
	}

	tokenIds := make([]string, 0, len(nftResp.Nfts))
	for _, ownerNFT := range nftResp.Nfts {
		tokenIds = append(tokenIds, ownerNFT.Id)
	}

	return &types.QueryClassAddressResponse{
		Exists:      true,
		TotalSupply: k.GetClassTotalSupplyCap(ctx, class.Id),
		Nfts:        tokenIds,
		Pagination:  nftResp.Pagination,
	}, nil
}

func (k Keeper) NftFilter(goCtx context.Context, r *types.QueryNftFilterRequest) (*types.QueryNftFilterResponse, error) {
	if r == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	var list []*types.NftList
	pageRes := &query.PageResponse{}

	// determine query type based on request parameters
	if r.TokenId != "" && r.ClassId != "" && r.Owner == "" {
		// query individual nft information
		_, has := k.GetClass(ctx, r.ClassId)
		if !has {
			return nil, nil
		}

		nftInfo, ok := k.GetNFT(ctx, r.ClassId, r.TokenId)
		if !ok {
			return nil, nil
		}

		owner := k.GetOwner(ctx, r.ClassId, r.TokenId)

		list = append(list, &types.NftList{
			ClassId: nftInfo.ClassId,
			TokenId: nftInfo.Id,
			Owner:   owner.String(),
			Uri:     nftInfo.Uri,
		})

		return &types.QueryNftFilterResponse{Nfts: list, Pagination: pageRes}, nil
	} else if r.ClassId != "" && r.Owner != "" && r.TokenId == "" {
		// query the holdings of a specific type of nft
		_, has := k.GetClass(ctx, r.ClassId)
		if !has {
			return nil, nil
		}
		nftResp, err := k.Keeper.NFTs(goCtx, &nft.QueryNFTsRequest{
			ClassId:    r.ClassId,
			Owner:      r.Owner,
			Pagination: r.Pagination,
		})
		if err != nil {
			return nil, err
		}

		for _, nftInfo := range nftResp.Nfts {
			list = append(list, &types.NftList{
				ClassId: nftInfo.ClassId,
				TokenId: nftInfo.Id,
				Owner:   r.Owner,
				Uri:     nftInfo.Uri,
			})
		}

		return &types.QueryNftFilterResponse{Nfts: list, Pagination: nftResp.Pagination}, nil
	} else if r.Owner != "" && r.TokenId == "" && r.ClassId == "" {
		// query the nft information held by the address
		nftResp, err := k.Keeper.NFTs(goCtx, &nft.QueryNFTsRequest{
			Owner:      r.Owner,
			Pagination: r.Pagination,
		})
		if err != nil {
			return nil, err
		}

		for _, nftInfo := range nftResp.Nfts {
			owner := k.GetOwner(ctx, nftInfo.ClassId, nftInfo.Id)
			list = append(list, &types.NftList{
				ClassId: nftInfo.ClassId,
				TokenId: nftInfo.Id,
				Owner:   owner.String(),
				Uri:     nftInfo.Uri,
			})
		}

		return &types.QueryNftFilterResponse{Nfts: list, Pagination: nftResp.Pagination}, nil
	}
	return nil, nil
}
