package keeper

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/gravity/types"
)

// AttestationHandler Handle is the entry point for Attestation processing.
//
//gocyclo:ignore
func (k Keeper) AttestationHandler(ctx sdk.Context, externalClaim types.ExternalClaim) error {
	switch claim := externalClaim.(type) {
	case *types.MsgSendToMeClaim:
		bridgeToken, _ := k.GetBridgeTokenByContract(ctx, claim.TokenContract)
		if bridgeToken == nil {
			return errorsmod.Wrap(types.ErrInvalid, "bridge token is not exist")
		}

		coin := sdk.NewCoin(bridgeToken.Denom, claim.Amount)
		receiveAddr, err := sdk.AccAddressFromBech32(claim.Receiver)
		if err != nil {
			return errorsmod.Wrap(types.ErrInvalid, "receiver address")
		}

		if err := k.bankKeeper.MintCoins(ctx, k.moduleName, sdk.NewCoins(coin)); err != nil {
			return errorsmod.Wrapf(err, "mint vouchers coins")
		}
		if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, receiveAddr, sdk.NewCoins(coin)); err != nil {
			return errorsmod.Wrap(err, "transfer vouchers")
		}
		// record supply so we can withdraw it later
		bridgeToken.Supply = bridgeToken.Supply.Add(claim.Amount)
		k.SetBridgeToken(ctx, bridgeToken)

	case *types.MsgSendToExternalClaim:
		k.OutgoingTxBatchExecuted(ctx, claim.TokenContract, claim.BatchNonce)

	case *types.MsgBridgeTokenClaim:
		// Check if it already exists
		isExist := k.HasBridgeToken(ctx, claim.TokenContract)
		if isExist {
			return errorsmod.Wrap(types.ErrInvalid, "bridge token is exist")
		}
		bridgeToken := types.BridgeToken{
			ContractAddress: claim.TokenContract,
			Denom:           strings.ToLower(claim.Symbol),
			Name:            claim.Name,
			Symbol:          claim.Symbol,
			Decimal:         claim.Decimals,
			Supply:          sdkmath.ZeroInt(),
		}
		k.SetBridgeToken(ctx, &bridgeToken)
		k.Logger(ctx).Info("add bridge token success", "symbol", claim.Symbol, "token", claim.TokenContract, "denom", bridgeToken.Denom)

	case *types.MsgRelayerSetUpdateClaim:
		observedRelayerSet := &types.RelayerSet{
			Nonce:   claim.RelayerSetNonce,
			Members: claim.Members,
		}
		// check the contents of the validator set against the store
		if claim.RelayerSetNonce != 0 {
			trustedRelayerSet := k.GetRelayerSet(ctx, claim.RelayerSetNonce)
			if trustedRelayerSet == nil {
				ctx.Logger().Error("Received attestation for a relayer set which does not exist in store", "relayerSetNonce", claim.RelayerSetNonce, "claim", claim)
				return errorsmod.Wrapf(types.ErrInvalid, "attested relayerSet (%v) does not exist in store", claim.RelayerSetNonce)
			}

			// overwrite the height, since it's not part of the claim
			observedRelayerSet.Height = trustedRelayerSet.Height
			if _, err := trustedRelayerSet.Equal(observedRelayerSet); err != nil {
				panic(fmt.Sprintf("Potential bridge highjacking: observed relayerSet (%+v) does not match stored relayerSet (%+v)! %s", observedRelayerSet, trustedRelayerSet, err.Error()))
			}
		}
		k.SetLastObservedRelayerSet(ctx, observedRelayerSet)
	default:
		return errorsmod.Wrapf(types.ErrInvalid, "event type: %s", claim.GetType())
	}
	return nil
}
