package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

	address, err := sdk.AccAddressFromBech32(r.Address)
	if err != nil {
		return nil, err
	}

	nfts := k.GetNFTsOfClassByOwner(ctx, r.ClassId, address)

	var tokenIds []string
	for _, nft := range nfts {
		tokenIds = append(tokenIds, nft.Id)
	}

	return &types.QueryClassAddressResponse{
		Exists:      true,
		TotalSupply: class.TotalSupply,
		Nfts:        tokenIds,
	}, nil
}

func (k Keeper) NftFilter(goCtx context.Context, r *types.QueryNftFilterRequest) (*types.QueryNftFilterResponse, error) {
	if r == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var list []*types.NftList

	var owner sdk.AccAddress
	if r.Owner != "" {
		var err error
		owner, err = sdk.AccAddressFromBech32(r.Owner)
		if err != nil {
			return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address %q: %v", r.Owner, err)
		}
	}

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

		return &types.QueryNftFilterResponse{
			Nfts: list,
		}, nil

	} else if r.ClassId != "" && r.Owner != "" && r.TokenId == "" {
		// query the holdings of a specific type of nft
		_, has := k.GetClass(ctx, r.ClassId)
		if !has {
			return nil, nil
		}

		nftInfos := k.GetNFTsOfClassByOwner(ctx, r.ClassId, owner)
		for _, nftInfo := range nftInfos {
			list = append(list, &types.NftList{
				ClassId: nftInfo.ClassId,
				TokenId: nftInfo.Id,
				Owner:   r.Owner,
				Uri:     nftInfo.Uri,
			})
		}

		return &types.QueryNftFilterResponse{
			Nfts: list,
		}, nil

	} else if r.Owner != "" && r.TokenId == "" && r.ClassId == "" {
		// query the nft information held by the address
		classes := k.GetClasses(ctx)
		for _, class := range classes {
			nftInfos := k.GetNFTsOfClassByOwner(ctx, class.Id, owner)
			for _, nftInfo := range nftInfos {
				nftOwner := k.GetOwner(ctx, nftInfo.ClassId, nftInfo.Id)
				list = append(list, &types.NftList{
					ClassId: nftInfo.ClassId,
					TokenId: nftInfo.Id,
					Owner:   nftOwner.String(),
					Uri:     nftInfo.Uri,
				})
			}
		}

		return &types.QueryNftFilterResponse{
			Nfts: list,
		}, nil

	}
	return nil, sdkerrors.ErrInvalidRequest.Wrap("invalid query parameters")
}
