package keeper

import (
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

func paginate[T any](items []T, req *query.PageRequest) ([]T, *query.PageResponse) {
	if req == nil {
		return items, nil
	}
	paged := append([]T(nil), items...)
	if req.Reverse {
		for i, j := 0, len(paged)-1; i < j; i, j = i+1, j-1 {
			paged[i], paged[j] = paged[j], paged[i]
		}
	}
	limit := req.Limit
	if limit == 0 {
		limit = query.DefaultLimit
	}
	start := req.Offset
	if start > uint64(len(paged)) {
		start = uint64(len(paged))
	}
	end := start + limit
	if end > uint64(len(paged)) {
		end = uint64(len(paged))
	}
	resp := &query.PageResponse{}
	if req.CountTotal {
		resp.Total = uint64(len(items))
	}
	return paged[start:end], resp
}

func filterProposals(proposals []types.Proposal, status types.ProposalStatus, kind types.ProposalKind) []types.Proposal {
	if status == types.ProposalStatusUnspecified && kind == types.ProposalKindUnspecified {
		return proposals
	}
	filtered := make([]types.Proposal, 0, len(proposals))
	for _, proposal := range proposals {
		if status != types.ProposalStatusUnspecified && proposal.Status != status {
			continue
		}
		if kind != types.ProposalKindUnspecified && proposal.Kind != kind {
			continue
		}
		filtered = append(filtered, proposal)
	}
	return filtered
}

func filterActions(actions []types.PendingAction, status types.PendingActionStatus) []types.PendingAction {
	if status == types.PendingActionStatusUnspecified {
		return actions
	}
	filtered := make([]types.PendingAction, 0, len(actions))
	for _, action := range actions {
		if action.Status == status {
			filtered = append(filtered, action)
		}
	}
	return filtered
}
