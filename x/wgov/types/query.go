package types

import (
	"context"

	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

type QueryMeTallyResultRequest struct {
	ProposalId uint64
}

type QueryMeTallyResultResponse struct {
	Tally *v1.TallyResult
}

type QueryServer interface {
	MeTallyResult(context.Context, *QueryMeTallyResultRequest) (*QueryMeTallyResultResponse, error)
}
