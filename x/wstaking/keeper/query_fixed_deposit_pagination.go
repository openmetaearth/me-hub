package keeper

import (
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func normalizeFixedDepositByRegionPagination(req *types.QueryFixedDepositByRegionRequest) {
	if req.Pagination == nil {
		req.Pagination = &query.PageRequest{}
	}
	if req.Pagination.Limit == 0 {
		req.Pagination.Limit = query.DefaultLimit
		req.Pagination.CountTotal = true
	}
}
