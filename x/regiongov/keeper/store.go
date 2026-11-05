package keeper

import (
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if len(bz) == 0 {
		return types.DefaultParams()
	}
	var params types.Params
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ParamsKey, k.cdc.MustMarshal(&params))
	return nil
}

func (k Keeper) GetProposalSeq(ctx sdk.Context) uint64 {
	return k.getUint64(ctx, types.ProposalSeqKey)
}

func (k Keeper) SetProposalSeq(ctx sdk.Context, seq uint64) {
	k.setUint64(ctx, types.ProposalSeqKey, seq)
}

func (k Keeper) GetPendingActionSeq(ctx sdk.Context) uint64 {
	return k.getUint64(ctx, types.PendingActionSeqKey)
}

func (k Keeper) SetPendingActionSeq(ctx sdk.Context, seq uint64) {
	k.setUint64(ctx, types.PendingActionSeqKey, seq)
}

func (k Keeper) SetProposal(ctx sdk.Context, proposal types.Proposal) {
	k.mustSet(ctx, types.ProposalKey(proposal.Id), &proposal)
}

func (k Keeper) GetProposal(ctx sdk.Context, id uint64) (types.Proposal, bool) {
	var proposal types.Proposal
	found := k.get(ctx, types.ProposalKey(id), &proposal)
	return proposal, found
}

func (k Keeper) GetAllProposals(ctx sdk.Context) []types.Proposal {
	var proposals []types.Proposal
	k.iterate(ctx, types.ProposalPrefix, func(_, value []byte) bool {
		var proposal types.Proposal
		k.cdc.MustUnmarshal(value, &proposal)
		proposals = append(proposals, proposal)
		return false
	})
	return proposals
}

func (k Keeper) SetElector(ctx sdk.Context, elector types.Elector) {
	acc, ok := operatorAccount(elector.Operator)
	if !ok {
		panic(fmt.Sprintf("elector %d has invalid operator %s", elector.ProposalId, elector.Operator))
	}
	k.mustSet(ctx, types.ElectorKey(elector.ProposalId, acc.String()), &elector)
}

func (k Keeper) GetElectors(ctx sdk.Context, proposalID uint64) []types.Elector {
	var electors []types.Elector
	k.iterate(ctx, types.ElectorPrefixByProposal(proposalID), func(_, value []byte) bool {
		var elector types.Elector
		k.cdc.MustUnmarshal(value, &elector)
		electors = append(electors, elector)
		return false
	})
	return electors
}

func (k Keeper) GetAllElectors(ctx sdk.Context) []types.Elector {
	var electors []types.Elector
	k.iterate(ctx, types.ElectorPrefix, func(_, value []byte) bool {
		var elector types.Elector
		k.cdc.MustUnmarshal(value, &elector)
		electors = append(electors, elector)
		return false
	})
	return electors
}

func (k Keeper) GetElector(ctx sdk.Context, proposalID uint64, address string) (types.Elector, bool) {
	var elector types.Elector
	found := k.get(ctx, types.ElectorKey(proposalID, address), &elector)
	return elector, found
}

func (k Keeper) SetVote(ctx sdk.Context, vote types.Vote) {
	k.mustSet(ctx, types.VoteKey(vote.ProposalId, vote.Voter), &vote)
}

func (k Keeper) GetVote(ctx sdk.Context, proposalID uint64, voter string) (types.Vote, bool) {
	var vote types.Vote
	found := k.get(ctx, types.VoteKey(proposalID, voter), &vote)
	return vote, found
}

func (k Keeper) GetVotes(ctx sdk.Context, proposalID uint64) []types.Vote {
	var votes []types.Vote
	k.iterate(ctx, types.VotePrefixByProposal(proposalID), func(_, value []byte) bool {
		var vote types.Vote
		k.cdc.MustUnmarshal(value, &vote)
		votes = append(votes, vote)
		return false
	})
	return votes
}

func (k Keeper) GetAllVotes(ctx sdk.Context) []types.Vote {
	var votes []types.Vote
	k.iterate(ctx, types.VotePrefix, func(_, value []byte) bool {
		var vote types.Vote
		k.cdc.MustUnmarshal(value, &vote)
		votes = append(votes, vote)
		return false
	})
	return votes
}

func (k Keeper) SetEquity(ctx sdk.Context, equity types.RegionEquity) {
	k.mustSet(ctx, types.EquityKey(equity.RegionId), &equity)
}

func (k Keeper) GetEquity(ctx sdk.Context, regionID string) (types.RegionEquity, bool) {
	var equity types.RegionEquity
	found := k.get(ctx, types.EquityKey(regionID), &equity)
	return equity, found
}

