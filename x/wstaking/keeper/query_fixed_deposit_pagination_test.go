package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"github.com/stretchr/testify/require"
)

func TestNormalizeFixedDepositByRegionPagination(t *testing.T) {
	tests := []struct {
		name       string
		pagination *query.PageRequest
		wantLimit  uint64
		wantTotal  bool
	}{
		{
			name:      "missing pagination uses sdk default",
			wantLimit: query.DefaultLimit,
			wantTotal: true,
		},
		{
			name:       "zero limit uses sdk default",
			pagination: &query.PageRequest{},
			wantLimit:  query.DefaultLimit,
			wantTotal:  true,
		},
		{
			name:       "explicit limit remains unchanged",
			pagination: &query.PageRequest{Limit: 5},
			wantLimit:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &types.QueryFixedDepositByRegionRequest{
				Pagination: tt.pagination,
			}

			normalizeFixedDepositByRegionPagination(req)

			require.NotNil(t, req.Pagination)
			require.Equal(t, tt.wantLimit, req.Pagination.Limit)
			require.Equal(t, tt.wantTotal, req.Pagination.CountTotal)
		})
	}
}
