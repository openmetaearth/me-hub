package keeper

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sort"

	"github.com/openmetaearth/me-hub/x/gravity/types"
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

// ClearGenesis clears module state just for test environment
func (k Keeper) ClearGenesis(ctx sdk.Context) {
	//genesis := gravitykeeper.ExportGenesis(ctx, k)
	k.IterateOutgoingTxBatches(ctx, func(batch *types.OutgoingTxBatch) bool {
		k.DeleteBatch(ctx, batch)
		return false
	})
	k.SetLastObservedEventNonce(ctx, 0)
	k.SetLastObservedBlockHeight(ctx, 0, 0)

	claimMap := make(map[uint64][]types.ExternalClaim)
	var nonces []uint64
	k.IterateAttestationAndClaim(ctx, func(att *types.Attestation, claim types.ExternalClaim) bool {
		if v, ok := claimMap[claim.GetEventNonce()]; !ok {
			claimMap[claim.GetEventNonce()] = []types.ExternalClaim{claim}
			nonces = append(nonces, claim.GetEventNonce())
		} else {
			claimMap[claim.GetEventNonce()] = append(v, claim)
		}
		return false
	})
	// Then we sort it
	sort.Slice(nonces, func(i, j int) bool {
		return nonces[i] < nonces[j]
	})

	// This iterates over all keys (event nonces) in the attestation mapping. Each value contains
	// a slice with one or more attestations at that event nonce. There can be multiple attestations
	// at one event nonce when Relayers disagree about what event happened at that nonce.
	for _, nonce := range nonces {
		// This iterates over all attestations at a particular event nonce.
		// They are ordered by when the first attestation at the event nonce was received.
		// This order is not important.
		for _, claim := range claimMap[nonce] {
			k.DeleteAttestation(ctx, claim)
		}
	}

	k.IterateUnbatchedTransactions(ctx, "", func(tx *types.OutgoingTransferTx) bool {
		err := k.DelUnbatchedTx(ctx, tx.Fee, tx.Id)
		if err != nil {
			panic(err)
		}
		return false
	})

	relayerSets := []types.RelayerSet{}
	k.IterateRelayerSets(ctx, false, func(relayerSet *types.RelayerSet) bool {
		relayerSets = append(relayerSets, *relayerSet)
		return false
	})

	for _, rs := range relayerSets {
		k.DeleteRelayerSetConfirm(ctx, rs.Nonce)
	}
	nextID := k.AutoIncrementID(ctx, types.KeyLastOutgoingBatchID)
	k.IterateBridgeTokenByDenom(ctx, func(token *types.BridgeToken) bool {
		k.DelBridgeToken(ctx, token)
		for i := uint64(0); i < nextID; i++ {
			k.DeleteBatchConfirm(ctx, i, token.ContractAddress)
		}
		return false
	})
	k.ClearAutoIncrementID(ctx)

	if lastObserved := k.GetLastObservedRelayerSet(ctx); lastObserved != nil {
		k.DelLastObservedRelayerSet(ctx)
	}
	relayers := k.GetAllRelayers(ctx, false)
	for _, relayer := range relayers {
		k.DelLastEventNonceByRelayer(ctx, sdk.MustAccAddressFromBech32(relayer.RelayerAddress))
	}
	return
}

