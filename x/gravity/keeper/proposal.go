package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

func (k Keeper) UpdateProposalRelayers(ctx sdk.Context, relayers []string) error {
	maxRelayers := k.GetMaxRelayers(ctx)
	if len(relayers) > int(maxRelayers) {
		return errorsmod.Wrapf(types.ErrInvalid, "relayer length must be less than or equal: %d", maxRelayers)
	}

	newRelayerMap := make(map[string]bool, len(relayers))
	for _, relayer := range relayers {
		newRelayerMap[relayer] = true
	}

	var unbondedRelayerList []types.Relayer
	totalPower, deleteTotalPower := sdkmath.ZeroInt(), sdkmath.ZeroInt()

	allRelayers := k.GetAllRelayers(ctx, false)
	proposalRelayer, _ := k.GetProposalRelayer(ctx)
	oldRelayerMap := make(map[string]bool, len(relayers))
	for _, relayer := range proposalRelayer.Relayers {
		oldRelayerMap[relayer] = true
	}

	for _, relayer := range allRelayers {
		if relayer.Online {
			totalPower = totalPower.Add(relayer.GetPower())
		}
		// relayer in new proposal
		if _, ok := newRelayerMap[relayer.RelayerAddress]; ok {
			continue
		}
		// relayer not in new proposal and relayer in old proposal
		if _, ok := oldRelayerMap[relayer.RelayerAddress]; ok {
			unbondedRelayerList = append(unbondedRelayerList, relayer)
			if relayer.Online {
				deleteTotalPower = deleteTotalPower.Add(relayer.GetPower())
			}
		}
	}

	maxChangePowerThreshold := types.AttestationProposalRelayerChangePowerThreshold.Mul(totalPower).Quo(sdkmath.NewInt(int64(types.PowerBase)))
	if deleteTotalPower.GT(sdk.ZeroInt()) && deleteTotalPower.GTE(maxChangePowerThreshold) {
		return errorsmod.Wrapf(types.ErrMaxChangePowerLimitExceeded,
			"maxChangePowerThreshold: %s, deleteTotalPower: %s", maxChangePowerThreshold.String(), deleteTotalPower.String())
	}

	// update proposal relayer
	k.SetProposalRelayer(ctx, &types.ProposalRelayer{Relayers: relayers})
	for _, unbondedRelayer := range unbondedRelayerList {
		if err := k.UnbondedRelayerFromProposal(ctx, unbondedRelayer); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) UnbondedRelayerFromProposal(ctx sdk.Context, relayer types.Relayer) error {
	relayerAddress := sdk.MustAccAddressFromBech32(relayer.RelayerAddress)
	//if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, relayerAddress, sdk.NewCoins(sdk.NewCoin(params.BaseDenom, relayer.DelegateAmount))); err != nil {
	//	return nil
	//}
	relayer.Online = false
	k.SetRelayer(ctx, relayerAddress, relayer)
	return nil
}
