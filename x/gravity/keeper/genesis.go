package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/st-chain/me-hub/x/gravity/types"
)

// InitGenesis import module genesis
//
//gocyclo:ignore
func InitGenesis(ctx sdk.Context, k Keeper, state *types.GenesisState) {
	if err := k.SetParams(ctx, &state.Params); err != nil {
		panic(err)
	}

	k.SetLastObservedEventNonce(ctx, state.LastObservedEventNonce)
	k.SetLastObservedBlockHeight(ctx, state.LastObservedBlockHeight.ExternalBlockHeight, state.LastObservedBlockHeight.BlockHeight)
	k.SetProposalRelayer(ctx, &state.ProposalRelayer)
	k.SetLastObservedRelayerSet(ctx, &state.LastObservedRelayerSet)
	k.SetLastSlashedRelayerSetNonce(ctx, state.LastSlashedRelayerSetNonce)
	k.SetLastSlashedBatchBlock(ctx, state.LastSlashedBatchBlock)
	for _, relayer := range state.Relayers {
		relayerAddress := sdk.MustAccAddressFromBech32(relayer.RelayerAddress)
		k.SetRelayer(ctx, relayerAddress, relayer)
		k.SetRelayerByExternalAddress(ctx, relayer.ExternalAddress, relayerAddress)
	}
	k.SetLastTotalPower(ctx)

	latestRelayerSetNonce := uint64(0)
	for i := 0; i < len(state.RelayerSets); i++ {
		set := state.RelayerSets[i]
		if set.Nonce > latestRelayerSetNonce {
			latestRelayerSetNonce = set.Nonce
		}
		k.StoreRelayerSet(ctx, &set)
	}
	k.SetLastRelayerSetNonce(ctx, latestRelayerSetNonce)

	for _, bridgeToken := range state.BridgeTokens {
		k.SetBridgeToken(ctx, &bridgeToken)
	}

	for i := 0; i < len(state.BatchConfirms); i++ {
		confirm := state.BatchConfirms[i]
		for _, relayer := range state.Relayers {
			if confirm.RelayerAddress == relayer.RelayerAddress {
				k.SetBatchConfirm(ctx, relayer.GetRelayer(), &confirm)
			}
		}
	}
	for i := 0; i < len(state.RelayerSetConfirms); i++ {
		confirm := state.RelayerSetConfirms[i]
		for _, relayer := range state.Relayers {
			if confirm.RelayerAddress == relayer.GetRelayerAddress() {
				k.SetRelayerSetConfirm(ctx, relayer.GetRelayer(), &confirm)
			}
		}
	}

	for i := 0; i < len(state.UnbatchedTransfers); i++ {
		transfer := state.UnbatchedTransfers[i]
		if err := k.AddUnbatchedTx(ctx, &transfer); err != nil {
			panic(err)
		}
	}

	for i := 0; i < len(state.Batches); i++ {
		batch := state.Batches[i]
		if err := k.StoreBatch(ctx, &batch); err != nil {
			panic(err)
		}
	}

	// reset attestations in state
	for i := 0; i < len(state.Attestations); i++ {
		att := state.Attestations[i]
		claim, err := types.UnpackAttestationClaim(k.cdc, &att)
		if err != nil {
			panic("couldn't cast to claim")
		}

		k.SetAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash(), &att)
	}

	// reset attestation state of specific validators
	// this must be done after the above to be correct
	for i := 0; i < len(state.Attestations); i++ {
		att := state.Attestations[i]
		claim, err := types.UnpackAttestationClaim(k.cdc, &att)
		if err != nil {
			panic("couldn't cast to claim")
		}
		// reconstruct the latest event nonce for every validator
		// if somehow this genesis state is saved when all attestations
		// have been cleaned up GetLastEventNonceByRelayer handles that case
		//
		// if we where to save and load the last event nonce for every validator
		// then we would need to carry that state forever across all chain restarts
		// but since we've already had to handle the edge case of new validators joining
		// while all attestations have already been cleaned up we can do this instead and
		// not carry around every validators event nonce counter forever.
		for _, vote := range att.Votes {
			relayer := sdk.MustAccAddressFromBech32(vote)
			last := k.GetLastEventNonceByRelayer(ctx, relayer)
			if claim.GetEventNonce() > last {
				k.SetLastEventNonceByRelayer(ctx, relayer, claim.GetEventNonce())
				k.SetLastEventBlockHeightByRelayer(ctx, relayer, claim.GetBlockHeight())
			}
		}
	}
}

// ExportGenesis export module status
func ExportGenesis(ctx sdk.Context, k Keeper) *types.GenesisState {
	state := &types.GenesisState{
		Params:                  k.GetParams(ctx),
		LastObservedEventNonce:  k.GetLastObservedEventNonce(ctx),
		LastObservedBlockHeight: k.GetLastObservedBlockHeight(ctx),
	}
	k.IterateRelayer(ctx, func(re types.Relayer) bool {
		state.Relayers = append(state.Relayers, re)
		return false
	})
	k.IterateRelayerSets(ctx, false, func(relayerSet *types.RelayerSet) bool {
		state.RelayerSets = append(state.RelayerSets, *relayerSet)
		return false
	})
	k.IterateOutgoingTxBatches(ctx, func(batch *types.OutgoingTxBatch) bool {
		state.Batches = append(state.Batches, *batch)
		return false
	})
	k.IterateAttestations(ctx, func(attestation *types.Attestation) bool {
		state.Attestations = append(state.Attestations, *attestation)
		return false
	})
	k.IterateUnbatchedTransactions(ctx, "", func(tx *types.OutgoingTransferTx) bool {
		state.UnbatchedTransfers = append(state.UnbatchedTransfers, *tx)
		return false
	})
	for _, vs := range state.RelayerSets {
		k.IterateRelayerSetConfirmByNonce(ctx, vs.Nonce, func(confirm *types.MsgRelayerSetConfirm) bool {
			state.RelayerSetConfirms = append(state.RelayerSetConfirms, *confirm)
			return false
		})
	}
	for _, batch := range state.Batches {
		k.IterateBatchConfirmByNonceAndTokenContract(ctx, batch.BatchNonce, batch.TokenContract, func(confirm *types.MsgConfirmBatch) bool {
			state.BatchConfirms = append(state.BatchConfirms, *confirm)
			return false
		})
	}
	k.IterateBridgeTokenByDenom(ctx, func(erc20ToDenom *types.BridgeToken) bool {
		state.BridgeTokens = append(state.BridgeTokens, *erc20ToDenom)
		return false
	})
	state.ProposalRelayer, _ = k.GetProposalRelayer(ctx)
	if lastObserved := k.GetLastObservedRelayerSet(ctx); lastObserved != nil {
		state.LastObservedRelayerSet = *lastObserved
	}
	state.LastSlashedBatchBlock = k.GetLastSlashedBatchBlock(ctx)
	state.LastSlashedRelayerSetNonce = k.GetLastSlashedRelayerSetNonce(ctx)
	return state
}