func (k Keeper) ResetGenesis(ctx sdk.Context) {
	k.SetLastObservedEventNonce(ctx, 317)
	k.SetLastObservedBlockHeight(ctx, 77790848, 0)

	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyLastOutgoingBatchID, sdk.Uint64ToBigEndian(20))
	store.Set(types.KeyLastTxPoolID, sdk.Uint64ToBigEndian(61))
	store.Set(types.LastSlashedBatchBlock, sdk.Uint64ToBigEndian(9580055))
	store.Set(types.LastSlashedRelayerSetNonce, sdk.Uint64ToBigEndian(1))

	set := k.GetCurrentRelayerSet(ctx)
	k.SetLastObservedRelayerSet(ctx, set)

	relayers := k.GetAllRelayers(ctx, false)
	for _, relayer := range relayers {
		if relayer.RelayerAddress == "mechain1qql8qg3k7f5g0j6x5m6j4x5l5p5u3z5f4h3g4" {
			k.SetLastEventNonceByRelayer(ctx, sdk.MustAccAddressFromBech32(relayer.RelayerAddress), 271)
		} else {
			k.SetLastEventNonceByRelayer(ctx, sdk.MustAccAddressFromBech32(relayer.RelayerAddress), 317)
		}
	}

	bridgeTokens := `{
	"bridge_tokens": [{
		"contract_address": "0x676E1ba786f36cd8fB06d2C6332Eb0cd3f1737f9",
		"denom": "usdd",
		"name": "MyToken",
		"symbol": "USDD",
		"decimal": 6,
		"supply": "1000000"
	}, {
		"contract_address": "0xB9cdEc4F2938Bd0447ffE65fDba30f987D77D85e",
		"denom": "usdt",
		"name": "Tether USD",
		"symbol": "USDT",
		"decimal": 6,
		"supply": "1331503666"
	}, {
		"contract_address": "0x8c2ee7E028b6cf64cdE59D64b95FdB3150Afad12",
		"denom": "usdt_5over26yqfevdp27vus2sp",
		"name": "MyToken",
		"symbol": "usdt_5oveR26yqfEVDP27Vus2sP",
		"decimal": 6,
		"supply": "1000000"
	}, {
		"contract_address": "0x48EFE309ad9c2cb62Ac7c45ab6102d7E38B1243f",
		"denom": "usdt_eusg6vr2r9mbat95wazqyu",
		"name": "MyToken",
		"symbol": "usdt_eUSG6vR2R9mbAT95WaZQYU",
		"decimal": 6,
		"supply": "1000000"
	}, {
		"contract_address": "0x58D47096700b01b275FFF39C9D0A642950a0D793",
		"denom": "usdt_gmhtbb5cgvahwn9kp7mlhd",
		"name": "MyToken",
		"symbol": "usdt_gmhtbb5cgvahwn9kp7mlhd",
		"decimal": 6,
		"supply": "1000000"
	}, {
		"contract_address": "0x46A7Ef398d415722eA2b619AA9bd73B4A334c886",
		"denom": "usdt_juxtbtednswqwvdwnfezrn",
		"name": "MyToken",
		"symbol": "usdt_juxtBTEDNsWQWVdWNFezRN",
		"decimal": 6,
		"supply": "1000000"
	}, {
		"contract_address": "0x7616d918F3775c7AB8Dd3d2F188dc65D55e33b5c",
		"denom": "usdt_nm5bzv6xvxwfbssbws4sxt",
		"name": "MyToken",
		"symbol": "usdt_Nm5BzV6XvXWfbsSbWS4Sxt",
		"decimal": 6,
		"supply": "1000000"
	}, {
		"contract_address": "0x3825E5c0BaE86971c1Ec29947fBeD0D9EEa7C088",
		"denom": "usdx",
		"name": "MyToken",
		"symbol": "USDX",
		"decimal": 6,
		"supply": "0"
	}, {
		"contract_address": "0x91346e814f34462A59aF61ED0139Aa5312489c19",
		"denom": "uusdc",
		"name": "USDC Coin",
		"symbol": "USDC",
		"decimal": 6,
		"supply": "0"
	}, {
		"contract_address": "0x09F9629a56B179c5977485ba1d74b98c00300bbB",
		"denom": "uusdt",
		"name": "Tether USD",
		"symbol": "USDT",
		"decimal": 6,
		"supply": "70000000000"
	}]
}`

	bridgeTokensStruct := struct {
		BridgeTokens []struct {
			ContractAddress string `json:"contract_address"`
			Denom           string `json:"denom"`
			Name            string `json:"name"`
			Symbol          string `json:"symbol"`
			Decimal         uint64 `json:"decimal"`
			Supply          string `json:"supply"`
		} `json:"bridge_tokens"`
	}{}
	err := json.Unmarshal([]byte(bridgeTokens), &bridgeTokensStruct)
	if err != nil {
		panic(err)
	}

	for _, bt := range bridgeTokensStruct.BridgeTokens {
		bridgeToken := types.BridgeToken{
			ContractAddress: bt.ContractAddress,
			Denom:           bt.Denom,
			Name:            bt.Name,
			Symbol:          bt.Symbol,
			Decimal:         bt.Decimal,
		}
		ok := false
		bridgeToken.Supply, ok = sdk.NewIntFromString(bt.Supply)
		if !ok {
			panic("invalid supply")
		}
		k.SetBridgeToken(ctx, &bridgeToken)
	}
	return
}