func (k Keeper) GetAllEquities(ctx sdk.Context) []types.RegionEquity {
	var equities []types.RegionEquity
	k.iterate(ctx, types.EquityPrefix, func(_, value []byte) bool {
		var equity types.RegionEquity
		k.cdc.MustUnmarshal(value, &equity)
		equities = append(equities, equity)
		return false
	})
	return equities
}

func (k Keeper) SetRegionFreeze(ctx sdk.Context, freeze types.RegionFreeze) {
	k.mustSet(ctx, types.RegionFreezeKey(freeze.RegionId), &freeze)
}

func (k Keeper) GetRegionFreeze(ctx sdk.Context, regionID string) (types.RegionFreeze, bool) {
	var freeze types.RegionFreeze
	found := k.get(ctx, types.RegionFreezeKey(regionID), &freeze)
	return freeze, found
}

func (k Keeper) GetAllRegionFreezes(ctx sdk.Context) []types.RegionFreeze {
	var freezes []types.RegionFreeze
	k.iterate(ctx, types.RegionFreezePrefix, func(_, value []byte) bool {
		var freeze types.RegionFreeze
		k.cdc.MustUnmarshal(value, &freeze)
		freezes = append(freezes, freeze)
		return false
	})
	return freezes
}

func (k Keeper) GetEquityFreeze(ctx sdk.Context, regionID, address string) (types.EquityFreeze, bool) {
	var freeze types.EquityFreeze
	found := k.get(ctx, types.EquityFreezeKey(regionID, address), &freeze)
	return freeze, found
}

func (k Keeper) SetEquityFreeze(ctx sdk.Context, freeze types.EquityFreeze) {
	k.mustSet(ctx, types.EquityFreezeKey(freeze.RegionId, freeze.Address), &freeze)
}

func (k Keeper) GetEquityFreezes(ctx sdk.Context, regionID string) []types.EquityFreeze {
	var freezes []types.EquityFreeze
	k.iterate(ctx, types.EquityFreezePrefixByRegion(regionID), func(_, value []byte) bool {
		var freeze types.EquityFreeze
		k.cdc.MustUnmarshal(value, &freeze)
		freezes = append(freezes, freeze)
		return false
	})
	return freezes
}

func (k Keeper) GetAllEquityFreezes(ctx sdk.Context) []types.EquityFreeze {
	var freezes []types.EquityFreeze
	k.iterate(ctx, types.EquityFreezePrefix, func(_, value []byte) bool {
		var freeze types.EquityFreeze
		k.cdc.MustUnmarshal(value, &freeze)
		freezes = append(freezes, freeze)
		return false
	})
	return freezes
}

func (k Keeper) SetPendingAction(ctx sdk.Context, action types.PendingAction) {
	k.mustSet(ctx, types.PendingActionKey(action.Id), &action)
}

func (k Keeper) GetPendingAction(ctx sdk.Context, id uint64) (types.PendingAction, bool) {
	var action types.PendingAction
	found := k.get(ctx, types.PendingActionKey(id), &action)
	return action, found
}

func (k Keeper) GetAllPendingActions(ctx sdk.Context) []types.PendingAction {
	var actions []types.PendingAction
	k.iterate(ctx, types.PendingActionPrefix, func(_, value []byte) bool {
		var action types.PendingAction
		k.cdc.MustUnmarshal(value, &action)
		actions = append(actions, action)
		return false
	})
	return actions
}

func (k Keeper) mustSet(ctx sdk.Context, key []byte, msg codec.ProtoMarshaler) {
	ctx.KVStore(k.storeKey).Set(key, k.cdc.MustMarshal(msg))
}

func (k Keeper) get(ctx sdk.Context, key []byte, msg codec.ProtoMarshaler) bool {
	bz := ctx.KVStore(k.storeKey).Get(key)
	if len(bz) == 0 {
		return false
	}
	k.cdc.MustUnmarshal(bz, msg)
	return true
}

func (k Keeper) iterate(ctx sdk.Context, prefix []byte, cb func(key, value []byte) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		if cb(iterator.Key(), iterator.Value()) {
			break
		}
	}
}

func (k Keeper) setUint64(ctx sdk.Context, key []byte, value uint64) {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, value)
	ctx.KVStore(k.storeKey).Set(key, bz)
}

func (k Keeper) getUint64(ctx sdk.Context, key []byte) uint64 {
	bz := ctx.KVStore(k.storeKey).Get(key)
	if len(bz) == 0 {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}
