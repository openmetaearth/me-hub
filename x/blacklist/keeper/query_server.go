package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/st-chain/me-hub/x/blacklist/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Blacklist(goCtx context.Context, req *types.QueryBlacklistRequest) (*types.QueryBlacklistResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	blacklist, found := k.GetBlacklist(ctx)
	if !found {
		blacklist = types.Blacklist{Addresses: []string{}}
	}

	// If a specific address is requested, check if it exists
	if req.Address != "" {
		var respAddresses []string
		for _, addr := range blacklist.Addresses {
			if addr == req.Address {
				respAddresses = []string{req.Address}
				break
			}
		}
		return &types.QueryBlacklistResponse{
			Blacklist: types.Blacklist{Addresses: respAddresses},
		}, nil
	}

	// If no specific address is requested, return all with pagination
	addresses := blacklist.Addresses

	pageReq := req.Pagination
	if pageReq == nil {
		pageReq = &query.PageRequest{}
	}

	offset := pageReq.Offset
	limit := pageReq.Limit
	if limit == 0 {
		limit = 100 // Default limit
	}

	total := uint64(len(addresses))
	if offset >= total {
		return &types.QueryBlacklistResponse{
			Blacklist:  types.Blacklist{Addresses: []string{}},
			Pagination: &query.PageResponse{Total: total},
		}, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	paginatedAddresses := addresses[offset:end]

	return &types.QueryBlacklistResponse{
		Blacklist:  types.Blacklist{Addresses: paginatedAddresses},
		Pagination: &query.PageResponse{Total: total},
	}, nil
}
