package regiongov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/regiongov/keeper"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

// InitGenesis initializes the module state from genesis.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	if err := genState.Validate(); err != nil {
		panic(err)
	}
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}
	k.SetProposalSeq(ctx, genState.ProposalSeq)
	k.SetPendingActionSeq(ctx, genState.PendingActionSeq)
	for _, proposal := range genState.Proposals {
		k.SetProposal(ctx, proposal)
		if proposal.Status == types.ProposalStatusVoting {
			k.InsertVotingQueue(ctx, proposal.VotingEndTime, proposal.Id)
		}
	}
	for _, elector := range genState.Electors {
		k.SetElector(ctx, elector)
	}
	for _, vote := range genState.Votes {
		k.SetVote(ctx, vote)
	}
	for _, equity := range genState.Equities {
		k.SetEquity(ctx, equity)
	}
	for _, freeze := range genState.RegionFreezes {
		k.SetRegionFreeze(ctx, freeze)
	}
	for _, freeze := range genState.EquityFreezes {
		k.SetEquityFreeze(ctx, freeze)
	}
	for _, action := range genState.PendingActions {
		k.SetPendingAction(ctx, action)
	}
}

// ExportGenesis returns the current module state.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params:           k.GetParams(ctx),
		ProposalSeq:      k.GetProposalSeq(ctx),
		PendingActionSeq: k.GetPendingActionSeq(ctx),
		Proposals:        k.GetAllProposals(ctx),
		Electors:         k.GetAllElectors(ctx),
		Votes:            k.GetAllVotes(ctx),
		Equities:         k.GetAllEquities(ctx),
		RegionFreezes:    k.GetAllRegionFreezes(ctx),
		EquityFreezes:    k.GetAllEquityFreezes(ctx),
		PendingActions:   k.GetAllPendingActions(ctx),
	}
}
