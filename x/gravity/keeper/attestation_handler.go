package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/openmetaearth/me-hub/utils"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

// AttestationHandler Handle is the entry point for Attestation processing.
//
//gocyclo:ignore
func (k Keeper) AttestationHandler(ctx sdk.Context, externalClaim types.ExternalClaim) error {
	switch claim := externalClaim.(type) {
	case *types.MsgSendToMeClaim:
		bridgeToken, err := k.GetBridgeTokenByContract(ctx, claim.TokenContract)
		if err != nil {
			return errorsmod.Wrapf(types.ErrInvalid, "bridge token does not exist: %s", claim.TokenContract)
		}

		receiveAddr, err := sdk.AccAddressFromBech32(claim.Receiver)
		if err != nil {
			return errorsmod.Wrap(types.ErrInvalid, "receiver address")
		}

		mintAmount := types.GetMintCoin(claim.Amount, claim.ChainName, bridgeToken)
		if err := k.bankKeeper.MintCoins(ctx, k.moduleName, sdk.NewCoins(mintAmount)); err != nil {
			return errorsmod.Wrapf(err, "mint vouchers coins")
		}
		if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, receiveAddr, sdk.NewCoins(mintAmount)); err != nil {
			return errorsmod.Wrap(err, "transfer vouchers")
		}
		// record supply so we can withdraw it later
		bridgeToken.Supply = bridgeToken.Supply.Add(mintAmount.Amount)
		k.SetBridgeToken(ctx, bridgeToken)

	case *types.MsgSendToExternalClaim:
		k.OutgoingTxBatchExecuted(ctx, claim.TokenContract, claim.BatchNonce)

	case *types.MsgBridgeTokenClaim:
		// Check if it already exists
		exists := k.HasBridgeToken(ctx, claim.TokenContract)
		if exists {
			return errorsmod.Wrap(types.ErrInvalid, "bridge token already exists")
		}
		denom := utils.GetDenom(claim.Symbol)
		if err := sdk.ValidateDenom(denom); err != nil {
			return errorsmod.Wrapf(types.ErrInvalid, "invalid denom derived from symbol: %v", err)
		}
		// This requires determining whether the same denom exists on the same chain, because different chains share the same denom.
		if existing, err := k.GetBridgeTokenByDenom(ctx, denom); err == nil {
			return errorsmod.Wrapf(
				types.ErrInvalid,
				"token %s already registered on %s chain (contract %s)",
				denom, claim.ChainName, existing.ContractAddress,
			)
		} else if !errorsmod.IsOf(err, types.ErrNotFound) {
			return errorsmod.Wrapf(err, "failed to look up bridge token by denom %s", denom)
		}

		bridgeToken := types.BridgeToken{
			ContractAddress: claim.TokenContract,
			Denom:           denom,
			Name:            claim.Name,
			Symbol:          claim.Symbol,
			Decimal:         claim.Decimals,
			Supply:          sdkmath.ZeroInt(),
		}
		k.SetBridgeToken(ctx, &bridgeToken)
		denomMeta, found := k.bankKeeper.GetDenomMetaData(ctx, bridgeToken.Denom)
		if !found {
			k.bankKeeper.SetDenomMetaData(ctx, banktypes.Metadata{
				Description: fmt.Sprintf("%s/%s", claim.ChainName, claim.TokenContract),
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    bridgeToken.Denom,
						Exponent: uint32(0),
					},
					{
						Denom:    claim.Symbol,
						Exponent: types.GetDecimals(claim),
					},
				},
				Base:    bridgeToken.Denom,
				Display: claim.Symbol,
				Name:    claim.Name,
				Symbol:  claim.Symbol,
				URI:     "",
				URIHash: "",
			})
		} else {
			denomMeta.Description = fmt.Sprintf("%s,%s/%s", denomMeta.Description, claim.ChainName, claim.TokenContract)
			k.bankKeeper.SetDenomMetaData(ctx, denomMeta)
		}
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
			match, err := trustedRelayerSet.Equal(observedRelayerSet)
			if err != nil {
				// this indicates that the members of the two sets are not equal
				return errorsmod.Wrapf(types.ErrInvalid, "potential bridge hijacking: observed relayerSet (%+v) does not match stored relayerSet (%+v)! %s", observedRelayerSet, trustedRelayerSet, err.Error())
			}
			if !match {
				return errorsmod.Wrapf(types.ErrInvalid, "potential bridge hijacking: observed relayerSet (%d) does not match stored relayerSet (%d)", observedRelayerSet.Nonce, trustedRelayerSet.Nonce)
			}
		}
		k.SetLastObservedRelayerSet(ctx, observedRelayerSet)
	default:
		return errorsmod.Wrapf(types.ErrInvalid, "event type: %s", claim.GetType())
	}
	return nil
}
