package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/gravity/types"
)

// claimHandlerCommon is an internal function that provides common code for processing claims once they are
// translated from the message to the Ethereum claim interface
func (s MsgServer) claimHandlerCommon(ctx sdk.Context, msg types.ExternalClaim) (err error) {
	bridgerAddr := msg.GetClaimer()
	if err := s.checkIsRelayer(ctx, bridgerAddr); err != nil {
		return err
	}

	// Add the claim to the store
	if _, err := s.Attest(ctx, bridgerAddr, msg); err != nil {
		return err
	}

	// Emit the handle message event
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeySender, bridgerAddr.String()),
	))

	return nil
}

func (s MsgServer) confirmHandlerCommon(ctx sdk.Context, signatureAddr, signature string, checkpoint []byte) (sdk.AccAddress, error) {
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalid, "signature decoding")
	}

	relayerAddr, found := s.GetRelayerByExternalAddress(ctx, signatureAddr)
	if !found {
		return nil, types.ErrNotFoundRelayer
	}

	relayer, found := s.GetRelayer(ctx, relayerAddr)
	if !found {
		return nil, types.ErrNotFoundRelayer
	}

	if relayer.ExternalAddress != signatureAddr {
		return nil, errorsmod.Wrapf(types.ErrInvalid, "got %s, expected %s", signatureAddr, relayer.ExternalAddress)
	}
	if err = types.ValidateEthereumSignature(checkpoint, sigBytes, relayer.ExternalAddress); err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with checkpoint %s found %s", relayer.ExternalAddress, hex.EncodeToString(checkpoint), signature))
	}
	return relayerAddr, nil
}
